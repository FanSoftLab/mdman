<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { processMarkdownBlocks, diffBlocks, type MarkdownBlock } from '../markdown';
	import mermaid from 'mermaid';
	import * as platform from '../platform';
	import { appState } from '../appState.svelte';
	import { t } from '../i18n/index.svelte';
	import markdownLightCss from 'github-markdown-css/github-markdown-light.css?inline';
	import hljsLightCss from 'highlight.js/styles/github.css?inline';

	interface Props {
		content?: string;
		/** Absolute path of the markdown file — used to resolve relative image paths */
		baseFilePath?: string;
		onscroll?: (scrollInfo: { scrollTop: number; scrollHeight: number; clientHeight: number }) => void;
		ondimensionschange?: () => void; // Called when section dimensions are measured/updated
		class?: string;
	}

	let { content = '', baseFilePath = '', onscroll, ondimensionschange, class: className = '' }: Props = $props();

	let previewContainer: HTMLDivElement;
	let blocks = $state<MarkdownBlock[]>([]);
	let isSyncingScroll = false;
	let isRendering = false; // Flag to block scroll sync during re-render
	let processingTimeout: ReturnType<typeof setTimeout> | null = null;

	// Section dimensions for smart scroll sync
	interface SectionDimension {
		startOffset: number;  // Pixel offset from top of container
		endOffset: number;    // End pixel offset
		height: number;       // Height in pixels
	}

	interface SectionInfo {
		block: MarkdownBlock;
		previewDimension: SectionDimension;
	}

	let sectionInfoList = $state<SectionInfo[]>([]);

	async function resolveImagesInHtml(html: string, base: string): Promise<string> {
		if (!html || !base) return html;
		// Match src="..." in img tags (markdown and raw HTML)
		const re = /(<img\b[^>]*?\bsrc\s*=\s*)(["'])([^"']+)(\2)/gi;
		const matches = [...html.matchAll(re)];
		if (!matches.length) return html;
		let out = html;
		// resolve unique srcs
		const cache = new Map<string, string>();
		for (const m of matches) {
			const src = m[3];
			if (cache.has(src)) continue;
			const lower = src.toLowerCase();
			if (lower.startsWith('http://') || lower.startsWith('https://') ||
				lower.startsWith('data:') || lower.startsWith('blob:') || lower.startsWith('wails://')) {
				cache.set(src, src);
				continue;
			}
			try {
				const url = await platform.resolveMediaURL(base, src);
				cache.set(src, url || src);
			} catch {
				cache.set(src, src);
			}
		}
		out = out.replace(re, (_full, pre, quote, src, q2) => {
			return pre + quote + (cache.get(src) || src) + q2;
		});
		// also markdown-style already converted; handle src='...'
		return out;
	}

	// Light stylesheets for the preview (github-markdown + highlight.js), toggled by preview theme
	let markdownLightStyle: HTMLStyleElement | null = null;
	let hljsLightStyle: HTMLStyleElement | null = null;

	function ensureLightStylesheets() {
		if (markdownLightStyle && hljsLightStyle) return;
		if (!markdownLightStyle) {
			markdownLightStyle = document.createElement('style');
			markdownLightStyle.media = 'not all';
			markdownLightStyle.textContent = markdownLightCss;
			document.head.appendChild(markdownLightStyle);
		}
		if (!hljsLightStyle) {
			hljsLightStyle = document.createElement('style');
			hljsLightStyle.media = 'not all';
			hljsLightStyle.textContent = hljsLightCss;
			document.head.appendChild(hljsLightStyle);
		}
	}

	function applyMermaidTheme(theme: 'dark' | 'light') {
		mermaid.initialize({
			startOnLoad: false,
			theme: theme === 'light' ? 'default' : 'dark',
			securityLevel: 'loose',
			fontFamily: 'inherit',
			htmlLabels: false,
			markdownAutoWrap: true,
			wrap: true,
		});
	}

	// Initialize mermaid and set up resize observer for dimension changes
	onMount(() => {
		applyMermaidTheme(appState.resolvedPreviewTheme);

		// Use ResizeObserver to detect any size changes (images, dynamic content, etc.)
		let resizeTimeout: ReturnType<typeof setTimeout> | null = null;
		const resizeObserver = new ResizeObserver(() => {
			// Debounce re-measurement to batch multiple resize events
			if (resizeTimeout) clearTimeout(resizeTimeout);
			resizeTimeout = setTimeout(() => {
				measureSectionDimensions();
			}, 50);
		});

		// Observe the container for any size changes within it
		if (previewContainer) {
			resizeObserver.observe(previewContainer);
		}

		return () => {
			resizeObserver.disconnect();
			if (resizeTimeout) clearTimeout(resizeTimeout);
		};
	});

	// Toggle light preview stylesheets when the preview theme changes
	$effect(() => {
		const isLight = appState.resolvedPreviewTheme === 'light';
		ensureLightStylesheets();
		if (markdownLightStyle) markdownLightStyle.media = isLight ? 'all' : 'not all';
		if (hljsLightStyle) hljsLightStyle.media = isLight ? 'all' : 'not all';
	});

	// Re-initialize mermaid and re-render diagrams when the preview theme changes
	$effect(() => {
		const theme = appState.resolvedPreviewTheme;
		if (!previewContainer) return;
		applyMermaidTheme(theme);
		const rendered = previewContainer.querySelectorAll('.mermaid[data-rendered]');
		if (rendered.length > 0) {
			rendered.forEach((div) => div.removeAttribute('data-rendered'));
			void renderMermaidDiagrams();
		}
	});

	// Debounced block processing when content changes
	$effect(() => {
		const currentContent = content;
		const _baseTrack = baseFilePath; // re-resolve images when file path changes
		
		// Clear any pending processing
		if (processingTimeout) {
			clearTimeout(processingTimeout);
		}

		// Debounce processing to avoid excessive updates
		processingTimeout = setTimeout(async () => {
			try {
				const newBlocks = await processMarkdownBlocks(currentContent || '');
				const base = baseFilePath;
				let resolvedBlocks = newBlocks;
				if (base) {
					// Resolve images sequentially with a shared cache (faster than N parallel Go calls)
					const cache = new Map<string, string>();
					resolvedBlocks = [];
					for (const b of newBlocks) {
						resolvedBlocks.push({
							...b,
							html: await resolveImagesInHtml(b.html, base),
						});
					}
				}
				blocks = resolvedBlocks;
			} catch (err) {
				console.error('preview render failed', err);
				// Fallback: show raw content so the pane is never blank
				blocks = [{
					id: 'fallback',
					source: currentContent || '',
					html: `<pre>${(currentContent || '').replace(/[&<>]/g, (c) => ({'&':'&amp;','<':'&lt;','>':'&gt;'}[c] as string))}</pre>`,
					startLine: 1,
					endLine: 1,
				} as any];
			}
		}, 250);
	});

	// Render mermaid diagrams for newly added/updated blocks
	$effect(() => {
		if (blocks.length > 0 && previewContainer) {
			// Block scroll sync during rendering
			isRendering = true;
			
			// Use tick to ensure DOM is updated, then render mermaid and add copy buttons
			tick().then(async () => {
				await renderMermaidDiagrams();
				addCopyButtons();
				
				// Wait for all images and dynamic content to load before measuring
				await waitForContentToLoad();
				
				// Measure section dimensions after all rendering is complete
				measureSectionDimensions();
				
				// Re-enable scroll sync after a short delay to let things settle
				setTimeout(() => {
					isRendering = false;
				}, 100);
			});
		}
	});

	// Wait for all images, iframes, and dynamic content to fully load
	async function waitForContentToLoad(): Promise<void> {
		if (!previewContainer) return;
		
		// Collect all loading promises
		const loadingPromises: Promise<void>[] = [];
		
		// Wait for all images to load
		const images = previewContainer.querySelectorAll('img');
		images.forEach((img) => {
			if (!img.complete) {
				loadingPromises.push(new Promise((resolve) => {
					const onLoad = () => { img.removeEventListener('error', onLoad); resolve(); };
					const onError = () => { img.removeEventListener('load', onLoad); resolve(); };
					img.addEventListener('load', onLoad, { once: true });
					img.addEventListener('error', onError, { once: true });
				}));
			}
		});
		
		// Wait for SVGs in mermaid diagrams to be ready (they might resize)
		const mermaidSvgs = previewContainer.querySelectorAll('.mermaid svg');
		mermaidSvgs.forEach((svg) => {
			// SVGs are synchronous but might have fonts loading
			loadingPromises.push(new Promise((resolve) => {
				// Use requestAnimationFrame to ensure layout is complete
				requestAnimationFrame(() => {
					requestAnimationFrame(() => resolve());
				});
			}));
		});
		
		// Wait for KaTeX elements to render (they use fonts that might load async)
		const katexElements = previewContainer.querySelectorAll('.katex');
		if (katexElements.length > 0) {
			// Wait for fonts to be ready
			loadingPromises.push(
				document.fonts.ready.then(() => {
					// Extra frame to ensure layout recalc
					return new Promise((resolve) => {
						requestAnimationFrame(() => resolve());
					});
				})
			);
		}
		
		// Wait for all promises with a timeout
		if (loadingPromises.length > 0) {
			await Promise.race([
				Promise.all(loadingPromises),
				new Promise(resolve => setTimeout(resolve, 2000)) // 2s timeout
			]);
		}
		
		// Final layout settle
		await new Promise(resolve => requestAnimationFrame(resolve));
	}

	// Measure the dimensions of each section in the preview
	function measureSectionDimensions() {
		if (!previewContainer) return;

		const blockElements = previewContainer.querySelectorAll('.markdown-block');
		const newSectionInfoList: SectionInfo[] = [];
		
		// Get container's scroll position and bounding rect for accurate measurements
		const containerRect = previewContainer.getBoundingClientRect();
		const scrollTop = previewContainer.scrollTop;

		blocks.forEach((block, index) => {
			const element = blockElements[index] as HTMLElement;
			if (element) {
				// For display:contents elements, we need to measure ALL their children
				// to get accurate bounding box including wrapped text
				const children = element.children;
				
				let startOffset = Infinity;
				let endOffset = 0;
				
				if (children.length > 0) {
					// Iterate through all children to find true bounds
					// This handles wrapped text and varying element heights
					for (let i = 0; i < children.length; i++) {
						const child = children[i] as HTMLElement;
						const rect = child.getBoundingClientRect();
						
						// Convert to scroll-relative positions
						const childTop = rect.top - containerRect.top + scrollTop;
						const childBottom = rect.bottom - containerRect.top + scrollTop;
						
						startOffset = Math.min(startOffset, childTop);
						endOffset = Math.max(endOffset, childBottom);
					}
					
					// Handle case where startOffset is still Infinity
					if (startOffset === Infinity) {
						startOffset = 0;
					}
				} else if (element.offsetHeight > 0) {
					// Fallback for non-contents elements
					const rect = element.getBoundingClientRect();
					startOffset = rect.top - containerRect.top + scrollTop;
					endOffset = rect.bottom - containerRect.top + scrollTop;
				} else {
					// Empty block - use previous section's end if available
					if (newSectionInfoList.length > 0) {
						const prevSection = newSectionInfoList[newSectionInfoList.length - 1];
						startOffset = prevSection.previewDimension.endOffset;
						endOffset = startOffset;
					} else {
						startOffset = 0;
						endOffset = 0;
					}
				}

				newSectionInfoList.push({
					block,
					previewDimension: {
						startOffset,
						endOffset,
						height: Math.max(0, endOffset - startOffset)
					}
				});
			}
		});

		sectionInfoList = newSectionInfoList;
		
		// Notify parent that dimensions have been updated
		ondimensionschange?.();
	}

	async function renderMermaidDiagrams() {
		if (!previewContainer) return;
		
		// Ensure fonts are loaded before rendering to prevent text cutting
		await document.fonts.ready;
		
		const mermaidDivs = previewContainer.querySelectorAll('.mermaid:not([data-rendered])');
		for (const div of mermaidDivs) {
			const code = div.getAttribute('data-mermaid');
			if (code) {
				try {
					const id = `mermaid-${Math.random().toString(36).substr(2, 9)}`;
					// Pass the container element so mermaid can inherit its font-family
					const { svg } = await mermaid.render(id, code, div as HTMLElement);
					div.innerHTML = svg;
					div.setAttribute('data-rendered', 'true');
				} catch (err) {
					console.error('Mermaid render error:', err);
					div.innerHTML = `<pre class="mermaid-error">${t('mermaidRenderError')}: ${err}</pre>`;
					div.setAttribute('data-rendered', 'true');
				}
			}
		}
	}

	// Add copy buttons to code blocks
	function addCopyButtons() {
		if (!previewContainer) return;
		
		const codeBlocks = previewContainer.querySelectorAll('pre:not([data-copy-added])');
		for (const pre of codeBlocks) {
			// Skip mermaid blocks
			if (pre.closest('.mermaid')) continue;
			
			pre.setAttribute('data-copy-added', 'true');
			
			// Create wrapper for positioning
			const wrapper = document.createElement('div');
			wrapper.className = 'code-block-wrapper';
			pre.parentNode?.insertBefore(wrapper, pre);
			wrapper.appendChild(pre);
			
			// Create copy button
			const copyBtn = document.createElement('button');
			copyBtn.className = 'copy-code-btn';
			copyBtn.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`;
			copyBtn.title = t('copyCode');
			
			copyBtn.addEventListener('click', async () => {
				const code = pre.querySelector('code');
				const text = code?.textContent || pre.textContent || '';
				
				try {
					await navigator.clipboard.writeText(text);
					copyBtn.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
					copyBtn.classList.add('copied');
					
					setTimeout(() => {
						copyBtn.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`;
						copyBtn.classList.remove('copied');
					}, 2000);
				} catch (err) {
					console.error('Failed to copy:', err);
				}
			});
			
			wrapper.appendChild(copyBtn);
		}
	}

	// Handle scroll events for scroll sync (section-based)
	function handleScroll() {
		// Don't emit scroll events while rendering or syncing
		if (!previewContainer || !onscroll || isSyncingScroll || isRendering) return;

		onscroll({
			scrollTop: previewContainer.scrollTop,
			scrollHeight: previewContainer.scrollHeight,
			clientHeight: previewContainer.clientHeight
		});
	}

	// Get scroll position as a continuous value across all sections
	// Returns fractional section index for smooth interpolation
	export function getScrollPosition(): { sectionIdx: number; posInSection: number } | null {
		if (!previewContainer || sectionInfoList.length === 0) return null;

		const scrollTop = previewContainer.scrollTop;
		const maxScroll = previewContainer.scrollHeight - previewContainer.clientHeight;
		
		// Handle edge cases
		if (maxScroll <= 0) return { sectionIdx: 0, posInSection: 0 };
		if (scrollTop <= 0) return { sectionIdx: 0, posInSection: 0 };
		if (scrollTop >= maxScroll) {
			return { sectionIdx: sectionInfoList.length - 1, posInSection: 1 };
		}
		
		// Find which section we're in
		for (let i = 0; i < sectionInfoList.length; i++) {
			const section = sectionInfoList[i];
			const dim = section.previewDimension;
			
			// Check if scroll is within this section's range
			if (scrollTop >= dim.startOffset && scrollTop < dim.endOffset) {
				const posInSection = dim.height > 0 
					? (scrollTop - dim.startOffset) / dim.height 
					: 0;
				return {
					sectionIdx: i,
					posInSection: Math.max(0, Math.min(1, posInSection))
				};
			}
			
			// Check if we're in the gap between this section and the next
			if (i < sectionInfoList.length - 1) {
				const nextSection = sectionInfoList[i + 1];
				const gapStart = dim.endOffset;
				const gapEnd = nextSection.previewDimension.startOffset;
				
				if (scrollTop >= gapStart && scrollTop < gapEnd) {
					// Interpolate through the gap - treat it as end of current section
					return {
						sectionIdx: i,
						posInSection: 1 // At end of section during gap
					};
				}
			}
		}
		
		// Fallback: use last section
		return { sectionIdx: sectionInfoList.length - 1, posInSection: 1 };
	}

	// Animation state for smooth scrolling
	let animationFrameId: number | null = null;
	let animationStartTime: number = 0;
	let animationStartScroll: number = 0;
	let animationTargetScroll: number = 0;
	let lastTargetScroll: number = 0;
	const ANIMATION_DURATION = 50; // ms - fast for responsive feel
	
	// Dynamic scroll threshold based on viewport size (0.5% of client height, min 2px, max 10px)
	function getScrollThreshold(): number {
		if (!previewContainer) return 5;
		const clientHeight = previewContainer.clientHeight;
		return Math.max(2, Math.min(10, clientHeight * 0.005));
	}

	// Smooth easing function (ease-out-quad for natural deceleration)
	function easeOutQuad(t: number): number {
		return t * (2 - t);
	}

	// Animate scroll to target position
	function animateScrollTo(targetScroll: number) {
		if (!previewContainer) return;
		const threshold = getScrollThreshold();

		// Skip if target is very close to current position (prevents micro-jitter)
		const currentScroll = previewContainer.scrollTop;
		if (Math.abs(targetScroll - currentScroll) < threshold) {
			return;
		}

		// If we're already animating toward a similar target, don't restart
		if (animationFrameId !== null && Math.abs(targetScroll - lastTargetScroll) < threshold) {
			return;
		}

		// Cancel any ongoing animation
		if (animationFrameId !== null) {
			cancelAnimationFrame(animationFrameId);
		}

		animationStartTime = performance.now();
		animationStartScroll = currentScroll;
		animationTargetScroll = targetScroll;
		lastTargetScroll = targetScroll;
		isSyncingScroll = true;

		function animate(currentTime: number) {
			const elapsed = currentTime - animationStartTime;
			const progress = Math.min(elapsed / ANIMATION_DURATION, 1);
			const easedProgress = easeOutQuad(progress);

			const newScroll = animationStartScroll + 
				(animationTargetScroll - animationStartScroll) * easedProgress;
			
			if (previewContainer) {
				previewContainer.scrollTop = Math.round(newScroll);
			}

			if (progress < 1) {
				animationFrameId = requestAnimationFrame(animate);
			} else {
				animationFrameId = null;
				// Keep isSyncingScroll true longer to prevent feedback loops
				setTimeout(() => {
					isSyncingScroll = false;
				}, 100);
			}
		}

		animationFrameId = requestAnimationFrame(animate);
	}

	// Scroll to a specific section and position within it (for section-based sync)
	export function scrollToSection(sectionIdx: number, posInSection: number, animate: boolean = true) {
		// Don't scroll while rendering - dimensions are unstable
		if (!previewContainer || sectionInfoList.length === 0 || isRendering) return;

		const section = sectionInfoList[Math.min(sectionIdx, sectionInfoList.length - 1)];
		if (!section) return;

		const dim = section.previewDimension;
		const targetScroll = dim.startOffset + (dim.height * posInSection);
		const clampedTarget = Math.max(0, Math.min(targetScroll, 
			previewContainer.scrollHeight - previewContainer.clientHeight));
		
		// Direct scroll - feedback prevention handled by scrollSource in parent
		isSyncingScroll = true;
		previewContainer.scrollTop = clampedTarget;
		isSyncingScroll = false;
	}

	// Method to scroll to a specific pixel offset (for anchor-based sync)
	export function scrollToOffset(offset: number) {
		// Don't scroll while rendering - dimensions are unstable
		if (!previewContainer || isRendering) return;
		const clampedOffset = Math.max(0, Math.min(offset, 
			previewContainer.scrollHeight - previewContainer.clientHeight));
		isSyncingScroll = true;
		previewContainer.scrollTop = clampedOffset;
		isSyncingScroll = false;
	}

	// Get section info for scroll sync (blocks with their line ranges)
	export function getSectionInfo(): Array<{ startLine: number; endLine: number; previewDimension: SectionDimension }> {
		return sectionInfoList.map(s => ({
			startLine: s.block.startLine,
			endLine: s.block.endLine,
			previewDimension: s.previewDimension
		}));
	}

	// Get scroll dimensions for viewport-aware sync
	export function getScrollDimensions(): { scrollTop: number; scrollHeight: number; clientHeight: number } | null {
		if (!previewContainer) return null;
		return {
			scrollTop: previewContainer.scrollTop,
			scrollHeight: previewContainer.scrollHeight,
			clientHeight: previewContainer.clientHeight
		};
	}

	// Method to scroll to a specific source line
	export function scrollToLine(lineNumber: number) {
		if (!previewContainer) return;

		const elements = previewContainer.querySelectorAll('[data-source-line]');
		let targetElement: Element | null = null;
		let closestLine = 0;

		for (const el of elements) {
			const elLine = parseInt(el.getAttribute('data-source-line') || '0', 10);
			if (elLine <= lineNumber && elLine > closestLine) {
				closestLine = elLine;
				targetElement = el;
			}
		}

		if (targetElement) {
			targetElement.scrollIntoView({ behavior: 'smooth', block: 'start' });
		}
	}

	// Method to scroll by percentage (fallback for proportional scroll sync)
	export function scrollToPercent(percent: number) {
		// Don't scroll while rendering - dimensions are unstable
		if (!previewContainer || isRendering) return;
		const maxScroll = previewContainer.scrollHeight - previewContainer.clientHeight;
		// Direct scroll - feedback prevention handled by scrollSource in parent
		isSyncingScroll = true;
		previewContainer.scrollTop = maxScroll * percent;
		isSyncingScroll = false;
	}

	// Get current scroll percentage
	export function getScrollPercent(): number {
		if (!previewContainer) return 0;
		const maxScroll = previewContainer.scrollHeight - previewContainer.clientHeight;
		return maxScroll > 0 ? previewContainer.scrollTop / maxScroll : 0;
	}

	// Get rendered HTML content for printing (with light theme mermaid)
	export async function getHTMLForPrint(): Promise<string> {
		if (!previewContainer) return blocks.map(b => b.html).join('\n');
		
		// Ensure fonts are loaded before rendering for print
		await document.fonts.ready;
		
		// Clone the content
		const clone = previewContainer.cloneNode(true) as HTMLDivElement;
		
		// Re-render mermaid diagrams with light theme for printing
		const mermaidDivs = clone.querySelectorAll('.mermaid');
		for (const div of mermaidDivs) {
			const code = div.getAttribute('data-mermaid');
			if (code) {
				try {
					// Temporarily switch to light theme with wrapping enabled
					mermaid.initialize({
						startOnLoad: false,
						theme: 'default',
						securityLevel: 'loose',
						fontFamily: 'inherit',
						htmlLabels: false,
						markdownAutoWrap: true,
						wrap: true
					});
					const id = `mermaid-print-${Math.random().toString(36).substr(2, 9)}`;
					// Pass the original div for font inheritance context
					const { svg } = await mermaid.render(id, code, div as HTMLElement);
					div.innerHTML = svg;
					
					// Normalize the SVG for better print output
					const svgElement = div.querySelector('svg');
					if (svgElement) {
						// Remove fixed height/width, let CSS control sizing
						svgElement.removeAttribute('height');
						svgElement.removeAttribute('width');
						
						// Remove inline max-width style that mermaid sets (often huge values)
						svgElement.style.removeProperty('max-width');
						
						// Set proper responsive sizing
						svgElement.style.maxWidth = '100%';
						svgElement.style.height = 'auto';
						svgElement.style.display = 'block';
						svgElement.style.margin = '0 auto';
						
						// Fix viewBox if it has excessive padding
						// Mermaid sometimes sets viewBox with negative coords or large padding
						const viewBox = svgElement.getAttribute('viewBox');
						if (viewBox) {
							const parts = viewBox.split(' ').map(Number);
							if (parts.length === 4) {
								// If viewBox starts with negative coordinates, it might cause offset
								// Keep the viewBox but ensure the SVG displays properly
								const [minX, minY, width, height] = parts;
								
								// If the viewBox has significant negative offset, adjust
								if (minY < -50) {
									// Adjust viewBox to start from a reasonable position
									const newMinY = Math.max(minY, -20);
									const newHeight = height + (minY - newMinY);
									svgElement.setAttribute('viewBox', `${minX} ${newMinY} ${width} ${newHeight}`);
								}
							}
						}
					}
					
					// Restore the current preview theme after printing
					mermaid.initialize({
						startOnLoad: false,
						theme: appState.resolvedPreviewTheme === 'light' ? 'default' : 'dark',
						securityLevel: 'loose',
						fontFamily: 'inherit',
						htmlLabels: false,
						markdownAutoWrap: true,
						wrap: true
					});
				} catch (err) {
					console.error('Mermaid print render error:', err);
				}
			}
		}
		
		// Extract content from markdown-block wrappers to get clean HTML
		const blockDivs = clone.querySelectorAll('.markdown-block');
		if (blockDivs.length > 0) {
			// Collect innerHTML from each block wrapper
			const htmlParts: string[] = [];
			blockDivs.forEach(blockDiv => {
				htmlParts.push(blockDiv.innerHTML);
			});
			return htmlParts.join('\n');
		}
		
		return clone.innerHTML;
	}

	// Get rendered HTML content (includes rendered mermaid SVGs - dark theme)
	export function getHTML(): string {
		// Return the actual DOM content which includes rendered mermaid diagrams
		return previewContainer?.innerHTML || blocks.map(b => b.html).join('\n');
	}
</script>

<div
	bind:this={previewContainer}
	class="preview-container markdown-body {className}"
	onscroll={handleScroll}
>
	{#each blocks as block (block.id)}
		<div class="markdown-block" data-block-id={block.id} data-start-line={block.startLine}>
			{@html block.html}
		</div>
	{/each}
</div>

<style>
	.preview-container {
		height: 100%;
		width: 100%;
		overflow-y: auto;
		padding: 20px;
		padding-top: 0;
	}

	/* Block wrapper - seamless rendering */
	.markdown-block {
		display: contents;
	}

	/* Scrollbar - follows preview theme */
	.preview-container::-webkit-scrollbar {
		width: 12px;
		height: 12px;
	}

	.preview-container::-webkit-scrollbar-track {
		background: var(--bg, #0d1117);
	}

	.preview-container::-webkit-scrollbar-thumb {
		background: var(--border, #30363d);
		border-radius: 6px;
		border: 3px solid var(--bg, #0d1117);
	}

	.preview-container::-webkit-scrollbar-thumb:hover {
		background: var(--scrollbar-thumb-hover, #484f58);
	}

	:global(html.preview-theme-light) .preview-container::-webkit-scrollbar-track {
		background: #f6f8fa;
	}

	:global(html.preview-theme-light) .preview-container::-webkit-scrollbar-thumb {
		background: #d0d7de;
		border: 3px solid #f6f8fa;
	}

	:global(html.preview-theme-light) .preview-container::-webkit-scrollbar-thumb:hover {
		background: #c8ced4;
	}

	/* Mermaid diagram layout (colors come from mermaid theme) */
	.preview-container :global(.mermaid) {
		padding: 0;
		border-radius: 6px;
		margin: 0 0 16px 0;
		text-align: center;
	}

	.preview-container :global(.mermaid svg) {
		overflow: visible !important;
	}

	.preview-container :global(.mermaid-error) {
		color: #f85149;
		background: var(--bg-elevated, #161b22);
		padding: 16px;
		border-radius: 6px;
		border: 1px solid #f85149;
	}

	/* KaTeX layout */
	/* Display math ($$...$$) - centered on its own line */
	.preview-container :global(.katex-display) {
		display: block;
		text-align: center;
		margin: 1em 0;
		overflow-x: auto;
		overflow-y: hidden;
		padding: 8px 0;
	}

	.preview-container :global(.katex-display > .katex) {
		display: inline-block;
		text-align: center;
	}

	/* Inline math ($...$) - flows with text */
	.preview-container :global(.katex:not(.katex-display .katex)) {
		font-size: 1.1em;
	}

	/* Task list styling */
	.preview-container :global(input[type='checkbox']) {
		margin-right: 8px;
	}

	/* Code block copy button */
	.preview-container :global(.code-block-wrapper) {
		position: relative;
		margin: 16px 0;
	}

	.preview-container :global(.code-block-wrapper pre) {
		margin: 0;
	}

	.preview-container :global(.copy-code-btn) {
		position: absolute;
		top: 8px;
		right: 8px;
		padding: 6px 8px;
		background: var(--bg-elevated, #30363d);
		border: 1px solid var(--border, #484f58);
		border-radius: 6px;
		color: var(--fg-muted, #8b949e);
		cursor: pointer;
		opacity: 0;
		transition: opacity 0.2s, background 0.2s, color 0.2s;
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 10;
	}

	.preview-container :global(.code-block-wrapper:hover .copy-code-btn) {
		opacity: 1;
	}

	.preview-container :global(.copy-code-btn:hover) {
		background: var(--bg-hover, #484f58);
		color: var(--fg, #c9d1d9);
	}

	.preview-container :global(.copy-code-btn.copied) {
		background: #238636;
		color: #ffffff;
		border-color: #3fb950;
	}

	.preview-container :global(.copy-code-btn svg) {
		width: 16px;
		height: 16px;
	}
</style>
