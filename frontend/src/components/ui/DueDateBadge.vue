<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Calendar, AlertTriangle, Clock } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  date: {
    type: String,
    required: true
  },
  showIcon: {
    type: Boolean,
    default: true
  }
})

const isOverdue = computed(() => {
  const due = new Date(props.date)
  const now = new Date()
  now.setHours(0, 0, 0, 0)
  due.setHours(0, 0, 0, 0)
  return due < now
})

const isToday = computed(() => {
  const due = new Date(props.date)
  const now = new Date()
  return due.toDateString() === now.toDateString()
})

const isTomorrow = computed(() => {
  const due = new Date(props.date)
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return due.toDateString() === tomorrow.toDateString()
})

const formattedDate = computed(() => {
  const due = new Date(props.date)
  return due.toLocaleDateString('pt-BR', { day: '2-digit', month: 'short' })
})

const statusLabel = computed(() => {
  if (isOverdue.value) return t('tasks.overdue')
  if (isToday.value) return t('tasks.dueToday')
  if (isTomorrow.value) return t('tasks.dueTomorrow')
  return formattedDate.value
})

const statusClasses = computed(() => {
  if (isOverdue.value) return 'bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400'
  if (isToday.value) return 'bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400'
  if (isTomorrow.value) return 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400'
  return 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400'
})

const StatusIcon = computed(() => {
  if (isOverdue.value) return AlertTriangle
  if (isToday.value) return Clock
  return Calendar
})
</script>

<template>
  <span
    :class="[
      statusClasses,
      'inline-flex items-center gap-1 text-xs font-medium px-2 py-0.5 rounded-full'
    ]"
  >
    <component v-if="showIcon" :is="StatusIcon" :size="12" />
    {{ statusLabel }}
  </span>
</template>
