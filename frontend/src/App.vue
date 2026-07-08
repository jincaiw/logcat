<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import {
  darkTheme,
  dateEnUS,
  dateZhCN,
  enUS,
  zhCN,
  NConfigProvider,
  NDialogProvider,
  NMessageProvider,
  NNotificationProvider,
  type GlobalThemeOverrides,
} from 'naive-ui'
import { useI18n } from '@/i18n'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'

const { initI18n, locale } = useI18n()
const themeStore = useThemeStore()
const authStore = useAuthStore()

function handleAuthExpired() {
  authStore.clearToken()
}

onMounted(() => {
  initI18n()
  window.addEventListener('auth:expired', handleAuthExpired)
})

onUnmounted(() => {
  window.removeEventListener('auth:expired', handleAuthExpired)
})

const theme = computed(() => (themeStore.isDark ? darkTheme : undefined))

const sharedCommon = {
  primaryColor: '#2563eb',
  primaryColorHover: '#1d4ed8',
  primaryColorPressed: '#1e40af',
  primaryColorSuppl: '#2563eb',
  borderRadius: '14px',
  borderRadiusSmall: '10px',
  fontFamily: "-apple-system, BlinkMacSystemFont, 'SF Pro Text', 'Segoe UI', sans-serif",
  fontSize: '14px',
  fontSizeMedium: '14px',
  fontSizeSmall: '13px',
  fontSizeMini: '12px',
  heightMedium: '36px',
  heightSmall: '32px',
  heightLarge: '42px',
}

const sharedComponents: Pick<GlobalThemeOverrides, 'Button' | 'Tag' | 'Message' | 'Notification' | 'Alert'> = {
  Button: {
    borderRadiusMedium: '999px',
    borderRadiusSmall: '999px',
    fontWeightStrong: '600',
  },
  Tag: { borderRadius: '999px' },
  Message: { borderRadius: '12px' },
  Notification: { borderRadius: '16px' },
  Alert: { borderRadius: '14px' },
}

function createThemeOverrides(isDark: boolean): GlobalThemeOverrides {
  const accent = isDark ? '#60a5fa' : '#2563eb'
  const accentHover = isDark ? '#93c5fd' : '#1d4ed8'
  const accentSoft = isDark ? 'rgba(96, 165, 250, 0.12)' : 'rgba(37, 99, 235, 0.08)'
  const accentSoftHover = isDark ? 'rgba(96, 165, 250, 0.18)' : 'rgba(37, 99, 235, 0.12)'
  const surface = isDark ? '#111827' : '#ffffff'
  const surfaceAlt = isDark ? '#162033' : '#f8fafc'
  const surfaceSoft = isDark ? '#0f172a' : '#eef2f7'
  const border = isDark ? 'rgba(148, 163, 184, 0.14)' : 'rgba(15, 23, 42, 0.08)'
  const borderStrong = isDark ? 'rgba(148, 163, 184, 0.2)' : 'rgba(15, 23, 42, 0.14)'
  const textSecondary = isDark ? '#b1b9c8' : '#556070'
  const textMuted = isDark ? '#7f8aa2' : '#8893a6'

  return {
    common: {
      ...sharedCommon,
      primaryColor: accent,
      primaryColorHover: accentHover,
      primaryColorPressed: isDark ? '#2563eb' : '#1e40af',
      primaryColorSuppl: accent,
      bodyColor: surface,
      baseColor: surface,
      cardColor: surface,
      modalColor: surface,
      popoverColor: surface,
      tableColor: surface,
      textColorBase: isDark ? '#f8fafc' : '#0f172a',
      textColor1: isDark ? '#f8fafc' : '#0f172a',
      textColor2: textSecondary,
      textColor3: textMuted,
      textColorDisabled: textMuted,
      placeholderColor: textMuted,
      dividerColor: border,
      borderColor: border,
      borderRadius: '14px',
      borderRadiusSmall: '10px',
      boxShadow1: isDark ? '0 1px 2px rgba(0, 0, 0, 0.22)' : '0 1px 2px rgba(15, 23, 42, 0.04)',
      boxShadow2: isDark ? '0 10px 30px rgba(0, 0, 0, 0.24)' : '0 10px 30px rgba(15, 23, 42, 0.06)',
      boxShadow3: isDark ? '0 18px 48px rgba(0, 0, 0, 0.34)' : '0 18px 48px rgba(15, 23, 42, 0.1)',
      opacityDisabled: '0.5',
    },
    ...sharedComponents,
    Card: {
      color: surface,
      colorModal: surface,
      borderColor: border,
      borderRadius: '18px',
      titleFontSizeMedium: '15px',
      titleFontWeight: '650',
      paddingMedium: '20px',
    },
    DataTable: {
      thColor: surfaceAlt,
      tdColor: surface,
      tdColorHover: accentSoftHover,
      borderColor: isDark ? 'rgba(255, 255, 255, 0.06)' : 'rgba(15, 23, 42, 0.06)',
      thFontWeight: '600',
      tdColorStriped: surfaceSoft,
      fontSize: '13px',
    },
    Input: {
      color: surface,
      colorFocus: surface,
      borderColor: isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(15, 23, 42, 0.12)',
      borderColorHover: accentSoftHover,
      borderColorFocus: accent,
      colorPlaceholder: textMuted,
      borderRadius: '12px',
      caretColor: accent,
    },
    Select: {
      peers: {
        InternalSelection: {
          color: surface,
          borderRadius: '12px',
        },
        InternalSelectMenu: {
          color: surface,
          borderRadius: '14px',
          optionTextColor: textSecondary,
          optionTextColorActive: isDark ? '#f8fafc' : '#0f172a',
          optionColorPending: accentSoft,
          optionColorActive: accentSoft,
        },
      },
    },
    Modal: {
      color: surface,
      borderColor: border,
      borderRadius: '20px',
      titleFontSize: '16px',
      titleFontWeight: '650',
    },
    Tabs: {
      tabTextColorLine: textMuted,
      tabTextColorActiveLine: accent,
      tabTextColorHoverLine: isDark ? '#e2e8f0' : '#334155',
      barColor: accent,
    },
    Form: {
      labelTextColor: textSecondary,
      feedbackTextColor: textMuted,
    },
    Menu: {
      itemTextColor: textSecondary,
      itemTextColorHover: isDark ? '#f8fafc' : '#0f172a',
      itemTextColorActive: accent,
      itemTextColorActiveHover: accent,
      itemTextColorChildActive: accent,
      itemColorHover: accentSoft,
      itemColorActive: accentSoftHover,
      itemColorActiveHover: accentSoftHover,
      borderRadius: '12px',
    },
    Dropdown: {
      color: surface,
      borderRadius: '14px',
      optionTextColor: textSecondary,
      optionTextColorHover: isDark ? '#f8fafc' : '#0f172a',
      optionColorHover: accentSoft,
    },
    Pagination: {
      itemTextColor: textSecondary,
      itemTextColorHover: isDark ? '#f8fafc' : '#0f172a',
      itemTextColorActive: '#ffffff',
      itemBorderRadius: '999px',
      itemColor: 'transparent',
      itemColorHover: accentSoft,
      itemColorActive: accent,
    },
    Switch: { railColorActive: accent },
    Popover: { color: surface, borderRadius: '14px' },
    Tooltip: { color: isDark ? '#0f172a' : '#0f172a', borderRadius: '10px', textColor: '#f8fafc' },
    Scrollbar: {
      color: isDark ? 'rgba(255, 255, 255, 0.12)' : 'rgba(15, 23, 42, 0.18)',
      colorHover: isDark ? 'rgba(255, 255, 255, 0.2)' : 'rgba(15, 23, 42, 0.28)',
      width: '6px',
      widthHover: '6px',
      borderRadius: '999px',
      borderRadiusHover: '999px',
    },
    Descriptions: {
      thColor: surfaceAlt,
      tdColor: surface,
      borderColor: border,
    },
    Divider: {
      color: borderStrong,
    },
  }
}

const lightThemeOverrides = createThemeOverrides(false)
const darkThemeOverrides = createThemeOverrides(true)
const themeOverrides = computed(() => (themeStore.isDark ? darkThemeOverrides : lightThemeOverrides))
const naiveLocale = computed(() => (locale.value === 'zh-CN' ? zhCN : enUS))
const naiveDateLocale = computed(() => (locale.value === 'zh-CN' ? dateZhCN : dateEnUS))
</script>

<template>
  <NConfigProvider :theme="theme" :locale="naiveLocale" :date-locale="naiveDateLocale" :theme-overrides="themeOverrides">
    <NMessageProvider>
      <NDialogProvider>
        <NNotificationProvider>
          <router-view />
        </NNotificationProvider>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>
