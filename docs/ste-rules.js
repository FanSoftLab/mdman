/**
 * Custom markdownlint rules for ASD-STE100 compliance
 * Usage: markdownlint --rules docs/ste-rules.js docs/**/*.md
 */

const VAGUE_TERMS = [
  'appropriate',
  'as needed',
  'some time',
  'some',
  'various',
  'certain',
  'particular',
  'specific',
  'proper',
  'correct',
  'suitable',
  'adequate',
  'sufficient',
  'reasonable',
  'normal',
  'standard',
  'typical',
  'usual',
  'common',
  'general',
  'basic',
  'simple',
  'easy',
  'fast',
  'quick',
  'robust',
  'secure',
  'reliable',
  'efficient',
  'optimal',
  'best practice',
  'industry standard'
];

const PASSIVE_PATTERNS = [
  /\bis\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bare\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bwas\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bwere\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bhas\s+been\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bhave\s+been\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bhad\s+been\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bwill\s+be\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bshould\s+be\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bmust\s+be\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i,
  /\bcan\s+be\s+(done|performed|executed|completed|finished|processed|handled|managed|created|generated|produced|made|built|constructed|formed|established|implemented|applied|used|utilized|employed|adopted)\b/i
];

function countWords(text) {
  return text.trim().split(/\s+/).filter(w => w.length > 0).length;
}

function splitSentences(text) {
  return text.split(/(?<=[.!?])\s+(?=[A-Z])/).filter(s => s.trim().length > 0);
}

function isProceduralLine(line) {
  return /^\d+\.\s/.test(line.trim()) || /^-\s/.test(line.trim());
}

function checkProceduralSentenceLength(params, onError) {
  const { lines, frontMatterLines } = params;
  let inProcedure = false;
  let procedureStart = -1;

  lines.forEach((line, index) => {
    const trimmed = line.trim();

    if (/^### Procedure$|^## Procedure$|^# Procedure$/.test(trimmed)) {
      inProcedure = true;
      procedureStart = index;
      return;
    }

    if (inProcedure && /^##|^#/.test(trimmed) && index > procedureStart) {
      inProcedure = false;
      return;
    }

    if (inProcedure && isProceduralLine(line)) {
      const sentences = splitSentences(trimmed.replace(/^\d+\.\s+|^-\s+/, ''));
      sentences.forEach(sentence => {
        const words = countWords(sentence);
        if (words > 20) {
          onError({
            lineNumber: index + 1,
            detail: `Procedural sentence has ${words} words (max 20): "${sentence.substring(0, 80)}..."`,
            range: [1, line.length]
          });
        }
      });
    }
  });
}

function checkDescriptiveSentenceLength(params, onError) {
  const { lines, frontMatterLines } = params;
  let inProcedure = false;
  let procedureStart = -1;

  lines.forEach((line, index) => {
    const trimmed = line.trim();

    if (/^### Procedure$|^## Procedure$|^# Procedure$/.test(trimmed)) {
      inProcedure = true;
      procedureStart = index;
      return;
    }

    if (inProcedure && /^##|^#/.test(trimmed) && index > procedureStart) {
      inProcedure = false;
      return;
    }

    if (!inProcedure && trimmed.length > 0 && !/^#/.test(trimmed) && !/^\|/.test(trimmed) && !/^```/.test(trimmed)) {
      const sentences = splitSentences(trimmed);
      sentences.forEach(sentence => {
        const words = countWords(sentence);
        if (words > 25) {
          onError({
            lineNumber: index + 1,
            detail: `Descriptive sentence has ${words} words (max 25): "${sentence.substring(0, 80)}..."`,
            range: [1, line.length]
          });
        }
      });
    }
  });
}

function checkPassiveVoice(params, onError) {
  const { lines, frontMatterLines } = params;
  let inProcedure = false;
  let procedureStart = -1;

  lines.forEach((line, index) => {
    const trimmed = line.trim();

    if (/^### Procedure$|^## Procedure$|^# Procedure$/.test(trimmed)) {
      inProcedure = true;
      procedureStart = index;
      return;
    }

    if (inProcedure && /^##|^#/.test(trimmed) && index > procedureStart) {
      inProcedure = false;
      return;
    }

    if (inProcedure && isProceduralLine(line)) {
      PASSIVE_PATTERNS.forEach(pattern => {
        if (pattern.test(trimmed)) {
          onError({
            lineNumber: index + 1,
            detail: `Passive voice detected in procedural step: "${trimmed.substring(0, 80)}..."`,
            range: [1, line.length]
          });
        }
      });
    }
  });
}

function checkVagueTerms(params, onError) {
  const { lines, frontMatterLines } = params;

  lines.forEach((line, index) => {
    const trimmed = line.trim();
    if (trimmed.length === 0 || /^#/.test(trimmed) || /^```/.test(trimmed)) return;

    VAGUE_TERMS.forEach(term => {
      const regex = new RegExp(`\\b${term}\\b`, 'i');
      if (regex.test(trimmed)) {
        onError({
          lineNumber: index + 1,
          detail: `Vague term "${term}" detected: "${trimmed.substring(0, 80)}..."`,
          range: [1, line.length]
        });
      }
    });
  });
}

function checkGlossaryTerms(params, onError) {
  const { lines, frontMatterLines } = params;
  const disallowedVariants = {
    'directory': 'folder',
    'file': 'note',
    'document': 'note',
    'markdown file': 'note',
    'zen mode': 'auto-hide UI',
    'distraction-free mode': 'auto-hide UI',
    'export': 'backup',
    'data dump': 'backup',
    'right-click menu': 'context menu',
    'Mermaid chart': 'diagram',
    'graph': 'diagram',
    'editing pane': 'editor',
    'markdown editor': 'editor',
    'dir': 'folder',
    'header': 'heading',
    'title': 'heading',
    'code fence': 'code block',
    'fenced code': 'code block',
    'code span': 'inline code',
    'math renderer': 'KaTeX',
    'hotkey': 'keyboard shortcut',
    'key binding': 'keyboard shortcut',
    'locale': 'language',
    'UI language': 'language',
    'hyperlink': 'link',
    'display math': 'math block',
    'formula block': 'math block',
    'diagram tool': 'Mermaid',
    'preview pane': 'preview',
    'rendered view': 'preview',
    'formatting bar': 'toolbar',
    'action bar': 'toolbar',
    'explorer': 'sidebar',
    'file explorer': 'sidebar',
    'strikeout': 'strikethrough',
    'crossed out': 'strikethrough',
    'checklist': 'task list',
    'todo list': 'task list',
    'outline': 'table of contents (TOC)',
    'editor tab': 'tab',
    'read-only mode': 'view-only mode',
    'preview mode': 'view-only mode',
    'root folder': 'workspace',
    'project folder': 'workspace'
  };

  lines.forEach((line, index) => {
    const trimmed = line.trim();
    if (trimmed.length === 0 || /^#/.test(trimmed) || /^```/.test(trimmed)) return;

    Object.entries(disallowedVariants).forEach(([variant, preferred]) => {
      const regex = new RegExp(`\\b${variant}\\b`, 'i');
      if (regex.test(trimmed)) {
        onError({
          lineNumber: index + 1,
          detail: `Disallowed variant "${variant}" — use "${preferred}" instead: "${trimmed.substring(0, 80)}..."`,
          range: [1, line.length]
        });
      }
    });
  });
}

module.exports = {
  checkProceduralSentenceLength,
  checkDescriptiveSentenceLength,
  checkPassiveVoice,
  checkVagueTerms,
  checkGlossaryTerms
};