<script lang="ts">
  import { t, getLocale, setLocale, type Locale } from '../i18n/index.svelte';

  interface Props {
    open?: boolean;
    themeApp?: 'dark' | 'light' | 'system';
    themePreview?: 'dark' | 'light' | 'system';
    recentMax?: number;
    onclose?: () => void;
    onapply?: (v: {
      themeApp: 'dark' | 'light' | 'system';
      themePreview: 'dark' | 'light' | 'system';
      locale: Locale;
      recentMax: number;
    }) => void;
  }

  let {
    open = false,
    themeApp = $bindable('system'),
    themePreview = $bindable('system'),
    recentMax = $bindable(20),
    onclose,
    onapply,
  }: Props = $props();

  let locale = $state<Locale>(getLocale());

  $effect(() => {
    if (open) locale = getLocale();
  });

  function apply() {
    setLocale(locale);
    onapply?.({ themeApp, themePreview, locale, recentMax });
  }
</script>

{#if open}
  <aside class="side-panel settings-panel">
    <div class="panel-header">
      <span>{t('settings')}</span>
      <button class="toc-close" onclick={onclose} aria-label={t('close')}>
        <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
          <path d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"/>
        </svg>
      </button>
    </div>
    <div class="panel-body">
      <label>
        <span>{t('themeApp')}</span>
        <select bind:value={themeApp} onchange={apply}>
          <option value="dark">{t('themeDark')}</option>
          <option value="light">{t('themeLight')}</option>
          <option value="system">{t('themeSystem')}</option>
        </select>
      </label>
      <label>
        <span>{t('themePreview')}</span>
        <select bind:value={themePreview} onchange={apply}>
          <option value="dark">{t('themeDark')}</option>
          <option value="light">{t('themeLight')}</option>
          <option value="system">{t('themeSystem')}</option>
        </select>
      </label>
      <label>
        <span>{t('language')}</span>
        <select bind:value={locale} onchange={apply}>
          <option value="en">{t('langEnglish')}</option>
          <option value="uk">{t('langUkrainian')}</option>
        </select>
      </label>
      <label>
        <span>{t('recentMax')}</span>
        <input type="number" min="5" max="100" bind:value={recentMax} onchange={apply} />
      </label>
    </div>
  </aside>
{/if}

<style>
  .side-panel {
    position: relative;
    flex-shrink: 0;
    width: 280px;
    height: 100%;
    background: var(--bg-elevated, #161b22);
    border-left: 1px solid var(--border, #30363d);
    display: flex;
    flex-direction: column;
    color: var(--fg, #e6edf3);
  }
  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 12px;
    border-bottom: 1px solid var(--border, #30363d);
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--fg-muted, #8b949e);
  }
  .toc-close {
    background: transparent;
    border: none;
    color: var(--fg-muted, #8b949e);
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    display: flex;
  }
  .toc-close:hover { background: var(--bg-hover, #21262d); color: var(--fg, #e6edf3); }
  .panel-body {
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 14px;
    overflow: auto;
  }
  label {
    display: flex;
    flex-direction: column;
    gap: 6px;
    font-size: 12px;
    color: var(--fg-muted, #8b949e);
  }
  select, input {
    background: var(--bg, #0d1117);
    color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d);
    border-radius: 6px;
    padding: 7px 10px;
    font-size: 13px;
  }
</style>
