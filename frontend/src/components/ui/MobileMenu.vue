<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { 
  Menu, X, Home, FolderKanban, User, LogOut, 
  Sun, Moon, Languages, Settings 
} from 'lucide-vue-next'
import UserAvatar from './UserAvatar.vue'
import ThemeToggle from './ThemeToggle.vue'
import LanguageToggle from './LanguageToggle.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()

const isOpen = ref(false)

const menuItems = [
  { path: '/', icon: Home, label: 'nav.dashboard' },
  { path: '/projects', icon: FolderKanban, label: 'nav.projects' },
  { path: '/profile', icon: User, label: 'common.profile' },
]

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const closeMenu = () => {
  isOpen.value = false
}

const navigateTo = (path) => {
  router.push(path)
  closeMenu()
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
  closeMenu()
}

watch(() => route.path, () => {
  closeMenu()
})
</script>

<template>
  <div class="lg:hidden">
    <button
      @click="toggleMenu"
      class="p-2 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 rounded-lg transition-colors"
    >
      <Menu :size="24" />
    </button>

    <Teleport to="body">
      <Transition
        enter-active-class="transition-opacity duration-300"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-opacity duration-300"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="isOpen"
          class="fixed inset-0 z-50 bg-black/50 lg:hidden"
          @click="closeMenu"
        />
      </Transition>

      <Transition
        enter-active-class="transition-transform duration-300"
        enter-from-class="-translate-x-full"
        enter-to-class="translate-x-0"
        leave-active-class="transition-transform duration-300"
        leave-from-class="translate-x-0"
        leave-to-class="-translate-x-full"
      >
        <div
          v-if="isOpen"
          class="fixed inset-y-0 left-0 z-50 w-72 bg-white dark:bg-dark-800 shadow-2xl lg:hidden"
        >
          <div class="flex flex-col h-full">
            <div class="flex items-center justify-between p-4 border-b border-gray-100 dark:border-dark-700">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                  <FolderKanban :size="18" class="text-white" />
                </div>
                <span class="font-bold text-gray-900 dark:text-white">TaskManager</span>
              </div>
              <button
                @click="closeMenu"
                class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-700"
              >
                <X :size="20" />
              </button>
            </div>

            <div class="p-4 border-b border-gray-100 dark:border-dark-700">
              <div class="flex items-center gap-3">
                <UserAvatar :user="authStore.user" size="md" />
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ authStore.user?.name }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{{ authStore.user?.email }}</p>
                </div>
              </div>
            </div>

            <nav class="flex-1 p-4 space-y-1 overflow-y-auto">
              <button
                v-for="item in menuItems"
                :key="item.path"
                @click="navigateTo(item.path)"
                :class="[
                  'w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                  route.path === item.path
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                    : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-dark-700'
                ]"
              >
                <component :is="item.icon" :size="20" />
                <span class="font-medium">{{ t(item.label) }}</span>
              </button>
            </nav>

            <div class="p-4 border-t border-gray-100 dark:border-dark-700 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('settings.appearance') }}</span>
                <ThemeToggle />
              </div>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-600 dark:text-gray-400">{{ t('settings.language') }}</span>
                <LanguageToggle />
              </div>
            </div>

            <div class="p-4 border-t border-gray-100 dark:border-dark-700">
              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-4 py-3 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
              >
                <LogOut :size="20" />
                <span class="font-medium">{{ t('common.logout') }}</span>
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>
