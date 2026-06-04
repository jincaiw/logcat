import { ref, computed } from 'vue'
import { getSystemConfigs } from '@/api/system'
import type { SystemConfig } from '@/types'

const timezone = ref('Asia/Shanghai')
const timeFormat = ref('YYYY-MM-DD HH:mm:ss')
const loaded = ref(false)

const COMMON_TIMEZONES = [
  { label: 'UTC', value: 'UTC' },
  { label: 'Asia/Shanghai (北京)', value: 'Asia/Shanghai' },
  { label: 'Asia/Tokyo (东京)', value: 'Asia/Tokyo' },
  { label: 'Asia/Seoul (首尔)', value: 'Asia/Seoul' },
  { label: 'Asia/Singapore (新加坡)', value: 'Asia/Singapore' },
  { label: 'Asia/Kolkata (孟买)', value: 'Asia/Kolkata' },
  { label: 'Asia/Dubai (迪拜)', value: 'Asia/Dubai' },
  { label: 'Europe/London (伦敦)', value: 'Europe/London' },
  { label: 'Europe/Paris (巴黎)', value: 'Europe/Paris' },
  { label: 'Europe/Berlin (柏林)', value: 'Europe/Berlin' },
  { label: 'Europe/Moscow (莫斯科)', value: 'Europe/Moscow' },
  { label: 'America/New_York (纽约)', value: 'America/New_York' },
  { label: 'America/Chicago (芝加哥)', value: 'America/Chicago' },
  { label: 'America/Denver (丹佛)', value: 'America/Denver' },
  { label: 'America/Los_Angeles (洛杉矶)', value: 'America/Los_Angeles' },
  { label: 'Pacific/Auckland (奥克兰)', value: 'Pacific/Auckland' },
  { label: 'Australia/Sydney (悉尼)', value: 'Australia/Sydney' },
]

const TIME_FORMAT_OPTIONS = [
  { label: 'YYYY-MM-DD HH:mm:ss', value: 'YYYY-MM-DD HH:mm:ss' },
  { label: 'YYYY/MM/DD HH:mm:ss', value: 'YYYY/MM/DD HH:mm:ss' },
  { label: 'DD-MM-YYYY HH:mm:ss', value: 'DD-MM-YYYY HH:mm:ss' },
  { label: 'MM/DD/YYYY HH:mm:ss', value: 'MM/DD/YYYY HH:mm:ss' },
  { label: 'YYYY-MM-DD HH:mm', value: 'YYYY-MM-DD HH:mm' },
]

function pad(n: number): string {
  return n < 10 ? '0' + n : String(n)
}

function formatToPattern(date: Date, pattern: string): string {
  const y = date.getFullYear()
  const M = date.getMonth() + 1
  const d = date.getDate()
  const H = date.getHours()
  const m = date.getMinutes()
  const s = date.getSeconds()

  return pattern
    .replace('YYYY', String(y))
    .replace('MM', pad(M))
    .replace('DD', pad(d))
    .replace('HH', pad(H))
    .replace('mm', pad(m))
    .replace('ss', pad(s))
}

function getDateInTimezone(dateInput: string | Date | null | undefined): Date | null {
  if (!dateInput) return null
  const d = typeof dateInput === 'string' ? new Date(dateInput) : dateInput
  if (isNaN(d.getTime())) return null
  return d
}

function formatTime(dateInput: string | Date | null | undefined): string {
  const d = getDateInTimezone(dateInput)
  if (!d) return '--'

  try {
    const tz = timezone.value
    const parts = new Intl.DateTimeFormat('en-US', {
      timeZone: tz,
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).formatToParts(d)

    const get = (type: string) => parts.find(p => p.type === type)?.value || ''

    const reconstructed = new Date(
      Number(get('year')),
      Number(get('month')) - 1,
      Number(get('day')),
      Number(get('hour')),
      Number(get('minute')),
      Number(get('second')),
    )

    return formatToPattern(reconstructed, timeFormat.value)
  } catch {
    return formatToPattern(d, timeFormat.value)
  }
}

function formatTimeShort(dateInput: string | Date | null | undefined): string {
  const d = getDateInTimezone(dateInput)
  if (!d) return '--'

  try {
    const tz = timezone.value
    const parts = new Intl.DateTimeFormat('en-US', {
      timeZone: tz,
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).formatToParts(d)

    const get = (type: string) => parts.find(p => p.type === type)?.value || ''
    return `${get('month')}-${get('day')} ${get('hour')}:${get('minute')}`
  } catch {
    const m = d.getMonth() + 1
    const day = d.getDate()
    const h = d.getHours()
    const min = d.getMinutes()
    return `${pad(m)}-${pad(day)} ${pad(h)}:${pad(min)}`
  }
}

async function loadTimezoneConfig() {
  if (loaded.value) return
  try {
    const res = await getSystemConfigs()
    if (res.data) {
      for (const cfg of res.data as SystemConfig[]) {
        if (cfg.configKey === 'timezone' && cfg.configValue) {
          timezone.value = cfg.configValue
        }
        if (cfg.configKey === 'timeFormat' && cfg.configValue) {
          timeFormat.value = cfg.configValue
        }
      }
    }
  } catch { /* ignore */ }
  loaded.value = true
}

function setTimezone(tz: string) {
  timezone.value = tz
  localStorage.setItem('timezone', tz)
}

function setTimeFormat(fmt: string) {
  timeFormat.value = fmt
  localStorage.setItem('timeFormat', fmt)
}

function initFromLocalStorage() {
  const savedTz = localStorage.getItem('timezone')
  const savedFmt = localStorage.getItem('timeFormat')
  if (savedTz) timezone.value = savedTz
  if (savedFmt) timeFormat.value = savedFmt
}

initFromLocalStorage()

export function useTimeFormat() {
  return {
    timezone: computed(() => timezone.value),
    timeFormat: computed(() => timeFormat.value),
    commonTimezones: COMMON_TIMEZONES,
    timeFormatOptions: TIME_FORMAT_OPTIONS,
    formatTime,
    formatTimeShort,
    loadTimezoneConfig,
    setTimezone,
    setTimeFormat,
  }
}
