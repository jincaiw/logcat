<script setup lang="ts">
import {
  NConfigProvider, NMessageProvider, NDialogProvider, NNotificationProvider,
  darkTheme, enUS, dateEnUS, zhCN, dateZhCN, type GlobalThemeOverrides
} from 'naive-ui'
import { computed, onMounted, onUnmounted } from 'vue'
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

const theme = computed(() => {
  return themeStore.isDark ? darkTheme : undefined
})

const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#0071e3',
    primaryColorHover: '#0077ED',
    primaryColorPressed: '#0060b3',
    primaryColorSuppl: '#0071e3',
    borderRadius: '8px',
    borderRadiusSmall: '6px',
    fontFamily: "-apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', sans-serif",
    fontSize: '14px',
    fontSizeMedium: '14px',
    fontSizeSmall: '13px',
    fontSizeMini: '12px',
    heightMedium: '36px',
    heightSmall: '32px',
    heightLarge: '40px',
  },
  Card: {
    color: '#ffffff',
    colorModal: '#ffffff',
    borderColor: 'rgba(0, 0, 0, 0.08)',
    borderRadius: '10px',
    titleFontSizeMedium: '15px',
    titleFontWeight: '600',
    paddingMedium: '20px',
  },
  DataTable: {
    thColor: '#fafafa',
    tdColor: '#ffffff',
    tdColorHover: 'rgba(0, 113, 227, 0.04)',
    borderColor: 'rgba(0, 0, 0, 0.06)',
    thFontWeight: '600',
    tdColorStriped: 'rgba(0, 0, 0, 0.02)',
    fontSize: '13px',
  },
  Input: {
    color: '#ffffff',
    colorFocus: '#ffffff',
    borderColor: 'rgba(0, 0, 0, 0.12)',
    borderColorHover: 'rgba(0, 0, 0, 0.2)',
    borderColorFocus: '#0071e3',
    colorPlaceholder: '#aeaeb2',
    borderRadius: '8px',
    caretColor: '#0071e3',
  },
  Button: {
    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    fontWeightStrong: '600',
  },
  Select: {
    peers: {
      InternalSelection: {
        color: '#ffffff',
        borderRadius: '8px',
      },
      InternalSelectMenu: {
        color: '#ffffff',
        borderRadius: '8px',
        optionTextColor: '#86868b',
        optionTextColorActive: '#1d1d1f',
        optionColorPending: 'rgba(0, 113, 227, 0.08)',
        optionColorActive: 'rgba(0, 113, 227, 0.08)',
      },
    },
  },
  Modal: {
    color: '#ffffff',
    borderColor: 'rgba(0, 0, 0, 0.08)',
    borderRadius: '20px',
    titleFontSize: '16px',
    titleFontWeight: '600',
  },
  Tabs: {
    tabTextColorLine: '#aeaeb2',
    tabTextColorActiveLine: '#0071e3',
    tabTextColorHoverLine: '#86868b',
    barColor: '#0071e3',
  },
  Tag: { borderRadius: '6px' },
  Form: {
    labelTextColor: '#86868b',
    feedbackTextColor: '#aeaeb2',
  },
  Menu: {
    itemTextColor: '#86868b',
    itemTextColorHover: '#1d1d1f',
    itemTextColorActive: '#0071e3',
    itemTextColorActiveHover: '#0071e3',
    itemTextColorChildActive: '#0071e3',
    itemColorHover: 'rgba(0, 113, 227, 0.06)',
    itemColorActive: 'rgba(0, 113, 227, 0.08)',
    itemColorActiveHover: 'rgba(0, 113, 227, 0.12)',
    borderRadius: '8px',
  },
  Dropdown: {
    color: '#fafafa',
    borderRadius: '8px',
    optionTextColor: '#86868b',
    optionTextColorHover: '#1d1d1f',
    optionColorHover: 'rgba(0, 113, 227, 0.06)',
  },
  Pagination: {
    itemTextColor: '#86868b',
    itemTextColorHover: '#1d1d1f',
    itemTextColorActive: '#ffffff',
    itemBorderRadius: '6px',
    itemColor: 'transparent',
    itemColorHover: 'rgba(0, 113, 227, 0.06)',
    itemColorActive: '#0071e3',
  },
  Switch: { railColorActive: '#0071e3' },
  Popover: { color: '#fafafa', borderRadius: '8px' },
  Tooltip: { color: '#1c1c1e', borderRadius: '6px', textColor: '#f5f5f7' },
  Message: { borderRadius: '8px' },
  Dialog: { borderRadius: '20px', color: '#ffffff' },
  Notification: { borderRadius: '10px' },
  Scrollbar: {
    color: 'rgba(0,0,0,0.15)',
    colorHover: 'rgba(0,0,0,0.25)',
    width: '6px',
    widthHover: '6px',
    borderRadius: '3px',
    borderRadiusHover: '3px',
  },
  Descriptions: {
    thColor: '#fafafa',
    tdColor: '#ffffff',
    borderColor: 'rgba(0, 0, 0, 0.06)',
  },
  Divider: {
    color: 'rgba(0, 0, 0, 0.08)',
  },
  Alert: {
    borderRadius: '8px',
  },
}

const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#0a84ff',
    primaryColorHover: '#409cff',
    primaryColorPressed: '#006fdd',
    primaryColorSuppl: '#0a84ff',
    borderRadius: '8px',
    borderRadiusSmall: '6px',
    fontFamily: "-apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', sans-serif",
    fontSize: '14px',
    fontSizeMedium: '14px',
    fontSizeSmall: '13px',
    fontSizeMini: '12px',
    heightMedium: '36px',
    heightSmall: '32px',
    heightLarge: '40px',
  },
  Card: {
    color: '#1c1c1e',
    colorModal: '#1c1c1e',
    borderColor: 'rgba(255, 255, 255, 0.08)',
    borderRadius: '10px',
    titleFontSizeMedium: '15px',
    titleFontWeight: '600',
    paddingMedium: '20px',
  },
  DataTable: {
    thColor: '#2c2c2e',
    tdColor: '#1c1c1e',
    tdColorHover: 'rgba(255, 255, 255, 0.04)',
    borderColor: 'rgba(255, 255, 255, 0.06)',
    thFontWeight: '600',
    tdColorStriped: 'rgba(255, 255, 255, 0.02)',
    fontSize: '13px',
  },
  Input: {
    color: '#1c1c1e',
    colorFocus: '#2c2c2e',
    borderColor: 'rgba(255, 255, 255, 0.1)',
    borderColorHover: 'rgba(255, 255, 255, 0.2)',
    borderColorFocus: '#0a84ff',
    colorPlaceholder: '#636366',
    borderRadius: '8px',
    caretColor: '#0a84ff',
  },
  Button: {
    borderRadiusMedium: '8px',
    borderRadiusSmall: '6px',
    fontWeightStrong: '600',
  },
  Select: {
    peers: {
      InternalSelection: {
        color: '#1c1c1e',
        borderRadius: '8px',
      },
      InternalSelectMenu: {
        color: '#2c2c2e',
        borderRadius: '8px',
        optionTextColor: '#98989d',
        optionTextColorActive: '#f5f5f7',
        optionColorPending: 'rgba(10, 132, 255, 0.12)',
        optionColorActive: 'rgba(10, 132, 255, 0.12)',
      },
    },
  },
  Modal: {
    color: '#1c1c1e',
    borderColor: 'rgba(255, 255, 255, 0.08)',
    borderRadius: '20px',
    titleFontSize: '16px',
    titleFontWeight: '600',
  },
  Tabs: {
    tabTextColorLine: '#636366',
    tabTextColorActiveLine: '#0a84ff',
    tabTextColorHoverLine: '#98989d',
    barColor: '#0a84ff',
  },
  Tag: { borderRadius: '6px' },
  Form: {
    labelTextColor: '#98989d',
    feedbackTextColor: '#636366',
  },
  Menu: {
    itemTextColor: '#98989d',
    itemTextColorHover: '#f5f5f7',
    itemTextColorActive: '#0a84ff',
    itemTextColorActiveHover: '#0a84ff',
    itemTextColorChildActive: '#0a84ff',
    itemColorHover: 'rgba(255, 255, 255, 0.06)',
    itemColorActive: 'rgba(10, 132, 255, 0.12)',
    itemColorActiveHover: 'rgba(10, 132, 255, 0.18)',
    borderRadius: '8px',
  },
  Dropdown: {
    color: '#2c2c2e',
    borderRadius: '8px',
    optionTextColor: '#98989d',
    optionTextColorHover: '#f5f5f7',
    optionColorHover: 'rgba(255, 255, 255, 0.06)',
  },
  Pagination: {
    itemTextColor: '#98989d',
    itemTextColorHover: '#f5f5f7',
    itemTextColorActive: '#ffffff',
    itemBorderRadius: '6px',
    itemColor: 'transparent',
    itemColorHover: 'rgba(255, 255, 255, 0.06)',
    itemColorActive: '#0a84ff',
  },
  Switch: { railColorActive: '#0a84ff' },
  Popover: { color: '#2c2c2e', borderRadius: '8px' },
  Tooltip: { color: '#2c2c2e', borderRadius: '6px', textColor: '#f5f5f7' },
  Message: { borderRadius: '8px' },
  Dialog: { borderRadius: '20px', color: '#1c1c1e' },
  Notification: { borderRadius: '10px' },
  Scrollbar: {
    color: 'rgba(255, 255, 255, 0.12)',
    colorHover: 'rgba(255, 255, 255, 0.2)',
    width: '6px',
    widthHover: '6px',
    borderRadius: '3px',
    borderRadiusHover: '3px',
  },
  Descriptions: {
    thColor: '#2c2c2e',
    tdColor: '#1c1c1e',
    borderColor: 'rgba(255, 255, 255, 0.06)',
  },
  Divider: {
    color: 'rgba(255, 255, 255, 0.08)',
  },
  Alert: {
    borderRadius: '8px',
  },
}

const themeOverrides = computed(() => {
  return themeStore.isDark ? darkThemeOverrides : lightThemeOverrides
})

const naiveLocale = computed(() => locale.value === 'zh-CN' ? zhCN : enUS)
const naiveDateLocale = computed(() => locale.value === 'zh-CN' ? dateZhCN : dateEnUS)
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
