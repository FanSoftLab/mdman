import en from '../../locales/en';
import uk from '../../locales/uk';

export type Locale = 'en' | 'uk';

const dictionaries: Record<Locale, Record<string, string>> = { en, uk };

// Use Svelte 5 runes for reactivity
let current = $state<Locale>('en');
const listeners = new Set<() => void>();

export function initLocale(code?: string) {
  if (code === 'uk' || code === 'en') current = code;
  else if (typeof navigator !== 'undefined' && navigator.language?.startsWith('uk')) current = 'uk';
  else current = 'en';
}

export function setLocale(code: Locale) {
  current = code;
  listeners.forEach((fn) => fn());
  try {
    localStorage.setItem('mdman.locale', code);
  } catch { /* ignore */ }
}

export function getLocale(): Locale {
  return current;
}

export function onLocaleChange(fn: () => void) {
  listeners.add(fn);
  return () => listeners.delete(fn);
}

export function t(key: string, vars?: Record<string, string | number>): string {
  let s = dictionaries[current][key] ?? dictionaries.en[key] ?? key;
  if (vars) {
    for (const [k, v] of Object.entries(vars)) {
      s = s.replace(new RegExp(`\\{${k}\\}`, 'g'), String(v));
    }
  }
  return s;
}
