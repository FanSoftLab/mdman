# Linux WebView

## Purpose
Explain the WebKitGTK-based webview on Linux and its configuration.

## What It Is
On Linux, mdman uses WebKitGTK (via Wails v2) to render the frontend. This is a system-provided library (`libwebkit2gtk-4.0` or `4.1`).

## Key Behaviors
- **GPU policy**: `WebviewGpuPolicyOnDemand` (set in `main.go`) — avoids black screen on NVIDIA/Wayland by deferring GPU process launch.
- **Dependencies**: Requires `libwebkit2gtk-4.1-dev` (or `4.0`), `libgtk-3-dev`, `libayatana-appindicator3-dev`.
- **AppImage**: Bundles WebKitGTK helper processes (`WebProcess`, `NetworkProcess`, `GPUProcess`) and libraries to run on any distribution.
- **Wayland vs X11**: Works on both. Wayland + NVIDIA may need `GDK_BACKEND=x11` or `WEBKIT_DISABLE_GPU_PROCESS=1`.

## When to Use
- Building or running on Linux.
- Troubleshooting display issues (black screen, crashes).

## Limitations
- WebKitGTK version tied to distribution (Debian Bookworm = 4.1, Bullseye = 4.0).
- No WebView2 equivalent (Windows uses WebView2, macOS uses WebKit).
- GPU acceleration may cause issues on some driver/hardware combos.

## Related Topics
- [Troubleshooting: App Won't Start on Linux](../troubleshooting/app-wont-start-linux.md)
- [Troubleshooting: Preview Shows Black Screen](../troubleshooting/preview-black-screen.md)
- [Install and Run](../getting-started/install-and-run.md)
- [Docker Build Reference](../reference/cli-build-reference.md)