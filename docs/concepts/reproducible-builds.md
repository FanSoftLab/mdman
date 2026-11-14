# Reproducible Builds

## Purpose
Produce identical binaries from the same source using Docker and BuildKit.

## What It Is
The `build.sh` script uses Docker with BuildKit to build mdman in a controlled environment. All dependencies (Go, Node.js, WebKitGTK, linuxdeploy) are pinned in the `Dockerfile`.

## Key Behaviors
- **Base image**: `golang:1.27-bookworm` (fixed Go version, Debian Bookworm).
- **Node.js**: Version 20 (via NodeSource repo).
- **WebKitGTK**: From Debian Bullseye (backported for 4.0 compatibility).
- **Build args**: Version, time, commit, source embedded via `-ldflags`.
- **Output**: Flat artifacts in `./dist/`.
- **Cache**: BuildKit layer cache enabled by default (`-no-cache` to disable).

## Version Embedding
Embedded at link time:
| Variable | Value |
|---|---|
| `main.version` | Git tag or `dev-<short-SHA>` |
| `main.buildTime` | ISO 8601 UTC |
| `main.buildCommit` | Full SHA |
| `main.buildSource` | Git remote URL or `local` |

## When to Use
- CI/CD pipelines.
- Release builds.
- Cross-compilation (Linux host → Windows/macOS targets).

## Limitations
- Requires Docker with BuildKit.
- Linux host required (cross-compiles to other OSes).
- AppImage only on Linux target.
- Not bit-for-bit reproducible across Docker versions (timestamps in layers).

## Related Topics
- [CLI and Build Reference](../reference/cli-build-reference.md)
- [Install and Run](../getting-started/install-and-run.md)
- [Dockerfile](../Dockerfile)