// HTTP 客户端基础封装：统一请求/响应处理、错误转换、认证注入。
import type { ApiSuccess } from '@/types'

const API_BASE = '/api'
const TOKEN_STORAGE_KEY = 'logcat-token'

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly endpoint: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/** 读取本地存储的认证 token */
export function getAuthToken(): string | null {
  return localStorage.getItem(TOKEN_STORAGE_KEY)
}

/** 保存认证 token */
export function setAuthToken(token: string): void {
  localStorage.setItem(TOKEN_STORAGE_KEY, token)
}

/** 清除认证 token */
export function clearAuthToken(): void {
  localStorage.removeItem(TOKEN_STORAGE_KEY)
}

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: unknown
  query?: Record<string, string | number | boolean | undefined>
  /** 是否跳过认证头注入（用于登录接口） */
  skipAuth?: boolean
}

/** 构造查询字符串 */
function buildQueryString(params?: Record<string, string | number | boolean | undefined>): string {
  if (!params) return ''
  const entries = Object.entries(params).filter(([, v]) => v !== undefined && v !== '')
  if (entries.length === 0) return ''
  const search = new URLSearchParams()
  for (const [k, v] of entries) {
    search.set(k, String(v))
  }
  return `?${search.toString()}`
}

let isHandlingUnauthorized = false

/** 401 处理：清除 token、通知 store 重置状态、跳转登录页 */
function handleUnauthorized(): void {
  if (isHandlingUnauthorized) return
  isHandlingUnauthorized = true
  clearAuthToken()
  window.dispatchEvent(new CustomEvent('auth:expired'))
  if (!window.location.hash.startsWith('#/login')) {
    window.location.hash = '#/login'
  }
  setTimeout(() => { isHandlingUnauthorized = false }, 1000)
}

/** 发送 HTTP 请求并解析 JSON 响应 */
async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, query, skipAuth = false } = options
  const url = `${API_BASE}/${endpoint}${buildQueryString(query)}`

  const headers: Record<string, string> = { 'Content-Type': 'application/json' }

  // 注入认证头
  if (!skipAuth) {
    const token = getAuthToken()
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const init: RequestInit = {
    method,
    headers,
  }

  if (body !== undefined) {
    init.body = typeof body === 'string' ? body : JSON.stringify(body)
  }

  let response: Response
  try {
    response = await fetch(url, init)
  } catch (err) {
    throw new ApiError(
      `Network error: ${err instanceof Error ? err.message : String(err)}`,
      0,
      endpoint,
    )
  }

  // 401 未认证：清除 token 并跳转登录页
  if (response.status === 401 && !skipAuth) {
    handleUnauthorized()
    throw new ApiError('未认证或会话已过期', 401, endpoint)
  }

  // 处理空响应
  const text = await response.text()
  if (!text) {
    if (!response.ok) {
      throw new ApiError(`HTTP ${response.status}`, response.status, endpoint)
    }
    return undefined as T
  }

  let data: unknown
  try {
    data = JSON.parse(text)
  } catch {
    if (!response.ok) {
      throw new ApiError(`HTTP ${response.status}: ${text}`, response.status, endpoint)
    }
    return text as unknown as T
  }

  if (!response.ok) {
    const errMsg =
      (typeof data === 'object' && data !== null && 'error' in data
        ? String((data as { error: string }).error)
        : `HTTP ${response.status}`)
    throw new ApiError(errMsg, response.status, endpoint)
  }

  return data as T
}

// ==================== 便捷方法 ====================

export const http = {
  get: <T>(endpoint: string, query?: RequestOptions['query']) =>
    request<T>(endpoint, { method: 'GET', query }),

  post: <T>(endpoint: string, body?: unknown, skipAuth = false) =>
    request<T>(endpoint, { method: 'POST', body, skipAuth }),

  put: <T>(endpoint: string, body?: unknown) =>
    request<T>(endpoint, { method: 'PUT', body }),

  delete: <T>(endpoint: string, query?: RequestOptions['query']) =>
    request<T>(endpoint, { method: 'DELETE', query }),
}

// ==================== 通用响应类型便捷别名 ====================

export type { ApiSuccess }
