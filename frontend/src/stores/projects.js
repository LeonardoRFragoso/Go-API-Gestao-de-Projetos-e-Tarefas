import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useProjectsStore = defineStore('projects', () => {
  const projects = ref([])
  const currentProject = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 10, total: 0, totalPages: 0 })

  async function fetchProjects(page = 1, limit = 10) {
    loading.value = true
    try {
      const response = await api.get('/projects', { params: { page, limit } })
      projects.value = response.data.data
      pagination.value = response.data.pagination
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar projetos'
    } finally {
      loading.value = false
    }
  }

  async function fetchProject(id) {
    loading.value = true
    try {
      const response = await api.get(`/projects/${id}`)
      currentProject.value = response.data.data
      return currentProject.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar projeto'
      return null
    } finally {
      loading.value = false
    }
  }

  async function createProject(data) {
    loading.value = true
    try {
      const response = await api.post('/projects', data)
      projects.value.unshift(response.data.data)
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao criar projeto'
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateProject(id, data) {
    loading.value = true
    try {
      const response = await api.put(`/projects/${id}`, data)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = response.data.data
      }
      if (currentProject.value?.id === id) {
        currentProject.value = response.data.data
      }
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao atualizar projeto'
      return null
    } finally {
      loading.value = false
    }
  }

  async function deleteProject(id) {
    loading.value = true
    try {
      await api.delete(`/projects/${id}`)
      projects.value = projects.value.filter(p => p.id !== id)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao excluir projeto'
      return false
    } finally {
      loading.value = false
    }
  }

  async function fetchProjectMembers(projectId) {
    loading.value = true
    try {
      const response = await api.get(`/projects/${projectId}/members`)
      return response.data.data || []
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao carregar membros'
      return []
    } finally {
      loading.value = false
    }
  }

  async function addProjectMember(projectId, memberData) {
    loading.value = true
    try {
      const response = await api.post(`/projects/${projectId}/members`, memberData)
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao adicionar membro'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function removeProjectMember(projectId, memberId) {
    loading.value = true
    try {
      await api.delete(`/projects/${projectId}/members/${memberId}`)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao remover membro'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    projects,
    currentProject,
    loading,
    error,
    pagination,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject,
    fetchProjectMembers,
    addProjectMember,
    removeProjectMember
  }
})
