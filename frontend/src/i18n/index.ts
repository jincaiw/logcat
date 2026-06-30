// 轻量级 i18n 实现：基于 Vue reactive API，无外部依赖。
// 支持语言切换、持久化、嵌套键访问（如 "menu.dashboard"）。
import { reactive, computed, type ComputedRef } from 'vue'
import zhCN from './zh-CN'
import enUS from './en-US'

export type Locale = 'zh-CN' | 'en-US'

type MessageDictionary = typeof zhCN

const messages: Record<Locale, MessageDictionary> = {
  'zh-CN': zhCN,
  'en-US': enUS,
}

const STORAGE_KEY = 'logcat-locale'

interface I18nState {
  locale: Locale
}

const state = reactive<I18nState>({
  locale: detectLocale(),
})

function detectLocale(): Locale {
  // 1. localStorage
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'zh-CN' || saved === 'en-US') {
    return saved
  }
  // 2. 系统默认语言保持中文
  return 'zh-CN'
}

function getNestedValue(obj: unknown, path: string): unknown {
  const keys = path.split('.')
  let current: unknown = obj
  for (const key of keys) {
    if (current === null || current === undefined || typeof current !== 'object') {
      return undefined
    }
    current = (current as Record<string, unknown>)[key]
  }
  return current
}

/** 翻译函数：t('menu.dashboard') => '系统状态'
 *  支持模板变量：t('parseTemplate.parseSuccess', { count: 5 }) => '解析成功 - 5 个字段'
 */
function t(key: string, params?: Record<string, string | number>): string {
  const value = getNestedValue(messages[state.locale], key)
  let result: string
  if (typeof value === 'string') {
    result = value
  } else {
    // 回退到中文
    const fallback = getNestedValue(messages['zh-CN'], key)
    if (typeof fallback === 'string') {
      result = fallback
    } else {
      // 找不到则返回 key 本身
      return key
    }
  }
  // 替换模板变量 {{var}}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      result = result.replace(new RegExp(`\\{\\{${k}\\}\\}`, 'g'), String(v))
    }
  }
  return result
}

/** 切换语言 */
function setLocale(locale: Locale): void {
  state.locale = locale
  localStorage.setItem(STORAGE_KEY, locale)
  document.documentElement.setAttribute('lang', locale)
}

/** 初始化 i18n（在应用启动时调用） */
function initI18n(): void {
  document.documentElement.setAttribute('lang', state.locale)
}

/** 响应式 locale（用于模板绑定） */
const locale: ComputedRef<Locale> = computed(() => state.locale)

/** 可用语言列表 */
const availableLocales: { value: Locale; label: string }[] = [
  { value: 'zh-CN', label: '简体中文' },
  { value: 'en-US', label: 'English' },
]

export function useI18n() {
  return {
    locale,
    t,
    setLocale,
    initI18n,
    availableLocales,
  }
}
