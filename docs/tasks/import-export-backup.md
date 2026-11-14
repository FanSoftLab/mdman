# Import, Export, and Backup

## Goal
Create a **backup** of all **notes**, **folders**, and settings; restore from a **backup**; export individual **notes** as Markdown.

## Prerequisites
- A **workspace** is open (for export/import).
- **Editor** pane visible (for **Download as Markdown**).

## Procedure

### Export All Data (Backup)
1. Open **Help** panel (`Ctrl+H`) → **Licenses** tab (or use empty state).
2. Click **Download Backup** (or from empty state: **Import Backup** → select file → **Export**).
3. Browser downloads `mdman-backup-<timestamp>.json`.
4. Save this file for recovery or migration.

### Import Backup
1. Open **Help** panel → **Licenses** tab (or empty state).
2. Click **Import Backup**.
3. Select a `mdman-backup-*.json` file.
4. Confirm: **Import** replaces current **workspace** data (notes, folders, settings).
5. Application reloads with imported data.

### Export Single Note as Markdown
1. Open the **note** in the **editor**.
2. In **sidebar**, right-click the **note**.
3. Select **Download as Markdown**.
4. Browser downloads `<note-name>.md`.

### Import Markdown Files
1. Drag `.md` files into **sidebar** folder (not supported directly).
2. Alternative: Use **Open File** (`Ctrl+O`) to load, then **Save** into **workspace**.

## Expected Result
- **Backup**: JSON file contains all **notes**, **folder** structure, and settings.
- **Import**: **Workspace** matches backup state; **notes**, **folders**, settings restored.
- **Download as Markdown**: Single `.md` file with **note** content.

## Verification
- Open backup JSON in text editor — verify `notes[]`, `settings`, `version` fields.
- After import — verify **sidebar** shows imported **folders** and **notes**.
- After download — open `.md` file — content matches **note**.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| `Failed to import backup` | Invalid JSON or version mismatch | Check file is valid JSON; ensure version `1.0` |
| Import loses current data | Import replaces all data | Export current data first (create backup) |
| `Download as Markdown` fails | Browser blocks download | Allow downloads for site; check popup blocker |
| Imported notes missing | Backup missing `notes[]` array | Verify backup file structure |

## Limitations
- **Backup** includes only **notes** in **workspace** (not external files).
- **Import** is destructive — replaces all current **notes**, **folders**, settings.
- No selective import (all or nothing).
- **Download as Markdown** exports single **note** only.
- Settings imported: `themeApp`, `themePreview`, `locale`, `recentFileMax`.

## JSON Backup Schema Reference
See [Reference: JSON Backup Schema](../reference/json-backup-schema.md).

## Related Topics
- [Reference: JSON Backup Schema](../reference/json-backup-schema.md)
- [Concept: Workspaces and Notes](../concepts/workspaces-notes.md)
- [Troubleshooting: Import/Backup Fails](../troubleshooting/import-backup-fails.md)