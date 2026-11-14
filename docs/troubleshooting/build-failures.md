# Build Failures

## Symptom
`./build.sh`, `wails build`, or Docker build fails with errors.

## Conditions
- Building from source or via Docker.
- Any platform.

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| Missing Go/Node.js | `go: command not found` or `npm: command not found` | Verified |
| Wrong Go version | `go.mod` requires 1.27; older Go fails | Verified |
| Missing Linux deps | `libwebkit2gtk` not found at link time | Verified |
| Docker BuildKit disabled | `DOCKER_BUILDKIT=1` not set | Verified |
| Network issues | `git clone`, `npm install`, `go mod download` fail | Verified |
| Disk space | `no space left on device` | Verified |
| Architecture mismatch | Cross-compile without proper toolchain | Verified |

## Diagnostic Checks
1. Check versions: `go version`, `node --version`, `npm --version`.
2. Check Docker: `docker version`, `docker buildx version`.
3. Run `wails doctor` — checks environment.
4. Check disk: `df -h`.
5. Read error output — identify failing stage.

## Corrective Actions

### Action 1: Install Required Tools
- Go 1.27+: `go install golang.org/dl/go1.27@latest && go1.27 download`
- Node.js 20+: Use NodeSource repo or `nvm`.
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Action 2: Install Linux Dependencies
- Debian/Ubuntu: `sudo apt install libwebkit2gtk-4.1-dev libgtk-3-dev libayatana-appindicator3-dev gcc g++ pkg-config`
- Fedora: `sudo dnf install webkit2gtk4.1-devel gtk3-devel libayatana-appindicator-gtk3-devel gcc g++ pkg-config`
- Arch: `sudo pacman -S webkit2gtk-4.1 gtk3 libayatana-appindicator base-devel`

### Action 3: Enable BuildKit
```bash
export DOCKER_BUILDKIT=1
./build.sh
```

### Action 4: Fix Network/Proxy
- Set `HTTP_PROXY`, `HTTPS_PROXY` for `go mod download`, `npm install`.
- Use `git config --global url."https://".insteadOf git://` for git.

### Action 5: Free Disk Space
```bash
docker system prune -a
rm -rf ~/.cache/go-build
npm cache clean --force
```

### Action 6: Clean Build
```bash
./build.sh -no-cache
```
Or manually:
```bash
rm -rf frontend/node_modules frontend/dist
go clean -cache -modcache
./build.sh
```

## Expected Result After Correction
- Build completes.
- Artifacts in `./dist/` (Docker) or `build/bin/` (wails).

## Escalation
If the problem persists:
- Collect: Full build log, `go version`, `node --version`, `docker version`, OS.
- Run `wails doctor` and include output.
- Report at [GitHub Issues](https://github.com/FanSoftLab/mdman/issues).

## Prevention
- Use Docker build (`./build.sh`) for consistent environment.
- Keep Go, Node.js, Docker updated.
- Run `wails doctor` before building.

## Related Topics
- [CLI and Build Reference](../reference/cli-build-reference.md)
- [Install and Run](../getting-started/install-and-run.md)
- [Concept: Reproducible Builds](../concepts/reproducible-builds.md)