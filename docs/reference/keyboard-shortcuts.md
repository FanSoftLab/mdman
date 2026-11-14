# Keyboard Shortcuts Reference

## Scope
Complete list of keyboard shortcuts in mdman. Source: `frontend/src/lib/components/Editor.svelte` keymaps and `frontend/src/App.svelte` global handlers. macOS users: substitute `Cmd` for `Ctrl`.

## Global Shortcuts (Anywhere in App)

| Shortcut | Action | `en.ts` Key | Notes |
|---|---|---|---|
| `Ctrl+S` | Save current **note** | `shortcutSave` | Fails silently if **note** has no path |
| `Ctrl+Shift+S` | Save **note** as… | `saveAs` | Opens system save dialog |
| `Ctrl+O` | Open file dialog | `shortcutOpen` | Loads file into new **tab** |
| `Ctrl+T` | New **tab** (untitled **note**) | `shortcutNewTab` | Creates new **tab** |
| `Ctrl+W` | Close current **tab** | `shortcutCloseTab` | Prompts if unsaved |
| `Ctrl+F` | Find | `shortcutFind` | Opens find bar |
| `Ctrl+H` | Find & Replace | `shortcutFindReplace` | Opens find/replace dialog |
| `Ctrl+P` | Print | `print` | Opens print dialog |
| `Ctrl+H` (UI focused) | Toggle **Help** Panel | `helpTitle` | Dual-purpose: see Editor shortcuts |

## Editor Shortcuts (Editor Focused)

### Formatting (Toggle on Selection)

| Shortcut | Action | Markdown Syntax | CodeMirror Binding |
|---|---|---|---|
| `Ctrl+B` | **Bold** | `**text**` | Custom: `toggleInlineFormat` |
| `Ctrl+I` | *Italic* | `_text_` | Custom: `toggleInlineFormat` |
| `Ctrl+~` | ~~Strikethrough~~ | `~~text~~` | Custom: `handleStrikethroughSelection` |
| `Ctrl+\`` | `Inline Code` / Code Block | `` `text` `` / `` ``` `` | Smart: single-line → inline, multi-line → block |
| `Ctrl+K` | Insert Link | `[text](url)` | Custom: `handleLinkSelection` |

### Headings (Line-Based)

| Shortcut | Action | Markdown Syntax |
|---|---|---|
| `Ctrl+Alt+1` | Heading 1 | `# ` |
| `Ctrl+Alt+2` | Heading 2 | `## ` |
| `Ctrl+Alt+3` | Heading 3 | `### ` |
| `Ctrl+Alt+4` | Heading 4 | `#### ` |
| `Ctrl+Alt+5` | Heading 5 | `##### ` |
| `Ctrl+Alt+6` | Heading 6 | `###### ` |

### Lists and Blocks (Line-Based)

| Shortcut | Action | Markdown Syntax |
|---|---|---|
| `Ctrl+.` | Bullet List | `- ` |
| `Ctrl+Shift+.` | Numbered List | `1. ` |
| `Ctrl+Q` | Blockquote | `> ` |

### Multi-Selection

| Shortcut | Action |
|---|---|
| `Ctrl+D` | Select next occurrence of selection/word |
| `Ctrl+Shift+L` | Select all occurrences of selection/word |

### Navigation and Editing

| Shortcut | Action | Source |
|---|---|---|
| `Ctrl+Z` | Undo | CodeMirror `historyKeymap` |
| `Ctrl+Y` / `Ctrl+Shift+Z` | Redo | CodeMirror `historyKeymap` |
| `Ctrl+A` | Select All | CodeMirror `defaultKeymap` |
| `Ctrl+X` | Cut | Native |
| `Ctrl+C` | Copy | Native |
| `Ctrl+V` | Paste | Native |
| `Ctrl+F` | Find | CodeMirror `searchKeymap` |
| `F3` / `Enter` | Find Next | CodeMirror `searchKeymap` |
| `Shift+F3` / `Shift+Enter` | Find Previous | CodeMirror `searchKeymap` |
| `Esc` | Close **Help** Panel / Context Menu | Custom |

### Help Overlay

| Shortcut | Action |
|---|---|
| `Ctrl+H` (Editor focused) | Toggle **Help** Overlay in **Editor** |

## Context Menu Shortcuts (Right-Click Menu)

The context menu shows these shortcuts next to each action (from `en.ts`):

| Action | Shortcut Label (`en.ts` Key) |
|---|---|
| **Bold** | `boldShortcut` (`Ctrl+B`) |
| **Italic** | `italicShortcut` (`Ctrl+I`) |
| **Strikethrough** | — (`Ctrl+~`) |
| **Inline Code** | — (`Ctrl+\``) |
| **Bullet List** | — (`Ctrl+.`) |
| **Numbered List** | — (`Ctrl+Shift+.` ) |
| **Heading 1–6** | — (`Ctrl+Alt+1`–`6`) |
| **Blockquote** | — (`Ctrl+Q`) |
| **Link** | — (`Ctrl+K`) |

## macOS Differences

| Windows/Linux | macOS |
|---|---|
| `Ctrl` | `Cmd` (⌘) |
| `Ctrl+Alt` | `Cmd+Option` (⌘⌥) |
| `Ctrl+Shift` | `Cmd+Shift` (⌘⇧) |

Example: `Ctrl+B` → `Cmd+B`, `Ctrl+Alt+1` → `Cmd+Option+1`.

## Version Notes

- Shortcuts defined in `Editor.svelte` custom keymap (highest precedence).
- CodeMirror default keymap filtered to remove conflicts (`Ctrl+B`, `Ctrl+I`, `Ctrl+\``).
- Global shortcuts in `App.svelte` `onKey` handler.
- `Ctrl+H` dual behavior: **Help** panel when UI focused, Find & Replace when **editor** focused.

## Related Topics
- [Use Keyboard Shortcuts](../tasks/use-keyboard-shortcuts.md)
- [Use the Toolbar](../tasks/use-toolbar.md)
- [Use the Context Menu](../tasks/use-context-menu.md)
- [Write and Format Markdown](../tasks/write-format-markdown.md)