# Import/Backup Fails

## Symptom
**Import Backup** shows error or does not restore data. **Export Backup** fails to download.

## Conditions
- **Help** panel → **Licenses** tab or empty state
- Valid or invalid `.json` backup file

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| Invalid JSON syntax | `json.Unmarshal` fails in `importBackup()` | Verified |
| Schema version mismatch | Backup `version` not `1.0` | Verified |
| Missing required fields | `notes[]`, `settings` absent | Verified |
| Browser download blocked | Browser popup blocker prevents `.json` download | Verified |
| File too large | Browser memory limit on large **workspace** | Hypothesis |

## Diagnostic Checks
1. **Export Backup** — if download fails, check browser download permissions.
2. **Import Backup** — open backup file in text editor; validate JSON syntax.
3. Check `version` field — must be `"1.0"`.
4. Check `notes[]` array exists.
5. Check `settings` object exists with required keys.

## Corrective Actions

### Action 1: Allow Downloads (Export Fails)
1. Browser address bar → site permissions → **Downloads** → **Allow**.
2. Retry **Export Backup**.

### Action 2: Validate and Repair Backup JSON (Import Fails)
1. Open backup file in editor (VS Code, `jq`).
2. Run `jq . backup.json` — validates syntax.
3. Fix missing fields:
   ```json
   {
     "version": "1.0",
     "exportedAt": "2026-01-15T10:30:00.000Z",
     "workspace": "/path/to/workspace",
     "notes": [],
     "settings": {
       "themeApp": "system",
       "themePreview": "system",
       "locale": "en",
       "recentFileMax": 20
     }
   }
   ```
4. Retry **Import Backup**.

### Action 3: Manual Restore (If Import Fails)
1. Create new **workspace** folder.
2. Extract `notes[]` from backup — write each `path`/`content` as `.md` files.
3. Copy `settings` values to **Settings** panel manually.

## Expected Result After Correction
- **Export Backup** downloads `mdman-backup-<timestamp>.json`.
- **Import Backup** restores **notes**, **folders**, settings; app reloads.

## Escalation
If the problem persists:
- Collect: Backup file (sanitized), browser console errors, mdman version.
- Report at [GitHub Issues](https://github.com/FanSoftLab/mdman/issues).

## Prevention
- Export backup regularly (weekly).
- Store backups outside **workspace** (cloud, external drive).
- Test import periodically in a test **workspace**.

## Related Topics
- [Import, Export, and Backup](../tasks/import-export-backup.md)
- [Reference: JSON Backup Schema](../reference/json-backup-schema.md)
- [Concept: Workspaces and Notes](../concepts/workspaces-notes.md)