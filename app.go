// Package main implements the Go backend for mdman, a privacy-first desktop
// Markdown editor built with Wails v2 (Go + system webview).
//
// The backend provides:
//   - Workspace and file management (open, create, rename, delete, move)
//   - Settings persistence (theme, language, recent files limit)
//   - Window state persistence (position, size, maximized)
//   - Media file resolution and caching (LRU, mtime-validated) for local images
//   - Recent files tracking
//   - Session persistence (open tabs, scroll positions)
//   - Single-instance enforcement with second-instance file opening
//
// The frontend (Svelte 5 + CodeMirror 6) handles all UI, editing, preview
// rendering, Mermaid diagrams, and KaTeX math. Communication is via Wails
// bind: Go methods are callable from JavaScript as `window.go.main.App.Method()`.
//
// Build with: `wails build` (or `./build.sh` for cross-platform Docker builds).
// Run dev: `wails dev` (live reload).
package main

import (
	"container/list"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// maxMediaBytes caps files that get read and base64-encoded for the webview.
const maxMediaBytes = 32 << 20 // 32 MiB

// mediaCache caches resolved data URLs for media files (LRU, mtime-validated).
// Without it, every preview re-render re-read every image from disk and
// re-encoded it to base64 — several copies of the image bytes per render.
type mediaCacheEntry struct {
	path    string
	dataURL string
	size    int64
	modNs   int64
}

type mediaCache struct {
	mu         sync.Mutex
	entries    map[string]*list.Element // abs path -> element (front = most recent)
	order      *list.List
	totalBytes int64
	maxEntries int
	maxBytes   int64
}

func newMediaCache(maxEntries int, maxBytes int64) *mediaCache {
	return &mediaCache{
		entries:    make(map[string]*list.Element),
		order:      list.New(),
		maxEntries: maxEntries,
		maxBytes:   maxBytes,
	}
}

func (c *mediaCache) get(path string, fi os.FileInfo) (string, bool) {
	if c == nil {
		return "", false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.entries[path]
	if !ok {
		return "", false
	}
	entry, ok := el.Value.(*mediaCacheEntry)
	if !ok || entry == nil {
		// Should never happen; drop the corrupt slot instead of panicking.
		c.order.Remove(el)
		delete(c.entries, path)
		return "", false
	}
	if entry.size != fi.Size() || entry.modNs != fi.ModTime().UnixNano() {
		// file changed on disk — drop stale entry
		c.order.Remove(el)
		delete(c.entries, path)
		c.totalBytes -= entry.size
		return "", false
	}
	c.order.MoveToFront(el)
	return entry.dataURL, true
}

func (c *mediaCache) put(path, dataURL string, size int64, modNs int64) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.entries[path]; ok {
		if old, ok := el.Value.(*mediaCacheEntry); ok {
			c.totalBytes -= old.size
		}
		c.order.Remove(el)
		delete(c.entries, path)
	}
	entry := &mediaCacheEntry{path: path, dataURL: dataURL, size: size, modNs: modNs}
	c.entries[path] = c.order.PushFront(entry)
	c.totalBytes += size
	// Evict even if this single entry alone exceeds maxBytes — only the very
	// first (empty-cache) insert is exempt, so the cache can never grow unbounded.
	for c.order.Len() > 1 && (c.order.Len() > c.maxEntries || c.totalBytes > c.maxBytes) {
		oldest := c.order.Back()
		if oldest == nil {
			break
		}
		evicted, ok := oldest.Value.(*mediaCacheEntry)
		c.order.Remove(oldest)
		if !ok {
			continue
		}
		delete(c.entries, evicted.path)
		c.totalBytes -= evicted.size
	}
}

type WindowState struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Maximised bool `json:"maximised"`
}

type RecentFile struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	OpenedAt int64  `json:"openedAt"`
}

type FileNode struct {
	Name        string     `json:"name"`
	Path        string     `json:"path"`
	IsDir       bool       `json:"isDir"`
	Size        int64      `json:"size,omitempty"`
	Modified    int64      `json:"modified,omitempty"`
	Children    []FileNode `json:"children,omitempty"`
	HasChildren bool       `json:"hasChildren,omitempty"`
}

type App struct {
	ctxMu         sync.Mutex
	ctx           context.Context
	configDir     string
	workspace     string
	settings      map[string]interface{}
	settingsMu    sync.RWMutex
	recentFiles   []RecentFile
	recentFilesMu sync.RWMutex
	recentFileMax int
	windowStateMu sync.Mutex
	lastWindow    WindowState
	lastNormalPos struct {
		X, Y  int
		valid bool
	}
	windowLoaded bool
	pendingMu    sync.Mutex
	pendingFiles []string
	media        *mediaCache
}

// clampRecentFileMax converts a settings float64 to a sane recent-file limit,
// rejecting non-finite, sub-normal, and out-of-range values instead of
// truncating them into a bogus int.
func clampRecentFileMax(v float64) int {
	if v < 1 || v > 1000 || v != v { // NaN guard
		return 0
	}
	return int(v + 0.5)
}

// clampRecentMax clamps the stored maximum so AddRecent can read it without
// additional locking beyond recentFilesMu.
func (a *App) clampRecentMax() {
	a.recentFilesMu.Lock()
	defer a.recentFilesMu.Unlock()
	if a.recentFileMax < 1 {
		a.recentFileMax = 1
	}
	if a.recentFileMax > 1000 {
		a.recentFileMax = 1000
	}
}

// NewApp constructs the application backend, initializes config directories,
// and loads persisted state (settings, recent files, window geometry).
//
// The config directory is $XDG_CONFIG_HOME/mdman (or ~/.config/mdman on Linux,
// ~/Library/Application Support/mdman on macOS, %APPDATA%/mdman on Windows).
// If the directory cannot be determined, "." is used as a fallback.
//
// Persisted files:
//   - settings.json: app theme, preview theme, language, recent file limit, workspace, session
//   - recent.json: list of recently opened files (max recentFileMax entries)
//   - window.json: last window position, size, maximized state
//
// Returns a fully initialized *App ready for Wails binding.
func NewApp() *App {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	configDir = filepath.Join(configDir, "mdman")
	_ = os.MkdirAll(configDir, 0755)

	a := &App{
		configDir:     configDir,
		settings:      make(map[string]interface{}),
		recentFiles:   []RecentFile{},
		recentFileMax: 20,
		media:         newMediaCache(32, 64<<20),
	}
	a.loadSettings()
	a.loadRecentFiles()
	a.loadWindowState()
	if v, ok := a.settings["workspace"].(string); ok && v != "" {
		a.workspace = v
	}
	if v, ok := a.settings["recentFileMax"].(float64); ok {
		if clamped := clampRecentFileMax(v); clamped > 0 {
			a.recentFileMax = clamped
		}
	}
	a.clampRecentMax()
	return a
}

// GetStartupTheme returns the persisted app theme ("dark" or "light") so the
// native window background can be set before the webview loads, preventing a
// visible flash on startup. Returns empty string if theme is "system" or unset.
func (a *App) GetStartupTheme() string {
	pref, _ := a.settings["themeApp"].(string)
	switch pref {
	case "dark":
		return "dark"
	case "light":
		return "light"
	default:
		return ""
	}
}

func (a *App) setContext(ctx context.Context) {
	a.ctxMu.Lock()
	a.ctx = ctx
	a.ctxMu.Unlock()
}

func (a *App) currentContext() context.Context {
	a.ctxMu.Lock()
	defer a.ctxMu.Unlock()
	return a.ctx
}

func (a *App) startup(ctx context.Context) {
	a.setContext(ctx)
}

func (a *App) DomReady(ctx context.Context) {
	a.setContext(ctx)
	// Do not touch window geometry here — GTK widget may not be ready yet.
	// Restore is scheduled from main after delays.
}

func (a *App) Shutdown(ctx context.Context) {
	a.persistWindowState()
	a.saveSettings()
	a.saveRecentFiles()
}

func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	a.persistWindowState()
	a.saveSettings()
	a.saveRecentFiles()
	return false
}

func (a *App) settingsPath() string { return filepath.Join(a.configDir, "settings.json") }
func (a *App) recentPath() string   { return filepath.Join(a.configDir, "recent.json") }
func (a *App) windowPath() string   { return filepath.Join(a.configDir, "window.json") }

func (a *App) loadSettings() {
	data, err := os.ReadFile(a.settingsPath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("read settings: %v", err)
		}
		return
	}
	if err := json.Unmarshal(data, &a.settings); err != nil {
		log.Printf("parse settings %s: %v (starting with defaults)", a.settingsPath(), err)
	}
}

func (a *App) saveSettings() {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		log.Printf("marshal settings: %v", err)
		return
	}
	if err := os.WriteFile(a.settingsPath(), data, 0644); err != nil {
		log.Printf("write settings: %v", err)
	}
}

// copySettingsValue deep-copies the nested map values that the frontend can
// store in settings, so returned snapshots never share mutable state with
// the app's internal map. Leaf scalar values are immutable by value.
func copySettingsValue(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, inner := range t {
			out[k] = copySettingsValue(inner)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, inner := range t {
			out[i] = copySettingsValue(inner)
		}
		return out
	default:
		return v
	}
}

// GetSettings returns a deep copy of all persisted settings.
// The returned map is safe to mutate without affecting internal state.
// Keys include: themeApp, themePreview, language, recentFileMax, workspace, session.
func (a *App) GetSettings() map[string]interface{} {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	out := make(map[string]interface{}, len(a.settings))
	for k, v := range a.settings {
		out[k] = copySettingsValue(v)
	}
	return out
}

// SaveSettings merges the provided settings into the persisted configuration
// and writes to disk. Only keys present in the input are updated; others
// retain their current values.
//
// Recognized keys with special handling:
//   - recentFileMax: clamped to range [1, 1000], NaN rejected
//   - workspace: validated as absolute directory path on SetWorkspace
//
// The settings are deep-copied to prevent external mutation of internal state.
// On failure to write, an error is logged but not returned (fire-and-forget
// persistence suitable for UI-driven saves).
func (a *App) SaveSettings(settings map[string]interface{}) {
	a.settingsMu.Lock()
	if a.settings == nil {
		a.settings = map[string]interface{}{}
	}
	for k, v := range settings {
		a.settings[k] = copySettingsValue(v)
	}
	if v, ok := settings["recentFileMax"].(float64); ok {
		if clamped := clampRecentFileMax(v); clamped > 0 {
			a.recentFileMax = clamped
		}
	}
	if v, ok := settings["workspace"].(string); ok {
		a.workspace = v
	}
	a.settingsMu.Unlock()
	a.saveSettings()
}

func (a *App) loadRecentFiles() {
	data, err := os.ReadFile(a.recentPath())
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("read recent files: %v", err)
		}
		return
	}
	var files []RecentFile
	if err := json.Unmarshal(data, &files); err != nil {
		log.Printf("parse recent files %s: %v (starting empty)", a.recentPath(), err)
		return
	}
	a.recentFiles = files
	a.clampRecentMax()
}

func (a *App) saveRecentFilesLocked() {
	data, err := json.MarshalIndent(a.recentFiles, "", "  ")
	if err != nil {
		log.Printf("marshal recent files: %v", err)
		return
	}
	if err := os.WriteFile(a.recentPath(), data, 0644); err != nil {
		log.Printf("write recent files: %v", err)
	}
}

func (a *App) saveRecentFiles() {
	a.recentFilesMu.RLock()
	defer a.recentFilesMu.RUnlock()
	a.saveRecentFilesLocked()
}

// ListRecent returns a copy of the recent files list, most recent first.
// The returned slice is safe to mutate without affecting internal state.
// Each entry contains the absolute path, base name, and open timestamp (ms since epoch).
func (a *App) ListRecent() []RecentFile {
	a.recentFilesMu.RLock()
	defer a.recentFilesMu.RUnlock()
	out := make([]RecentFile, len(a.recentFiles))
	copy(out, a.recentFiles)
	return out
}

// AddRecent records a file as recently opened, moving it to the top of the
// list. If the file already exists in the list, it is removed before re-adding
// to avoid duplicates. The list is truncated to recentFileMax entries.
// The path is stored as-is (absolute paths from dialogs, relative from CLI).
// No-op if path is empty.
func (a *App) AddRecent(path string) {
	if path == "" {
		return
	}
	name := filepath.Base(path)
	a.recentFilesMu.Lock()
	defer a.recentFilesMu.Unlock()
	max := a.recentFileMax
	a.recentFiles = slices.DeleteFunc(a.recentFiles, func(f RecentFile) bool {
		return f.Path == path
	})
	a.recentFiles = append([]RecentFile{{
		Path:     path,
		Name:     name,
		OpenedAt: time.Now().UnixMilli(),
	}}, a.recentFiles...)
	if len(a.recentFiles) > max {
		a.recentFiles = a.recentFiles[:max]
	}
	a.saveRecentFilesLocked()
}

// ClearRecent removes all entries from the recent files list and persists the
// empty list to disk.
func (a *App) ClearRecent() {
	a.recentFilesMu.Lock()
	a.recentFiles = nil
	a.saveRecentFilesLocked()
	a.recentFilesMu.Unlock()
}

func (a *App) loadWindowState() {
	data, err := os.ReadFile(a.windowPath())
	if err != nil {
		return
	}
	var st WindowState
	if err := json.Unmarshal(data, &st); err != nil {
		log.Printf("parse window state %s: %v (using defaults)", a.windowPath(), err)
		return
	}
	if st.Width >= 200 && st.Height >= 200 {
		a.windowStateMu.Lock()
		a.lastWindow = st
		a.windowLoaded = true
		a.windowStateMu.Unlock()
	}
}

// GetWindowState returns the last known window geometry with sensible defaults.
// If the saved width is < 900 or height < 600, defaults of 1280x800 are applied
// to ensure a usable initial window size on first run or after config loss.
func (a *App) GetWindowState() WindowState {
	a.windowStateMu.Lock()
	defer a.windowStateMu.Unlock()
	st := a.lastWindow
	if st.Width < 900 {
		st.Width = 1280
	}
	if st.Height < 600 {
		st.Height = 800
	}
	return st
}

func (a *App) persistWindowState() {
	ctx := a.currentContext()
	if ctx == nil {
		return
	}
	// Guard against GTK criticals when the native window is not ready.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("persistWindowState recovered: %v", r)
		}
	}()

	maximised := runtime.WindowIsMaximised(ctx)
	w, h := runtime.WindowGetSize(ctx)
	x, y := runtime.WindowGetPosition(ctx)
	if w < 200 || h < 200 {
		return
	}

	state := WindowState{
		Maximised: maximised,
		Width:     w,
		Height:    h,
		X:         x,
		Y:         y,
	}
	if maximised {
		a.windowStateMu.Lock()
		if a.lastNormalPos.valid {
			state.X = a.lastNormalPos.X
			state.Y = a.lastNormalPos.Y
			if a.lastWindow.Width > 0 {
				state.Width = a.lastWindow.Width
			}
			if a.lastWindow.Height > 0 {
				state.Height = a.lastWindow.Height
			}
		}
		a.windowStateMu.Unlock()
	} else {
		a.windowStateMu.Lock()
		a.lastNormalPos.X = x
		a.lastNormalPos.Y = y
		a.lastNormalPos.valid = true
		a.windowStateMu.Unlock()
	}

	a.windowStateMu.Lock()
	a.lastWindow = state
	a.windowStateMu.Unlock()

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Printf("marshal window state: %v", err)
		return
	}
	if err := os.WriteFile(a.windowPath(), data, 0644); err != nil {
		log.Printf("write window state: %v", err)
	}
}

func (a *App) restoreWindowState() {
	ctx := a.currentContext()
	if ctx == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("restoreWindowState recovered: %v", r)
		}
	}()

	a.windowStateMu.Lock()
	st := a.lastWindow
	loaded := a.windowLoaded
	a.windowStateMu.Unlock()
	if !loaded || st.Width < 200 || st.Height < 200 {
		return
	}

	// Apply size first; position only when not maximised.
	runtime.WindowSetSize(ctx, st.Width, st.Height)
	if st.Maximised {
		runtime.WindowMaximise(ctx)
		return
	}
	runtime.WindowUnmaximise(ctx)
	// Skip 0,0 only when it likely means "unset"; still apply multi-monitor negatives.
	runtime.WindowSetPosition(ctx, st.X, st.Y)
}

// GetWorkspace returns the currently open workspace (root folder) path, or empty
// string if no workspace has been set.
func (a *App) GetWorkspace() string {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	return a.workspace
}

// SetWorkspace sets the workspace root folder. The path is resolved to absolute
// and must exist as a directory. On success, the workspace is persisted to
// settings.json and the recent files list is NOT cleared (user may switch
// workspaces and return).
// Returns os.ErrNotExist if the path does not exist or is not a directory.
func (a *App) SetWorkspace(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fi, err := os.Stat(abs)
	if err != nil || !fi.IsDir() {
		return os.ErrNotExist
	}
	a.settingsMu.Lock()
	a.workspace = abs
	if a.settings == nil {
		a.settings = map[string]interface{}{}
	}
	a.settings["workspace"] = abs
	a.settingsMu.Unlock()
	a.saveSettings()
	return nil
}

// SelectWorkspace opens a native directory picker dialog and sets the chosen
// folder as the workspace. Returns the selected path on success, or an error
// if the dialog is cancelled or the workspace cannot be set.
// Returns os.ErrInvalid if the Wails context is not yet available.
func (a *App) SelectWorkspace() (string, error) {
	ctx := a.currentContext()
	if ctx == nil {
		return "", os.ErrInvalid
	}
	path, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{
		Title: "Open Folder",
	})
	if err != nil || path == "" {
		return "", err
	}
	if err := a.SetWorkspace(path); err != nil {
		return "", err
	}
	return path, nil
}

func isMarkdown(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".txt"
}

// ListFiles returns the directory tree for the workspace or a subdirectory.
// The `rel` parameter is a slash-separated path relative to the workspace root
// (empty string or "." means the workspace root).
//
// Only Markdown files (extensions .md, .markdown, .mdown, .txt) and directories
// are returned. Hidden files (starting with ".") are excluded. The result is
// sorted: directories first (alphabetically), then files (alphabetically).
// Each FileNode includes name, relative path, isDir flag, size (files only),
// modification time (ms since epoch), and HasChildren hint for lazy loading.
//
// Returns an error if the workspace is not set or the resolved path escapes
// the workspace boundary (path traversal protection).
func (a *App) ListFiles(rel string) ([]FileNode, error) {
	ws := a.GetWorkspace()
	if ws == "" {
		return []FileNode{}, nil
	}
	var root string
	if rel == "" || rel == "." {
		root = ws
	} else {
		// Containment check so ".." or absolute rel paths cannot escape the workspace.
		var err error
		root, err = a.resolvePath(rel)
		if err != nil {
			return nil, err
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var nodes []FileNode
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		full := filepath.Join(root, name)
		relPath, err := filepath.Rel(ws, full)
		if err != nil || relPath == ".." || strings.HasPrefix(relPath, ".."+string(os.PathSeparator)) {
			continue
		}
		relPath = filepath.ToSlash(relPath)
		info, _ := e.Info()
		node := FileNode{
			Name:  name,
			Path:  relPath,
			IsDir: e.IsDir(),
		}
		if info != nil {
			node.Modified = info.ModTime().UnixMilli()
			if !e.IsDir() {
				node.Size = info.Size()
			}
		}
		if e.IsDir() {
			sub, _ := os.ReadDir(full)
			for _, s := range sub {
				if strings.HasPrefix(s.Name(), ".") {
					continue
				}
				if s.IsDir() || isMarkdown(s.Name()) {
					node.HasChildren = true
					break
				}
			}
		} else if !isMarkdown(name) {
			continue
		}
		nodes = append(nodes, node)
	}
	slices.SortFunc(nodes, func(a, b FileNode) int {
		if a.IsDir != b.IsDir {
			if a.IsDir {
				return -1
			}
			return 1
		}
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
	return nodes, nil
}

func (a *App) resolvePath(rel string) (string, error) {
	ws := a.GetWorkspace()
	if ws == "" {
		return "", os.ErrInvalid
	}
	full := filepath.Join(ws, filepath.FromSlash(rel))
	abs, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if abs != ws && !strings.HasPrefix(abs, ws+string(os.PathSeparator)) {
		return "", os.ErrPermission
	}
	return abs, nil
}

// ReadFile loads a file's contents as a UTF-8 string. Relative paths are
// resolved within the workspace via resolvePath (path traversal protected).
// Absolute paths are accepted directly (e.g., from file dialogs or CLI args).
// On success, the file is added to the recent files list.
// Returns the file content or an error if the file cannot be read.
func (a *App) ReadFile(rel string) (string, error) {
	var full string
	if filepath.IsAbs(rel) {
		full = rel
	} else {
		var err error
		full, err = a.resolvePath(rel)
		if err != nil {
			return "", err
		}
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	a.AddRecent(full)
	return string(data), nil
}

// WriteFile saves content to a file. Relative paths are resolved within the
// workspace (path traversal protected); absolute paths are accepted directly
// (e.g., from Save As dialog). Parent directories are created as needed (0755).
// On success, the file is added to the recent files list.
// Returns an error if the path cannot be resolved or the write fails.
func (a *App) WriteFile(path string, content string) error {
	var full string
	if filepath.IsAbs(path) {
		full = path
	} else {
		var err error
		full, err = a.resolvePath(path)
		if err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		return err
	}
	a.AddRecent(full)
	return nil
}

// SaveFileAs opens a native Save As dialog with the suggested filename,
// writes the content to the chosen path, and adds it to recent files.
// If the user cancels, returns ("", nil). If suggestedName is empty,
// "untitled.md" is used. If the chosen path has no extension, ".md" is appended.
// Returns the saved file path on success, or an error if the write fails.
// Returns os.ErrInvalid if the Wails context is not yet available.
func (a *App) SaveFileAs(content string, suggestedName string) (string, error) {
	ctx := a.currentContext()
	if ctx == nil {
		return "", os.ErrInvalid
	}
	if suggestedName == "" {
		suggestedName = "untitled.md"
	}
	path, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		Title:           "Save As",
		DefaultFilename: suggestedName,
		Filters: []runtime.FileFilter{
			{DisplayName: "Markdown", Pattern: "*.md;*.markdown"},
			{DisplayName: "All files", Pattern: "*"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}
	if filepath.Ext(path) == "" {
		path += ".md"
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	a.AddRecent(path)
	return path, nil
}

// OpenFileDialog opens a native file picker for Markdown and text files,
// reads the selected file, and returns its path, name, and content.
// If the user cancels, returns (nil, nil). The file is added to recent files.
// Returns an error if the file cannot be read.
// Returns os.ErrInvalid if the Wails context is not yet available.
func (a *App) OpenFileDialog() (map[string]string, error) {
	ctx := a.currentContext()
	if ctx == nil {
		return nil, os.ErrInvalid
	}
	path, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Open File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Markdown", Pattern: "*.md;*.markdown;*.txt"},
			{DisplayName: "All files", Pattern: "*"},
		},
	})
	if err != nil || path == "" {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	a.AddRecent(path)
	return map[string]string{
		"path":    path,
		"name":    filepath.Base(path),
		"content": string(data),
	}, nil
}

// CreateFile creates a new file with the given content at the relative path
// within the workspace. Delegates to WriteFile (handles directory creation).
// Path traversal is prevented by resolvePath.
func (a *App) CreateFile(rel string, content string) error {
	return a.WriteFile(rel, content)
}

// CreateFolder creates a directory at the relative path within the workspace.
// Parent directories are created as needed (0755). Path traversal is prevented
// by resolvePath.
func (a *App) CreateFolder(rel string) error {
	full, err := a.resolvePath(rel)
	if err != nil {
		return err
	}
	return os.MkdirAll(full, 0755)
}

// DeletePath removes a file or directory (recursively) at the relative path
// within the workspace. Path traversal is prevented by resolvePath.
// Returns an error if the path cannot be resolved or the removal fails.
func (a *App) DeletePath(rel string) error {
	full, err := a.resolvePath(rel)
	if err != nil {
		return err
	}
	return os.RemoveAll(full)
}

// RenamePath moves/renames a file or directory from oldRel to newRel within
// the workspace. Both paths are resolved and validated against the workspace
// boundary. Uses os.Rename (atomic on same filesystem). Returns an error if
// either path cannot be resolved or the rename fails (e.g., cross-device).
func (a *App) RenamePath(oldRel, newRel string) error {
	oldFull, err := a.resolvePath(oldRel)
	if err != nil {
		return err
	}
	newFull, err := a.resolvePath(newRel)
	if err != nil {
		return err
	}
	return os.Rename(oldFull, newFull)
}

// GetSession returns a deep copy of the persisted session state (open tabs,
// scroll positions, active tab). Returns an empty map if no session is stored.
// The returned map is safe to mutate without affecting internal state.
func (a *App) GetSession() map[string]interface{} {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	if s, ok := a.settings["session"].(map[string]interface{}); ok {
		out := make(map[string]interface{}, len(s))
		for k, v := range s {
			out[k] = copySettingsValue(v)
		}
		return out
	}
	return map[string]interface{}{}
}

// SaveSession persists the session state (open tabs, scroll positions, active
// tab) to settings.json. The input map is deep-copied to prevent external
// mutation of internal state. Overwrites any existing session data.
func (a *App) SaveSession(session map[string]interface{}) {
	snapshot := make(map[string]interface{}, len(session))
	for k, v := range session {
		snapshot[k] = copySettingsValue(v)
	}
	a.settingsMu.Lock()
	if a.settings == nil {
		a.settings = map[string]interface{}{}
	}
	a.settings["session"] = snapshot
	a.settingsMu.Unlock()
	a.saveSettings()
}

// GetContext returns the current Wails runtime context, or nil if the app
// hasn't started up yet. Used by the frontend to emit events or call runtime
// functions that require a context.
func (a *App) GetContext() context.Context {
	return a.currentContext()
}

// SetPendingFiles stores file paths passed via CLI arguments or second-instance
// launch for later consumption by the frontend. Called from main.go before
// wails.Run. The list replaces any existing pending files.
func (a *App) SetPendingFiles(files []string) {
	a.pendingMu.Lock()
	a.pendingFiles = append([]string{}, files...)
	a.pendingMu.Unlock()
}

// ConsumePendingFiles returns and clears the pending file paths list.
// Called by the frontend on DomReady to open files passed via CLI or
// second-instance launch. Returns empty slice if no pending files.
func (a *App) ConsumePendingFiles() []string {
	a.pendingMu.Lock()
	defer a.pendingMu.Unlock()
	out := append([]string{}, a.pendingFiles...)
	a.pendingFiles = nil
	return out
}

// SnapshotWindow persists the current window geometry (position, size, maximized)
// to window.json. Called periodically from the frontend (on blur, resize) and
// on shutdown. Panic-safe: recovers and logs instead of crashing the UI thread.
func (a *App) SnapshotWindow() {
	// Called periodically from the frontend; never panic into the UI thread.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("SnapshotWindow recovered: %v", r)
		}
	}()
	a.persistWindowState()
}

// ResolveMediaURL resolves a media reference (image, video, etc.) from a
// Markdown file and returns a data: URL for the webview to display.
//
// The `mediaPath` can be:
//   - Absolute URL (http://, https://, data:, blob:, wails://) — returned as-is
//   - Absolute file path — used directly
//   - Relative path — resolved against the Markdown file's directory (or workspace
//     root if markdownFile is empty)
//
// The file is read, base64-encoded, and cached (LRU, mtime-validated) to avoid
// re-reading on every preview render. Files larger than 32 MiB are rejected
// to prevent OOM in the webview (base64 adds ~33% overhead).
//
// Supported MIME types by extension: png, jpg/jpeg, gif, webp, svg, bmp, ico.
// Other extensions default to application/octet-stream.
//
// Returns the data URL on success, or an error if the file cannot be read,
// is too large, or the path cannot be resolved.
func (a *App) ResolveMediaURL(markdownFile string, mediaPath string) (string, error) {
	mediaPath = strings.TrimSpace(mediaPath)
	if mediaPath == "" {
		return "", os.ErrInvalid
	}
	// leave absolute http(s), data:, blob: alone
	lower := strings.ToLower(mediaPath)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "data:") || strings.HasPrefix(lower, "blob:") ||
		strings.HasPrefix(lower, "wails://") {
		return mediaPath, nil
	}

	var full string
	if filepath.IsAbs(mediaPath) {
		full = mediaPath
	} else {
		base := markdownFile
		if base == "" {
			base = a.GetWorkspace()
		}
		if base == "" {
			return "", os.ErrNotExist
		}
		dir := base
		if fi, err := os.Stat(base); err == nil && !fi.IsDir() {
			dir = filepath.Dir(base)
		}
		// strip leading ./
		rel := strings.TrimPrefix(mediaPath, "./")
		full = filepath.Join(dir, filepath.FromSlash(rel))
	}

	fi, err := os.Stat(full)
	if err != nil {
		return "", err
	}
	// Cap decode size: a huge file would be fully read, base64-encoded (~1.33x)
	// and held in the cache; refuse instead of OOMing the webview.
	if fi.Size() > maxMediaBytes {
		return "", fmt.Errorf("media file %s is %d bytes, over the %d byte limit", full, fi.Size(), maxMediaBytes)
	}
	// Cache hit avoids re-reading + re-encoding the file on every preview re-render.
	if dataURL, ok := a.media.get(full, fi); ok {
		return dataURL, nil
	}

	data, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	mime := "application/octet-stream"
	ext := strings.ToLower(filepath.Ext(full))
	switch ext {
	case ".png":
		mime = "image/png"
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".gif":
		mime = "image/gif"
	case ".webp":
		mime = "image/webp"
	case ".svg":
		mime = "image/svg+xml"
	case ".bmp":
		mime = "image/bmp"
	case ".ico":
		mime = "image/x-icon"
	}
	enc := base64.StdEncoding.EncodeToString(data)
	dataURL := "data:" + mime + ";base64," + enc
	a.media.put(full, dataURL, fi.Size(), fi.ModTime().UnixNano())
	return dataURL, nil
}

// RestoreWindow restores the persisted window geometry (position, size, maximized)
// from window.json. Called from main.go after DomReady with delays to ensure
// the GTK widget is ready. Panic-safe: recovers and logs on failure.
func (a *App) RestoreWindow() {
	a.restoreWindowState()
}
