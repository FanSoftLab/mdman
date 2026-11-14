# Changelog

## Format
Based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Categories: **Added**, **Changed**, **Fixed**, **Removed**, **Security**.

## Unreleased

## [1.0.0] - 2026-01-15
### Added
- Rename application to **mdman** (from mdmark).
- Persist pane layout state (`sidebarWidth`, `editorWidthPercent`, `tocWidth`, `showDesktopSidebar`) in `localStorage` with `settings.json` fallback.
- Startup theme flash fix: window background matches saved theme before frontend loads.
- **LICENSES.md** with full third-party attribution.
- Ukrainian (`uk`) localization for all UI strings.
- **Settings** panel: App theme, Preview theme, Language, Recent files limit.
- **Help** panel: Overview, Shortcuts, About, Licenses tabs with search.
- **Recent Files** panel with clear action.
- **Table of Contents** (TOC) panel with heading sync.
- **Find & Replace** dialog with case-sensitive, whole-word, regex options.
- **Print** dialog with light-theme rendering for diagrams and math.
- **Auto-hide UI** mode (hides sidebar/toolbar until pointer nears edge).
- **View-only mode** (preview only, no editor).
- **Word Wrap** toggle.
- **Sync Scroll** toggle (section-based mapping).
- Keyboard shortcuts for all major actions.
- Context menu in editor with formatting, headings, blocks, extras.
- Toolbar with formatting, diagram, math, view toggles.
- Tabbed interface (new, close, switch, dirty indicator).
- Drag-and-drop folder/file organization in sidebar.
- File operations: New File, New Folder, Rename, Delete, Download as Markdown.
- Open Folder / Open File / Save / Save As.
- Backup/Import (JSON) for full workspace migration.
- Mermaid diagrams (flowchart, sequence, class, state, pie, gantt, gitgraph, ER, journey).
- KaTeX math (inline `$...$`, block `$$...$$`).
- GitHub Flavored Markdown (GFM): tables, task lists, strikethrough, autolinks.
- Syntax highlighting (highlight.js) for code blocks.
- Dark/light/system themes for app and preview independently.
- Window state persistence (position, size, maximized).
- Single-instance lock with second-instance file opening.
- CLI file arguments: `mdman file1.md file2.md`.
- Cross-platform build via Docker (`build.sh`): Linux, Windows, macOS, AMD64, ARM64.
- Linux AppImage with bundled WebKitGTK helpers.
- Reproducible builds with version, time, commit, source embedding.

### Fixed
- AppImage black screen on Linux (bundled WebKitGTK helpers, `WebviewGpuPolicyOnDemand`).
- Startup theme flash (window background set before show).
- Settings persistence across restarts (dual `localStorage`/`settings.json`).
- Window restore on Wayland/NVIDIA (delayed `WindowShow`).
- Memory leaks in scroll sync (RAF cleanup, debounce).

## [0.9.0] - 2025-12-20
### Added
- Initial Markdown editor with CodeMirror 6.
- Live preview with GitHub Markdown CSS.
- Basic toolbar (bold, italic, headings, lists, code).
- Sidebar file tree.
- Dark theme only.

### Fixed
- Basic rendering and editing.

## [0.1.0] - 2025-11-01
### Added
- Project initialization.
- Wails v2 + Svelte 5 scaffold.
- Basic Go backend with file operations.

---

## Version Numbering
- **Major**: Breaking changes to data format, CLI, or plugin API.
- **Minor**: New features, UI additions, new settings.
- **Patch**: Bug fixes, performance, minor corrections.

## Release Process
1. Tag commit: `git tag v1.0.0`
2. Build: `./build.sh -v 1.0.0`
3. Artifacts in `./dist/`
4. Publish GitHub Release with artifacts.

## Related Topics
- [CLI and Build Reference](../reference/cli-build-reference.md)
- [Install and Run](../getting-started/install-and-run.md)