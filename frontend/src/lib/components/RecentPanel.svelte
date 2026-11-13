<script lang="ts">
  import { t } from '../i18n/index.svelte';
  import type { RecentFile } from '../platform';

  interface Props {
    open?: boolean;
    files?: RecentFile[];
    workspace?: string;
    onclose?: () => void;
    onopen?: (path: string) => void;
    onclear?: () => void;
  }
  let { open = false, files = [], workspace = '', onclose, onopen, onclear }: Props = $props();

  function rel(path: string): string {
    if (workspace && path.startsWith(workspace)) {
      let r = path.slice(workspace.length);
      if (r.startsWith('/') || r.startsWith('\\')) r = r.slice(1);
      return r || path;
    }
    return path;
  }
  function fmt(ts: number): string {
    try { return new Date(ts).toLocaleString(); } catch { return ''; }
  }
</script>

{#if open}
  <aside class="side-panel recent-panel">
    <div class="panel-header">
      <span>{t('recent')}</span>
      <div class="hdr-actions">
        {#if files.length}
          <button class="link" onclick={onclear}>{t('clearRecent')}</button>
        {/if}
        <button class="toc-close" onclick={onclose} aria-label={t('close')}>
          <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14">
            <path d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"/>
          </svg>
        </button>
      </div>
    </div>
    <div class="panel-body">
      {#if !files.length}
        <p class="empty">{t('emptyRecent')}</p>
      {:else}
        <ul>
          {#each files as f}
            <li>
              <button class="item" onclick={() => onopen?.(f.path)}>
                <span class="name">{f.name}</span>
                <span class="meta">{fmt(f.openedAt)}</span>
                <span class="path" title={f.path}>{rel(f.path)}</span>
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </div>
  </aside>
{/if}

<style>
  .side-panel {
    flex-shrink: 0;
    width: 300px;
    height: 100%;
    background: var(--bg-elevated, #161b22);
    border-left: 1px solid var(--border, #30363d);
    display: flex;
    flex-direction: column;
    color: var(--fg, #e6edf3);
  }
  .panel-header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 10px 12px; border-bottom: 1px solid var(--border, #30363d);
    font-size: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em;
    color: var(--fg-muted, #8b949e);
  }
  .hdr-actions { display: flex; align-items: center; gap: 8px; }
  .link { background: none; border: none; color: var(--fg-muted, #8b949e); cursor: pointer; font-size: 11px; }
  .link:hover { color: #58a6ff; }
  .toc-close { background: transparent; border: none; color: var(--fg-muted, #8b949e); cursor: pointer; padding: 4px; border-radius: 4px; display: flex; }
  .toc-close:hover { background: var(--bg-hover, #21262d); }
  .panel-body { overflow: auto; flex: 1; padding: 6px; }
  ul { list-style: none; margin: 0; padding: 0; }
  .item {
    width: 100%; text-align: left; background: transparent; border: none; color: inherit;
    padding: 8px 10px; border-radius: 6px; cursor: pointer;
    display: grid; grid-template-columns: 1fr auto; gap: 2px 8px;
  }
  .item:hover { background: var(--bg-hover, #21262d); }
  .name { font-size: 13px; font-weight: 500; }
  .meta { font-size: 11px; color: var(--fg-muted, #8b949e); }
  .path { grid-column: 1 / -1; font-size: 11px; color: var(--fg-muted, #8b949e); font-family: ui-monospace, monospace;
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .empty { padding: 20px; text-align: center; color: var(--fg-muted, #8b949e); font-size: 13px; }
</style>
