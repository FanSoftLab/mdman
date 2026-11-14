# Diagrams Not Rendering

## Symptom
Mermaid **diagram** code block in **preview** shows `Error rendering diagram` or raw code.

## Conditions
- **Note** contains ````mermaid` ... ```` block.
- **Preview** pane visible.

## Verified Cause(s)

| Cause | Evidence | Confidence |
|---|---|---|
| Invalid Mermaid syntax | Mermaid parser throws error; shown in **preview** | Verified |
| Unsupported diagram type | Mermaid 11.x does not support all types | Verified |
| Diagram too large | SVG generation exceeds internal limits | Hypothesis |
| Mermaid version mismatch | Bundled Mermaid 11.x vs. newer syntax | Verified |

## Diagnostic Checks
1. Copy diagram source (hover **preview** code block → **Copy code**).
2. Paste into [Mermaid Live Editor](https://mermaid.live) — validates syntax.
3. Check **preview** error message — often indicates line/column.

## Corrective Actions

### Action 1: Fix Syntax Errors
1. Open [Mermaid Live Editor](https://mermaid.live).
2. Paste diagram source.
3. Fix reported errors (missing semicolons, invalid keywords, mismatched brackets).
4. Copy corrected source back to **editor**.

### Action 2: Use Supported Diagram Types
Supported in Mermaid 11.x:
- `flowchart` / `graph`
- `sequenceDiagram`
- `classDiagram`
- `stateDiagram-v2`
- `pie`
- `gantt`
- `gitgraph`
- `erDiagram`
- `journey`
- `mindmap`
- `timeline`
- `quadrantChart`
- `requirementDiagram`

Unsupported: `C4Context`, `C4Container`, `C4Component`, `C4Deployment` (require plugins).

### Action 3: Simplify Large Diagrams
- Split into multiple smaller diagrams.
- Reduce node count.
- Avoid extremely long labels.

## Expected Result After Correction
- **Preview** renders diagram as SVG.
- No error message.

## Escalation
If valid syntax still fails:
- Collect: Diagram source, Mermaid Live Editor result, mdman version.
- Report at [GitHub Issues](https://github.com/localtools/mdman/issues).

## Prevention
- Validate syntax in Mermaid Live Editor before pasting.
- Use **toolbar** → **More Options** → **Insert Diagram** templates as starting point.
- Keep diagrams under 50 nodes for reliable rendering.

## Related Topics
- [Work with Diagrams](../tasks/work-with-diagrams.md)
- [Concept: Mermaid Diagrams](../concepts/mermaid-diagrams.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)