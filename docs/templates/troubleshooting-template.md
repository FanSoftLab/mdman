# Troubleshooting Template

**Content Type:** Troubleshooting (Problem Resolution)
**Standard:** ISO/IEC/IEEE 26514 §2.9 + ASD-STE100
**Max sentence length:** 20 words (procedural steps), 25 words (descriptive)
**Max sentences per paragraph:** 6

---

## [Problem Title: Observable Symptom]

*Example: "App Won't Start on Linux", "Preview Shows Black Screen", "Files Not Saving"*

### Symptom
One sentence: what the user sees (error message, blank screen, missing data, crash).

### Conditions
- When does this occur? (OS, version, configuration, user action)
- Be specific: "Linux with Wayland and NVIDIA GPU" not "some Linux systems."

### Verified Cause(s)
| Cause | Evidence | Confidence |
|---|---|---|
| [Root cause] | [Source: code, logs, reproduction] | Verified / Hypothesis |

- Distinguish **confirmed causes** from **hypotheses**.
- Mark hypotheses with `[HYPOTHESIS]`.
- Do not present guesses as confirmed diagnoses.

### Diagnostic Checks
Steps the user can take to confirm the cause.
1. [Check 1 — e.g., "Run `ldd mdman` and verify `libwebkit2gtk-4.0.so` is found."]
2. [Check 2]
3. [Check N]

### Corrective Actions
Ordered from least to most invasive. Each action is a numbered procedure.

#### Action 1: [Least Invasive Fix]
1. [Step 1]
2. [Step 2]
**Expected Result:** [What the user sees if this works.]

#### Action 2: [Next Fix]
1. [Step 1]
2. [Step 2]
**Expected Result:** [What the user sees if this works.]

#### Action N: [Most Invasive / Last Resort]
1. [Step 1]
2. [Step 2]
**Expected Result:** [What the user sees if this works.]

### Expected Result After Correction
The symptom is resolved. Describe the normal behavior (e.g., "The app starts and shows the **welcome screen**.").

### Escalation
If the problem persists after all actions:
- Collect: OS, version (`mdman --version` or About dialog), steps to reproduce, logs.
- Report: [GitHub Issues URL] with the collected information.

### Prevention
- Configuration or usage changes that avoid the problem.

### Related Topics
- [Link to concept topic]
- [Link to reference topic (settings, CLI)]
- [Link to task topic]

---

## Writing Rules (ASD-STE100)

- **State observable symptom first** — not the internal error.
- **Verified causes only** — mark hypotheses clearly.
- **Diagnostic checks** are testable (commands, UI observations).
- **Corrective actions** use task-template format (numbered steps, ≤20 words/sentence).
- **Distinguish severity** — escalation path for unresolved issues.
- **Glossary terms** from `appendix/glossary.md` used consistently.
- **Abbreviations** defined at first use.
- **UI labels** match `en.ts` exactly.

## Review Checklist

- [ ] Symptom is user-observable (not internal error)
- [ ] Conditions are specific (OS, version, config)
- [ ] Causes marked Verified / Hypothesis
- [ ] Diagnostic checks are actionable
- [ ] Corrective actions are numbered procedures
- [ ] Expected result after each action
- [ ] Escalation path provided
- [ ] Prevention advice included
- [ ] Sentences ≤ 20 words (steps), ≤ 25 words (descriptive)
- [ ] Glossary terms used correctly
- [ ] Cross-references valid