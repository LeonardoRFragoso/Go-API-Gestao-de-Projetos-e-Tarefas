<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import { Bell, Check, CheckCheck, Trash2, X, AlertCircle, MessageSquare, Zap, Users } from 'lucide-vue-next'

const notificationsStore = useNotificationsStore()
const isOpen = ref(false)
const dropdownRef = ref(null)

onMounted(() => {
  notificationsStore.fetchUnreadCount()
  notificationsStore.fetchNotifications(1, 10)
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

const toggleDropdown = () => {
  isOpen.value = !isOpen.value
  if (isOpen.value) {
    notificationsStore.fetchNotifications(1, 10)
  }
}

const formatTime = (date) => {
  const now = new Date()
  const past = new Date(date)
  const diffMs = now - past
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)
  
  if (diffMins < 1) return 'agora'
  if (diffMins < 60) return `há ${diffMins} min`
  if (diffHours < 24) return `há ${diffHours}h`
  if (diffDays < 7) return `há ${diffDays} dias`
  return past.toLocaleDateString('pt-BR')
}

const getNotificationIcon = (type) => {
  const iconMap = {
    task_assigned: Zap,
    task_unassigned: Zap,
    task_updated: AlertCircle,
    task_comment: MessageSquare,
    task_due_soon: AlertCircle,
    task_overdue: AlertCircle,
    project_invite: AlertCircle,
    team_invite: Users,
    mention: MessageSquare
  }
  return iconMap[type] || Bell
}

const getNotificationColor = (type) => {
  const colorMap = {
    task_assigned: 'text-blue-600 dark:text-blue-400',
    task_unassigned: 'text-blue-600 dark:text-blue-400',
    task_updated: 'text-yellow-600 dark:text-yellow-400',
    task_comment: 'text-green-600 dark:text-green-400',
    task_due_soon: 'text-orange-600 dark:text-orange-400',
    task_overdue: 'text-red-600 dark:text-red-400',
    project_invite: 'text-purple-600 dark:text-purple-400',
    team_invite: 'text-purple-600 dark:text-purple-400',
    mention: 'text-green-600 dark:text-green-400'
  }
  return colorMap[type] || 'text-gray-600 dark:text-gray-400'
}

const getNotificationBg = (type) => {
  const bgMap = {
    task_assigned: 'bg-blue-100 dark:bg-blue-900/30',
    task_unassigned: 'bg-blue-100 dark:bg-blue-900/30',
    task_updated: 'bg-yellow-100 dark:bg-yellow-900/30',
    task_comment: 'bg-green-100 dark:bg-green-900/30',
    task_due_soon: 'bg-orange-100 dark:bg-orange-900/30',
    task_overdue: 'bg-red-100 dark:bg-red-900/30',
    project_invite: 'bg-purple-100 dark:bg-purple-900/30',
    team_invite: 'bg-purple-100 dark:bg-purple-900/30',
    mention: 'bg-green-100 dark:bg-green-900/30'
  }
  return bgMap[type] || 'bg-gray-100 dark:bg-gray-900/30'
}

const handleMarkAsRead = async (id) => {
  await notificationsStore.markAsRead(id)
}

const handleMarkAllAsRead = async () => {
  await notificationsStore.markAllAsRead()
}

const handleDelete = async (id) => {
  await notificationsStore.deleteNotification(id)
}
</script>

<template>
  <div ref="dropdownRef" class="relative">
    <button @click="toggleDropdown" class="relative p-2 text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">
      <Bell :size="20" />
      <span v-if="notificationsStore.hasUnread" 
        class="absolute top-0 right-0 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
        {{ notificationsStore.unreadCount > 9 ? '9+' : notificationsStore.unreadCount }}
      </span>
    </button>

    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div v-if="isOpen" class="absolute right-0 mt-2 w-80 md:w-96 bg-white dark:bg-dark-800 rounded-xl shadow-xl border border-gray-100 dark:border-dark-700 z-50 max-h-[70vh] flex flex-col overflow-hidden">
        <div class="p-4 border-b border-gray-100 dark:border-dark-600 flex items-center justify-between bg-gradient-to-r from-primary-50 to-primary-100 dark:from-primary-900/20 dark:to-primary-900/10">
          <div class="flex items-center gap-2">
            <Bell :size="20" class="text-primary-600 dark:text-primary-400" />
            <h3 class="font-semibold text-gray-900 dark:text-white">Notificações</h3>
            <span v-if="notificationsStore.unreadCount > 0" class="ml-2 px-2 py-1 bg-red-500 text-white text-xs rounded-full font-semibold">
              {{ notificationsStore.unreadCount }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <button v-if="notificationsStore.hasUnread" @click="handleMarkAllAsRead" 
              class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 flex items-center gap-1 transition-colors">
              <CheckCheck :size="16" />
              Marcar todas
            </button>
            <button @click="isOpen = false" class="p-1 hover:bg-gray-200 dark:hover:bg-dark-700 rounded transition-colors">
              <X :size="16" class="text-gray-600 dark:text-gray-400" />
            </button>
          </div>
        </div>

        <div class="overflow-y-auto flex-1">
          <div v-if="notificationsStore.loading" class="p-8 text-center text-gray-500">
            Carregando...
          </div>
          <div v-else-if="notificationsStore.notifications.length === 0" class="p-8 text-center text-gray-500 dark:text-gray-400">
            <Bell :size="32" class="mx-auto mb-2 opacity-50" />
            <p>Nenhuma notificação</p>
          </div>
          <div v-else class="divide-y divide-gray-100 dark:divide-dark-600">
            <div v-for="notification in notificationsStore.notifications" :key="notification.id"
              class="p-4 hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors cursor-pointer group"
              :class="{ 'bg-primary-50 dark:bg-primary-900/10 border-l-4 border-primary-500': !notification.read }">
              <div class="flex items-start gap-3">
                <div :class="[getNotificationBg(notification.type), 'p-2 rounded-lg flex-shrink-0']">
                  <component :is="getNotificationIcon(notification.type)" :size="18" :class="getNotificationColor(notification.type)" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-medium text-gray-900 dark:text-white text-sm">{{ notification.title }}</p>
                  <p class="text-sm text-gray-600 dark:text-gray-400 mt-0.5 line-clamp-2">{{ notification.message }}</p>
                  <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ formatTime(notification.created_at) }}</p>
                </div>
                <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
                  <button v-if="!notification.read" @click="handleMarkAsRead(notification.id)" 
                    class="p-1 text-gray-400 hover:text-green-600 dark:hover:text-green-400 transition-colors" title="Marcar como lida">
                    <Check :size="16" />
                  </button>
                  <button @click="handleDelete(notification.id)" 
                    class="p-1 text-gray-400 hover:text-red-600 dark:hover:text-red-400 transition-colors" title="Excluir">
                    <Trash2 :size="16" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="notificationsStore.notifications.length > 0" class="p-3 border-t border-gray-100 dark:border-dark-600 text-center bg-gray-50 dark:bg-dark-700/50">
          <RouterLink to="/notifications" @click="isOpen = false" 
            class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium transition-colors">
            Ver todas as notificações →
          </RouterLink>
        </div>
      </div>
    </Transition>
  </div>
</template>
