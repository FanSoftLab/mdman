import { resolveMediaURL } from './platform';

/**
 * Cross-render cache for resolved media data URLs.
 *
 * Keyed by "markdownFile|mediaSrc". Without this, every preview re-render
 * re-read every image from disk through the Go IPC bridge and re-encoded it
 * to base64 (2x the image bytes, per render, per image). With it, unchanged
 * images cost one Map lookup per render.
 *
 * Tradeoff: an image edited on disk mid-session keeps its stale data URL
 * until the app restarts. Bounded at MAX_ENTRIES (oldest evicted first).
 */
const cache = new Map<string, string>();
const MAX_ENTRIES = 256;

export async function resolveMediaCached(markdownFile: string, mediaPath: string): Promise<string> {
	const key = markdownFile + '|' + mediaPath;
	const hit = cache.get(key);
	if (hit !== undefined) {
		// re-insert to mark as most recently used
		cache.delete(key);
		cache.set(key, hit);
		return hit;
	}
	let url: string;
	try {
		url = (await resolveMediaURL(markdownFile, mediaPath)) || mediaPath;
	} catch {
		return mediaPath;
	}
	if (cache.size >= MAX_ENTRIES) {
		const oldest = cache.keys().next().value;
		if (oldest !== undefined) cache.delete(oldest);
	}
	cache.set(key, url);
	return url;
}
