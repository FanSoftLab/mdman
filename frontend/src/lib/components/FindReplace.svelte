<script lang="ts">
  import { t } from '../i18n/index.svelte';

  interface Props {
    open?: boolean;
    onclose?: () => void;
    onfind?: (query: string, opts: { caseSensitive: boolean; backwards?: boolean }) => void;
    onreplace?: (query: string, replacement: string, opts: { caseSensitive: boolean; all?: boolean }) => void;
  }

  let { open = false, onclose, onfind, onreplace }: Props = $props();
  let query = $state('');
  let replacement = $state('');
  let caseSensitive = $state(false);
  let showReplace = $state(true);

  function find(backwards = false) {
    if (!query) return;
    onfind?.(query, { caseSensitive, backwards });
  }
  function replaceOne() {
    if (!query) return;
    onreplace?.(query, replacement, { caseSensitive, all: false });
  }
  function replaceAll() {
    if (!query) return;
    onreplace?.(query, replacement, { caseSensitive, all: true });
  }

  function onKey(e: KeyboardEvent) {
    if (e.key === 'Escape') onclose?.();
    if (e.key === 'Enter') {
      e.preventDefault();
      if (e.shiftKey) find(true);
      else find(false);
    }
  }
</script>

{#if open}
  <div class="bar">
    <input
      class="q"
      placeholder={t('search')}
      bind:value={query}
      onkeydown={onKey}
    />
    {#if showReplace}
      <input
        class="q"
        placeholder={t('replace')}
        bind:value={replacement}
        onkeydown={onKey}
      />
    {/if}
    <label class="chk"><input type="checkbox" bind:checked={caseSensitive} /> {t('caseSensitive')}</label>
    <button onclick={() => find(true)}>{t('findPrev')}</button>
    <button onclick={() => find(false)}>{t('findNext')}</button>
    {#if showReplace}
      <button onclick={replaceOne}>{t('replaceOne')}</button>
      <button onclick={replaceAll}>{t('replaceAll')}</button>
    {/if}
    <button class="icon" onclick={onclose} aria-label={t('close')}>✕</button>
  </div>
{/if}

<style>
  .bar {
    display: flex; flex-wrap: wrap; gap: 6px; align-items: center;
    padding: 6px 10px; background: var(--bg-elevated, #161b22);
    border-bottom: 1px solid var(--border, #30363d); font-size: 12px;
  }
  .q {
    background: var(--bg, #0d1117); color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d); border-radius: 5px;
    padding: 5px 8px; min-width: 140px;
  }
  .chk { display: flex; align-items: center; gap: 4px; color: var(--fg-muted, #8b949e); }
  button {
    background: var(--bg, #0d1117); color: var(--fg, #e6edf3);
    border: 1px solid var(--border, #30363d); border-radius: 5px;
    padding: 4px 8px; cursor: pointer; font-size: 12px;
  }
  button.icon { border: none; background: transparent; color: var(--fg-muted, #8b949e); }
</style>
