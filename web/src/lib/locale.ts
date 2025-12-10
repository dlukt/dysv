export const supportedLocales = ['de', 'en', 'hr'] as const
export type Locale = (typeof supportedLocales)[number]

// Returns best-effort locale inferred from the first path segment.
export function getLocaleFromPath(pathname: string): Locale {
  const segment = pathname.split('/').filter(Boolean)[0]
  if (supportedLocales.includes(segment as Locale)) {
    return segment as Locale
  }
  return 'de'
}

export function switchLocale(pathname: string, targetLocale: Locale): string {
  const segments = pathname.split('/').filter(Boolean)
  if (supportedLocales.includes(segments[0] as Locale)) {
    segments.shift()
  }
  if (targetLocale !== 'de') {
    segments.unshift(targetLocale)
  }
  return '/' + segments.join('/')
}
