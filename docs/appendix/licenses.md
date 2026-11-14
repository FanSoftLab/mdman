# Licenses

## Scope
Third-party licenses for libraries and tools used in mdman. Source: `frontend/src/lib/components/HelpPanel.svelte` `licenses` array. mdman itself is MIT (see `LICENSE`).

## Dependencies

| Name | Version | Description | License | URL |
|---|---|---|---|---|
| Svelte | 5.x | Cybernetically enhanced web apps | MIT | https://github.com/sveltejs/svelte |
| SvelteMark | upstream | Privacy-first local Markdown editor (UI source) | MIT | https://github.com/MasFana/sveltemark |
| CodeMirror | 6.x | Extensible code editor for the browser | MIT | https://github.com/codemirror/dev |
| Mermaid | 11.x | Diagramming and charting tool | MIT | https://github.com/mermaid-js/mermaid |
| KaTeX | 0.16.x | Fast math typesetting for the web | MIT | https://github.com/KaTeX/KaTeX |
| highlight.js | 11.x | Syntax highlighting for code blocks | BSD-3-Clause | https://github.com/highlightjs/highlight.js |
| Font Awesome Free | 7.x | Icon set used across the UI | CC BY 4.0 / OFL / MIT | https://github.com/fortawesome/font-awesome |
| github-markdown-css | 5.x | GitHub Markdown CSS styles | MIT | https://github.com/sindresorhus/github-markdown-css |
| Dexie.js | 4.x | Minimalistic IndexedDB wrapper | Apache-2.0 | https://github.com/dexie/Dexie.js |
| JSZip | 3.x | Create and read ZIP archives | MIT / GPL-3.0 | https://github.com/Stuk/jszip |
| lucide-svelte | 0.5x | Lucide icons for Svelte | ISC | https://github.com/lucide-icons/lucide |
| github-slugger | 2.x | Slug generation for headings | ISC | https://github.com/Flet/github-slugger |
| unified | 11.x | Content transformation platform (Markdown pipeline) | MIT | https://github.com/unifiedjs/unified |
| remark-gfm | 4.x | GitHub Flavored Markdown support | MIT | https://github.com/remarkjs/remark-gfm |
| remark-math | 6.x | Math blocks support | MIT | https://github.com/remarkjs/remark-math |
| rehype-highlight | 7.x | Syntax highlighting in HAST | MIT | https://github.com/rehypejs/rehype-highlight |
| rehype-katex | 7.x | Render math via KaTeX | MIT | https://github.com/remarkjs/remark-math |
| rehype-slug | 6.x | Add ids to headings | MIT | https://github.com/rehypejs/rehype-slug |
| Wails | 2.x | Desktop apps with Go and web technologies | MIT | https://github.com/wailsapp/wails |
| Vite | 6.x | Frontend build tool | MIT | https://github.com/vitejs/vite |
| TypeScript | 5.x | Typed JavaScript | Apache-2.0 | https://github.com/microsoft/TypeScript |

## Go Dependencies (Indirect)
All Go dependencies are transitive via `github.com/wailsapp/wails/v2`. Key indirect deps include:
- `github.com/wailsapp/go-webview2` (Windows WebView2 loader)
- `github.com/go-ole/go-ole` (Windows COM)
- `github.com/godbus/dbus/v5` (Linux D-Bus)
- `github.com/pkg/errors` (Error wrapping)
- `golang.org/x/sys`, `golang.org/x/text`, `golang.org/x/crypto` (Standard library extensions)

Full list: `go mod graph` or `go.sum`.

## License Texts

### MIT License (mdman and most deps)
```
Copyright (c) <year> <copyright holders>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

### Apache-2.0 (Dexie.js, TypeScript)
See: https://www.apache.org/licenses/LICENSE-2.0

### BSD-3-Clause (highlight.js)
See: https://opensource.org/licenses/BSD-3-Clause

### ISC (lucide-svelte, github-slugger)
See: https://opensource.org/licenses/ISC

### CC BY 4.0 / OFL / MIT (Font Awesome Free)
See: https://fontawesome.com/license

### GPL-3.0 (JSZip dual license)
See: https://www.gnu.org/licenses/gpl-3.0.html

## Compliance
- All licenses are permissive or copyleft-compatible for desktop application distribution.
- No license requires source distribution for binary-only distribution (GPL-3.0 on JSZip is dual-licensed MIT).
- Attribution provided in **Help** panel and this document.

## Related Topics
- [Help Panel](../concepts/help-panel.md)
- [README](../README.md)