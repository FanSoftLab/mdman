# Configure Settings

## Goal
Change application behavior via the **Settings** panel.

## Prerequisites
- mdman running.
- **Settings** panel accessible (click **Settings** in **toolbar**).

## Procedure

### Open Settings
1. Click **Settings** (gear icon) in **toolbar**.
2. **Settings** panel opens on right.

### Change App Theme
1. In **Settings**, find **App theme** dropdown.
2. Select: `Dark`, `Light`, or `System`.
3. Change applies immediately.

### Change Preview Theme
1. In **Settings**, find **Preview theme** dropdown.
2. Select: `Dark`, `Light`, or `System`.
3. **Preview** pane updates.

### Change Language
1. In **Settings**, find **Language** dropdown.
2. Select: `English` or `Ukrainian`.
3. UI text updates immediately.

### Change Recent Files Limit
1. In **Settings**, find **Recent files to keep** input.
2. Enter number (5–100).
3. Press `Enter` or click outside.
4. Excess entries trimmed from **Recent Files** panel.

### Apply and Close
1. Click **Apply** (saves to `settings.json` and `localStorage`).
2. Click **Close** (X) or press `Esc` to dismiss panel.

## Expected Result
- Selected values persist across restarts.
- UI updates immediately for theme/language.
- **Recent Files** panel respects new limit.

## Verification
- Restart mdman — settings unchanged.
- Check `settings.json` in config directory — values match.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Settings revert after restart | `settings.json` not writable | Fix config directory permissions (see [Settings Not Persisting](../troubleshooting/settings-not-persisting.md)) |
| Language not fully translated | Missing `uk.ts` keys | Report missing keys; fallback to English |
| Theme not applied | CSS variable issue | Restart mdman; check `localStorage` |

## Limitations
- No custom theme colors.
- No font size setting (use OS scaling).
- No auto-save interval setting (debounced 50 ms).

## Related Topics
- [Reference: Settings Reference](../reference/settings-reference.md)
- [Concept: Themes and Appearance](../concepts/themes-appearance.md)
- [Concept: Auto-hide UI](../concepts/auto-hide-ui.md)
- [Troubleshooting: Settings Not Persisting](../troubleshooting/settings-not-persisting.md)