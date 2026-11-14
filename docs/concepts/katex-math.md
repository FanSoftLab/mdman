# KaTeX Math

## Purpose
Write and render mathematical formulas in **notes** using LaTeX-like syntax. KaTeX provides fast, high-quality math typesetting in the **preview**.

## What It Is
KaTeX is a JavaScript library that renders LaTeX math notation to HTML/MathML. In mdmath, it processes two syntax forms:
- **Inline Math**: `$...$` — renders within a text line.
- **Math Block**: `$$...$$` — renders as a centered display block.

The **editor** provides syntax highlighting for math delimiters. The **preview** renders formatted formulas.

## Key Behaviors
- **Syntax**: Standard LaTeX math mode commands (`\alpha`, `\sum`, `\int`, `\frac`, `\sqrt`, matrices, align, etc.).
- **Rendering**: Occurs in **preview** only. **Editor** shows raw LaTeX.
- **Live Update**: Changes reflect in **preview** after debounce (50 ms).
- **Themes**: Math color adapts to **preview theme** (CSS `color` inheritance).
- **Print**: Math renders in print output (KaTeX CSS included).
- **Copy Code**: **Preview** math blocks have copy button for LaTeX source.

## Supported Features (KaTeX 0.16.x)
| Feature | Example |
|---|---|
| Greek letters | `\alpha`, `\beta`, `\Gamma` |
| Operators | `\sum`, `\int`, `\prod`, `\lim` |
| Fractions | `\frac{a}{b}`, `\dfrac{a}{b}` |
| Roots | `\sqrt{x}`, `\sqrt[n]{x}` |
| Sub/superscript | `x_i`, `x^2`, `x_i^2` |
| Matrices | `\begin{matrix} a & b \\ c & d \end{matrix}` |
| Align | `\begin{align} a &= b \\ c &= d \end{align}` |
| Text in math | `\text{some text}` |
| Delimiters | `\left(`, `\right)`, `\langle`, `\rvert` |
| Functions | `\sin`, `\cos`, `\log`, `\exp` |
| Accents | `\hat{x}`, `\vec{v}`, `\bar{y}` |

Full support: [KaTeX Supported Functions](https://katex.org/docs/supported.html).

## When to Use
- Technical **notes** requiring formulas (math, physics, engineering, CS).
- Inline variables and equations in prose.
- Display equations for derivations, proofs, or references.

## Related Settings
None specific to KaTeX. Uses **preview theme** for color.

## Limitations
- KaTeX is not full LaTeX; some packages/commands unsupported (e.g., `\usepackage`, custom macros).
- Very long formulas may overflow **preview** width (horizontal scroll appears).
- **Editor** does not validate LaTeX; errors appear in **preview** as raw text or KaTeX error.
- KaTeX version: 0.16.x (bundled).

## Related Topics
- [Work with Math](../tasks/work-with-math.md)
- [Write and Format Markdown](../tasks/write-format-markdown.md)
- [Troubleshooting: Math Not Rendering](../troubleshooting/math-not-rendering.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)