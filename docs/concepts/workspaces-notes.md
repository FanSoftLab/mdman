# Workspaces and Notes

## Purpose
Organize **notes** in a hierarchical **folder** structure within a **workspace**.

## What It Is
A **workspace** is the root **folder** you open in mdman. All **notes** and sub**folders** live inside it. The **sidebar** displays this tree.

## Key Behaviors
- **Single workspace**: Only one **workspace** open at a time.
- **Persistence**: Last **workspace** restores on startup (stored in `settings.json`).
- **File system sync**: **Sidebar** mirrors actual directory structure. Operations (create, rename, delete, move) apply to disk immediately.
- **Note discovery**: Only Markdown files (`.md`, `.markdown`, `.mdown`, `.txt`) appear. Hidden files (`.name`) excluded.
- **Sort order**: **Folders** first, then **notes**; alphabetical, case-insensitive.

## When to Use
- Group related **notes** by project, topic, or date.
- Keep personal and work **notes** in separate **workspaces** (switch via **Open Folder**).

## Related Settings
| Setting | Key | Values | Default |
|---|---|---|---|
| Workspace | `workspace` | Absolute path | (last opened) |
| Recent files limit | `recentFileMax` | 5–100 | 20 |

## Limitations
- **Notes** outside **workspace** cannot be opened via **sidebar** (use **Open File**).
- Symbolic links followed; target must be inside **workspace**.
- No virtual folders or tags (only file system hierarchy).

## Related Topics
- [Open Your First Folder](../getting-started/open-first-folder.md)
- [Organize Folders and Files](../tasks/organize-folders-files.md)
- [Import, Export, and Backup](../tasks/import-export-backup.md)
- [Reference: JSON Backup Schema](../reference/json-backup-schema.md)