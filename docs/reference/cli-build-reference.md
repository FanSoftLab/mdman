# CLI and Build Reference

## Scope
Command-line interface for `build.sh`, `wails` commands, Docker build arguments, and version embedding. Source: `build.sh`, `main.go`, `Dockerfile`, `wails.json`.

## build.sh Flags

| Flag | Values | Default | Description |
|---|---|---|---|
| `-mode` | `release`, `debug` | `release` | Build optimization. `debug` disables `-s -w` ldflags. |
| `-os` | `linux`, `windows`, `darwin`, `all` (comma-separated) | Host OS | Target operating system(s). |
| `-arch` | `amd64`, `arm64`, `all` (comma-separated) | Host arch | Target architecture(s). |
| `-run-tests` | (boolean flag) | `false` | Runs `go test ./...` and `go vet ./...` (best-effort). |
| `-no-cache` | (boolean flag) | `false` | Disables Docker BuildKit cache (`--no-cache`). |
| `-no-appimage` | (boolean flag) | `false` | Skips Linux AppImage packaging. |
| `-v`, `--version` | SemVer string (e.g., `1.0.0`) | Git tag or `dev-<short-SHA>` | Application version embedded in binary. |

### Examples
```bash
# Default: current OS/arch, release, AppImage on Linux
./build.sh

# Linux AMD64 only, version 1.0.0
./build.sh -os linux -arch amd64 -v 1.0.0

# Cross-platform: Linux and Windows, AMD64, no AppImage
./build.sh -os linux,windows -arch amd64 -no-appimage -v 1.2.0

# Debug build with tests
./build.sh -mode debug -run-tests
```

### Output
- Artifacts written to `./dist/` (flat).
- Naming: `mdman-<os>-<arch>` (`.exe` on Windows).
- Linux AppImage: `mdman-linux-<arch>.AppImage`.

## Wails CLI Commands

| Command | Description |
|---|---|
| `wails dev` | Development mode with live reload (frontend + Go). |
| `wails build` | Production build for current platform. |
| `wails build -platform <os/arch>` | Cross-compile (requires Go toolchain for target). |
| `wails doctor` | Check environment for Wails requirements. |

### Wails Build Flags (via `build.sh` ldflags)
Embedded at link time via `-X main.<var>=<value>`:

| Variable | Source | Description |
|---|---|---|
| `main.version` | `-v` flag / git tag | Application version string. |
| `main.buildTime` | `date -u` | ISO 8601 UTC timestamp. |
| `main.buildCommit` | `git rev-parse HEAD` | Full commit SHA. |
| `main.buildSource` | `git remote get-url origin` | Repository URL or `local`. |

## Docker Build Arguments (Dockerfile)

| Arg | Default | Description |
|---|---|---|
| `GO_VERSION` | `1.27` | Go toolchain version (base image tag). |
| `NODE_VERSION` | `20` | Node.js version for frontend build. |
| `RUN_TESTS` | `false` | Run tests in builder stage. |
| `BUILD_MODE` | `release` | `release` or `debug`. |
| `APP_VERSION` | `dev` | Version string for binary. |
| `BUILD_TIME` | `unknown` | Build timestamp. |
| `BUILD_COMMIT` | `unknown` | Git commit SHA. |
| `BUILD_SOURCE` | `local` | Repository URL. |
| `PLATFORMS` | `linux/amd64` | Comma-separated `os/arch` targets. |
| `BUILD_APPIMAGE` | `true` | Build AppImage on Linux. |

### Docker Build Command (via build.sh)
```bash
docker build \
  --build-arg RUN_TESTS=true \
  --build-arg BUILD_MODE=release \
  --build-arg APP_VERSION=1.0.0 \
  --build-arg BUILD_TIME="2026-01-15T10:30:00Z" \
  --build-arg BUILD_COMMIT="abc123..." \
  --build-arg BUILD_SOURCE="https://github.com/localtools/mdman" \
  --build-arg PLATFORMS="linux/amd64,windows/amd64" \
  --build-arg BUILD_APPIMAGE=true \
  --target exporter \
  --output "type=local,dest=.buildout" \
  .
```

## wails.json Configuration

```json
{
  "name": "mdman",
  "outputfilename": "mdman",
  "frontend:install": "npm install",
  "frontend:build": "npm run build:wails",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": { "name": "mdman", "email": "" },
  "info": {
    "companyName": "mdman",
    "productName": "mdman",
    "productVersion": "1.0.0",
    "copyright": "Copyright © 2026",
    "comments": "Desktop port of SvelteMark — privacy-first local Markdown editor"
  }
}
```

### Key Fields
| Field | Purpose |
|---|---|
| `name` / `outputfilename` | Binary name (`mdman`). |
| `frontend:install` | Runs in Docker builder stage. |
| `frontend:build` | `vite build` output to `frontend/dist`. |
| `info.productVersion` | Fallback version (overridden by ldflags). |

## Go Module (go.mod)

```go
module github.com/localtools/mdman

go 1.27

require github.com/wailsapp/wails/v2 v2.17.1
```

- Single dependency: `wails/v2` v2.17.1.
- All other deps are indirect (transitive).

## Version Embedding Details

In `main.go`:
```go
var (
	version     = "dev"
	buildTime   = "unknown"
	buildCommit = "unknown"
	buildSource = "local"
)
```

Overridden at build via `-ldflags`:
```
-s -w -X main.version=1.0.0 -X main.buildTime=2026-01-15T10:30:00Z -X main.buildCommit=abc123 -X main.buildSource=https://github.com/localtools/mdman
```

- `-s -w`: Strip symbol table and DWARF (release only).
- Accessible at runtime via `main.version` etc.
- Displayed in **Help** → **About** panel.

## File Arguments (CLI Launch)
- `mdman file1.md file2.md` — opens files in new **tabs** on startup.
- Handled in `main.go` `flag.Args()` → `app.SetPendingFiles()` → `app.ConsumePendingFiles()` on **DomReady**.

## Single Instance Lock
- `UniqueId: "com.localtools.mdman"` prevents multiple instances.
- Second launch: activates existing window, emits `open-files` event with arguments.

## Related Topics
- [Install and Run](../getting-started/install-and-run.md)
- [Troubleshooting: Build Failures](../troubleshooting/build-failures.md)
- [Concept: Reproducible Builds](../concepts/reproducible-builds.md)