# Settings Not Persisting

## Symptom
Changes in **Settings** panel (theme, language, recent limit) revert after restart.

## Conditions
- Any platform
- **Settings** panel used or **toolbar** toggles

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| `settings.json` write permission denied | `os.WriteFile` fails silently in `saveSettings()` | Verified |
| Config directory not writable | `os.MkdirAll` or `os.WriteFile` fails | Verified |
| `localStorage` quota exceeded | Browser storage limit (rare) | Hypothesis |
| Dual persistence conflict | `localStorage` (UI) vs `settings.json` (Go) priority | Verified (code) |

## Diagnostic Checks
1. Change a setting (e.g., **App theme**).
2. Restart mdman.
3. Check if setting reverted.
4. Inspect config directory:
   - Linux: `~/.config/mdman/settings.json`
   - macOS: `~/Library/Application Support/mdman/settings.json`
   - Windows: `%APPDATA%\mdman\settings.json`
5. Verify file exists and contains new value.

## Corrective Actions

### Action 1: Fix Config Directory Permissions
1. Locate config directory (see above).
2. Ensure user owns it:
   ```bash
   sudo chown -R $USER:<group> <config-dir>
   chmod -R u+rw <config-dir>
   ```
3. Restart mdman; change setting; restart again.

### Action 2: Clear localStorage (If Corrupted)
1. Open browser DevTools (F12) → Application → Local Storage → `mdman` keys.
2. Delete all `mdman` keys.
3. Restart mdman; reconfigure settings.

### Action 3: Understand Persistence Priority
- **Go-backed settings** (`themeApp`, `themePreview`, `locale`, `recentFileMax`, `workspace`, `viewOnlyMode`, `session`): `settings.json` is source of truth.
- **UI-only settings** (`sidebarWidth`, `editorWidthPercent`, `tocWidth`, `showTOC`, `showDesktopSidebar`, `autoHideUI`, `wordWrap`, `syncScrollEnabled`): `localStorage` is source of truth (written on every change).
- On startup: `localStorage` wins for UI settings; `settings.json` is fallback.

## Expected Result After Correction
- Settings persist across restarts.
- **Settings** panel shows saved values.

## Escalation
If the problem persists:
- Collect: Config directory path, `ls -la` output, `settings.json` content, `localStorage` dump.
- Report at [GitHub Issues](https://github.com/FanSoftLab/mdman/issues).

## Prevention
- Run mdman as same user consistently.
- Do not manually edit `settings.json` while app runs.
- Avoid clearing browser data for mdman origin.

## Related Topics
- [Configure Settings](../tasks/configure-settings.md)
- [Reference: Settings Reference](../reference/settings-reference.md)
- [Concept: Themes and Appearance](../concepts/themes-appearance.md)