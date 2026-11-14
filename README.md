# mdman

**Desktop port of [SvelteMark](https://github.com/MasFana/sveltemark).**

Privacy-first, local-only Markdown editor. Same UI and behaviour as SvelteMark, running as a native desktop app via Wails v2 (Go + system webview).

- Notes live in IndexedDB (Dexie) inside the webview — nothing leaves the device
- CodeMirror 6, live preview, Mermaid, KaTeX, GFM
- Hierarchical folders/files, backup/import JSON
- Dark GitHub-style UI

## Documentation

- **Getting Started**
  - [Install and Run](docs/getting-started/install-and-run.md)
  - [Open Your First Folder](docs/getting-started/open-first-folder.md)
  - [Create Your First Note](docs/getting-started/create-first-note.md)
- **Tasks**
  - [Write and Format Markdown](docs/tasks/write-format-markdown.md)
  - [Use the Toolbar](docs/tasks/use-toolbar.md)
  - [Use Keyboard Shortcuts](docs/tasks/use-keyboard-shortcuts.md)
  - [Use the Context Menu](docs/tasks/use-context-menu.md)
  - [Work with Diagrams](docs/tasks/work-with-diagrams.md)
  - [Work with Math](docs/tasks/work-with-math.md)
  - [Organize Folders and Files](docs/tasks/organize-folders-files.md)
  - [Import, Export, and Backup](docs/tasks/import-export-backup.md)
  - [Configure Settings](docs/tasks/configure-settings.md)
- **Concepts**
  - [Workspaces and Notes](docs/concepts/workspaces-notes.md)
  - [Themes and Appearance](docs/concepts/themes-appearance.md)
  - [Auto-hide UI](docs/concepts/auto-hide-ui.md)
  - [Mermaid Diagrams](docs/concepts/mermaid-diagrams.md)
  - [KaTeX Math](docs/concepts/katex-math.md)
  - [Linux WebView](docs/concepts/linux-webview.md)
  - [Reproducible Builds](docs/concepts/reproducible-builds.md)
  - [Help Panel](docs/concepts/help-panel.md)
- **Reference**
  - [Keyboard Shortcuts](docs/reference/keyboard-shortcuts.md)
  - [Settings Reference](docs/reference/settings-reference.md)
  - [Markdown Syntax Supported](docs/reference/markdown-syntax.md)
  - [JSON Backup Schema](docs/reference/json-backup-schema.md)
  - [CLI and Build Reference](docs/reference/cli-build-reference.md)
- **Troubleshooting**
  - [App Won't Start on Linux](docs/troubleshooting/app-wont-start-linux.md)
  - [Preview Shows Black Screen](docs/troubleshooting/preview-black-screen.md)
  - [Files Not Saving](docs/troubleshooting/files-not-saving.md)
  - [Diagrams Not Rendering](docs/troubleshooting/diagrams-not-rendering.md)
  - [Math Not Rendering](docs/troubleshooting/math-not-rendering.md)
  - [Settings Not Persisting](docs/troubleshooting/settings-not-persisting.md)
  - [Import/Backup Fails](docs/troubleshooting/import-backup-fails.md)
  - [Build Failures](docs/troubleshooting/build-failures.md)
- **Appendix**
  - [Glossary](docs/appendix/glossary.md)
  - [Licenses](docs/appendix/licenses.md)
  - [Changelog](docs/appendix/changelog.md)

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