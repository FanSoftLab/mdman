# Work with Math

## Goal
Insert, edit, and render mathematical formulas in **notes** using KaTeX.

## Prerequisites
- A **note** is open in the **editor**.
- **Preview** pane is visible.

## Procedure

### Insert Inline Math
1. Place cursor in text where formula should appear.
2. Click **Inline Math** in **toolbar** → **More Options** → **Inline Math**.
3. Type LaTeX between `$` delimiters (e.g., `$E = mc^2$`).
4. The **preview** renders inline.

### Insert Math Block
1. Place cursor on a blank line.
2. Click **Math Block** in **toolbar** → **More Options** → **Math Block**.
3. Type LaTeX between `$$` delimiters on separate lines:
   ```latex
   $$
   \int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
   $$
   ```
4. The **preview** renders as centered display block.

### Edit Math
1. Locate `$...$` or `$$...$$` in the **editor**.
2. Modify LaTeX. The **preview** updates live.
3. Use **editor** syntax highlighting (math delimiters colored).

### Common Patterns

#### Inline Variables
- `$x_i$`, `$n^2$`, `$\alpha$`, `$\beta$`

#### Fractions and Roots
- `$\frac{a}{b}$`, `$\sqrt{x}$`, `$\sqrt[3]{y}$`

#### Summation and Integrals
- `$\sum_{i=1}^{n} i$`, `$\int_{0}^{1} f(x) dx$`

#### Greek Letters
- `$\alpha, \beta, \gamma, \Gamma, \Delta, \Omega$`

#### Matrices
```latex
$$
\begin{pmatrix}
a & b \\
c & d
\end{pmatrix}
$$
```

#### Aligned Equations
```latex
$$
\begin{align}
x &= y + z \\
  &= 2y
\end{align}
$$
```

#### Text in Math
- `$v_{\text{initial}}$`, `$\text{If } x > 0 \text{ then } y = 1$`

### Copy Math Source
1. In **preview**, hover the rendered math block.
2. Click **Copy code** button (top-right).
3. LaTeX source copies to clipboard.

### Print Math
1. Click **Print** in **toolbar** (`Ctrl+P`).
2. Math renders with KaTeX CSS for print quality.

## Expected Result
- Formulas render in **preview** as formatted math.
- **Editor** shows syntax-highlighted LaTeX.
- Inline math flows with text; blocks center on own lines.

## Verification
- Verify **preview** shows formatted math (not raw LaTeX).
- Verify inline math stays on text baseline.
- Verify print output includes math.

## Recovery / Common Errors

| Error / Situation | Cause | Corrective Action |
|---|---|---|
| Raw LaTeX shows in **preview** | Unclosed delimiter (`$` or `$$`) | Ensure matching pairs; each `$` opens/closes |
| `KaTeX parse error` in **preview** | Invalid LaTeX command | Check [KaTeX Supported Functions](https://katex.org/docs/supported.html); fix syntax |
| Formula overflows **preview** width | Very wide expression | Use `\\` line breaks in `align`; or reduce font size via CSS |
| Inline math not rendering | Missing `$` or nested `$` | Use single `$` for inline; never nest |

## Limitations
- **Editor** does not validate LaTeX; only **preview** shows errors.
- No math-specific shortcuts (use **toolbar** menu).
- Unsupported LaTeX: `\usepackage`, custom macros, `tikz`, `pgfplots`.
- KaTeX 0.16.x — some newer LaTeX commands unavailable.

## Related Topics
- [Concept: KaTeX Math](../concepts/katex-math.md)
- [Write and Format Markdown](../tasks/write-format-markdown.md)
- [Troubleshooting: Math Not Rendering](../troubleshooting/math-not-rendering.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)