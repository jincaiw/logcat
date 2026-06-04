import { ref, onMounted, onUnmounted } from 'vue'

const isMobile = ref(false)
const isTablet = ref(false)

let initialized = false
const listeners: (() => void)[] = []

function update() {
  const w = window.innerWidth
  isMobile.value = w < 768
  isTablet.value = w >= 768 && w < 1024
}

let resizeTimer: ReturnType<typeof setTimeout> | null = null
function onResize() {
  if (resizeTimer) clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => {
    update()
    listeners.forEach(fn => fn())
  }, 100)
}

function init() {
  if (initialized) return
  initialized = true
  update()
  window.addEventListener('resize', onResize, { passive: true })
}

export function useIsMobile() {
  init()

  const localIsMobile = ref(isMobile.value)
  const localIsTablet = ref(isTablet.value)

  const sync = () => {
    localIsMobile.value = isMobile.value
    localIsTablet.value = isTablet.value
  }

  onMounted(() => {
    sync()
    listeners.push(sync)
  })

  onUnmounted(() => {
    const idx = listeners.indexOf(sync)
    if (idx > -1) listeners.splice(idx, 1)
  })

  return { isMobile: localIsMobile, isTablet: localIsTablet }
}

export function responsiveWidth(maxWidth: number): string {
  return `min(${maxWidth}px, calc(100vw - 32px))`
}
