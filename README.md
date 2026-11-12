# mdmark

**Desktop port of [SvelteMark](https://github.com/MasFana/sveltemark).**

Privacy-first, local-only Markdown editor. Same UI and behaviour as SvelteMark, running as a native desktop app via Wails v2 (Go + system webview).

- Notes live in IndexedDB (Dexie) inside the webview — nothing leaves the device
- CodeMirror 6, live preview, Mermaid, KaTeX, GFM
- Hierarchical folders/files, backup/import JSON
- Dark GitHub-style UI

## Requirements

- Go 1.27+
- Node.js 20+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- Linux: `libwebkit2gtk-4.1-dev` (or 4.0), `libgtk-3-dev`, `libayatana-appindicator3-dev`

## Develop

```bash
cd frontend && npm install && cd ..
wails dev
```

## Build

```bash
cd frontend && npm install && npm run build && cd ..
wails build
```

Binary appears under `build/bin/`.

## Relation to SvelteMark

This is a direct desktop packaging of the SvelteMark application UI and data model (IndexedDB notes, folders, editor, preview, toolbar). The Go side only provides the native window and embeds the built frontend.

Upstream: https://github.com/MasFana/sveltemark (MIT)

## License

MIT (same as SvelteMark). See [LICENSE](LICENSE).

## Docker / reproducible build

```bash
./build.sh                              # current OS/arch, AppImage on Linux
./build.sh -os linux -arch amd64
./build.sh -os linux,windows -arch amd64 -v 1.0.0
./build.sh -no-appimage -no-cache
```

Requires Docker with BuildKit. Output goes to `./dist/`.
