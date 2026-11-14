# Help Panel

## Purpose
In-app documentation accessible via **Help** button (`Ctrl+H`).

## What It Is
A modal dialog with four tabs: **Overview**, **Shortcuts**, **About**, **Licenses**. Includes search filter.

## Key Behaviors
- **Open**: Click **Help** in **toolbar** or press `Ctrl+H` (when UI focused).
- **Search**: Filters content across all tabs in real time.
- **Tabs**:
  - **Overview**: App description, key features.
  - **Shortcuts**: Complete keyboard shortcut table.
  - **About**: Version, license, upstream link.
  - **Licenses**: Third-party license list (searchable).
- **Close**: Click **Close**, press `Esc`, or click overlay.
- **Localization**: All content uses `en.ts`/`uk.ts` strings.

## When to Use
- Learn shortcuts and features.
- Check version and licenses.
- Troubleshoot without leaving app.

## Limitations
- No deep linking to specific sections.
- Search only filters visible text (not code blocks).
- Content updates require app rebuild.

## Related Topics
- [Use Keyboard Shortcuts](../tasks/use-keyboard-shortcuts.md)
- [Appendix: Licenses](../appendix/licenses.md)
- [Appendix: Changelog](../appendix/changelog.md)