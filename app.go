package main

import (
	"container/list"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
	entry := el.Value.(*mediaCacheEntry)
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
		old := el.Value.(*mediaCacheEntry)
		c.totalBytes -= old.size
		c.order.Remove(el)
		delete(c.entries, path)
	}
	entry := &mediaCacheEntry{path: path, dataURL: dataURL, size: size, modNs: modNs}
	c.entries[path] = c.order.PushFront(entry)
	c.totalBytes += size
	for c.order.Len() > 1 && (c.order.Len() > c.maxEntries || c.totalBytes > c.maxBytes) {
		oldest := c.order.Back()
		if oldest == nil {
			break
		}
		evicted := oldest.Value.(*mediaCacheEntry)
		c.order.Remove(oldest)
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
	if v, ok := a.settings["recentFileMax"].(float64); ok && v > 0 {
		a.recentFileMax = int(v)
	}
	return a
}

// GetStartupTheme reports the theme the window background should use before
// the frontend loads, to avoid a colour flash on startup.
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

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DomReady(ctx context.Context) {
	a.ctx = ctx
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
		return
	}
	_ = json.Unmarshal(data, &a.settings)
}

func (a *App) saveSettings() {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(a.settingsPath(), data, 0644)
}

func (a *App) GetSettings() map[string]interface{} {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	out := make(map[string]interface{}, len(a.settings))
	for k, v := range a.settings {
		out[k] = v
	}
	return out
}

func (a *App) SaveSettings(settings map[string]interface{}) {
	a.settingsMu.Lock()
	if a.settings == nil {
		a.settings = map[string]interface{}{}
	}
	for k, v := range settings {
		a.settings[k] = v
	}
	if v, ok := settings["recentFileMax"].(float64); ok && v > 0 {
		a.recentFileMax = int(v)
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
		return
	}
	var files []RecentFile
	if json.Unmarshal(data, &files) == nil {
		a.recentFiles = files
	}
}

func (a *App) saveRecentFiles() {
	a.recentFilesMu.RLock()
	defer a.recentFilesMu.RUnlock()
	data, err := json.MarshalIndent(a.recentFiles, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(a.recentPath(), data, 0644)
}

func (a *App) ListRecent() []RecentFile {
	a.recentFilesMu.RLock()
	defer a.recentFilesMu.RUnlock()
	out := make([]RecentFile, len(a.recentFiles))
	copy(out, a.recentFiles)
	return out
}

func (a *App) AddRecent(path string) {
	if path == "" {
		return
	}
	name := filepath.Base(path)
	a.recentFilesMu.Lock()
	a.recentFiles = slices.DeleteFunc(a.recentFiles, func(f RecentFile) bool {
		return f.Path == path
	})
	a.recentFiles = append([]RecentFile{{
		Path:     path,
		Name:     name,
		OpenedAt: time.Now().UnixMilli(),
	}}, a.recentFiles...)
	if len(a.recentFiles) > a.recentFileMax {
		a.recentFiles = a.recentFiles[:a.recentFileMax]
	}
	a.recentFilesMu.Unlock()
	a.saveRecentFiles()
}

func (a *App) ClearRecent() {
	a.recentFilesMu.Lock()
	a.recentFiles = nil
	a.recentFilesMu.Unlock()
	a.saveRecentFiles()
}

func (a *App) loadWindowState() {
	data, err := os.ReadFile(a.windowPath())
	if err != nil {
		return
	}
	var st WindowState
	if json.Unmarshal(data, &st) == nil && st.Width >= 200 && st.Height >= 200 {
		a.windowStateMu.Lock()
		a.lastWindow = st
		a.windowLoaded = true
		a.windowStateMu.Unlock()
	}
}

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
	if a.ctx == nil {
		return
	}
	// Guard against GTK criticals when the native window is not ready.
	defer func() { recover() }()

	maximised := runtime.WindowIsMaximised(a.ctx)
	w, h := runtime.WindowGetSize(a.ctx)
	x, y := runtime.WindowGetPosition(a.ctx)
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
		return
	}
	_ = os.WriteFile(a.windowPath(), data, 0644)
}

func (a *App) restoreWindowState() {
	if a.ctx == nil {
		return
	}
	defer func() { recover() }()

	a.windowStateMu.Lock()
	st := a.lastWindow
	loaded := a.windowLoaded
	a.windowStateMu.Unlock()
	if !loaded || st.Width < 200 || st.Height < 200 {
		return
	}

	// Apply size first; position only when not maximised.
	runtime.WindowSetSize(a.ctx, st.Width, st.Height)
	if st.Maximised {
		runtime.WindowMaximise(a.ctx)
		return
	}
	runtime.WindowUnmaximise(a.ctx)
	// Skip 0,0 only when it likely means "unset"; still apply multi-monitor negatives.
	runtime.WindowSetPosition(a.ctx, st.X, st.Y)
}

func (a *App) GetWorkspace() string {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	return a.workspace
}

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

func (a *App) SelectWorkspace() (string, error) {
	if a.ctx == nil {
		return "", os.ErrInvalid
	}
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
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

func (a *App) ListFiles(rel string) ([]FileNode, error) {
	ws := a.GetWorkspace()
	if ws == "" {
		return []FileNode{}, nil
	}
	root := ws
	if rel != "" && rel != "." {
		root = filepath.Join(ws, rel)
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
		relPath, _ := filepath.Rel(ws, full)
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

func (a *App) SaveFileAs(content string, suggestedName string) (string, error) {
	if a.ctx == nil {
		return "", os.ErrInvalid
	}
	if suggestedName == "" {
		suggestedName = "untitled.md"
	}
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
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

func (a *App) OpenFileDialog() (map[string]string, error) {
	if a.ctx == nil {
		return nil, os.ErrInvalid
	}
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
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

func (a *App) CreateFile(rel string, content string) error {
	return a.WriteFile(rel, content)
}

func (a *App) CreateFolder(rel string) error {
	full, err := a.resolvePath(rel)
	if err != nil {
		return err
	}
	return os.MkdirAll(full, 0755)
}

func (a *App) DeletePath(rel string) error {
	full, err := a.resolvePath(rel)
	if err != nil {
		return err
	}
	return os.RemoveAll(full)
}

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

func (a *App) GetSession() map[string]interface{} {
	a.settingsMu.RLock()
	defer a.settingsMu.RUnlock()
	if s, ok := a.settings["session"].(map[string]interface{}); ok {
		return s
	}
	return map[string]interface{}{}
}

func (a *App) SaveSession(session map[string]interface{}) {
	a.settingsMu.Lock()
	if a.settings == nil {
		a.settings = map[string]interface{}{}
	}
	a.settings["session"] = session
	a.settingsMu.Unlock()
	a.saveSettings()
}

func (a *App) GetContext() context.Context {
	return a.ctx
}

func (a *App) SetPendingFiles(files []string) {
	a.pendingMu.Lock()
	a.pendingFiles = append([]string{}, files...)
	a.pendingMu.Unlock()
}

// ConsumePendingFiles returns and clears CLI / second-instance file paths.
func (a *App) ConsumePendingFiles() []string {
	a.pendingMu.Lock()
	defer a.pendingMu.Unlock()
	out := append([]string{}, a.pendingFiles...)
	a.pendingFiles = nil
	return out
}

// SnapshotWindow is called from the frontend periodically / on blur.
func (a *App) SnapshotWindow() {
	// Called periodically from the frontend; never panic into the UI thread.
	defer func() { recover() }()
	a.persistWindowState()
}

// ResolveMediaURL resolves a relative media path against the markdown file directory
// and returns a data URL so the webview can display local images.
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

// RestoreWindow is a public alias used from main after DomReady.
func (a *App) RestoreWindow() {
	a.restoreWindowState()
}
