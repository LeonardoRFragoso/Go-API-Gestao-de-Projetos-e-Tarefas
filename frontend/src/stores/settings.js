import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

export const useSettingsStore = defineStore('settings', () => {
  const darkMode = ref(localStorage.getItem('darkMode') === 'true')
  const locale = ref(localStorage.getItem('locale') || 'pt-BR')
  const sidebarCollapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true')
  const onboardingCompleted = ref(localStorage.getItem('onboardingCompleted') === 'true')

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
    applyDarkMode()
  }

  function setDarkMode(value) {
    darkMode.value = value
    applyDarkMode()
  }

  function applyDarkMode() {
    localStorage.setItem('darkMode', darkMode.value)
    if (darkMode.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  function setLocale(newLocale) {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
    localStorage.setItem('sidebarCollapsed', sidebarCollapsed.value)
  }

  function completeOnboarding() {
    onboardingCompleted.value = true
    localStorage.setItem('onboardingCompleted', 'true')
  }

  function resetOnboarding() {
    onboardingCompleted.value = false
    localStorage.removeItem('onboardingCompleted')
  }

  function initializeSettings() {
    applyDarkMode()
  }

  return {
    darkMode,
    locale,
    sidebarCollapsed,
    onboardingCompleted,
    toggleDarkMode,
    setDarkMode,
    setLocale,
    toggleSidebar,
    completeOnboarding,
    resetOnboarding,
    initializeSettings
  }
})
