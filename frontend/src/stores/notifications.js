import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios'

export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref([])
  const unreadCount = ref(0)
  const loading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0 })

  const hasUnread = computed(() => unreadCount.value > 0)

  async function fetchNotifications(page = 1, limit = 20, unreadOnly = false) {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/notifications', { 
        params: { page, limit, unread_only: unreadOnly } 
      })
      notifications.value = response.data.data
      pagination.value = response.data.pagination
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch notifications'
    } finally {
      loading.value = false
    }
  }

  async function fetchUnreadCount() {
    try {
      const response = await api.get('/notifications/count')
      unreadCount.value = response.data.unread_count
    } catch (err) {
      console.error('Failed to fetch unread count:', err)
    }
  }

  async function markAsRead(id) {
    try {
      await api.put(`/notifications/${id}/read`)
      const notification = notifications.value.find(n => n.id === id)
      if (notification && !notification.read) {
        notification.read = true
        notification.read_at = new Date().toISOString()
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to mark as read'
      return false
    }
  }

  async function markAllAsRead() {
    try {
      await api.put('/notifications/read-all')
      notifications.value.forEach(n => {
        if (!n.read) {
          n.read = true
          n.read_at = new Date().toISOString()
        }
      })
      unreadCount.value = 0
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to mark all as read'
      return false
    }
  }

  async function deleteNotification(id) {
    try {
      await api.delete(`/notifications/${id}`)
      const notification = notifications.value.find(n => n.id === id)
      if (notification && !notification.read) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      notifications.value = notifications.value.filter(n => n.id !== id)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete notification'
      return false
    }
  }

  async function deleteAll() {
    try {
      await api.delete('/notifications')
      notifications.value = []
      unreadCount.value = 0
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete all notifications'
      return false
    }
  }

  function addNotification(notification) {
    notifications.value.unshift(notification)
    if (!notification.read) {
      unreadCount.value++
    }
  }

  return {
    notifications,
    unreadCount,
    hasUnread,
    loading,
    error,
    pagination,
    fetchNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    deleteAll,
    addNotification
  }
})
