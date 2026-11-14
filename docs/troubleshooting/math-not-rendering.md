# Math Not Rendering

## Symptom
KaTeX math formulas (`$...$` or `$$...$$`) show as raw LaTeX in **preview** or display `KaTeX parse error`.

## Conditions
- **Note** contains math delimiters.
- **Preview** pane visible.

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| Unmatched delimiters | KaTeX requires balanced `$`/`$$`; single `$` breaks parsing | Verified |
| Unsupported LaTeX command | KaTeX 0.16.x supports subset of LaTeX | Verified |
| Nested math delimiters | `$ $inner$ $` invalid | Verified |
| Invalid syntax in math mode | Missing braces, unknown commands | Verified |

## Diagnostic Checks
1. Count `$` and `$$` in **note** — must be even.
2. Isolate problematic formula — comment out sections to locate.
3. Check [KaTeX Supported Functions](https://katex.org/docs/supported.html) for used commands.

## Corrective Actions

### Action 1: Balance Delimiters
- Inline: `$...$` — exactly two `$` per formula.
- Block: `$$...$$` — exactly two `$$` pairs.
- Never mix: `$ $$...$$ $` invalid.

### Action 2: Use Supported Commands Only
Common unsupported:
- `\usepackage{...}` — not supported.
- `\newcommand` / `\def` — not supported.
- `\begin{tikzpicture}` / `pgfplots` — not supported.
- Custom fonts / packages.

Supported: All [KaTeX Supported Functions](https://katex.org/docs/supported.html).

### Action 3: Fix Syntax Errors
- Braces: `\frac{a}{b}` not `\frac a b`.
- Subscripts: `x_i` not `x_i` (works), but `x_{i+1}` needs braces.
- Special chars: `\%`, `\$`, `\_`, `\&`, `\#`, `\{`, `\}`.

## Expected Result After Correction
- **Preview** renders math as formatted formulas.
- Inline math flows with text; blocks center.

## Escalation
If valid KaTeX still fails:
- Collect: Formula source, KaTeX.org test result, mdman version.
- Report at [GitHub Issues](https://github.com/localtools/mdman/issues).

## Prevention
- Test formulas at [KaTeX.org](https://katex.org/) before pasting.
- Use **toolbar** → **More Options** → **Inline Math** / **Math Block** templates.
- Keep formulas simple; split complex ones.

## Related Topics
- [Work with Math](../tasks/work-with-math.md)
- [Concept: KaTeX Math](../concepts/katex-math.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)