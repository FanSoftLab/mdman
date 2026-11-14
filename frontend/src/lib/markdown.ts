import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import remarkRehype from 'remark-rehype';
import rehypeKatex from 'rehype-katex';
import rehypeSlug from 'rehype-slug';
import rehypeHighlight from 'rehype-highlight';
import rehypeStringify from 'rehype-stringify';
import GithubSlugger from 'github-slugger';
import type { Root as MdastRoot, RootContent, Heading } from 'mdast';
import type { Root as HastRoot, Element } from 'hast';
import { visit } from 'unist-util-visit';

// Custom plugin to inject source line numbers into HTML elements
function rehypeSourceLines() {
    return (tree: HastRoot) => {
        visit(tree, 'element', (node: Element) => {
            if (node.position?.start?.line) {
                node.properties = node.properties || {};
                node.properties['data-source-line'] = node.position.start.line;
            }
        });
    };
}

// Transform mermaid node labels to use markdown strings for auto text wrapping
// This wraps text in backticks like "`text`" to enable automatic text wrapping
function transformMermaidForAutoWrap(code: string): string {

    // Process each line to handle node definitions
    const lines = code.split('\n');
    let inFrontmatter = false;

    const processedLines = lines.map((line) => {
        const trimmedLine = line.trim();

        // Track frontmatter state
        if (trimmedLine === '---') {
            inFrontmatter = !inFrontmatter;
            return line;
        }

        // Skip lines in frontmatter, comments, or already have markdown strings
        if (inFrontmatter || trimmedLine.startsWith('%%') || line.includes('"`')) {
            return line;
        }

        let processedLine = line;

        // Process square brackets [text] - most common
        // Match: A[text] but not A["`text`"] or A["text"]
        processedLine = processedLine.replace(/(\b\w+)\[([^\]"'`]+)\]/g, (match, id, text) => {
            // Skip if text is a link syntax or icon
            if (text.startsWith('fa:') || text.startsWith('fab:') || text.startsWith('fas:')) return match;
            return `${id}["\`${text}\`"]`;
        });

        // Process parentheses (text) for rounded nodes
        // Match: A(text) but not A((text)) or A("`text`")
        processedLine = processedLine.replace(/(\b\w+)\(([^()"'`]+)\)(?!\))/g, (match, id, text) => {
            // Skip if it looks like a subgraph reference
            if (text.startsWith('[') || text.includes('|')) return match;
            return `${id}("\`${text}\`")`;
        });

        // Process curly braces {text} for diamond/decision nodes
        // Match: A{text} but not A{{text}}
        processedLine = processedLine.replace(/(\b\w+)\{([^{}"'`]+)\}(?!\})/g, (match, id, text) => {
            return `${id}{"\`${text}\`"}`;
        });

        // Process stadium shape ([text])
        processedLine = processedLine.replace(/(\b\w+)\(\[([^\]"'`]+)\]\)/g, (match, id, text) => {
            return `${id}(["\`${text}\`"])`;
        });

        // Process subroutine [[text]]
        processedLine = processedLine.replace(/(\b\w+)\[\[([^\]"'`]+)\]\]/g, (match, id, text) => {
            return `${id}[["\`${text}\`"]]`;
        });

        // Process hexagon {{text}}
        processedLine = processedLine.replace(/(\b\w+)\{\{([^}"'`]+)\}\}/g, (match, id, text) => {
            return `${id}{{"\`${text}\`"}}`;
        });

        // Process double circle (((text)))
        processedLine = processedLine.replace(/(\b\w+)\(\(\(([^)"'`]+)\)\)\)/g, (match, id, text) => {
            return `${id}((("\`${text}\`")))`;
        });

        return processedLine;
    });

    const result = processedLines.join('\n');

    return result;
}

// Custom plugin to handle mermaid code blocks
function remarkMermaid() {
    return (tree: MdastRoot) => {
        visit(tree, 'code', (node: { lang?: string | null; meta?: string | null; value: string; type: string; data?: { hName?: string; hProperties?: Record<string, unknown> } }) => {
            if (node.lang === 'mermaid') {
                // Transform mermaid code to use markdown strings for auto text wrapping
                const transformedCode = transformMermaidForAutoWrap(node.value);

                // Transform to a div that will be processed by mermaid.js client-side
                node.type = 'code';
                node.data = node.data || {};
                node.data.hName = 'div';
                node.data.hProperties = {
                    className: ['mermaid'],
                    'data-mermaid': transformedCode
                };
            }
        });
    };
}

// Preprocess markdown to convert inline $$...$$ on its own line to proper block math
// and convert \[...\] bracket syntax to $$ block math
// remark-math treats inline $$...$$ as text math, but users often expect it to be display math
function preprocessDisplayMath(markdown: string): string {
    // Split into lines
    const lines = markdown.split('\n');
    const result: string[] = [];
    let inBracketMath = false;
    const bracketMathLines: string[] = [];
    let bracketMathIndent = '';

    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];
        const trimmed = line.trim();

        // Get the leading whitespace/indentation to preserve it
        const leadingWhitespace = line.match(/^(\s*)/)?.[1] || '';

        // Handle \[ ... \] bracket syntax for display math
        if (trimmed === '\\[') {
            inBracketMath = true;
            bracketMathLines.length = 0;
            bracketMathIndent = leadingWhitespace;
            continue;
        }

        if (inBracketMath) {
            if (trimmed === '\\]') {
                // End of bracket math block - convert to $$ block
                inBracketMath = false;
                result.push(bracketMathIndent + '$$');
                result.push(...bracketMathLines);
                result.push(bracketMathIndent + '$$');
                continue;
            } else {
                // Accumulate lines inside \[...\]
                bracketMathLines.push(line);
                continue;
            }
        }

        // Check if the line is primarily a $$...$$ math expression
        // Pattern: line that starts and ends with $$ with content in between
        // This handles cases like: $$ x = 1 $$ or $$\frac{1}{2}$$
        const inlineDisplayMatch = trimmed.match(/^\$\$(.+)\$\$$/);

        if (inlineDisplayMatch) {
            // Check if this is truly a single $$...$$ expression (exactly 2 $$ sequences)
            const dollarMatches = trimmed.match(/\$\$/g);
            if (dollarMatches && dollarMatches.length === 2) {
                // Convert to proper block math, preserving indentation
                const mathContent = inlineDisplayMatch[1].trim();
                result.push(leadingWhitespace + '$$');
                result.push(leadingWhitespace + mathContent);
                result.push(leadingWhitespace + '$$');
                continue;
            }
        }

        result.push(line);
    }

    return result.join('\n');
}

// Create the unified processor
const processor = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkMath)
    .use(remarkMermaid)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeKatex)
    .use(rehypeHighlight, { detect: true, ignoreMissing: true })
    .use(rehypeSlug)
    .use(rehypeSourceLines)
    .use(rehypeStringify, { allowDangerousHtml: true });

// Create a parser-only processor for splitting into blocks
const parser = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkMath);

// Block processor for individual blocks (WITHOUT rehypeSlug - slugs are applied from full document context)
const blockProcessor = unified()
    .use(remarkParse)
    .use(remarkGfm)
    .use(remarkMath)
    .use(remarkMermaid)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeKatex)
    .use(rehypeHighlight, { detect: true, ignoreMissing: true })
    .use(rehypeSourceLines)
    .use(rehypeStringify, { allowDangerousHtml: true });

/**
 * Represents a markdown block with its source and rendered HTML
 */
export interface MarkdownBlock {
    id: string;           // Unique identifier for diffing
    source: string;       // Original markdown source
    html: string;         // Rendered HTML
    startLine: number;    // Starting line number in source
    endLine: number;      // Ending line number in source
}

/**
 * Generate a unique ID for a block based on content and position
 */
function generateBlockId(source: string, startLine: number, index: number): string {
    // Use a combination of content hash, line number, and index for uniqueness
    let hash = 0;
    const str = source + '|' + startLine;
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash;
    }
    // Include index to guarantee uniqueness even for identical content
    return `b${index}-${Math.abs(hash).toString(36)}-${startLine}`;
}

/**
 * Split a parsed AST into top-level blocks (source slices + line ranges).
 * Reuses a single parse — callers that need both blocks and headings must
 * pass the same tree here instead of re-parsing the document.
 */
function blocksFromTree(tree: MdastRoot, markdown: string): { source: string; startLine: number; endLine: number }[] {
    const lines = markdown.split('\n');
    const blocks: { source: string; startLine: number; endLine: number }[] = [];
    for (const node of tree.children as RootContent[]) {
        if (node.position) {
            const startLine = node.position.start.line;
            const endLine = node.position.end.line;
            // Extract the source text for this block (1-indexed to 0-indexed)
            const source = lines.slice(startLine - 1, endLine).join('\n');
            blocks.push({ source, startLine, endLine });
        }
    }
    return blocks;
}

/**
 * Split markdown into top-level blocks (one full parse per call).
 */
export function splitIntoBlocks(markdown: string): { source: string; startLine: number; endLine: number }[] {
    if (!markdown.trim()) {
        return [];
    }
    return blocksFromTree(parser.parse(markdown) as MdastRoot, markdown);
}

/**
 * Process a single markdown block to HTML
 */
export async function processBlock(source: string): Promise<string> {
    if (!source.trim()) return '';
    // Preprocess to handle inline $$...$$ as display math
    const preprocessed = preprocessDisplayMath(source);
    const result = await blockProcessor.process(preprocessed);
    return String(result);
}

/**
 * Represents a heading with its metadata for slug tracking
 */
export interface HeadingInfo {
    level: number;
    text: string;
    slug: string;
    blockId?: string; // Block identifier for tracking which block this heading belongs to
    occurrenceInBlock: number; // Occurrence number within this block
    globalOccurrence: number; // Global occurrence number in the entire document
}

export interface TocEntry {
    level: number;
    text: string;
    id: string;
}

/**
 * Build heading slugs (block-context aware) and table of contents from an
 * already-parsed AST in a single pass. Operating on the shared tree instead
 * of re-parsing the document is the main rendering-pipeline optimization:
 * the previous implementation parsed the full document 2 extra times per render.
 * Uses GithubSlugger to handle duplicates correctly (same algorithm as rehype-slug).
 */
function buildSlugMapAndToc(tree: MdastRoot, blocks: { source: string; startLine: number; endLine: number }[]): {
    slugMap: Map<string, HeadingInfo>;
    toc: TocEntry[];
} {
    const headingInfoMap = new Map<string, HeadingInfo>();
    const toc: TocEntry[] = [];
    const slugger = new GithubSlugger();
    const headingOccurrences = new Map<string, number>(); // Track global occurrences

    // Visit all heading nodes
    visit(tree, 'heading', (node: Heading) => {
        // Extract text content from heading ONLY (not from nested code blocks, etc)
        let text = '';
        visit(node, 'text', (textNode: { value: string }) => {
            text += textNode.value;
        });

        if (text.trim()) {
            // Generate slug using GithubSlugger (same as rehype-slug)
            const slug = slugger.slug(text);

            // Track global occurrence count
            const occurrenceKey = `h${node.depth}-${text}`;
            const globalOccurrence = (headingOccurrences.get(occurrenceKey) || 0) + 1;
            headingOccurrences.set(occurrenceKey, globalOccurrence);

            // Find which block this heading belongs to
            const nodeStart = node.position?.start.line || 0;
            let blockId = 'unknown';
            for (let i = 0; i < blocks.length; i++) {
                if (nodeStart >= blocks[i].startLine && nodeStart <= blocks[i].endLine) {
                    blockId = `block-${i}`;
                    break;
                }
            }

            const headingInfo: HeadingInfo = {
                level: node.depth,
                text,
                slug,
                blockId,
                occurrenceInBlock: 1, // Will be updated if needed
                globalOccurrence
            };

            // Store with composite key for lookup
            const key = `h${node.depth}-${blockId}-${text}`;
            headingInfoMap.set(key, headingInfo);

            // Also store by slug for reverse lookup
            const slugKey = `slug-${slug}`;
            headingInfoMap.set(slugKey, headingInfo);

            toc.push({ level: node.depth, text: text.trim(), id: slug });
        }
    });

    return { slugMap: headingInfoMap, toc };
}

/**
 * Decode HTML entities to plain text
 * Handles common entities like &#x26; and &amp;
 */
function decodeHtmlEntities(text: string): string {
    return text
        .replace(/&#x26;/g, '&')
        .replace(/&#x27;/g, "'")
        .replace(/&#x22;/g, '"')
        .replace(/&#x3c;/g, '<')
        .replace(/&#x3e;/g, '>')
        .replace(/&amp;/g, '&')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
        .replace(/&quot;/g, '"')
        .replace(/&#39;/g, "'")
        .replace(/&#x2F;/g, '/');
}

/**
 * Apply slugs from full document to block HTML with proper duplicate handling
 * This ensures duplicate heading IDs are consistent across blocks
 * Uses block context to prevent collisions
 */
function applySlugsToBlockHTML(blockHTML: string, slugMap: Map<string, HeadingInfo>, blockId: string): string {
    // Find headings in block HTML and inject IDs from full document mapping
    const headingRegex = /<h([1-6])(\s+[^>]*)?>(.+?)<\/h\1>/g;
    let result = blockHTML;

    result = result.replace(headingRegex, (match, level, attrs, content) => {

        // Extract plain text from HTML content (strip tags and decode entities)
        let plainText = content.replace(/<[^>]+>/g, '');
        plainText = decodeHtmlEntities(plainText);

        // Look up this heading using block context
        const mapKey = `h${level}-${blockId}-${plainText}`;
        const headingInfo = slugMap.get(mapKey);

        if (headingInfo) {
            const slugId = headingInfo.slug;
            // Check if heading already has an id attribute
            if ((attrs || '').includes('id=')) {
                // Replace existing id
                return match.replace(/id="[^"]*"/, `id="${slugId}"`);
            } else {
                // Add id attribute
                return `<h${level}${attrs || ''} id="${slugId}">${content}</h${level}>`;
            }
        }

        return match;
    });

    return result;
}

/**
 * LRU-ish cache of pre-slug block HTML keyed by block source.
 * Unchanged blocks skip the (expensive) parse/render entirely on re-renders.
 * Slugs are applied on top of the cached HTML every render, so heading
 * edits above a block still shift its ids correctly.
 */
const blockHtmlCache = new Map<string, { html: string; size: number }>();
const BLOCK_CACHE_MAX_ENTRIES = 1000;
const BLOCK_CACHE_MAX_BYTES = 64 * 1024 * 1024;
let blockCacheBytes = 0;

function getCachedBlockHtml(source: string): string | null {
    const entry = blockHtmlCache.get(source);
    if (entry === undefined) return null;
    // re-insert to mark as most recently used
    blockHtmlCache.delete(source);
    blockHtmlCache.set(source, entry);
    return entry.html;
}

function putCachedBlockHtml(source: string, html: string): void {
    const prev = blockHtmlCache.get(source);
    if (prev !== undefined) {
        blockCacheBytes -= prev.size;
        blockHtmlCache.delete(source);
    }
    const size = source.length + html.length;
    blockHtmlCache.set(source, { html, size });
    blockCacheBytes += size;
    while (
        blockHtmlCache.size > 1 &&
        (blockHtmlCache.size > BLOCK_CACHE_MAX_ENTRIES || blockCacheBytes > BLOCK_CACHE_MAX_BYTES)
    ) {
        const oldest = blockHtmlCache.keys().next().value;
        if (oldest === undefined) break;
        const evicted = blockHtmlCache.get(oldest);
        blockCacheBytes -= evicted?.size ?? 0;
        blockHtmlCache.delete(oldest);
    }
}

/**
 * Result of an incremental document render.
 */
export interface RenderResult {
    blocks: MarkdownBlock[];
    /** Indexes of blocks whose source changed since prevBlocks (i.e. were re-rendered). */
    changedIndexes: number[];
    /** Table of contents extracted from the same AST pass (no extra parse). */
    toc: TocEntry[];
}

/**
 * Render markdown into blocks incrementally.
 *
 * OPTIMIZED pipeline: exactly ONE full-document parse per call (the previous
 * implementation performed ~N+3 parses: splitIntoBlocks, heading-slugs pass
 * plus a second splitIntoBlocks inside it, and one parse per block).
 * Unchanged blocks are served from the pre-slug HTML cache, so editing one
 * block of a large document costs O(changed block) instead of O(document).
 * Output is identical to a full render (slug handling included).
 */
export async function renderDocument(
    prevBlocks: MarkdownBlock[] | null,
    markdown: string
): Promise<RenderResult> {
    if (!markdown.trim()) {
        return { blocks: [], changedIndexes: [], toc: [] };
    }

    const tree = parser.parse(markdown) as MdastRoot;
    const rawBlocks = blocksFromTree(tree, markdown);
    const { slugMap, toc } = buildSlugMapAndToc(tree, rawBlocks);

    const blocks: MarkdownBlock[] = [];
    const changed = new Set<number>();

    for (let i = 0; i < rawBlocks.length; i++) {
        const raw = rawBlocks[i];
        const id = generateBlockId(raw.source, raw.startLine, i);
        const prev = prevBlocks && prevBlocks[i] && prevBlocks[i].id === id ? prevBlocks[i] : null;
        if (!prev) changed.add(i);

        let preSlugHtml = getCachedBlockHtml(raw.source);
        if (preSlugHtml === null) {
            preSlugHtml = await processBlock(raw.source);
            putCachedBlockHtml(raw.source, preSlugHtml);
            changed.add(i);
        }

        const html = applySlugsToBlockHTML(preSlugHtml, slugMap, `block-${i}`);

        blocks.push({
            id,
            source: raw.source,
            html,
            startLine: raw.startLine,
            endLine: raw.endLine
        });
    }

    return { blocks, changedIndexes: [...changed], toc };
}

/**
 * Process markdown into blocks with rendered HTML (full render, cold start).
 * Kept as the public full-render API; incremental updates should use
 * renderDocument(prevBlocks, markdown).
 */
export async function processMarkdownBlocks(markdown: string): Promise<MarkdownBlock[]> {
    return (await renderDocument(null, markdown)).blocks;
}

/**
 * Extract table of contents from markdown using AST (not regex)
 * This properly distinguishes between:
 * - Actual headings (# Comment)
 * - # symbols in code blocks (```bash # comment\n```)
 *
 * Uses GithubSlugger to generate consistent IDs, properly handling duplicates.
 * NOTE: callers that already render blocks should use renderDocument's toc
 * instead — it reuses the same AST pass (no extra parse).
 */
export function extractTableOfContents(markdown: string): TocEntry[] {
    if (!markdown.trim()) {
        return [];
    }

    try {
        const tree = parser.parse(markdown) as MdastRoot;
        const blocks = blocksFromTree(tree, markdown);
        return buildSlugMapAndToc(tree, blocks).toc;
    } catch (e) {
        console.error('Error extracting table of contents:', e);
        return [];
    }
}

/**
 * Process markdown string to HTML (legacy full render)
 */
export async function processMarkdown(markdown: string): Promise<string> {
    // Preprocess to handle inline $$...$$ as display math
    const preprocessed = preprocessDisplayMath(markdown);
    const result = await processor.process(preprocessed);
    return String(result);
}

/**
 * Synchronous version for when async isn't needed
 */
export function processMarkdownSync(markdown: string): string {
    // Preprocess to handle inline $$...$$ as display math
    const preprocessed = preprocessDisplayMath(markdown);
    const result = processor.processSync(preprocessed);
    return String(result);
}

// Export the processor for advanced use cases
export { processor };
