<script setup>
import { Languages } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'
import { ref } from 'vue'

const { locale } = useI18n()
const settingsStore = useSettingsStore()
const showMenu = ref(false)

const languages = [
  { code: 'pt-BR', name: 'Português', flag: '🇧🇷' },
  { code: 'en', name: 'English', flag: '🇺🇸' }
]

const currentLanguage = () => {
  return languages.find(l => l.code === locale.value) || languages[0]
}

const setLanguage = (code) => {
  locale.value = code
  settingsStore.setLocale(code)
  showMenu.value = false
}

const toggleLanguage = () => {
  const newLocale = locale.value === 'pt-BR' ? 'en' : 'pt-BR'
  setLanguage(newLocale)
}
</script>

<template>
  <div class="relative">
    <button
      @click="showMenu = !showMenu"
      class="flex items-center gap-2 px-3 py-2 rounded-lg text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
    >
      <span class="text-lg">{{ currentLanguage().flag }}</span>
      <span class="text-sm font-medium hidden sm:inline">{{ currentLanguage().code.toUpperCase() }}</span>
    </button>

    <Transition
      enter-active-class="transition-all duration-200"
      enter-from-class="opacity-0 scale-95 -translate-y-2"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition-all duration-200"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 -translate-y-2"
    >
      <div
        v-if="showMenu"
        class="absolute right-0 mt-2 w-40 bg-white dark:bg-dark-800 rounded-lg shadow-lg border border-gray-100 dark:border-dark-700 py-1 z-50"
      >
        <button
          v-for="lang in languages"
          :key="lang.code"
          @click="setLanguage(lang.code)"
          :class="[
            'w-full flex items-center gap-3 px-4 py-2 text-sm hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors',
            locale === lang.code ? 'text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/20' : 'text-gray-700 dark:text-gray-300'
          ]"
        >
          <span class="text-lg">{{ lang.flag }}</span>
          <span>{{ lang.name }}</span>
        </button>
      </div>
    </Transition>

    <div
      v-if="showMenu"
      class="fixed inset-0 z-40"
      @click="showMenu = false"
    />
  </div>
</template>
