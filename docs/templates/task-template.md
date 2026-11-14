# Task Template

**Content Type:** Task (Procedure)
**Standard:** ISO/IEC/IEEE 26514 §2.6 + ASD-STE100
**Max sentence length:** 20 words (procedural)
**Max sentences per paragraph:** 1 (per step)

---

## [Task Title: Verb + Object]

*Example: "Open a Folder", "Create a Note", "Export a Backup"*

### Goal
One sentence: what the user achieves by completing this task.

### Prerequisites
- List each prerequisite as a bullet.
- Include: workspace open, file permissions, specific settings, required file types.
- If none, write "None."

### Procedure

1. **Step 1:** [Action on specific UI element]. [Result if not obvious].
   - Use imperative mood. One main action per step.
   - Name the exact UI label from `en.ts` (e.g., "Click **Open Folder** in the **toolbar**").
   - Maximum 20 words per sentence.

2. **Step 2:** [Action]. [Result].

3. **Step N:** [Action]. [Result].

### Expected Result
What the user sees when the task succeeds. Be specific (e.g., "The folder appears in the **sidebar**. The first note opens in the **editor**.").

### Verification
How the user confirms success (e.g., "Verify the note title appears in the **tab bar**.").

### Recovery / Common Errors
| Error / Situation | Cause | Corrective Action |
|---|---|---|
| [Error message or symptom] | [Verified cause] | [Step to resolve] |

### Limitations
- Any constraints (e.g., "Only works when a **workspace** is open.")
- Platform-specific behavior.

### Related Topics
- [Link to concept topic]
- [Link to reference topic]
- [Link to troubleshooting topic]

---

## Writing Rules (ASD-STE100)

- **Active voice:** "Click **Save**" not "**Save** is clicked."
- **Concrete subject:** "The **editor** saves the file" not "It saves the file."
- **Explicit conditions:** "If the **workspace** is not set, the command fails" not "When possible."
- **One action per step** — split compound steps.
- **No vague terms:** Replace "appropriate," "as needed," "some time" with specifics or mark `[UNVERIFIED]`.
- **UI labels** must match `en.ts` exactly (capitalization, spacing).
- **Glossary terms** from `appendix/glossary.md` used consistently.
- **Abbreviations** defined at first use in this topic.

## Review Checklist

- [ ] Goal states user outcome
- [ ] All prerequisites listed
- [ ] Steps numbered, sequential, one action each
- [ ] UI labels match `en.ts`
- [ ] Expected result verifiable
- [ ] Recovery covers common failures
- [ ] Sentences ≤ 20 words
- [ ] No passive voice in steps
- [ ] Glossary terms used correctly
- [ ] Cross-references valid