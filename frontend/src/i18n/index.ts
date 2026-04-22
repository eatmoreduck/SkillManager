import { computed, ref } from 'vue'
import {
  fallbackLocale,
  localeLabels,
  messages,
  supportedLocales,
  type Locale,
  type MessageKey
} from './messages'

type TranslateParams = Record<string, string | number>

const locale = ref<Locale>(fallbackLocale)

function isSupportedLocale(value: string): value is Locale {
  return (supportedLocales as readonly string[]).includes(value)
}

export function normalizeLocale(value?: string | null): Locale {
  if (!value) {
    return fallbackLocale
  }

  if (isSupportedLocale(value)) {
    return value
  }

  const normalized = value.toLowerCase()
  if (normalized.startsWith('zh')) {
    return 'zh-CN'
  }
  if (normalized.startsWith('en')) {
    return 'en-US'
  }

  return fallbackLocale
}

export function setLocale(value?: string | null) {
  locale.value = normalizeLocale(value)
  if (typeof document !== 'undefined') {
    document.documentElement.lang = locale.value
  }
}

function interpolate(template: string, params?: TranslateParams): string {
  if (!params) {
    return template
  }

  return template.replace(/\{(\w+)\}/g, (_, key: string) => {
    const value = params[key]
    return value === undefined ? `{${key}}` : String(value)
  })
}

export function t(key: MessageKey, params?: TranslateParams): string {
  const currentMessages = messages[locale.value]
  const template = currentMessages[key] ?? messages[fallbackLocale][key] ?? key
  return interpolate(template, params)
}

export function formatDateTime(
  value: string | Date,
  options: Intl.DateTimeFormatOptions = { dateStyle: 'medium', timeStyle: 'short' }
): string {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) {
    return '-'
  }

  return new Intl.DateTimeFormat(locale.value, options).format(date)
}

export function formatRelativeTime(value: string | Date): string {
  const date = value instanceof Date ? value : new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }

  const diffMs = date.getTime() - Date.now()
  const absMs = Math.abs(diffMs)

  const units: Array<[Intl.RelativeTimeFormatUnit, number]> = [
    ['year', 1000 * 60 * 60 * 24 * 365],
    ['month', 1000 * 60 * 60 * 24 * 30],
    ['week', 1000 * 60 * 60 * 24 * 7],
    ['day', 1000 * 60 * 60 * 24],
    ['hour', 1000 * 60 * 60],
    ['minute', 1000 * 60]
  ]

  const formatter = new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' })

  for (const [unit, unitMs] of units) {
    if (absMs >= unitMs || unit === 'minute') {
      const valueForUnit = Math.round(diffMs / unitMs)
      return formatter.format(valueForUnit, unit)
    }
  }

  return formatter.format(0, 'minute')
}

export const languageOptions = supportedLocales.map((value) => ({
  label: localeLabels[value],
  value
}))

export function useI18n() {
  return {
    locale: computed(() => locale.value),
    languageOptions,
    setLocale,
    t
  }
}
