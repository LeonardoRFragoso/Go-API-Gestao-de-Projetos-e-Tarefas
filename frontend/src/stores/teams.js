import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useTeamsStore = defineStore('teams', () => {
  const teams = ref([])
  const currentTeam = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 10, total: 0 })

  async function fetchTeams(page = 1, limit = 10) {
    loading.value = true
    error.value = null
    try {
      const response = await api.get('/teams', { params: { page, limit } })
      teams.value = response.data.data
      pagination.value = response.data.pagination
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch teams'
    } finally {
      loading.value = false
    }
  }

  async function fetchTeam(id) {
    loading.value = true
    error.value = null
    try {
      const response = await api.get(`/teams/${id}`)
      currentTeam.value = response.data
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch team'
      return null
    } finally {
      loading.value = false
    }
  }

  async function createTeam(data) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/teams', data)
      teams.value.unshift(response.data)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create team'
      return null
    } finally {
      loading.value = false
    }
  }

  async function updateTeam(id, data) {
    loading.value = true
    error.value = null
    try {
      const response = await api.put(`/teams/${id}`, data)
      const index = teams.value.findIndex(t => t.id === id)
      if (index !== -1) {
        teams.value[index] = response.data
      }
      if (currentTeam.value?.id === id) {
        currentTeam.value = { ...currentTeam.value, ...response.data }
      }
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update team'
      return null
    } finally {
      loading.value = false
    }
  }

  async function deleteTeam(id) {
    loading.value = true
    error.value = null
    try {
      await api.delete(`/teams/${id}`)
      teams.value = teams.value.filter(t => t.id !== id)
      if (currentTeam.value?.id === id) {
        currentTeam.value = null
      }
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete team'
      return false
    } finally {
      loading.value = false
    }
  }

  async function addMember(teamId, userId, role = 'member') {
    try {
      await api.post(`/teams/${teamId}/members`, { user_id: userId, role })
      await fetchTeam(teamId)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to add member'
      return false
    }
  }

  async function removeMember(teamId, memberId) {
    try {
      await api.delete(`/teams/${teamId}/members/${memberId}`)
      await fetchTeam(teamId)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to remove member'
      return false
    }
  }

  async function updateMemberRole(teamId, memberId, role) {
    try {
      await api.put(`/teams/${teamId}/members/${memberId}`, { role })
      await fetchTeam(teamId)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update member role'
      return false
    }
  }

  async function addProject(teamId, projectId) {
    try {
      await api.post(`/teams/${teamId}/projects`, { project_id: projectId })
      await fetchTeam(teamId)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to add project'
      return false
    }
  }

  async function removeProject(teamId, projectId) {
    try {
      await api.delete(`/teams/${teamId}/projects/${projectId}`)
      await fetchTeam(teamId)
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to remove project'
      return false
    }
  }

  return {
    teams,
    currentTeam,
    loading,
    error,
    pagination,
    fetchTeams,
    fetchTeam,
    createTeam,
    updateTeam,
    deleteTeam,
    addMember,
    removeMember,
    updateMemberRole,
    addProject,
    removeProject
  }
})
