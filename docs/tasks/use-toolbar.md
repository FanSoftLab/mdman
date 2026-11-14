# Use the Toolbar

## Goal
Apply Markdown formatting and access editor actions using the **toolbar** buttons.

## Prerequisites
- A **note** is open in the **editor**.
- The **toolbar** is visible (not hidden by **auto-hide UI**).

## Procedure

### Formatting Buttons (Left Group)
1. **Bold** (`Ctrl+B`): Wraps selection in `**...**`.
2. **Italic** (`Ctrl+I`): Wraps selection in `_..._`.
3. **Strikethrough** (`Ctrl+~`): Wraps selection in `~~...~~`.
4. **Inline Code** (`Ctrl+\``): Wraps selection in `` `...` ``.
5. **Heading 1** (`Ctrl+Alt+1`): Adds `# ` prefix.
6. **Heading 2** (`Ctrl+Alt+2`): Adds `## ` prefix.
7. **Heading 3** (`Ctrl+Alt+3`): Adds `### ` prefix.
8. **Bullet List** (`Ctrl+.`): Adds `- ` prefix.
9. **Numbered List** (`Ctrl+Shift+.`): Adds `1. ` prefix.
10. **Task List**: Adds `- [ ] ` prefix.
11. **Insert Link** (`Ctrl+K`): Opens dialog for `[text](url)`.
12. **Insert Image**: Opens dialog for `![alt](url)`.
13. **Code Block** (`Ctrl+\``): Inserts `` ``` `` fenced block.
14. **Insert Table**: Inserts Markdown table template.
15. **Blockquote** (`Ctrl+Q`): Adds `> ` prefix.
16. **Horizontal Rule**: Inserts `---`.

### More Options Menu (Right Group)
Click **More Options** (three dots) to access:
1. **Insert Diagram** → **Basic Diagram**, **Flowchart**, **Sequence Diagram**, **Pie Chart**: Inserts Mermaid code block.
2. **Inline Math**: Wraps selection in `$...$`.
3. **Math Block**: Wraps selection in `$$...$$`.
4. **Print** (`Ctrl+P`): Opens print dialog with light-theme rendering.
5. **Auto-hide UI**: Toggles **auto-hide UI** mode.
6. **Word Wrap**: Toggles line wrapping in **editor**.
7. **Sync Scroll**: Toggles synchronized scrolling between **editor** and **preview**.
8. **Reset Layout**: Resets pane widths to defaults.

### View and File Actions (Rightmost Group)
1. **Table of Contents** (TOC): Toggles **TOC** panel on right.
2. **Recent Files**: Opens **Recent Files** panel.
3. **Settings**: Opens **Settings** panel.
4. **Help** (`Ctrl+H`): Opens **Help** panel.
5. **Find & Replace** (`Ctrl+H`): Opens find/replace dialog.
6. **Save** (`Ctrl+S`): Saves current **note**.
7. **Save As** (`Ctrl+Shift+S`): Saves current **note** with new name.
8. **View Toggle**: Switches between **edit mode** and **view-only mode**.
9. **New Tab** (`Ctrl+T`): Creates new untitled **tab**.
10. **Close Tab** (`Ctrl+W`): Closes current **tab**.

## Expected Result
- Each button performs its labeled action.
- Tooltips show keyboard shortcuts on hover.
- Disabled buttons (e.g., **TOC** when no headings) are inactive.

## Verification
- Hover each button — tooltip matches action label from `en.ts`.
- Press shortcut — same action occurs.
- **TOC** button enables when headings exist in **note**.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Toolbar hidden | **Auto-hide UI** enabled | Move pointer to top edge (10 px) to reveal |
| Button disabled | No valid context (e.g., no headings for TOC) | Add required content; button enables automatically |
| Shortcut conflicts | System or browser captures key | Use **toolbar** button instead; check OS shortcut settings |

## Limitations
- **Auto-hide UI** hides **toolbar** and **sidebar** until pointer nears edge.
- Some actions require **editor** focus (e.g., formatting shortcuts).
- **Print** uses light theme regardless of **app theme**.

## Related Topics
- [Write and Format Markdown](write-format-markdown.md)
- [Use Keyboard Shortcuts](use-keyboard-shortcuts.md)
- [Use the Context Menu](use-context-menu.md)
- [Concept: Auto-hide UI](../concepts/auto-hide-ui.md)
- [Settings Reference](../reference/settings-reference.md)