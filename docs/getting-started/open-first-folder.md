# Open Your First Folder

## Goal
Open a **workspace** (root folder) in the **sidebar** to browse and edit **notes**.

## Prerequisites
- mdman is running (see [Install and Run](install-and-run.md)).
- A folder exists on your file system that you want to use as your **workspace**.

## Procedure

1. Click **Open Folder** in the **toolbar** (folder icon) or press `Ctrl+O`.
2. In the system dialog, navigate to your desired folder.
3. Click **Select Folder** (or **Open** on some systems).
4. The **sidebar** displays the folder structure. **Notes** (Markdown files) and **folders** appear. Hidden files (starting with `.`) are excluded.

## Expected Result
- The **sidebar** shows the folder tree with the selected **workspace** as root.
- The **toolbar** displays the **workspace** path.
- The first **note** in the folder (alphabetically) opens in the **editor**.
- The **preview** renders the **note** content.

## Verification
- Verify the **sidebar** shows your folder structure.
- Verify a **note** opens in the **editor** pane.
- Verify the **preview** pane shows rendered Markdown.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| "Permission denied" when selecting folder | Insufficient file system permissions | Choose a folder you own; or run mdman with appropriate permissions |
| Folder opens but shows empty | No Markdown files in folder or subfolders | Create a **note** (see [Create Your First Note](create-first-note.md)) or add `.md` files |
| Sidebar does not update | UI refresh issue | Click **Open Folder** again; or restart mdman |

## Limitations
- Only one **workspace** can be open at a time.
- Symbolic links are followed; target must be readable.
- Network-mounted folders may have slower performance.

## Related Topics
- [Create Your First Note](create-first-note.md)
- [Organize with Folders & Files](../tasks/organize-folders-files.md)
- [Concept: Workspaces and Notes](../concepts/workspaces-notes.md)