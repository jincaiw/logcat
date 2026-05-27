import type { Directive, DirectiveBinding } from 'vue'
import { usePermissionStore } from '@/stores/permission'

const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const { value } = binding
    if (!value) return

    const permissionStore = usePermissionStore()

    if (typeof value === 'string') {
      if (!permissionStore.hasPermission(value)) {
        if (binding.modifiers.disable) {
          el.setAttribute('disabled', 'true')
          el.classList.add('n-button--disabled')
          el.classList.add('is-disabled')
        } else {
          el.parentNode?.removeChild(el)
        }
      }
    } else if (Array.isArray(value)) {
      const hasAny = value.some((code: string) => permissionStore.hasPermission(code))
      if (!hasAny) {
        if (binding.modifiers.disable) {
          el.setAttribute('disabled', 'true')
          el.classList.add('n-button--disabled')
          el.classList.add('is-disabled')
        } else {
          el.parentNode?.removeChild(el)
        }
      }
    }
  },
}

export default permission