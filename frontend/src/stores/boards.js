import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useBoardsStore = defineStore('boards', () => {
  const boards = ref([])
  const currentBoard = ref(null)
  const loading = ref(false)
  const error = ref(null)

  async function fetchBoardsByProject(projectId) {
    loading.value = true
    try {
      const response = await api.get(`/projects/${projectId}/boards`)
      boards.value = response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar quadros'
    } finally {
      loading.value = false
    }
  }

  async function fetchBoard(id) {
    loading.value = true
    try {
      const response = await api.get(`/boards/${id}`)
      currentBoard.value = response.data.data
      return currentBoard.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar quadro'
      return null
    } finally {
      loading.value = false
    }
  }

  async function createBoard(data) {
    loading.value = true
    try {
      const response = await api.post('/boards', data)
      boards.value.push(response.data.data)
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao criar quadro'
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateBoard(id, data) {
    loading.value = true
    try {
      const response = await api.put(`/boards/${id}`, data)
      const index = boards.value.findIndex(b => b.id === id)
      if (index !== -1) {
        boards.value[index] = response.data.data
      }
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao atualizar quadro'
      return null
    } finally {
      loading.value = false
    }
  }

  async function deleteBoard(id) {
    loading.value = true
    try {
      await api.delete(`/boards/${id}`)
      boards.value = boards.value.filter(b => b.id !== id)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao excluir quadro'
      return false
    } finally {
      loading.value = false
    }
  }

  async function reorderLists(boardId, positions) {
    try {
      await api.put(`/boards/${boardId}/lists/reorder`, { positions })
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao reordenar listas'
    }
  }

  return {
    boards,
    currentBoard,
    loading,
    error,
    fetchBoardsByProject,
    fetchBoard,
    createBoard,
    updateBoard,
    deleteBoard,
    reorderLists
  }
})
