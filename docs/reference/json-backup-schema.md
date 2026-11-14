# JSON Backup Schema

## Scope
Structure of the `mdman-backup-<timestamp>.json` file produced by **Export Backup** and consumed by **Import Backup**. Version 1.0.

## Root Object

| Field | Type | Required | Description |
|---|---|---|---|
| `version` | string | Yes | Schema version (`"1.0"`). |
| `exportedAt` | string (ISO 8601) | Yes | Export timestamp (UTC). |
| `workspace` | string | Yes | Absolute path of **workspace** at export time. |
| `notes` | array[Note] | Yes | All **notes** in **workspace**. |
| `settings` | Settings | Yes | User settings at export time. |

## Note Object

| Field | Type | Required | Description |
|---|---|---|---|
| `path` | string | Yes | Relative path from **workspace** root (POSIX separators `/`). |
| `content` | string | Yes | Full Markdown content (UTF-8). |
| `modified` | integer | Yes | Last modified timestamp (milliseconds since Unix epoch). |

## Settings Object

| Field | Type | Required | Description |
|---|---|---|---|
| `themeApp` | string | Yes | `dark`, `light`, or `system`. |
| `themePreview` | string | Yes | `dark`, `light`, or `system`. |
| `locale` | string | Yes | `en` or `uk`. |
| `recentFileMax` | integer | Yes | Maximum **recent files** entries (5–100). |

## Example

```json
{
  "version": "1.0",
  "exportedAt": "2026-01-15T10:30:00.000Z",
  "workspace": "/home/user/notes",
  "notes": [
    {
      "path": "project/readme.md",
      "content": "# Project\n\nDescription here.",
      "modified": 1705312200000
    },
    {
      "path": "project/todo.md",
      "content": "- [ ] Task 1\n- [x] Task 2",
      "modified": 1705315800000
    }
  ],
  "settings": {
    "themeApp": "system",
    "themePreview": "dark",
    "locale": "en",
    "recentFileMax": 20
  }
}
```

## Version History

| Version | Date | Changes |
|---|---|---|
| 1.0 | 2026-01-15 | Initial schema. |

## Validation Rules
- `version` must be `"1.0"` (current). Future versions may add fields.
- `exportedAt` must be valid ISO 8601 UTC.
- `workspace` must be absolute path (informational; import uses current **workspace**).
- `notes[].path` must be relative, POSIX-style, no leading `/`.
- `notes[].content` must be valid UTF-8.
- `notes[].modified` must be non-negative integer.
- `settings` fields must match allowed values.

## Import Behavior
- **Import** ignores `workspace` field — restores into current **workspace**.
- **Import** overwrites all **notes** and **folders** in current **workspace**.
- **Import** applies `settings` (theme, locale, recent limit).
- **Import** does not restore `recent.json`, `window.json`, or `localStorage` layout.

## Related Topics
- [Import, Export, and Backup](../tasks/import-export-backup.md)
- [Concept: Workspaces and Notes](../concepts/workspaces-notes.md)