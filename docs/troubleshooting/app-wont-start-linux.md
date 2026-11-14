# App Won't Start on Linux

## Symptom
Running `mdman` (or `wails dev`) produces no window, crashes immediately, or shows an error in the terminal.

## Conditions
- Linux (any distribution)
- Native binary or `wails dev`
- May occur on Wayland, X11, or specific GPU drivers (NVIDIA)

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| Missing WebKitGTK dependencies | `ldd mdman` shows `libwebkit2gtk-4.0.so` not found | Verified |
| Wayland + NVIDIA GPU: WebView GPU process fails | Wails v2 default GPU policy crashes on NVIDIA/Wayland | Verified (fixed in `main.go`: `WebviewGpuPolicyOnDemand`) |
| Incorrect `libwebkit2gtk` version | Package `libwebkit2gtk-4.0-dev` vs `4.1-dev` mismatch | Verified |
| Sandbox/namespace restrictions | Flatpak, AppImage, or kernel `unprivileged_userns_clone=0` | Hypothesis |

## Diagnostic Checks
1. Run from terminal: `./mdman` — observe error output.
2. Check dependencies: `ldd mdman | grep webkit` — verify `libwebkit2gtk-4.0.so` or `4.1.so` resolves.
3. Check Wayland: `echo $XDG_SESSION_TYPE` — if `wayland`, test with `GDK_BACKEND=x11 ./mdman`.
4. Check GPU: `lspci | grep -i nvidia` — if present, GPU policy may be cause.
5. Check kernel namespaces: `sysctl kernel.unprivileged_userns_clone` — if `0`, enable it.

## Corrective Actions

### Action 1: Install Correct WebKitGTK Dependencies
1. Identify distribution:
   - Debian/Ubuntu 22.04+: `sudo apt install libwebkit2gtk-4.1-dev libgtk-3-dev libayatana-appindicator3-dev`
   - Debian/Ubuntu 20.04: `sudo apt install libwebkit2gtk-4.0-dev libgtk-3-dev libayatana-appindicator3-dev`
   - Fedora: `sudo dnf install webkit2gtk4.1-devel gtk3-devel libayatana-appindicator-gtk3-devel`
   - Arch: `sudo pacman -S webkit2gtk-4.1 gtk3 libayatana-appindicator`
2. Re-run `mdman`.

### Action 2: Force X11 Backend (Wayland + NVIDIA)
1. Run with X11 backend:
   ```bash
   GDK_BACKEND=x11 ./mdman
   ```
2. If this works, create a launcher script or desktop entry with `Exec=env GDK_BACKEND=x11 mdman`.

### Action 3: Enable User Namespaces (If Sandbox Issue)
1. Check current value:
   ```bash
   sysctl kernel.unprivileged_userns_clone
   ```
2. If `0`, enable temporarily:
   ```bash
   sudo sysctl kernel.unprivileged_userns_clone=1
   ```
3. Persist across reboots:
   ```bash
   echo "kernel.unprivileged_userns_clone=1" | sudo tee /etc/sysctl.d/99-userns.conf
   sudo sysctl --system
   ```

### Action 4: Run with Debug Logging
1. Set Wails debug flag:
   ```bash
   WAILS_DEBUG=1 ./mdman
   ```
2. Observe console output for specific error.

## Expected Result After Correction
- Application window opens.
- **Welcome screen** or last **workspace** loads.
- No terminal errors.

## Escalation
If the problem persists after all actions:
- Collect: Distribution (`cat /etc/os-release`), kernel (`uname -r`), GPU (`lspci`), session type (`echo $XDG_SESSION_TYPE`), full terminal output.
- Report at [GitHub Issues](https://github.com/localtools/mdman/issues) with collected information.

## Prevention
- Install dependencies before building/running.
- On Wayland + NVIDIA, use `GDK_BACKEND=x11` or ensure `WebviewGpuPolicyOnDemand` is set (default in `main.go`).
- Keep WebKitGTK updated via system package manager.

## Related Topics
- [Preview Shows Black Screen](preview-black-screen.md)
- [Install and Run](../getting-started/install-and-run.md)
- [Concept: Linux WebView](../concepts/linux-webview.md)