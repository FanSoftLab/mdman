# Install and Run mdman

## Goal
Install mdman on your system and launch the application.

## Prerequisites
- Linux, macOS, or Windows
- Go 1.27 or later (for building from source)
- Node.js 20 or later (for building from source)
- Wails CLI v2.17.1 (for building from source)
- Linux only: `libwebkit2gtk-4.1-dev` (or 4.0), `libgtk-3-dev`, `libayatana-appindicator3-dev`

## Procedure

### Option 1: Download Prebuilt Binary (Recommended)
1. Go to the [Releases page](https://github.com/localtools/mdman/releases).
2. Download the archive for your operating system and architecture.
3. Extract the archive.
4. Run the `mdman` executable.

### Option 2: Build from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/localtools/mdman.git
   cd mdman
   ```
2. Install frontend dependencies:
   ```bash
   cd frontend && npm install && cd ..
   ```
3. Run in development mode:
   ```bash
   wails dev
   ```
   The application opens with live reload.
4. Build a release binary:
   ```bash
   cd frontend && npm run build && cd ..
   wails build
   ```
   The binary appears in `build/bin/`.

### Option 3: Docker Build (Linux, Cross-Platform)
1. Ensure Docker with BuildKit is installed.
2. Run the build script:
   ```bash
   ./build.sh
   ```
   Output appears in `./dist/`.
3. For specific targets:
   ```bash
   ./build.sh -os linux,windows -arch amd64 -v 1.0.0
   ```

## Expected Result
The mdman application window opens. The **welcome screen** appears with options to open a **folder** or create a new **note**.

## Verification
- Verify the window title shows "mdman".
- Verify the **toolbar** displays formatting buttons.
- Verify the **sidebar** shows "No folder open" message.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| `wails: command not found` | Wails CLI not installed or not in PATH | Run `go install github.com/wailsapp/wails/v2/cmd/wails@latest` and ensure `~/go/bin` is in PATH |
| `libwebkit2gtk-4.0-dev: package not found` (Linux) | Wrong package name for distribution | On Debian/Ubuntu: `apt install libwebkit2gtk-4.1-dev libgtk-3-dev libayatana-appindicator3-dev`. On Fedora: `dnf install webkit2gtk4.1-devel gtk3-devel libayatana-appindicator-gtk3-devel` |
| `npm install` fails | Network or permission issue | Run with `--legacy-peer-deps` flag; check network connectivity |
| App window opens but shows black screen | GPU acceleration issue on Linux | See [Troubleshooting: Preview Shows Black Screen](../troubleshooting/preview-black-screen.md) |

## Limitations
- Prebuilt binaries may not be available for all architectures.
- Docker build requires Linux host (cross-compiles to Windows/macOS).
- Development mode (`wails dev`) requires Go and Node.js toolchains.

## Related Topics
- [Open Your First Folder](open-first-folder.md)
- [Create Your First Note](create-first-note.md)
- [Troubleshooting: App Won't Start on Linux](../troubleshooting/app-wont-start-linux.md)