# Themes and Appearance

## Purpose
Customize the visual appearance of mdman: application chrome, preview rendering, and syntax highlighting.

## What It Is
mdman supports independent theming for the **app** (window chrome, sidebar, toolbar, editor) and **preview** (rendered Markdown). Each can be `dark`, `light`, or `system` (follows OS preference).

## Key Behaviors
- **App theme**: Controls editor background, sidebar, toolbar, dialogs. Applied at startup (window background) and runtime (CSS variables).
- **Preview theme**: Controls rendered Markdown colors (GitHub Markdown CSS). Independent of **app theme**.
- **System detection**: Uses `prefers-color-scheme` media query. Updates live on OS theme change.
- **Syntax highlighting**: CodeMirror theme matches **app theme** (GitHub Dark/Light).
- **Diagrams**: Mermaid inherits **preview theme** colors.
- **Math**: KaTeX inherits text color from **preview theme**.
- **Print**: Forces light theme for readability.

## Settings

| Setting | Key | Values | Default | Scope |
|---|---|---|---|---|
| App theme | `themeApp` | `dark`, `light`, `system` | `system` | `settings.json` + `localStorage` |
| Preview theme | `themePreview` | `dark`, `light`, `system` | `system` | `settings.json` + `localStorage` |

## When to Use
- Match OS theme: leave both as `system`.
- Prefer dark editor but light preview: set **App theme** = `dark`, **Preview theme** = `light`.
- Reduce eye strain: use `dark` for both in low light.

## Limitations
- No custom theme creation (only three presets).
- **App theme** changes apply instantly; **Preview theme** may need **preview** refresh.
- Print always uses light theme (not configurable).

## Related Topics
- [Configure Settings](../tasks/configure-settings.md)
- [Reference: Settings Reference](../reference/settings-reference.md)
- [Auto-hide UI](auto-hide-ui.md)