#!/usr/bin/env node
// Pipeline benchmark for mdmark markdown rendering. Runs in Docker (no DOM).
// Usage: node --expose-gc bench/bench.mjs --out /out/frontend-bench.json
//
// Scenarios:
//   fullRender-100k / fullRender-1M  — full document render via processMarkdownBlocks
//   tocExtract-100k / tocExtract-1M  — extractTableOfContents (full AST pass)
//   editOneBlock-200                 — re-render after editing 1 of 200 blocks
//                                      (uses renderDocument(prev, doc) when available,
//                                       falls back to full processMarkdownBlocks)
//   renderOneBlock                   — single block render
//
// Always writes golden.json (block HTML dump of synth-100k) next to --out
// so bench.sh can verify output equivalence between runs.

import { build } from 'esbuild';
import { mkdir, mkdtemp, writeFile } from 'node:fs/promises';
import { createHash } from 'node:crypto';
import { tmpdir } from 'node:os';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

const here = path.dirname(fileURLToPath(import.meta.url));
const srcRoot = path.resolve(here, '..');

const args = process.argv.slice(2);
function argOf(name, def) {
	const i = args.indexOf(name);
	return i >= 0 ? args[i + 1] : def;
}
const outPath = argOf('--out', 'frontend-bench.json');

// ---- esbuild bundle of src/lib/markdown.ts (pure JS deps, node-safe) ----
async function bundleMarkdown() {
	const dir = await mkdtemp(path.join(tmpdir(), 'mdmark-bench-'));
	const outfile = path.join(dir, 'markdown.bundle.mjs');
	await build({
		entryPoints: [path.join(srcRoot, 'src/lib/markdown.ts')],
		outfile,
		bundle: true,
		format: 'esm',
		platform: 'node',
		target: 'node20',
		sourcemap: false,
		logLevel: 'silent'
	});
	return import(pathToFileURL(outfile).href);
}

// ---- deterministic fixture generation ----
function mulberry32(seed) {
	let a = seed >>> 0;
	return function () {
		a |= 0;
		a = (a + 0x6d2b79f5) | 0;
		let t = Math.imul(a ^ (a >>> 15), 1 | a);
		t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
		return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
	};
}

const TITLES = ['Overview', 'Getting Started', 'Performance', 'Notes', 'Overview', 'Advanced Usage'];

const JS_SNIPPET = [
	'export function debounce(fn, ms) {',
	'  let t = null;',
	'  return (...args) => {',
	'    clearTimeout(t);',
	'    t = setTimeout(() => fn(...args), ms);',
	'  };',
	'}'
].join('\n');

function makeSection(i) {
	const parts = [];
	const title = TITLES[i % TITLES.length]; // deliberate duplicates for slug handling
	parts.push(`## ${title} (section ${i})`);
	parts.push(
		'Lorem ipsum **dolor** sit amet, *consectetur* adipiscing elit, sed do `eiusmod` tempor ' +
			'incididunt ut labore et [dolore](https://example.com/ Ut enim ad minim veniam, quis ' +
			'nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Inline ' +
			'math $a^2 + b^2 = c^2$ appears here.'
	);
	if (i % 3 === 0) parts.push('```js\n' + JS_SNIPPET + '\n```');
	if (i % 4 === 0) {
		parts.push(
			'| Column A | Column B | Column C |\n|----------|----------|----------|\n' +
				'| 1 | alpha | 0.25 |\n| 2 | beta  | 0.50 |\n| 3 | gamma | 0.75 |'
		);
	}
	if (i % 5 === 0) parts.push('$$\nE = mc^2 + \\int_0^\\infty e^{-x} dx\n$$');
	if (i % 7 === 0) {
		parts.push(
			'```mermaid\ngraph TD\nA[Start] --> B{Check}\nB -->|yes| C[Process]\nB -->|no| D[Skip]\nC --> E[End]\nD --> E\n```'
		);
	}
	if (i % 6 === 0) parts.push(`![diagram](./assets/img-${i}.png)`);
	return parts.join('\n\n');
}

// Generate a document whose length is at least targetBytes.
function genDoc(targetBytes) {
	const chunks = [];
	let len = 0;
	for (let i = 0; len < targetBytes; i++) {
		const s = makeSection(i);
		chunks.push(s);
		len += s.length + 2;
	}
	return chunks.join('\n\n');
}

// ---- measurement ----
let sink = 0;
async function bench(name, fn, runs) {
	for (let i = 0; i < 2; i++) await fn(); // warmup
	if (typeof global.gc === 'function') global.gc();
	const samples = [];
	let heapAfter = 0;
	for (let i = 0; i < runs; i++) {
		if (typeof global.gc === 'function') global.gc();
		const t0 = performance.now();
		await fn();
		samples.push(performance.now() - t0);
		if (typeof global.gc === 'function') global.gc();
		heapAfter = process.memoryUsage().heapUsed;
	}
	samples.sort((a, b) => a - b);
	const median = samples[Math.floor(samples.length / 2)];
	return {
		name,
		runs: samples.length,
		medianMs: +median.toFixed(2),
		minMs: +samples[0].toFixed(2),
		maxMs: +samples[samples.length - 1].toFixed(2),
		heapUsedMB: +(heapAfter / 1048576).toFixed(1)
	};
}

const md = await bundleMarkdown();
const hasIncremental = typeof md.renderDocument === 'function';

const doc100k = genDoc(100 * 1024);
const doc1M = genDoc(1024 * 1024);

const para200 = Array.from(
	{ length: 200 },
	(_, i) =>
		`Paragraph ${i}: lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`
);
const doc200 = para200.join('\n\n');
const edited200 = para200.map((b, i) => (i === 100 ? b + ' edited' : b)).join('\n\n');
const prev200 = await md.processMarkdownBlocks(doc200);

const scenarios = [];
scenarios.push(await bench('fullRender-100k', async () => { sink += (await md.processMarkdownBlocks(doc100k)).length; }, 5));
scenarios.push(await bench('fullRender-1M', async () => { sink += (await md.processMarkdownBlocks(doc1M)).length; }, 3));
scenarios.push(await bench('tocExtract-100k', async () => { sink += md.extractTableOfContents(doc100k).length; }, 5));
scenarios.push(await bench('tocExtract-1M', async () => { sink += md.extractTableOfContents(doc1M).length; }, 3));
scenarios.push(
	await bench(
		'editOneBlock-200',
		async () => {
			if (hasIncremental) {
				const r = await md.renderDocument(prev200, edited200);
				sink += r.changedIndexes.length;
			} else {
				sink += (await md.processMarkdownBlocks(edited200)).length;
			}
		},
		7
	)
);
scenarios.push(await bench('renderOneBlock', async () => { sink += (await md.processBlock(JS_SNIPPET)).length; }, 7));

const imageRe = /!\[[^\]]*\]\(([^)]+)\)/g;
let imageCount = 0;
for (const m of doc100k.matchAll(imageRe)) imageCount++;

const summary = {
	scenarios,
	meta: {
		incrementalApi: hasIncremental,
		uniqueImagesPerRender100k: imageCount,
		doc100kBytes: doc100k.length,
		doc1MBytes: doc1M.length,
		node: process.version,
		sink
	}
};

const outDir = path.dirname(path.resolve(outPath));
await mkdir(outDir, { recursive: true });
await writeFile(outPath, JSON.stringify(summary, null, 2));

// golden block dump for equivalence check across runs
const goldenBlocks = await md.processMarkdownBlocks(doc100k);
const golden = {
	fixture: 'synth-100k',
	fixtureHash: createHash('sha256').update(doc100k).digest('hex'),
	blocks: goldenBlocks.map((b) => ({ id: b.id, source: b.source, html: b.html, startLine: b.startLine, endLine: b.endLine }))
};
await writeFile(path.join(outDir, 'golden.json'), JSON.stringify(golden));

console.log(JSON.stringify(summary, null, 2));
