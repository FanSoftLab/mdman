# Auto-hide UI

## Purpose
Maximize editing space by hiding the **sidebar** and **toolbar** until needed.

## What It Is
When enabled, the **sidebar** and **toolbar** are hidden. They appear only when the pointer nears the window edge (left for **sidebar**, top for **toolbar**, 10 px threshold).

## Key Behaviors
- **Activation**: Toggle via **toolbar** → **More Options** → **Auto-hide UI**.
- **Sidebar reveal**: Pointer at left edge (≤10 px) or over **sidebar** area.
- **Toolbar reveal**: Pointer at top edge (≤10 px) or over **toolbar**/**dropdown**.
- **Persistence**: Stored in `localStorage` (key `autoHideUI`).
- **Mobile**: Disabled on mobile (viewport < 768 px); **sidebar** uses overlay instead.
- **View-only mode**: **Auto-hide UI** still applies to **toolbar**.

## Settings
| Setting | Key | Values | Default | Storage |
|---|---|---|---|---|
| Auto-hide UI | `autoHideUI` | `true`, `false` | `false` | `localStorage` only |

## When to Use
- Small screens or focused writing sessions.
- Distraction-free editing.

## Limitations
- **Toolbar** hidden means no access to buttons; use keyboard shortcuts.
- **Sidebar** hidden means no file tree access; use **Recent Files** panel or `Ctrl+O`.
- Does not affect **TOC** panel or **Settings**/**Help** panels.

## Related Topics
- [Use the Toolbar](../tasks/use-toolbar.md)
- [Use Keyboard Shortcuts](../tasks/use-keyboard-shortcuts.md)
- [Reference: Settings Reference](../reference/settings-reference.md)