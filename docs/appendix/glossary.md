# Glossary

**Source of truth:** `frontend/src/locales/en.ts` (UI strings) and `frontend/src/locales/uk.ts` (Ukrainian translations)

## Terms

| Preferred Term | Definition | Context | Disallowed Variants | Ukrainian |
|---|---|---|---|---|
| app theme | The color scheme of the application window and UI chrome. Values: `dark`, `light`, `system`. | Settings, toolbar | application theme, UI theme | тема додатку |
| auto-hide UI | A mode that hides the sidebar and toolbar until the pointer nears the window edge. | Settings, toolbar | zen mode, distraction-free mode | авто-приховування UI |
| backup | A JSON file containing all notes, folders, and settings for portability or recovery. | File menu, import/export | export, data dump | резервна копія |
| blockquote | A Markdown block element for quoted text, rendered with a left border. | Toolbar, context menu | quote, citation | цитата |
| code block | A fenced Markdown block for multi-line code with syntax highlighting. | Toolbar, context menu | fenced code, code fence | блок коду |
| context menu | Right-click menu in the editor showing formatting and editing actions. | Editor | right-click menu | контекстне меню |
| diagram | A visual chart rendered from Mermaid syntax in a code block. | Toolbar, preview | Mermaid chart, graph | діаграма |
| editor | The left pane (desktop) or full view (mobile) where you write Markdown. | UI layout, settings | editing pane, markdown editor | редактор |
| folder | A directory in the workspace that contains notes or other folders. | Explorer, file tree | directory, dir | тека |
| GFM | GitHub Flavored Markdown — the Markdown dialect mdman supports (tables, task lists, strikethrough, autolinks). | Concept, reference | GitHub Markdown | GFM |
| heading | A Markdown heading level 1–6 (`#` to `######`). | Toolbar, context menu | header, title | заголовок |
| help panel | The in-app documentation dialog with overview, shortcuts, about, and licenses. | Toolbar, shortcut `Ctrl+H` | help dialog, documentation | довідка |
| horizontal rule | A thematic break in Markdown (`---`). | Toolbar, context menu | divider, line | горизонтальна лінія |
| inline code | Monospaced text within a line, wrapped in backticks (`` `code` ``). | Toolbar, context menu | code span | вбудований код |
| KaTeX | A fast math typesetting library for rendering LaTeX-like formulas in the preview. | Concept, toolbar | math renderer | KaTeX |
| keyboard shortcut | A key combination that triggers an action (e.g., `Ctrl+S` to save). | Help panel, toolbar tooltips | hotkey, key binding | гаряча клавіша |
| language | The UI language. Values: `en` (English), `uk` (Ukrainian). | Settings | locale, UI language | мова |
| link | A Markdown hyperlink (`[text](url)`). | Toolbar, context menu | hyperlink | посилання |
| math block | A display math formula wrapped in `$$...$$`. | Toolbar, context menu | display math, formula block | блок формул |
| Mermaid | A diagramming tool that renders charts from text definitions in fenced code blocks. | Concept, toolbar | diagram tool | Mermaid |
| note | A single Markdown file (`.md`, `.markdown`, `.mdown`, `.txt`) edited in mdman. | User-facing: "Create a note", "Open a note" | file, document, markdown file | файл |
| preview | The right pane (desktop) or full view (mobile) that renders Markdown to HTML. | UI layout, settings | preview pane, rendered view | перегляд |
| preview theme | The color scheme of the rendered Markdown preview. Values: `dark`, `light`, `system`. | Settings | preview theme | тема перегляду |
| recent files | A list of recently opened notes, persisted across sessions. | Sidebar panel, help | history, recent | останні файлі |
| sidebar | The left panel showing the file tree, recent files, and settings. | UI layout, toolbar | explorer, file explorer | бічна панель |
| strikethrough | Text with a line through it (`~~text~~`). | Toolbar, context menu | strikeout, crossed out | закреслений |
| tab | An editor tab holding one open note. | Tab bar, shortcuts | editor tab | вкладка |
| table | A Markdown table with header row and aligned columns. | Toolbar, context menu | markdown table | таблиця |
| table of contents (TOC) | A navigable list of headings extracted from the current note. | Toolbar, right panel | TOC, outline | зміст |
| task list | A Markdown list with checkboxes (`- [ ]`, `- [x]`). | Toolbar, context menu | checklist, todo list | список завдань |
| toolbar | The top bar with formatting buttons, view toggles, and actions. | UI layout | formatting bar, action bar | панель інструментів |
| view-only mode | A read-only mode that shows the preview without the editor. | Toolbar, shortcut | read-only mode, preview mode | режим перегляду |
| workspace | The root folder opened in the sidebar; all notes live inside it. | Explorer, settings | root folder, project folder | робоча тека |

## Abbreviations

| Abbreviation | Full Form | First-Use Context |
|---|---|---|
| TOC | Table of Contents | "The table of contents (TOC) lists headings..." |
| GFM | GitHub Flavored Markdown | "mdman supports GitHub Flavored Markdown (GFM)..." |
| JSON | JavaScript Object Notation | "The backup is a JSON (JavaScript Object Notation) file..." |
| CLI | Command-Line Interface | "The build script provides a CLI (command-line interface)..." |
| GPU | Graphics Processing Unit | "On Linux, GPU acceleration may cause a black webview..." |
| UI | User Interface | "The UI (user interface) consists of..." |
| HTML | HyperText Markup Language | "The preview renders Markdown to HTML..." |
| CSS | Cascading Style Sheets | "Themes use CSS (Cascading Style Sheets) variables..." |
| SVG | Scalable Vector Graphics | "Icons are inline SVG (Scalable Vector Graphics)..." |
| URL | Uniform Resource Locator | "Enter the link URL (Uniform Resource Locator)..." |

## UI String Keys (from `en.ts`)

These keys are the single source for all user-facing text. Documentation must reference these exact keys when quoting UI labels.

**App-level:** `appName`, `loading`, `explorer`, `openFolder`, `newFile`, `newFolder`, `save`, `saveAs`, `open`, `settings`, `help`, `recent`, `toc`, `close`, `search`, `replace`, `findReplace`, `findNext`, `findPrev`, `replaceOne`, `replaceAll`, `caseSensitive`, `themeApp`, `themePreview`, `themeDark`, `themeLight`, `themeSystem`, `language`, `recentMax`, `tabs`, `newTab`, `unsaved`, `saving`, `saved`, `noFolder`, `noFolderHint`, `emptyRecent`, `clearRecent`, `relativePath`, `openedAt`, `cancel`, `apply`, `editMode`, `viewMode`

**Help:** `helpTitle`, `helpOverview`, `helpOverviewP1`, `helpOverviewP2`, `helpShortcuts`, `helpAbout`, `helpAboutP1`, `helpAboutP2`, `aboutTitle`, `aboutMeta`, `githubLink`, `helpLicenses`, `helpBody`, `shortcutSave`, `shortcutOpen`, `shortcutFind`, `shortcutFindReplace`, `shortcutNewTab`, `shortcutCloseTab`, `welcomeTitle`, `noFile`, `langEnglish`, `langUkrainian`

**Toolbar:** `bold`, `boldShortcut`, `italic`, `italicShortcut`, `strikethrough`, `inlineCode`, `heading1`–`heading3`, `bulletList`, `numberedList`, `taskList`, `insertLink`, `insertImage`, `codeBlock`, `insertTable`, `blockquote`, `horizontalRule`, `insertDiagram`, `basicDiagram`, `flowchart`, `sequenceDiagram`, `pieChart`, `inlineMath`, `mathBlock`, `moreOptions`, `print`, `autoHideUI`, `wordWrap`, `syncScroll`, `resetLayout`, `copyCode`, `mermaidRenderError`

**Editor Context Menu:** `cut`, `copy`, `paste`, `selectAll`, `find`, `findAndReplace`, `categoryHeading`, `categoryBlock`, `categoryExtra`, `categoryOther`

**Editor Formatting:** `fmtBold`, `fmtItalic`, `fmtStrikethrough`, `fmtCode`, `fmtBulletList`, `fmtNumberedList`, `fmtCodeBlock`, `fmtTask`, `fmtBlockquote`, `fmtHorizontalRule`, `fmtLink`, `fmtImage`, `fmtTable`, `fmtInlineMath`, `fmtMermaidDiagram`, `fmtHeading1`–`fmtHeading6`

**Help Overlay:** `helpKeyboardShortcuts`, `helpClose`, `helpFormatting`, `helpLists`, `helpHeadings`, `helpInsert`, `helpMultiSelect`, `helpOther`, `helpDescBold`, `helpDescItalic`, `helpDescStrikethrough`, `helpDescInlineCode`, `helpDescBlockquote`, `helpDescBulletList`, `helpDescNumberedList`, `helpDescHeading1`–`helpDescHeading6`, `helpDescLink`, `helpDescSelectNext`, `helpDescSelectAll`, `helpDescToggleHelp`, `helpDescFindReplace`, `helpDescUndo`, `helpDescRedo`, `helpFooter`

**File Tree:** `newFileMenu`, `newFolderMenu`, `downloadMarkdown`, `rename`, `delete`, `folderNamePlaceholder`, `fileNamePlaceholder`, `confirmDeleteFolder`, `confirmDeleteFile`, `emptyFolderAlert`

**App Toolbar/Status:** `toggleSidebar`, `showSidebar`, `hideSidebar`, `tableOfContents`, `closeTableOfContents`, `viewOnly`, `unsavedChanges`, `untitled`

## Usage Rules

1. **Use preferred terms exactly** — match spelling, capitalization, and spacing from this table.
2. **Do not substitute disallowed variants** — e.g., never write "directory" when "folder" is the user-facing term.
3. **Define abbreviations at first use** in each topic — then use the abbreviation alone.
4. **UI labels in docs** must match `en.ts` values exactly (including capitalization).
5. **Code identifiers** in docs use exact casing from source (`themeApp`, not "Theme App").
6. **Ukrainian terms** in `uk.ts` must stay in sync — any new English term requires a Ukrainian entry.