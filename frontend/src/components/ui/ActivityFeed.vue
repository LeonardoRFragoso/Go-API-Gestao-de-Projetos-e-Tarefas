<script setup>
import { useI18n } from 'vue-i18n'
import { CheckCircle, PlusCircle, Edit, Trash2, ArrowRight, Clock } from 'lucide-vue-next'
import UserAvatar from './UserAvatar.vue'

const { t } = useI18n()

defineProps({
  activities: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const getActivityIcon = (type) => {
  const icons = {
    task_created: PlusCircle,
    task_completed: CheckCircle,
    task_updated: Edit,
    task_deleted: Trash2,
    task_moved: ArrowRight
  }
  return icons[type] || Clock
}

const getActivityColor = (type) => {
  const colors = {
    task_created: 'text-green-500 bg-green-100 dark:bg-green-900/30',
    task_completed: 'text-blue-500 bg-blue-100 dark:bg-blue-900/30',
    task_updated: 'text-yellow-500 bg-yellow-100 dark:bg-yellow-900/30',
    task_deleted: 'text-red-500 bg-red-100 dark:bg-red-900/30',
    task_moved: 'text-purple-500 bg-purple-100 dark:bg-purple-900/30'
  }
  return colors[type] || 'text-gray-500 bg-gray-100 dark:bg-gray-800'
}

const formatTime = (date) => {
  const now = new Date()
  const activityDate = new Date(date)
  const diffMs = now - activityDate
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'Agora'
  if (diffMins < 60) return `${diffMins}m atrás`
  if (diffHours < 24) return `${diffHours}h atrás`
  if (diffDays < 7) return `${diffDays}d atrás`
  return activityDate.toLocaleDateString()
}
</script>

<template>
  <div class="space-y-4">
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="flex items-start gap-3 animate-pulse">
        <div class="w-8 h-8 bg-gray-200 dark:bg-dark-700 rounded-full"></div>
        <div class="flex-1">
          <div class="h-4 bg-gray-200 dark:bg-dark-700 rounded w-3/4 mb-2"></div>
          <div class="h-3 bg-gray-200 dark:bg-dark-700 rounded w-1/4"></div>
        </div>
      </div>
    </div>

    <div v-else-if="activities.length === 0" class="text-center py-8 text-gray-500 dark:text-gray-400">
      <Clock :size="32" class="mx-auto mb-2 opacity-50" />
      <p class="text-sm">{{ t('dashboard.noActivity') || 'Nenhuma atividade recente' }}</p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="activity in activities"
        :key="activity.id"
        class="flex items-start gap-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors"
      >
        <div :class="[getActivityColor(activity.type), 'p-2 rounded-full']">
          <component :is="getActivityIcon(activity.type)" :size="16" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm text-gray-900 dark:text-white">
            <span class="font-medium">{{ activity.user?.name || 'Usuário' }}</span>
            <span class="text-gray-600 dark:text-gray-400"> {{ activity.description }}</span>
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-500 mt-0.5">
            {{ formatTime(activity.created_at) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
