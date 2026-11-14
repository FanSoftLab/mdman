# Use Keyboard Shortcuts

## Goal
Perform common actions quickly using keyboard shortcuts.

## Prerequisites
- A **note** is open in the **editor**.
- The **editor** or **preview** pane has focus.

## Procedure

### File Operations
| Shortcut | Action | Context |
|---|---|---|
| `Ctrl+S` | Save current **note** | Editor focused, **note** has path |
| `Ctrl+Shift+S` | Save **note** as… | Editor focused |
| `Ctrl+O` | Open file dialog | Anywhere |
| `Ctrl+T` | New **tab** (untitled **note**) | Anywhere |
| `Ctrl+W` | Close current **tab** | Anywhere |

### Editing
| Shortcut | Action | Context |
|---|---|---|
| `Ctrl+Z` | Undo | Editor focused |
| `Ctrl+Y` / `Ctrl+Shift+Z` | Redo | Editor focused |
| `Ctrl+A` | Select all | Editor focused |
| `Ctrl+X` | Cut | Editor focused, selection exists |
| `Ctrl+C` | Copy | Editor focused, selection exists |
| `Ctrl+V` | Paste | Editor focused |

### Formatting (Selection Required)
| Shortcut | Action | Markdown Syntax |
|---|---|---|
| `Ctrl+B` | **Bold** | `**text**` |
| `Ctrl+I` | *Italic* | `_text_` |
| `Ctrl+~` | ~~Strikethrough~~ | `~~text~~` |
| `Ctrl+\`` | `Inline Code` | `` `text` `` |
| `Ctrl+K` | Insert Link | `[text](url)` |

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
| `Ctrl+\`` | Code Block / Inline Code | `` ``` `` / `` ` `` |

### Navigation and UI
| Shortcut | Action | Context |
|---|---|---|
| `Ctrl+F` | Find | Anywhere |
| `Ctrl+H` | Find & Replace | Anywhere |
| `F3` / `Enter` | Find Next | Find dialog open |
| `Shift+F3` / `Shift+Enter` | Find Previous | Find dialog open |
| `Ctrl+P` | Print | Anywhere |
| `Ctrl+H` | Toggle **Help** Panel | Anywhere |

### Multi-Selection (CodeMirror)
| Shortcut | Action |
|---|---|
| `Ctrl+D` | Select next occurrence of current selection/word |
| `Ctrl+Shift+L` | Select all occurrences of current selection/word |

## Expected Result
- Each shortcut triggers its action immediately.
- Formatting shortcuts toggle: apply if absent, remove if present.
- Line-based shortcuts affect the current line or selection lines.

## Verification
- Open **Help** panel (`Ctrl+H`) → **Shortcuts** tab — all listed shortcuts work.
- Test each shortcut in context — action matches table.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Shortcut does nothing | Wrong pane focused (e.g., **preview** for formatting) | Click **editor** pane; retry |
| `Ctrl+H` opens Find & Replace instead of **Help** | **Editor** focused | Click **toolbar** or **sidebar**; press `Ctrl+H` |
| `Ctrl+D` deletes line | Browser default behavior | Use in **editor** only; CodeMirror captures it there |
| Shortcut conflicts with OS | Window manager captures key | Change OS shortcut; use **toolbar** button |

## Limitations
- Shortcuts require **editor** focus for text operations.
- `Ctrl+H` dual-purpose: **Help** panel when UI focused, Find & Replace when **editor** focused.
- macOS: Use `Cmd` instead of `Ctrl` for most shortcuts.

## Related Topics
- [Write and Format Markdown](write-format-markdown.md)
- [Use the Toolbar](use-toolbar.md)
- [Reference: Keyboard Shortcuts](../reference/keyboard-shortcuts.md)
- [Use the Context Menu](use-context-menu.md)