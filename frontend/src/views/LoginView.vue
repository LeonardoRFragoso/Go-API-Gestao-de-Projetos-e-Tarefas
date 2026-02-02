<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toastification'
import { FolderKanban, Mail, Lock } from 'lucide-vue-next'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'
import LanguageToggle from '@/components/ui/LanguageToggle.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()

const email = ref('')
const password = ref('')

const handleSubmit = async () => {
  const success = await authStore.login(email.value, password.value)
  if (success) {
    toast.success(t('toast.loginSuccess'))
    router.push('/')
  }
}

onMounted(() => {
  settingsStore.initializeSettings()
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-dark-950 py-12 px-4 sm:px-6 lg:px-8 transition-colors">
    <div class="absolute top-4 right-4 flex items-center gap-2">
      <ThemeToggle />
      <LanguageToggle />
    </div>
    
    <div class="max-w-md w-full animate-fade-in">
      <div class="text-center mb-8">
        <div class="mx-auto w-16 h-16 bg-primary-600 rounded-2xl flex items-center justify-center shadow-lg">
          <FolderKanban :size="32" class="text-white" />
        </div>
        <h2 class="mt-4 text-3xl font-bold text-gray-900 dark:text-white">TaskManager</h2>
        <p class="mt-2 text-gray-600 dark:text-gray-400">{{ t('auth.loginSubtitle') }}</p>
      </div>
      
      <div class="card p-8">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div v-if="authStore.error" class="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 p-3 rounded-lg text-sm">
            {{ authStore.error }}
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('auth.email') }}</label>
            <div class="relative">
              <Mail :size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input id="email" v-model="email" type="email" required class="input pl-10" placeholder="seu@email.com" />
            </div>
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('auth.password') }}</label>
            <div class="relative">
              <Lock :size="18" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input id="password" v-model="password" type="password" required class="input pl-10" placeholder="••••••••" />
            </div>
          </div>
          <button type="submit" :disabled="authStore.loading" class="w-full btn-primary py-3">
            <span v-if="authStore.loading">{{ t('common.loading') }}</span>
            <span v-else>{{ t('auth.login') }}</span>
          </button>
        </form>
        <p class="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
          {{ t('auth.noAccount') }}
          <RouterLink to="/register" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium">
            {{ t('auth.register') }}
          </RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>
