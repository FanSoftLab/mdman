# Use the Context Menu

## Goal
Access formatting and editing actions via right-click in the **editor**.

## Prerequisites
- A **note** is open in the **editor**.
- The **editor** pane has focus.

## Procedure

### Open the Context Menu
1. Right-click anywhere in the **editor** pane.
2. The context menu appears at the pointer position.
3. If text is selected, the menu shows selection-aware actions.

### Menu Structure
The menu organizes actions into categories:

#### Standard Actions (Always Visible)
- **Cut** (`Ctrl+X`): Remove selection to clipboard.
- **Copy** (`Ctrl+C`): Copy selection to clipboard.
- **Paste** (`Ctrl+V`): Insert clipboard content.
- **Select All** (`Ctrl+A`): Select entire **note**.
- **Find** (`Ctrl+F`): Open find dialog.
- **Find & Replace** (`Ctrl+H`): Open find/replace dialog.

#### Formatting (Selection-Aware)
Appears when text is selected or cursor on a line:
- **Bold** (`Ctrl+B`): Toggle `**...**`.
- **Italic** (`Ctrl+I`): Toggle `_..._`.
- **Strikethrough** (`Ctrl+~`): Toggle `~~...~~`.
- **Inline Code** (`Ctrl+\``): Toggle `` `...` ``.

#### Heading (Line-Based)
- **Heading 1**–**Heading 6** (`Ctrl+Alt+1`–`6`): Set heading level on current line or selected lines.

#### Block (Line-Based)
- **Code Block** (`Ctrl+\``): Insert fenced code block.
- **Task (unchecked)**: Insert `- [ ] `.
- **Blockquote** (`Ctrl+Q`): Insert `> `.
- **Horizontal Rule**: Insert `---`.

#### Extra
- **Link** (`Ctrl+K`): Wrap selection in `[...](url)` or insert template.
- **Image**: Insert `![alt](url)` template.
- **Table**: Insert Markdown table template.

#### Other
- **Inline Math**: Wrap selection in `$...$`.
- **Mermaid Diagram**: Insert Mermaid code block template.

### Submenus
Hover **Heading**, **Block**, **Extra**, or **Other** to reveal submenu with all options.

### Smart Behavior
- **Formatting toggles**: If selection already wrapped, removes markers.
- **Line-based actions**: Apply to current line if no selection; to all selected lines if multi-line selection.
- **Code Block**: Single-line selection → inline code; multi-line → fenced block.
- **Heading**: Same level on selected line → removes heading; different level → changes level.

## Expected Result
- Menu appears at click position.
- Actions apply immediately.
- Menu closes after action or clicking outside.

## Verification
- Right-click with no selection — all categories visible.
- Right-click with word selected — formatting actions wrap the word.
- Right-click with multi-line selection — line-based actions apply to all lines.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Menu does not appear | Clicked outside **editor** (e.g., **preview**, **sidebar**) | Click inside **editor** pane |
| Menu clipped at screen edge | Click near right/bottom window edge | Menu auto-repositions; if not, click closer to center |
| Submenu does not open | Pointer moved too far | Hover category label directly; submenu opens on hover |

## Limitations
- Context menu only in **editor** pane (not **preview**).
- Some actions duplicate **toolbar** and shortcuts.
- Native browser context menu may appear on some elements (e.g., images in **preview**).

## Related Topics
- [Write and Format Markdown](write-format-markdown.md)
- [Use the Toolbar](use-toolbar.md)
- [Use Keyboard Shortcuts](use-keyboard-shortcuts.md)
- [Reference: Keyboard Shortcuts](../reference/keyboard-shortcuts.md)