# Preview Shows Black Screen

## Symptom
The **preview** pane (or entire window on mobile/view-only) displays a black rectangle. **Editor** may work normally.

## Conditions
- Linux (primary), possibly other platforms
- Wayland session with NVIDIA GPU
- After wake from sleep or display hot-plug
- AppImage on distributions with older `glibc`

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| WebKitGTK GPU process crash on NVIDIA/Wayland | `WebviewGpuPolicyOnDemand` mitigates; black screen persists on some drivers | Verified |
| Missing WebKitGTK helper processes in AppImage | `linuxdeploy` excludes `WebProcess`, `NetworkProcess`, `GPUProcess`; webview stays black | Verified (fixed in `Dockerfile` bundling) |
| `glibc` version mismatch (AppImage) | AppImage built on newer `glibc` than host; WebKitGTK fails to initialize | Verified |
| GTK theme engine conflict | Custom GTK themes (e.g., `adwaita` variants) cause WebView render failure | Hypothesis |

## Diagnostic Checks
1. Run from terminal: `./mdman` — check for `WEBKIT` or `GPU` errors in stderr.
2. Test **editor** pane — if **editor** works but **preview** is black, issue is WebView rendering.
3. Check session type: `echo $XDG_SESSION_TYPE` — Wayland vs X11.
4. Check GPU: `lspci | grep -i nvidia`.
5. Test with `GDK_BACKEND=x11 ./mdman` — if fixed, cause is Wayland-specific.
6. For AppImage: run with `--appimage-extract-and-run` to extract and test.

## Corrective Actions

### Action 1: Force X11 Backend (Wayland + NVIDIA)
1. Run:
   ```bash
   GDK_BACKEND=x11 ./mdman
   ```
2. If fixed, create persistent launcher:
   ```bash
   echo 'Exec=env GDK_BACKEND=x11 mdman' | sudo tee /usr/share/applications/mdman.desktop
   ```

### Action 2: Disable GPU Acceleration (WebView)
1. The app sets `WebviewGpuPolicyOnDemand` by default (`main.go`).
2. If still failing, set environment variable before launch:
   ```bash
   WEBKIT_DISABLE_GPU_PROCESS=1 ./mdman
   ```
   Or:
   ```bash
   WEBKIT_DISABLE_COMPOSITING_MODE=1 ./mdman
   ```

### Action 3: Use AppImage with Extracted Mode
1. Extract AppImage:
   ```bash
   ./mdman-linux-amd64.AppImage --appimage-extract
   ```
2. Run extracted:
   ```bash
   ./squid/AppRun
   ```
3. If this works, the issue is AppImage runtime integration.

### Action 4: Update System WebKitGTK
1. Update packages:
   - Debian/Ubuntu: `sudo apt update && sudo apt upgrade libwebkit2gtk-4.1-0`
   - Fedora: `sudo dnf update webkit2gtk4.1`
   - Arch: `sudo pacman -S webkit2gtk-4.1`
2. Restart mdman.

### Action 5: Switch GTK Theme
1. Temporarily use default Adwaita:
   ```bash
   GTK_THEME=Adwaita ./mdman
   ```
2. If fixed, the custom theme causes the issue.

## Expected Result After Correction
- **Preview** pane renders Markdown correctly.
- No black rectangles.
- Diagrams and math render in **preview**.

## Escalation
If the problem persists after all actions:
- Collect: Distribution, kernel, GPU, session type, `WEBKIT_DISABLE_GPU_PROCESS` test result, terminal output.
- Test with `wails dev` (development mode) — if it works, issue is production build/AppImage specific.
- Report at [GitHub Issues](https://github.com/localtools/mdman/issues).

## Prevention
- Build AppImage on oldest supported `glibc` (Docker `bookworm` base used in `Dockerfile`).
- Bundle WebKitGTK helpers in AppImage (`Dockerfile` copies `/usr/lib/*/webkit2gtk-4.0`).
- Default `WebviewGpuPolicyOnDemand` in `main.go` avoids eager GPU initialization.

## Related Topics
- [App Won't Start on Linux](app-wont-start-linux.md)
- [Install and Run](../getting-started/install-and-run.md)
- [Concept: Linux WebView](../concepts/linux-webview.md)
- [Docker Build Reference](../reference/cli-build-reference.md)