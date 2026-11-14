# Reference Template

**Content Type:** Reference (Technical Details)
**Standard:** ISO/IEC/IEEE 26514 §2.8 + ASD-STE100
**Max sentence length:** 25 words
**Max sentences per paragraph:** 6

---

## [Reference Title: Noun Phrase]

*Example: "Keyboard Shortcuts", "Settings Reference", "CLI Options", "JSON Backup Schema"*

### Scope
One sentence: what this reference covers and for whom.

### [Section Name: e.g., "Keyboard Shortcuts"]

#### Table: [Category]

| Shortcut / Option / Key | Action / Purpose / Description | Context / Conditions | Notes |
|---|---|---|---|
| `Ctrl+S` | Save current note | Editor focused, note has path | Fails silently if no path |
| `Ctrl+Shift+S` | Save note as… | Editor focused | Opens save dialog |
| `Ctrl+O` | Open file dialog | Anywhere | Loads file into new tab |

*Use exact key combinations, setting keys, CLI flags from source code.*

### [Section Name: e.g., "Settings"]

#### Setting: `themeApp`
- **UI Label:** App theme (from `en.ts`: `themeApp`)
- **Values:** `dark`, `light`, `system`
- **Default:** `system`
- **Description:** Controls the application window and UI chrome color scheme.
- **Persistence:** Stored in `settings.json`, synced to `localStorage`.
- **Version Added:** 1.0.0

#### Setting: `themePreview`
- **UI Label:** Preview theme (from `en.ts`: `themePreview`)
- **Values:** `dark`, `light`, `system`
- **Default:** `system`
- **Description:** Controls the rendered Markdown preview color scheme.
- **Persistence:** Stored in `settings.json`, synced to `localStorage`.
- **Version Added:** 1.0.0

*Repeat for each setting. Use exact keys from `app.go` and `SettingsPanel.svelte`.*

### [Section Name: e.g., "CLI Options"]

#### `build.sh` Flags
| Flag | Values | Default | Description |
|---|---|---|---|
| `-mode` | `release`, `debug` | `release` | Build optimization level |
| `-os` | `linux`, `windows`, `darwin`, `all` | Host OS | Target operating system(s) |
| `-arch` | `amd64`, `arm64`, `all` | Host arch | Target architecture(s) |
| `-run-tests` | (boolean) | `false` | Run `go test` and `go vet` |
| `-no-cache` | (boolean) | `false` | Disable Docker build cache |
| `-no-appimage` | (boolean) | `false` | Skip Linux AppImage packaging |
| `-v`, `--version` | `x.y.z` | Git tag or `dev-SHA` | Application version string |

*Use exact flags from `build.sh` usage function.*

### [Section Name: e.g., "File Formats"]

#### JSON Backup Schema
```json
{
  "version": "1.0",
  "exportedAt": "2026-01-15T10:30:00Z",
  "workspace": "/home/user/notes",
  "notes": [
    {
      "path": "folder/note.md",
      "content": "# Note\nContent here.",
      "modified": 1705312200000
    }
  ],
  "settings": {
    "themeApp": "system",
    "themePreview": "system",
    "recentFileMax": 20
  }
}
```
- **Version:** 1.0
- **Encoding:** UTF-8
- **Fields:** All required unless marked optional.

### Version-Specific Behavior
| Version | Change |
|---|---|
| 1.0.0 | Initial release |
| 1.1.0 | Added `themePreview` setting |

### Related Topics
- [Link to concept topic]
- [Link to task topic]
- [Link to troubleshooting topic]

---

## Writing Rules (ASD-STE100)

- **Exact identifiers** from implementation (setting keys, CLI flags, JSON fields).
- **Purpose, values, defaults, constraints** for each entry.
- **Examples** match actual interface (test them).
- **Mark version-specific behavior** when it matters.
- **No inferred defaults** from examples — state explicitly.
- **Glossary terms** from `appendix/glossary.md` used consistently.
- **Abbreviations** defined at first use.

## Review Checklist

- [ ] All entries use exact identifiers from source
- [ ] Values, defaults, constraints verified
- [ ] Examples tested against current build
- [ ] Version notes accurate
- [ ] Sentences ≤ 25 words
- [ ] Glossary terms used correctly
- [ ] Cross-references valid