<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { AlertCircle, ArrowUp, ArrowDown, Minus } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  priority: {
    type: String,
    required: true,
    validator: (value) => ['low', 'medium', 'high', 'urgent'].includes(value)
  },
  showIcon: {
    type: Boolean,
    default: true
  },
  size: {
    type: String,
    default: 'sm',
    validator: (value) => ['xs', 'sm', 'md'].includes(value)
  }
})

const config = {
  low: {
    bg: 'bg-gray-100 dark:bg-gray-800',
    text: 'text-gray-600 dark:text-gray-400',
    icon: ArrowDown
  },
  medium: {
    bg: 'bg-blue-100 dark:bg-blue-900/30',
    text: 'text-blue-600 dark:text-blue-400',
    icon: Minus
  },
  high: {
    bg: 'bg-orange-100 dark:bg-orange-900/30',
    text: 'text-orange-600 dark:text-orange-400',
    icon: ArrowUp
  },
  urgent: {
    bg: 'bg-red-100 dark:bg-red-900/30',
    text: 'text-red-600 dark:text-red-400',
    icon: AlertCircle
  }
}

const sizeClasses = {
  xs: 'text-xs px-1.5 py-0.5 gap-0.5',
  sm: 'text-xs px-2 py-0.5 gap-1',
  md: 'text-sm px-2.5 py-1 gap-1.5'
}

const iconSizes = {
  xs: 10,
  sm: 12,
  md: 14
}

const priorityConfig = computed(() => config[props.priority])
const label = computed(() => t(`tasks.priorities.${props.priority}`))
</script>

<template>
  <span
    :class="[
      priorityConfig.bg,
      priorityConfig.text,
      sizeClasses[size],
      'inline-flex items-center font-medium rounded-full transition-colors'
    ]"
  >
    <component
      v-if="showIcon"
      :is="priorityConfig.icon"
      :size="iconSizes[size]"
    />
    {{ label }}
  </span>
</template>
