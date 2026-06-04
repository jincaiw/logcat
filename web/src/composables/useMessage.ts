import { createDiscreteApi } from 'naive-ui'

const { message } = createDiscreteApi(['message'])

export function useAppMessage() {
  return message
}
