# Organize Folders and Files

## Goal
Manage **folder** structure and **notes** within the **workspace** using the **sidebar**.

## Prerequisites
- A **workspace** is open (see [Open Your First Folder](../getting-started/open-first-folder.md)).
- **Sidebar** is visible.

## Procedure

### Create a Folder
1. In **sidebar**, right-click **workspace** root or an existing **folder**.
2. Select **New Folder**.
3. Type folder name. Press `Enter`.
4. New **folder** appears in tree.

### Create a Note in a Folder
1. In **sidebar**, right-click target **folder**.
2. Select **New File**.
3. Type file name (`.md` added automatically). Press `Enter`.
4. **Note** opens in **editor**.

### Rename Folder or Note
1. In **sidebar**, right-click item.
2. Select **Rename**.
3. Edit name. Press `Enter`.
4. **Note** updates in **editor** if open; **folder** updates in tree.

### Move Notes and Folders
- **Drag and drop** in **sidebar**:
  - Drag **note** onto **folder** → moves into that **folder**.
  - Drag **folder** onto another **folder** → nests it.
  - Drag to **workspace** root → moves to top level.

### Delete Folder or Note
1. In **sidebar**, right-click item.
2. Select **Delete**.
3. Confirm dialog:
   - **Folder**: "Delete folder '{name}' and all its contents?"
   - **Note**: "Delete file '{name}'?"
4. Click **Delete**. Item removed permanently.

### Expand/Collapse Folders
- Click **folder** chevron (▶/▼) to toggle children visibility.
- **Folders** with **notes** or sub**folders** show chevron.

### Sort Order
- **Folders** appear before **notes** at each level.
- Within each group: alphabetical, case-insensitive.
- Hidden files (starting with `.`) are never shown.

### Download Note as Markdown
1. In **sidebar**, right-click **note**.
2. Select **Download as Markdown**.
3. Browser downloads `<name>.md`.

## Expected Result
- **Sidebar** reflects **folder** hierarchy and **note** locations.
- Operations persist to file system immediately.
- **Editor** updates if affected **note** is open.

## Verification
- Verify **folder** appears in correct parent.
- Verify **note** opens with correct content after move/rename.
- Verify deleted items no longer appear in **sidebar** or file system.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| "Failed to create folder" | Permission denied | Check **workspace** permissions; run mdman with write access |
| "Failed to rename" | Name conflict or invalid chars | Choose unique name; avoid `/`, `\`, `:`, `*`, `?`, `"`, `<`, `>`, `|` |
| "Delete folder failed" | Files locked by other process | Close other editors; retry |
| Drag-drop does nothing | Dropped on invalid target | Drop on **folder** or **workspace** root only |
| Hidden files not shown | By design | Hidden files (`.name`) excluded from **sidebar** |

## Limitations
- Operations apply to file system immediately (no undo in **sidebar**).
- Maximum path length OS-dependent.
- Symbolic links followed; target must be inside **workspace**.
- No batch operations (multi-select not supported in **sidebar**).

## Related Topics
- [Create Your First Note](../getting-started/create-first-note.md)
- [Import, Export, and Backup](../tasks/import-export-backup.md)
- [Concept: Workspaces and Notes](../concepts/workspaces-notes.md)
- [Troubleshooting: Files Not Saving](../troubleshooting/files-not-saving.md)