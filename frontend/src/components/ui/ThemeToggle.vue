<script setup>
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import { useSettingsStore } from '@/stores/settings'
import { ref, onMounted } from 'vue'

const settingsStore = useSettingsStore()
const showMenu = ref(false)

const themes = [
  { id: 'light', icon: Sun, label: 'Light' },
  { id: 'dark', icon: Moon, label: 'Dark' },
]

const setTheme = (theme) => {
  settingsStore.setDarkMode(theme === 'dark')
  showMenu.value = false
}

const toggleTheme = () => {
  settingsStore.toggleDarkMode()
}

onMounted(() => {
  settingsStore.initializeSettings()
})
</script>

<template>
  <button
    @click="toggleTheme"
    class="p-2 rounded-lg text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-dark-700 transition-all duration-200"
    :title="settingsStore.darkMode ? 'Modo Claro' : 'Modo Escuro'"
  >
    <Transition
      mode="out-in"
      enter-active-class="transition-all duration-200"
      enter-from-class="opacity-0 rotate-90 scale-0"
      enter-to-class="opacity-100 rotate-0 scale-100"
      leave-active-class="transition-all duration-200"
      leave-from-class="opacity-100 rotate-0 scale-100"
      leave-to-class="opacity-0 -rotate-90 scale-0"
    >
      <Moon v-if="!settingsStore.darkMode" :size="20" />
      <Sun v-else :size="20" />
    </Transition>
  </button>
</template>
