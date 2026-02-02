import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useStatsStore = defineStore('stats', () => {
  const taskStats = ref({ todo: 0, in_progress: 0, done: 0 })
  const weeklyData = ref([0, 0, 0, 0, 0, 0, 0])
  const projectsCount = ref(0)
  const loading = ref(false)
  const error = ref(null)

  async function fetchDashboardStats() {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/stats/dashboard')
      taskStats.value = response.data.task_stats
      weeklyData.value = response.data.weekly_data
      projectsCount.value = response.data.projects_count
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch stats'
      console.error('Failed to fetch dashboard stats:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchTaskStats() {
    try {
      const response = await api.get('/stats/tasks')
      taskStats.value = response.data
    } catch (err) {
      console.error('Failed to fetch task stats:', err)
    }
  }

  async function fetchWeeklyProgress() {
    try {
      const response = await api.get('/stats/weekly')
      weeklyData.value = response.data.data
    } catch (err) {
      console.error('Failed to fetch weekly progress:', err)
    }
  }

  return {
    taskStats,
    weeklyData,
    projectsCount,
    loading,
    error,
    fetchDashboardStats,
    fetchTaskStats,
    fetchWeeklyProgress
  }
})
