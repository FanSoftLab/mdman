# Write and Format Markdown

## Goal
Compose and format **notes** using Markdown syntax with the **toolbar**, keyboard shortcuts, or the **context menu**.

## Prerequisites
- A **note** is open in the **editor** (see [Create Your First Note](../getting-started/create-first-note.md)).
- The **editor** pane is visible (not in **view-only mode**).

## Procedure

### Basic Text Formatting
1. Select text in the **editor**.
2. Apply formatting using one of these methods:
   - **Toolbar**: Click **Bold** (`Ctrl+B`), **Italic** (`Ctrl+I`), **Strikethrough** (`Ctrl+~`), or **Inline Code** (`Ctrl+\``).
   - **Context Menu**: Right-click → select formatting option.
   - **Keyboard**: Press the shortcut keys.
3. The Markdown syntax wraps the selection. The **preview** updates.

### Headings
1. Place cursor on a line.
2. Apply heading level:
   - **Toolbar**: Click **Heading 1**–**Heading 3** (more in **More Options**).
   - **Keyboard**: `Ctrl+Alt+1` through `Ctrl+Alt+6`.
   - **Context Menu**: Right-click → **Heading** → choose level.
3. The line prefix changes to `#` through `######`.

### Lists
1. Place cursor at line start.
2. Start a list:
   - **Bullet List**: Click **Bullet List** (`Ctrl+.`) or type `- `.
   - **Numbered List**: Click **Numbered List** (`Ctrl+Shift+.`) or type `1. `.
   - **Task List**: Click **Task List** or type `- [ ] `.
3. Press `Enter` to continue the list. Press `Enter` twice to end.

### Links and Images
1. Select text for a link (or place cursor).
2. Click **Insert Link** (`Ctrl+K`) or **Insert Image**.
3. Enter the URL in the dialog.
4. For images, enter alt text and URL.

### Code Blocks
1. Place cursor on a blank line.
2. Click **Code Block** (`Ctrl+\``) or type `` ``` ``.
3. Optional: Type language name after backticks for syntax highlighting (e.g., `` ```python ``).

### Tables
1. Click **Insert Table**.
2. The **editor** inserts a template table with header and three columns.
3. Edit cell content. Use `|` to separate columns.

### Blockquotes and Horizontal Rules
1. **Blockquote**: Click **Blockquote** (`Ctrl+Q`) or type `> ` at line start.
2. **Horizontal Rule**: Click **Horizontal Rule** or type `---` on a blank line.

### Math Formulas
1. **Inline Math**: Click **Inline Math** or wrap text in `$...$`.
2. **Math Block**: Click **Math Block** or wrap lines in `$$...$$`.

## Expected Result
- The **editor** shows Markdown syntax.
- The **preview** renders formatted output in real time.
- Syntax highlighting appears in code blocks.
- Diagrams render in **preview** (see [Work with Diagrams](work-with-diagrams.md)).
- Math formulas render in **preview** (see [Work with Math](work-with-math.md)).

## Verification
- Toggle **Sync Scroll** in **toolbar** to verify **editor** and **preview** scroll together.
- Check **Table of Contents** (TOC) button — headings appear in the **TOC** panel.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Formatting not applied | No text selected, cursor not on line | Select text or position cursor correctly |
| Preview does not update | Large file, performance mode | Wait for debounce (50 ms); check **Sync Scroll** |
| Mermaid diagram shows error | Invalid Mermaid syntax | Check **Mermaid** syntax; see [Work with Diagrams](work-with-diagrams.md) |
| Math not rendering | Invalid LaTeX syntax | Check KaTeX syntax; see [Work with Math](work-with-math.md) |

## Limitations
- Maximum file size for smooth editing: ~1 MB.
- Very long lines may affect **editor** performance.
- Some GitHub Flavored Markdown (GFM) extensions may differ from GitHub rendering.

## Related Topics
- [Use the Toolbar](use-toolbar.md)
- [Use Keyboard Shortcuts](use-keyboard-shortcuts.md)
- [Use the Context Menu](use-context-menu.md)
- [Work with Diagrams](work-with-diagrams.md)
- [Work with Math](work-with-math.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)