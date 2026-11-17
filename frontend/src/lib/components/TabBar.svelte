<script lang="ts">
  import { t } from '../i18n/index.svelte';

  export type Tab = {
    id: string;
    title: string;
    path: string;
    dirty?: boolean;
    /** Last known disk mtime (ns) — used to detect external changes. */
    lastModNs?: number;
    /** Last known disk size — paired with lastModNs. */
    lastSize?: number;
  };

  interface Props {
    tabs?: Tab[];
    activeId?: string;
    onselect?: (id: string) => void;
    onclose?: (id: string) => void;
    onnew?: () => void;
  }

  let { tabs = [], activeId = '', onselect, onclose, onnew }: Props = $props();
</script>

<div class="tabbar">
  <div class="tabs">
    {#each tabs as tab (tab.id)}
      <button
        class="tab"
        class:active={tab.id === activeId}
        onclick={() => onselect?.(tab.id)}
        title={tab.path || tab.title}
      >
        <span class="title">{tab.title || t('noFile')}</span>
        {#if tab.dirty}<span class="dot">●</span>{/if}
        <span
          class="x"
          role="button"
          tabindex="0"
          onclick={(e) => { e.stopPropagation(); onclose?.(tab.id); }}
          onkeydown={(e) => { if (e.key === 'Enter') { e.stopPropagation(); onclose?.(tab.id); } }}
        >×</span>
      </button>
    {/each}
  </div>
  <button class="add" title={t('newTab')} onclick={onnew}>+</button>
</div>

<style>
  .tabbar {
    display: flex; align-items: stretch; background: var(--bg-elevated, #161b22);
    border-bottom: 1px solid var(--border, #30363d); height: 34px; min-height: 34px;
  }
  .tabs { display: flex; overflow-x: auto; flex: 1; }
  .tab {
    display: flex; align-items: center; gap: 6px; padding: 0 10px; border: none;
    background: transparent; color: var(--fg-muted, #8b949e); border-right: 1px solid var(--border, #30363d);
    cursor: pointer; font-size: 12px; max-width: 180px;
  }
  .tab.active { background: var(--bg, #0d1117); color: var(--fg, #e6edf3); }
  .title { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .dot { color: #d29922; font-size: 10px; }
  .x { opacity: 0.5; padding: 0 2px; border-radius: 3px; }
  .x:hover { opacity: 1; background: var(--bg-hover, #21262d); }
  .add {
    border: none; background: transparent; color: var(--fg-muted, #8b949e); padding: 0 12px; cursor: pointer; font-size: 16px;
  }
  .add:hover { color: var(--fg, #e6edf3); }
</style>
