<script lang="ts">
	import { t } from '../i18n/index.svelte';
	import * as platform from '../platform';
	import type { FileNode } from '../platform';

	interface Props {
		workspace?: string;
		onworkspace?: (path: string) => void;
		onopenfile?: (path: string, content: string) => void;
	}

	let { workspace = '', onworkspace, onopenfile }: Props = $props();

	let nodes = $state<FileNode[]>([]);
	let openDirs = $state<Record<string, boolean>>({});
	let loading = $state(false);

	async function refresh(rel = '') {
		loading = true;
		try {
			nodes = await platform.listFiles(rel);
		} catch {
			nodes = [];
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (workspace) refresh('');
		else nodes = [];
	});

	async function openFolder() {
		const path = await platform.selectWorkspace();
		if (path) {
			onworkspace?.(path);
			await refresh('');
		}
	}

	async function toggleDir(node: FileNode) {
		const key = node.path;
		if (openDirs[key]) {
			openDirs = { ...openDirs, [key]: false };
			return;
		}
		openDirs = { ...openDirs, [key]: true };
		if (!node.children?.length) {
			const children = await platform.listFiles(node.path);
			nodes = patchChildren(nodes, node.path, children);
		}
	}

	function patchChildren(list: FileNode[], path: string, children: FileNode[]): FileNode[] {
		return list.map((n) => {
			if (n.path === path) return { ...n, children };
			if (n.children) return { ...n, children: patchChildren(n.children, path, children) };
			return n;
		});
	}

	async function openFile(node: FileNode) {
		if (node.isDir) {
			await toggleDir(node);
			return;
		}
		try {
			const content = await platform.readFile(node.path);
			// Prefer absolute when possible
			const abs = workspace ? `${workspace.replace(/[/\\]$/, '')}/${node.path}` : node.path;
			onopenfile?.(abs, content);
		} catch (e) {
			console.error(e);
		}
	}
</script>

<aside class="sidebar">
	<div class="header">
		<span class="title">{t('explorer')}</span>
		<div class="actions">
			<button class="icon" title={t('openFolder')} onclick={openFolder} aria-label={t('openFolder')}>
				<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor">
					<path d="M1.75 2.5a.25.25 0 0 0-.25.25v10.5c0 .138.112.25.25.25h12.5a.25.25 0 0 0 .25-.25V4.75a.25.25 0 0 0-.25-.25h-6.386a.25.25 0 0 1-.2-.1l-.9-1.2a.25.25 0 0 0-.2-.1H1.75Z"></path>
				</svg>
			</button>
		</div>
	</div>
	<div class="body">
		{#if !workspace}
			<div class="empty">
				<p>{t('noFolder')}</p>
				<p class="hint">{t('noFolderHint')}</p>
				<button class="primary" onclick={openFolder}>{t('openFolder')}</button>
			</div>
		{:else if loading && !nodes.length}
			<div class="empty"><p>…</p></div>
		{:else}
			<ul class="tree">
				{#each nodes as node}
					{@render treeItem(node, 0)}
				{/each}
			</ul>
		{/if}
	</div>
</aside>

{#snippet treeItem(node: FileNode, depth: number)}
	<li>
		<button
			class="row"
			class:dir={node.isDir}
			style="padding-left: {8 + depth * 14}px"
			onclick={() => openFile(node)}
		>
			<span class="chev" class:open={openDirs[node.path]}>
				{#if node.isDir}▸{:else}·{/if}
			</span>
			<span class="name">{node.name}</span>
		</button>
		{#if node.isDir && openDirs[node.path] && node.children?.length}
			<ul>
				{#each node.children as child}
					{@render treeItem(child, depth + 1)}
				{/each}
			</ul>
		{/if}
	</li>
{/snippet}

<style>
	.sidebar {
		display: flex;
		flex-direction: column;
		height: 100%;
		width: 100%;
		background: var(--bg-elevated, #161b22);
		border-right: 1px solid var(--border, #30363d);
		color: var(--fg, #e6edf3);
	}
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 8px 12px;
		height: 41px;
		box-sizing: border-box;
		border-bottom: 1px solid var(--border, #30363d);
		font-size: 12px;
		font-weight: 600;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: var(--fg-muted, #8b949e);
	}
	.actions { display: flex; gap: 4px; }
	.icon {
		background: transparent; border: none; cursor: pointer; color: var(--fg-muted, #8b949e);
		font-size: 14px; padding: 2px 6px; border-radius: 4px;
	}
	.icon:hover { background: var(--bg-hover, #21262d); color: var(--fg, #e6edf3); }
	.body { flex: 1; overflow: auto; }
	.empty {
		display: flex; flex-direction: column; align-items: center; justify-content: center;
		gap: 8px; padding: 24px; color: var(--fg-muted, #8b949e); text-align: center; height: 100%;
	}
	.hint { font-size: 12px; opacity: 0.8; margin: 0; }
	.primary {
		margin-top: 8px; background: #238636; color: #fff; border: none; border-radius: 6px;
		padding: 8px 14px; cursor: pointer; font-size: 13px;
	}
	.tree, .tree ul { list-style: none; margin: 0; padding: 0; }
	.row {
		display: flex; align-items: center; gap: 4px; width: 100%;
		background: transparent; border: none; color: inherit; text-align: left;
		padding: 4px 8px; cursor: pointer; font-size: 13px;
	}
	.row:hover { background: var(--bg-hover, #21262d); }
	.chev { width: 12px; color: var(--fg-muted, #8b949e); display: inline-block; transition: transform .15s; }
	.chev.open { transform: rotate(90deg); }
	.name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
</style>
