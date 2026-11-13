/** Wails runtime bridge */

declare global {
  interface Window {
    go?: {
      main?: {
        App?: Record<string, (...args: unknown[]) => Promise<unknown>>;
      };
    };
  }
}

function api() {
  return window.go?.main?.App;
}

export type RecentFile = { path: string; name: string; openedAt: number };
export type FileNode = {
  name: string;
  path: string;
  isDir: boolean;
  size?: number;
  modified?: number;
  hasChildren?: boolean;
  children?: FileNode[];
};

export async function getSettings(): Promise<Record<string, unknown>> {
  const a = api();
  if (!a?.GetSettings) return {};
  return (await a.GetSettings()) as Record<string, unknown>;
}

export async function saveSettings(s: Record<string, unknown>): Promise<void> {
  const a = api();
  if (a?.SaveSettings) await a.SaveSettings(s);
}

export async function listRecent(): Promise<RecentFile[]> {
  const a = api();
  if (!a?.ListRecent) return [];
  return ((await a.ListRecent()) as RecentFile[]) ?? [];
}

export async function clearRecent(): Promise<void> {
  const a = api();
  if (a?.ClearRecent) await a.ClearRecent();
}

export async function getWorkspace(): Promise<string> {
  const a = api();
  if (!a?.GetWorkspace) return '';
  return ((await a.GetWorkspace()) as string) ?? '';
}

export async function selectWorkspace(): Promise<string> {
  const a = api();
  if (!a?.SelectWorkspace) return '';
  return ((await a.SelectWorkspace()) as string) ?? '';
}

export async function listFiles(rel = ''): Promise<FileNode[]> {
  const a = api();
  if (!a?.ListFiles) return [];
  return ((await a.ListFiles(rel)) as FileNode[]) ?? [];
}

export async function readFile(path: string): Promise<string> {
  const a = api();
  if (!a?.ReadFile) return '';
  return ((await a.ReadFile(path)) as string) ?? '';
}

export async function writeFile(path: string, content: string): Promise<void> {
  const a = api();
  if (a?.WriteFile) await a.WriteFile(path, content);
}

export async function saveFileAs(content: string, name?: string): Promise<string> {
  const a = api();
  if (!a?.SaveFileAs) return '';
  return ((await a.SaveFileAs(content, name ?? 'untitled.md')) as string) ?? '';
}

export async function openFileDialog(): Promise<{ path: string; name: string; content: string } | null> {
  const a = api();
  if (!a?.OpenFileDialog) return null;
  const r = (await a.OpenFileDialog()) as { path: string; name: string; content: string } | null;
  return r ?? null;
}

export async function createFile(rel: string, content: string): Promise<void> {
  const a = api();
  if (a?.CreateFile) await a.CreateFile(rel, content);
}

export async function createFolder(rel: string): Promise<void> {
  const a = api();
  if (a?.CreateFolder) await a.CreateFolder(rel);
}

export async function deletePath(rel: string): Promise<void> {
  const a = api();
  if (a?.DeletePath) await a.DeletePath(rel);
}

export async function renamePath(oldRel: string, newRel: string): Promise<void> {
  const a = api();
  if (a?.RenamePath) await a.RenamePath(oldRel, newRel);
}

export async function getSession(): Promise<Record<string, unknown>> {
  const a = api();
  if (!a?.GetSession) return {};
  return ((await a.GetSession()) as Record<string, unknown>) ?? {};
}

export async function saveSession(session: Record<string, unknown>): Promise<void> {
  const a = api();
  if (a?.SaveSession) await a.SaveSession(session);
}

export function isWails(): boolean {
  return typeof window !== 'undefined' && !!window.go?.main?.App;
}

export async function consumePendingFiles(): Promise<string[]> {
  const a = api();
  if (!a?.ConsumePendingFiles) return [];
  return ((await a.ConsumePendingFiles()) as string[]) ?? [];
}

/** Subscribe to Wails runtime events (e.g. open-files from second instance). */
export function onOpenFiles(handler: (paths: string[]) => void): () => void {
  const rt = (window as unknown as { runtime?: { EventsOn?: (e: string, cb: (...a: unknown[]) => void) => () => void } }).runtime;
  if (rt?.EventsOn) {
    return rt.EventsOn('open-files', (...args: unknown[]) => {
      const paths = (args[0] as string[]) || [];
      handler(paths);
    });
  }
  // Fallback: custom event
  const fn = (e: Event) => handler(((e as CustomEvent).detail as string[]) || []);
  window.addEventListener('mdman-open-files', fn);
  return () => window.removeEventListener('mdman-open-files', fn);
}

export async function snapshotWindow(): Promise<void> {
  const a = api();
  if (a?.SnapshotWindow) await a.SnapshotWindow();
}

export async function resolveMediaURL(markdownFile: string, mediaPath: string): Promise<string> {
  const a = api();
  if (!a?.ResolveMediaURL) return mediaPath;
  try {
    const url = (await a.ResolveMediaURL(markdownFile, mediaPath)) as string;
    return url || mediaPath;
  } catch {
    return mediaPath;
  }
}
