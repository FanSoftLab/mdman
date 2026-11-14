# Settings Reference

## Scope
All user-configurable settings in mdman. Source: `app.go` (Go backend), `SettingsPanel.svelte` (UI), `en.ts` (labels). Settings persist in `settings.json` (Go) and `localStorage` (UI layout).

## Settings Categories

### App Theme
- **Setting Key:** `themeApp`
- **UI Label:** App theme (`themeApp` in `en.ts`)
- **Values:** `dark`, `light`, `system`
- **Default:** `system`
- **Description:** Controls the application window and UI chrome color scheme. `system` follows OS preference.
- **Persistence:** `settings.json` (Go), synced to `localStorage` on UI change.
- **Version Added:** 1.0.0
- **Applied At:** Startup (window background), runtime (CSS variables).

### Preview Theme
- **Setting Key:** `themePreview`
- **UI Label:** Preview theme (`themePreview` in `en.ts`)
- **Values:** `dark`, `light`, `system`
- **Default:** `system`
- **Description:** Controls the rendered Markdown **preview** color scheme. Independent of **app theme**.
- **Persistence:** `settings.json` (Go), synced to `localStorage`.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (CSS class on `<html>`).

### Language
- **Setting Key:** `locale`
- **UI Label:** Language (`language` in `en.ts`)
- **Values:** `en` (English), `uk` (Ukrainian)
- **Default:** `en` (or `uk` if browser locale is Ukrainian)
- **Description:** UI language for all labels, dialogs, and help content.
- **Persistence:** `settings.json` (Go), `localStorage` key `mdman.locale`.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (reactive via Svelte store).

### Recent Files Limit
- **Setting Key:** `recentFileMax`
- **UI Label:** Recent files to keep (`recentMax` in `en.ts`)
- **Values:** Integer, 5–100
- **Default:** 20
- **Description:** Maximum number of entries in **Recent Files** panel.
- **Persistence:** `settings.json` (Go).
- **Version Added:** 1.0.0
- **Applied At:** On limit change (trims existing list).

### Workspace
- **Setting Key:** `workspace`
- **UI Label:** (Not directly editable; set via **Open Folder**)
- **Values:** Absolute file system path (string)
- **Default:** Empty (no **workspace** open)
- **Description:** Root **folder** opened in **sidebar**. All **notes** must reside within.
- **Persistence:** `settings.json` (Go).
- **Version Added:** 1.0.0
- **Applied At:** Startup (restores last **workspace**), runtime (on **Open Folder**).

### View-Only Mode
- **Setting Key:** `viewOnlyMode`
- **UI Label:** View Only / Edit Mode (`viewOnly` / `editMode` in `en.ts`)
- **Values:** `true` (view-only), `false` (edit)
- **Default:** `false`
- **Description:** Hides **editor**; shows only **preview**. Toggled via **toolbar** button.
- **Persistence:** `settings.json` (Go), `localStorage` key `viewOnlyMode`.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (UI layout).

### Auto-Hide UI
- **Setting Key:** `autoHideUI`
- **UI Label:** Auto-hide UI (`autoHideUI` in `en.ts`)
- **Values:** `true`, `false`
- **Default:** `false`
- **Description:** Hides **sidebar** and **toolbar** until pointer nears window edge (10 px).
- **Persistence:** `localStorage` only (UI preference).
- **Version Added:** 1.0.0
- **Applied At:** Runtime (CSS classes).

### Word Wrap
- **Setting Key:** `wordWrap`
- **UI Label:** Word Wrap (`wordWrap` in `en.ts`)
- **Values:** `true`, `false`
- **Default:** `false`
- **Description:** Enables line wrapping in **editor** (CodeMirror `lineWrapping`).
- **Persistence:** `localStorage` only.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (CodeMirror extension).

### Sync Scroll
- **Setting Key:** `syncScrollEnabled`
- **UI Label:** Sync Scroll (`syncScroll` in `en.ts`)
- **Values:** `true`, `false`
- **Default:** `true`
- **Description:** Synchronizes scroll position between **editor** and **preview** using section-based mapping.
- **Persistence:** `localStorage` only.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (scroll handlers).

### Show Table of Contents
- **Setting Key:** `showTOC`
- **UI Label:** Table of Contents (`toc` / `tableOfContents` in `en.ts`)
- **Values:** `true`, `false`
- **Default:** `false`
- **Description:** Shows **TOC** panel on right with heading links.
- **Persistence:** `localStorage` key `showTOC` (primary), fallback to `settings.json`.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (UI state).

### Sidebar Visibility (Desktop)
- **Setting Key:** `showDesktopSidebar`
- **UI Label:** Show/Hide Sidebar (`showSidebar` / `hideSidebar` in `en.ts`)
- **Values:** `true`, `false`
- **Default:** `true`
- **Description:** Shows/hides **sidebar** on desktop. Toggled via **toolbar** button.
- **Persistence:** `localStorage` key `showDesktopSidebar` (primary), fallback to `settings.json`.
- **Version Added:** 1.0.0
- **Applied At:** Runtime (CSS width).

### Pane Widths (localStorage Only)
These settings persist only in `localStorage`, written on every drag:

| Key | Default | Range | Description |
|---|---|---|---|
| `sidebarWidth` | 250 | 150–500 | **Sidebar** width in pixels |
| `editorWidthPercent` | 50 | 20–80 | **Editor** width as percentage of **editor-preview** container |
| `tocWidth` | 250 | 150–400 | **TOC** panel width in pixels |

### Session State (Restored on Startup)
Persisted in `settings.json` under `session` key:

| Field | Description |
|---|---|
| `tabs[]` | Array of `{ id, title, path }` for open **tabs** |
| `activeTabId` | ID of active **tab** |

Restored on startup: **tabs** recreated, active **tab** loaded.

## Settings File Locations

| File | Path | Contents |
|---|---|---|
| `settings.json` | `<configDir>/mdman/settings.json` | All Go-backed settings (`themeApp`, `themePreview`, `locale`, `recentFileMax`, `workspace`, `viewOnlyMode`, `session`) |
| `recent.json` | `<configDir>/mdman/recent.json` | **Recent Files** list |
| `window.json` | `<configDir>/mdman/window.json` | Window geometry (position, size, maximized) |
| `localStorage` | Browser `localStorage` | UI layout (`sidebarWidth`, `editorWidthPercent`, `tocWidth`, `showTOC`, `showDesktopSidebar`, `viewOnlyMode`, `autoHideUI`, `wordWrap`, `syncScrollEnabled`, `mdman.locale`) |

`<configDir>` = `$XDG_CONFIG_HOME/mdman` (Linux), `~/Library/Application Support/mdman` (macOS), `%APPDATA%/mdman` (Windows).

## Settings UI (SettingsPanel.svelte)

The **Settings** panel exposes these controls:

| Control | Setting | Type |
|---|---|---|
| **App theme** dropdown | `themeApp` | `select` |
| **Preview theme** dropdown | `themePreview` | `select` |
| **Language** dropdown | `locale` | `select` |
| **Recent files to keep** | `recentFileMax` | `number` (min 5, max 100) |

Other settings (pane widths, **auto-hide UI**, **word wrap**, **sync scroll**) are controlled via **toolbar** toggles or direct manipulation (drag handles).

## Version-Specific Behavior

| Version | Change |
|---|---|
| 1.0.0 | Initial settings set |
| 1.1.0 | Added `themePreview` (separate from `themeApp`) |
| 1.2.0 | Added `autoHideUI`, `wordWrap`, `syncScrollEnabled` to `localStorage` |

## Related Topics
- [Configure Settings](../tasks/configure-settings.md)
- [Concept: Themes and Appearance](../concepts/themes-appearance.md)
- [Concept: Auto-hide UI](../concepts/auto-hide-ui.md)
- [Troubleshooting: Settings Not Persisting](../troubleshooting/settings-not-persisting.md)