<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import { Bell, Check, CheckCheck, Trash2, X } from 'lucide-vue-next'

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
  const icons = {
    task_assigned: '📋',
    task_unassigned: '📋',
    task_updated: '✏️',
    task_comment: '💬',
    task_due_soon: '⏰',
    task_overdue: '⚠️',
    project_invite: '📁',
    team_invite: '👥',
    mention: '@'
  }
  return icons[type] || '🔔'
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
      <div v-if="isOpen" class="absolute right-0 mt-2 w-80 md:w-96 bg-white dark:bg-dark-800 rounded-xl shadow-lg border border-gray-100 dark:border-dark-700 z-50 max-h-[70vh] flex flex-col">
        <div class="p-4 border-b border-gray-100 dark:border-dark-600 flex items-center justify-between">
          <h3 class="font-semibold text-gray-900 dark:text-white">Notificações</h3>
          <div class="flex items-center gap-2">
            <button v-if="notificationsStore.hasUnread" @click="handleMarkAllAsRead" 
              class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 flex items-center">
              <CheckCheck :size="16" class="mr-1" />
              Marcar todas
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
              class="p-4 hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
              :class="{ 'bg-primary-50 dark:bg-primary-900/10': !notification.read }">
              <div class="flex items-start gap-3">
                <span class="text-xl">{{ getNotificationIcon(notification.type) }}</span>
                <div class="flex-1 min-w-0">
                  <p class="font-medium text-gray-900 dark:text-white text-sm">{{ notification.title }}</p>
                  <p class="text-sm text-gray-600 dark:text-gray-400 mt-0.5 line-clamp-2">{{ notification.message }}</p>
                  <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ formatTime(notification.created_at) }}</p>
                </div>
                <div class="flex items-center gap-1">
                  <button v-if="!notification.read" @click="handleMarkAsRead(notification.id)" 
                    class="p-1 text-gray-400 hover:text-green-600" title="Marcar como lida">
                    <Check :size="16" />
                  </button>
                  <button @click="handleDelete(notification.id)" 
                    class="p-1 text-gray-400 hover:text-red-600" title="Excluir">
                    <Trash2 :size="16" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="notificationsStore.notifications.length > 0" class="p-3 border-t border-gray-100 dark:border-dark-600 text-center">
          <RouterLink to="/notifications" @click="isOpen = false" 
            class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700">
            Ver todas as notificações
          </RouterLink>
        </div>
      </div>
    </Transition>
  </div>
</template>
