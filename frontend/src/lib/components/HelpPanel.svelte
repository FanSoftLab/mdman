<script lang="ts">
  import { t } from '../i18n/index.svelte';

  interface Props { open?: boolean; onclose?: () => void; }
  let { open = false, onclose }: Props = $props();

  type Section = 'overview' | 'shortcuts' | 'about' | 'licenses';
  let section = $state<Section>('overview');
  let searchQuery = $state('');

  const licenses = [
    { name: 'Svelte', version: '5.x', description: 'Cybernetically enhanced web apps', license: 'MIT', url: 'https://github.com/sveltejs/svelte' },
    { name: 'SvelteMark', version: 'upstream', description: 'Privacy-first local Markdown editor (UI source)', license: 'MIT', url: 'https://github.com/MasFana/sveltemark' },
    { name: 'CodeMirror', version: '6.x', description: 'Extensible code editor for the browser', license: 'MIT', url: 'https://github.com/codemirror/dev' },
    { name: 'Mermaid', version: '11.x', description: 'Diagramming and charting tool', license: 'MIT', url: 'https://github.com/mermaid-js/mermaid' },
    { name: 'KaTeX', version: '0.16.x', description: 'Fast math typesetting for the web', license: 'MIT', url: 'https://github.com/KaTeX/KaTeX' },
    { name: 'highlight.js', version: '11.x', description: 'Syntax highlighting for code blocks', license: 'BSD-3-Clause', url: 'https://github.com/highlightjs/highlight.js' },
    { name: 'Font Awesome Free', version: '7.x', description: 'Icon set used across the UI', license: 'CC BY 4.0 / OFL / MIT', url: 'https://github.com/fortawesome/font-awesome' },
    { name: 'github-markdown-css', version: '5.x', description: 'GitHub Markdown CSS styles', license: 'MIT', url: 'https://github.com/sindresorhus/github-markdown-css' },
    { name: 'Dexie.js', version: '4.x', description: 'Minimalistic IndexedDB wrapper', license: 'Apache-2.0', url: 'https://github.com/dexie/Dexie.js' },
    { name: 'JSZip', version: '3.x', description: 'Create and read ZIP archives', license: 'MIT / GPL-3.0', url: 'https://github.com/Stuk/jszip' },
    { name: 'lucide-svelte', version: '0.5x', description: 'Lucide icons for Svelte', license: 'ISC', url: 'https://github.com/lucide-icons/lucide' },
    { name: 'github-slugger', version: '2.x', description: 'Slug generation for headings', license: 'ISC', url: 'https://github.com/Flet/github-slugger' },
    { name: 'unified', version: '11.x', description: 'Content transformation platform (Markdown pipeline)', license: 'MIT', url: 'https://github.com/unifiedjs/unified' },
    { name: 'remark-gfm', version: '4.x', description: 'GitHub Flavored Markdown support', license: 'MIT', url: 'https://github.com/remarkjs/remark-gfm' },
    { name: 'remark-math', version: '6.x', description: 'Math blocks support', license: 'MIT', url: 'https://github.com/remarkjs/remark-math' },
    { name: 'rehype-highlight', version: '7.x', description: 'Syntax highlighting in HAST', license: 'MIT', url: 'https://github.com/rehypejs/rehype-highlight' },
    { name: 'rehype-katex', version: '7.x', description: 'Render math via KaTeX', license: 'MIT', url: 'https://github.com/remarkjs/remark-math' },
    { name: 'rehype-slug', version: '6.x', description: 'Add ids to headings', license: 'MIT', url: 'https://github.com/rehypejs/rehype-slug' },
    { name: 'Wails', version: '2.x', description: 'Desktop apps with Go and web technologies', license: 'MIT', url: 'https://github.com/wailsapp/wails' },
    { name: 'Vite', version: '6.x', description: 'Frontend build tool', license: 'MIT', url: 'https://github.com/vitejs/vite' },
    { name: 'TypeScript', version: '5.x', description: 'Typed JavaScript', license: 'Apache-2.0', url: 'https://github.com/microsoft/TypeScript' },
  ];

  const sections: { id: Section; labelKey: string }[] = [
    { id: 'overview', labelKey: 'helpOverview' },
    { id: 'shortcuts', labelKey: 'helpShortcuts' },
    { id: 'about', labelKey: 'helpAbout' },
    { id: 'licenses', labelKey: 'helpLicenses' },
  ];

  function match(text: string): boolean {
    if (!searchQuery.trim()) return true;
    return text.toLowerCase().includes(searchQuery.toLowerCase());
  }
</script>

{#if open}
  <div class="help-overlay" role="presentation" onclick={onclose}>
    <div class="help-dialog" role="dialog" aria-label={t('helpTitle')} onclick={(e) => e.stopPropagation()}>
      <header>
        <h2>{t('helpTitle')}</h2>
        <input class="search" type="search" placeholder={t('search')} bind:value={searchQuery} />
        <button class="icon-close" onclick={onclose} aria-label={t('close')}>✕</button>
      </header>
      <div class="help-body">
        <nav class="help-nav">
          {#each sections as s}
            <button class:active={section === s.id} onclick={() => (section = s.id)}>{t(s.labelKey)}</button>
          {/each}
        </nav>
        <div class="help-content">
          {#if section === 'overview'}
            <h3>{t('helpOverview')}</h3>
            {#if match(t('helpOverviewP1'))}<p>{t('helpOverviewP1')}</p>{/if}
            {#if match(t('helpOverviewP2'))}<p>{t('helpOverviewP2')}</p>{/if}
          {:else if section === 'shortcuts'}
            <h3>{t('helpShortcuts')}</h3>
            <table>
              <tbody>
                <tr><td><kbd>Ctrl</kbd>+<kbd>S</kbd></td><td>{t('shortcutSave')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>S</kbd></td><td>{t('saveAs')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>O</kbd></td><td>{t('shortcutOpen')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>F</kbd></td><td>{t('shortcutFind')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>H</kbd></td><td>{t('shortcutFindReplace')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>T</kbd></td><td>{t('shortcutNewTab')}</td></tr>
                <tr><td><kbd>Ctrl</kbd>+<kbd>W</kbd></td><td>{t('shortcutCloseTab')}</td></tr>
              </tbody>
            </table>
          {:else if section === 'about'}
            <h3>{t('helpAbout')}</h3>
            <div class="about-row">
              <div class="about-icon" aria-hidden="true">M</div>
              <div>
                <div class="about-title">mdman</div>
                <div class="about-meta">MIT · desktop SvelteMark</div>
              </div>
            </div>
            <p>{t('helpAboutP1')}</p>
            <p>{t('helpAboutP2')}</p>
            <p><a href="https://github.com/MasFana/sveltemark" target="_blank" rel="noreferrer">SvelteMark on GitHub</a></p>
          {:else}
            <h3>{t('helpLicenses')}</h3>
            <ul class="lic">
              {#each licenses.filter((l) => match(l.name + ' ' + l.description + ' ' + l.license)) as l}
                <li>
                  <div class="lic-main">
                    <a href={l.url} target="_blank" rel="noreferrer">{l.name}</a>
                    <span class="ver">{l.version}</span>
                    <span class="lic-name">{l.license}</span>
                  </div>
                  <p class="desc">{l.description}</p>
                </li>
              {/each}
            </ul>
          {/if}
        </div>
      </div>
      <footer>
        <button class="close-btn" onclick={onclose}>{t('close')}</button>
      </footer>
    </div>
  </div>
{/if}

<style>
  .help-overlay {
    position: fixed; inset: 0; background: rgba(0,0,0,.55); z-index: 1100;
    display: flex; align-items: center; justify-content: center; padding: 24px;
  }
  .help-dialog {
    width: min(960px, 96vw); height: min(680px, 92vh);
    background: var(--bg-elevated, #161b22);
    color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d);
    border-radius: 10px;
    display: flex; flex-direction: column; overflow: hidden;
    box-shadow: 0 16px 48px rgba(0,0,0,.45);
  }
  header {
    display: flex; align-items: center; gap: 12px;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border, #30363d);
    background: var(--bg-elevated, #161b22);
  }
  h2 { margin: 0; font-size: 15px; font-weight: 600; flex: 0 0 auto; }
  .search {
    flex: 1; max-width: 320px; margin-left: auto;
    background: var(--bg, #0d1117);
    color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d);
    border-radius: 6px; padding: 7px 12px; font-size: 13px;
  }
  .icon-close {
    background: transparent; border: none; color: var(--fg-muted, #8b949e);
    cursor: pointer; font-size: 16px; padding: 4px 8px; border-radius: 4px;
  }
  .icon-close:hover { background: var(--bg-hover, #21262d); color: var(--fg, #e6edf3); }
  .help-body { display: flex; flex: 1; min-height: 0; }
  .help-nav {
    width: 200px;
    border-right: 1px solid var(--border, #30363d);
    padding: 10px 8px;
    display: flex; flex-direction: column; gap: 2px; overflow: auto;
    background: var(--bg-elevated, #161b22);
  }
  .help-nav button {
    text-align: left; background: transparent; border: none;
    color: var(--fg-muted, #8b949e);
    padding: 9px 12px; border-radius: 6px; cursor: pointer; font-size: 13px;
  }
  .help-nav button:hover { background: var(--bg-hover, #21262d); color: var(--fg, #e6edf3); }
  .help-nav button.active {
    background: var(--bg-hover, #21262d);
    color: var(--fg, #e6edf3);
    font-weight: 600;
  }
  .help-content {
    flex: 1; overflow: auto; padding: 18px 22px;
    font-size: 13px; line-height: 1.6;
    background: var(--bg, #0d1117);
  }
  h3 { margin: 0 0 14px; font-size: 16px; font-weight: 600; }
  table { width: 100%; border-collapse: collapse; }
  td {
    padding: 8px 0;
    border-bottom: 1px solid var(--border, #30363d);
  }
  td:first-child { width: 40%; color: var(--fg-muted, #8b949e); }
  kbd {
    background: var(--bg-elevated, #161b22);
    border: 1px solid var(--border, #30363d);
    border-radius: 4px; padding: 1px 6px; font-size: 12px;
  }
  .about-row { display: flex; gap: 14px; align-items: center; margin-bottom: 16px; }
  .about-icon {
    width: 48px; height: 48px; border-radius: 10px;
    background: linear-gradient(135deg, #238636, #1f6feb);
    display: flex; align-items: center; justify-content: center;
    font-weight: 800; font-size: 20px; color: #fff;
  }
  .about-title { font-size: 18px; font-weight: 700; }
  .about-meta { font-size: 12px; color: var(--fg-muted, #8b949e); margin-top: 2px; }
  .lic { list-style: none; margin: 0; padding: 0; }
  .lic li {
    padding: 12px 0;
    border-bottom: 1px solid var(--border, #30363d);
  }
  .lic-main { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
  .lic a { color: #58a6ff; text-decoration: none; font-weight: 600; min-width: 140px; }
  .lic a:hover { text-decoration: underline; }
  .ver, .lic-name { color: var(--fg-muted, #8b949e); font-size: 12px; }
  .lic-name {
    background: var(--bg-elevated, #161b22);
    border: 1px solid var(--border, #30363d);
    border-radius: 4px; padding: 1px 6px;
  }
  .desc { margin: 6px 0 0; color: var(--fg-muted, #8b949e); font-size: 12px; }
  footer {
    display: flex; justify-content: flex-end;
    padding: 12px 16px;
    border-top: 1px solid var(--border, #30363d);
    background: var(--bg-elevated, #161b22);
  }
  .close-btn {
    background: var(--bg-hover, #21262d);
    color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d);
    border-radius: 6px; padding: 7px 16px; cursor: pointer; font-size: 13px;
  }
  .close-btn:hover { border-color: #58a6ff; }
  a { color: #58a6ff; }
</style>
