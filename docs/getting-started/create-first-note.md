# Create Your First Note

## Goal
Create a new **note** (Markdown file) in the open **workspace** and write content.

## Prerequisites
- A **workspace** is open (see [Open Your First Folder](open-first-folder.md)).
- The **editor** pane is visible (not in **view-only mode**).

## Procedure

### Create a Note in the Sidebar
1. In the **sidebar**, right-click the **workspace** root or a **folder**.
2. Select **New File** from the context menu.
3. Type a name (e.g., `my-note.md`). The `.md` extension is added automatically if omitted.
4. Press `Enter`. The new **note** opens in the **editor**.

### Create a Note via Toolbar
1. Click **New File** in the **toolbar** (plus icon) or press `Ctrl+T`.
2. A new untitled **note** opens in a new **tab**.
3. Click **Save** in the **toolbar** or press `Ctrl+S`.
4. In the save dialog, choose a location within the **workspace**.
5. Type a file name and click **Save**.

### Write Content
1. Type Markdown in the **editor** pane.
2. Use the **toolbar** buttons for formatting (see [Write & Format Markdown](../tasks/write-format-markdown.md)).
3. Use keyboard shortcuts (see [Keyboard Shortcuts Reference](../reference/keyboard-shortcuts.md)).
4. The **preview** pane updates in real time.

## Expected Result
- The new **note** appears in the **sidebar** tree.
- The **note** opens in a **tab** with the file name as title.
- The **editor** shows your typed content.
- The **preview** shows rendered HTML.
- The **tab** shows an unsaved indicator (`●`) until saved.

## Verification
- Verify the **note** appears in the **sidebar** under the correct **folder**.
- Verify the **tab** title matches the file name.
- Verify the **preview** renders your Markdown correctly.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| "Failed to create file" | Permission denied on target folder | Check folder permissions; choose a writable location |
| "File already exists" | Name conflict | Choose a different name |
| Unsaved changes lost on close | App closed before save | mdman prompts before closing unsaved **tabs**; click **Cancel** and save |

## Limitations
- **Notes** must be inside the open **workspace**.
- File names cannot contain `/` or `\`.
- Maximum path length depends on the operating system.

## Related Topics
- [Write & Format Markdown](../tasks/write-format-markdown.md)
- [Manage Tabs & Recent Files](../tasks/manage-tabs-recent.md)
- [Save and Export](../tasks/save-export-backup.md)