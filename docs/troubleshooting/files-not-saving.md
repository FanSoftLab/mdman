# Files Not Saving / Permission Errors

## Symptom
**Save** (`Ctrl+S`) or **Save As** fails silently or shows error. **Note** changes not persisted to disk.

## Conditions
- Any platform
- **Note** in **workspace** or outside
- File system permissions, disk full, or path issues

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| No write permission on **workspace** or **note** | `os.WriteFile` returns `EACCES` | Verified |
| **Workspace** not set (untitled **note**) | `Save` requires existing path; `Save As` works | Verified |
| Disk full | `ENOSPC` error | Verified |
| Path traversal blocked | `resolvePath` rejects paths outside **workspace** | Verified (`app.go`) |
| File locked by another process | `EBUSY` or `ETXTBSY` | Hypothesis |

## Diagnostic Checks
1. Try **Save As** — if it works, **note** has no path (untitled).
2. Check **workspace** permissions: `ls -ld <workspace>` — owner must have write.
3. Check **note** permissions: `ls -l <note>` — write bit required.
4. Check disk space: `df -h <workspace>`.
5. Open terminal, try `echo test > <workspace>/test.txt` — verifies write access.

## Corrective Actions

### Action 1: Use Save As for Untitled Notes
1. Press `Ctrl+Shift+S` (**Save As**).
2. Choose location in **workspace**.
3. File saves; subsequent `Ctrl+S` works.

### Action 2: Fix Workspace Permissions
1. Ensure you own the **workspace**:
   ```bash
   sudo chown -R $USER:<group> <workspace>
   chmod -R u+w <workspace>
   ```
2. Retry **Save**.

### Action 3: Free Disk Space
1. Check space: `df -h`.
2. Remove unnecessary files.
3. Retry **Save**.

### Action 4: Avoid Path Traversal
- **Save As** only allows paths inside **workspace**.
- Move **note** into **workspace** first (drag in **sidebar**).

## Expected Result After Correction
- **Save** (`Ctrl+S`) succeeds; **toolbar** shows **Saved** checkmark.
- **Note** content matches **editor** on reload.

## Escalation
If the problem persists:
- Collect: OS, **workspace** path, `ls -la` output, disk space, exact error message (if any).
- Report at [GitHub Issues](https://github.com/FanSoftLab/mdman/issues).

## Prevention
- Always open a **workspace** before creating **notes**.
- Keep disk space > 1 GB free.
- Run mdman as user who owns **workspace**.

## Related Topics
- [Create Your First Note](../getting-started/create-first-note.md)
- [Open Your First Folder](../getting-started/open-first-folder.md)
- [Organize Folders and Files](../tasks/organize-folders-files.md)