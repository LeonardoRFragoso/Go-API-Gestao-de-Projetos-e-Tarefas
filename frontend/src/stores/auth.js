import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const isAuthenticated = computed(() => !!user.value)

  function checkAuth() {
    const storedUser = localStorage.getItem('user')
    const token = localStorage.getItem('access_token')
    if (storedUser && token) {
      user.value = JSON.parse(storedUser)
    }
  }

  async function login(email, password) {
    loading.value = true
    error.value = null
    try {
      const response = await api.post('/auth/login', { email, password })
      const { user: userData, tokens } = response.data.data
      
      user.value = userData
      localStorage.setItem('user', JSON.stringify(userData))
      localStorage.setItem('access_token', tokens.access_token)
      localStorage.setItem('refresh_token', tokens.refresh_token)
      
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao fazer login'
      return false
    } finally {
      loading.value = false
    }
  }

  async function register(name, email, password) {
    loading.value = true
    error.value = null
    try {
      await api.post('/auth/register', { name, email, password })
      return await login(email, password)
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao criar conta'
      return false
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      user.value = null
      localStorage.removeItem('user')
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }

  async function updateProfile(data) {
    loading.value = true
    try {
      const response = await api.put('/users/me', data)
      user.value = response.data.data
      localStorage.setItem('user', JSON.stringify(user.value))
      return true
    } catch (err) {
      error.value = err.response?.data?.error || 'Erro ao atualizar perfil'
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    loading,
    error,
    isAuthenticated,
    checkAuth,
    login,
    register,
    logout,
    updateProfile
  }
})
