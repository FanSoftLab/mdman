# Mermaid Diagrams

## Purpose
Create and render diagrams from text definitions within **notes**. Mermaid enables flowcharts, sequence diagrams, class diagrams, state diagrams, pie charts, and more without leaving the **editor**.

## What It Is
Mermaid is a JavaScript diagramming library that parses a Markdown-like syntax inside fenced code blocks tagged `mermaid`. The **preview** renders the diagram as SVG. The **editor** provides syntax highlighting for Mermaid code blocks.

## Key Behaviors
- **Syntax**: Wrap diagram definition in ````mermaid` ... ````.
- **Rendering**: Occurs in **preview** only. **Editor** shows raw syntax.
- **Live Update**: Changes in **editor** reflect in **preview** after debounce (50 ms).
- **Themes**: Diagram colors adapt to **preview theme** (`dark`, `light`, `system`).
- **Print**: Diagrams render in print output (light theme forced for clarity).
- **Copy Code**: **Preview** code blocks have a copy button for the Mermaid source.

## Supported Diagram Types
| Type | Keyword | Example |
|---|---|---|
| Flowchart | `flowchart` or `graph` | `flowchart TD; A --> B` |
| Sequence | `sequenceDiagram` | `sequenceDiagram; A->>B: Hello` |
| Class | `classDiagram` | `classDiagram; ClassA <|-- ClassB` |
| State | `stateDiagram-v2` | `stateDiagram-v2; [*] --> State1` |
| Pie Chart | `pie` | `pie; "A" : 30; "B" : 70` |
| Gantt | `gantt` | `gantt; dateFormat YYYY-MM-DD` |
| Git Graph | `gitgraph` | `gitgraph; commit` |
| ER Diagram | `erDiagram` | `erDiagram; ENTITY ||--o{ RELATION : ""` |
| User Journey | `journey` | `journey; title My Day` |

Full list: [Mermaid Documentation](https://mermaid.js.org/syntax/flowchart.html).

## When to Use
- Visualize processes, workflows, or architectures in technical **notes**.
- Create diagrams that stay in sync with documentation (single source).
- Avoid external diagramming tools for simple charts.

## Related Settings
None specific to Mermaid. Uses **preview theme** for colors.

## Limitations
- Complex diagrams may render slowly in **preview**.
- Very large diagrams may exceed SVG size limits.
- Mermaid version: 11.x (bundled). Syntax matches v11.
- **Editor** does not validate Mermaid syntax; errors appear in **preview** as `Error rendering diagram`.

## Related Topics
- [Work with Diagrams](../tasks/work-with-diagrams.md)
- [Write and Format Markdown](../tasks/write-format-markdown.md)
- [Troubleshooting: Diagrams Not Rendering](../troubleshooting/diagrams-not-rendering.md)
- [Reference: Markdown Syntax Supported](../reference/markdown-syntax.md)