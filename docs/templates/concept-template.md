# Concept Template

**Content Type:** Concept (Explanatory)
**Standard:** ISO/IEC/IEEE 26514 §2.7 + ASD-STE100
**Max sentence length:** 25 words (descriptive)
**Max sentences per paragraph:** 6
**Max noun cluster:** 3 words

---

## [Concept Title: Noun Phrase]

*Example: "Workspaces and Notes", "Themes and Appearance", "Mermaid Diagrams"*

### Purpose
One sentence: what this feature is and why a user would use it.

### What It Is
- Clear, concise explanation of the concept.
- No procedures — only explanation.
- Use concrete nouns and active verbs.
- Maximum 25 words per sentence.

### Key Behaviors
- How the feature works from the user's perspective.
- Inputs, outputs, dependencies, limitations.
- Distinguish verified behavior from assumptions (`[UNVERIFIED]`).

### When to Use
- Scenarios where this feature helps the user.
- Contrast with alternatives if relevant.

### Related Settings
| Setting | Key (from `en.ts`) | Values | Default | Effect |
|---|---|---|---|---|
| [Setting name] | `settingKey` | `value1`, `value2` | `default` | [What it changes] |

### Limitations
- Known constraints, unsupported cases, platform differences.

### Related Topics
- [Link to task topic(s)]
- [Link to reference topic(s)]
- [Link to troubleshooting topic(s)]

---

## Writing Rules (ASD-STE100)

- **Active voice preferred:** "The **preview** renders Markdown" not "Markdown is rendered by the **preview**."
- **Concrete wording:** Name the actor, action, object, result.
- **One main idea per sentence** — split complex sentences.
- **No vague terms:** Replace "some," "various," "appropriate" with specifics.
- **Glossary terms** from `appendix/glossary.md` used consistently.
- **Abbreviations** defined at first use in this topic.
- **UI labels** match `en.ts` exactly when referenced.
- **Code identifiers** use exact casing from source (`themeApp`, `workspace`).

## Review Checklist

- [ ] Purpose states what and why
- [ ] Explanation is complete but minimal
- [ ] Key behaviors verified against source
- [ ] Settings table accurate (keys, values, defaults)
- [ ] Limitations included
- [ ] Sentences ≤ 25 words
- [ ] Paragraphs ≤ 6 sentences
- [ ] Noun clusters ≤ 3 words
- [ ] Glossary terms used correctly
- [ ] Cross-references valid