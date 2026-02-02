import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useTasksStore = defineStore('tasks', () => {
  const currentTask = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchTask(id) {
    loading.value = true
    try {
      const response = await api.get(`/tasks/${id}`)
      currentTask.value = response.data.data
      return currentTask.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar tarefa'
      return null
    } finally {
      loading.value = false
    }
  }

  async function createTask(data) {
    loading.value = true
    try {
      const response = await api.post('/tasks', data)
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao criar tarefa'
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateTask(id, data) {
    loading.value = true
    try {
      const response = await api.put(`/tasks/${id}`, data)
      currentTask.value = response.data.data
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao atualizar tarefa'
      return null
    } finally {
      loading.value = false
    }
  }

  async function deleteTask(id) {
    loading.value = true
    try {
      await api.delete(`/tasks/${id}`)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao excluir tarefa'
      return false
    } finally {
      loading.value = false
    }
  }

  async function moveTask(taskId, listId, position) {
    try {
      const response = await api.put(`/tasks/${taskId}/move`, {
        list_id: listId,
        position
      })
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao mover tarefa'
      return null
    }
  }

  async function reorderTasks(listId, positions) {
    try {
      await api.put(`/lists/${listId}/tasks/reorder`, { positions })
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao reordenar tarefas'
    }
  }

  async function addComment(taskId, content) {
    try {
      const response = await api.post('/comments', {
        task_id: taskId,
        content
      })
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao adicionar comentário'
      return null
    }
  }

  async function fetchComments(taskId, page = 1, limit = 20) {
    try {
      const response = await api.get(`/tasks/${taskId}/comments`, {
        params: { page, limit }
      })
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar comentários'
      return null
    }
  }

  return {
    currentTask,
    loading,
    error,
    fetchTask,
    createTask,
    updateTask,
    deleteTask,
    moveTask,
    reorderTasks,
    addComment,
    fetchComments
  }
})
