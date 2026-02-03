<script setup>
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toastification'
import { 
  FolderKanban, ChevronDown, User, LogOut, Settings,
  Home, Folder, Users
} from 'lucide-vue-next'
import NotificationDropdown from '@/components/NotificationDropdown.vue'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'
import LanguageToggle from '@/components/ui/LanguageToggle.vue'
import MobileMenu from '@/components/ui/MobileMenu.vue'
import BottomNavigation from '@/components/ui/BottomNavigation.vue'
import OnboardingModal from '@/components/ui/OnboardingModal.vue'
import KeyboardShortcuts from '@/components/ui/KeyboardShortcuts.vue'

const { t } = useI18n()
const toast = useToast()
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const settingsStore = useSettingsStore()

const showUserMenu = ref(false)
const keyboardShortcutsRef = ref(null)

const handleLogout = async () => {
  await authStore.logout()
  toast.success(t('toast.logoutSuccess'))
  router.push('/login')
}

const isActiveRoute = (path) => {
  return route.path === path
}

const handleNewTask = () => {
  toast.info('Use o atalho N dentro de um quadro para criar uma tarefa')
}

const handleSearch = () => {
  toast.info('Busca global em desenvolvimento')
}

onMounted(() => {
  settingsStore.initializeSettings()
})
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-dark-950 transition-colors duration-200">
    <nav class="bg-white dark:bg-dark-800 shadow-sm border-b border-gray-200 dark:border-dark-700 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <MobileMenu />
            
            <RouterLink to="/" class="flex items-center">
              <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                <FolderKanban :size="18" class="text-white" />
              </div>
              <span class="ml-2 text-xl font-bold text-gray-900 dark:text-white hidden sm:block">TaskManager</span>
            </RouterLink>
            
            <div class="hidden lg:ml-8 lg:flex lg:space-x-1">
              <RouterLink 
                to="/" 
                :class="[
                  'nav-link flex items-center gap-2',
                  isActiveRoute('/') ? 'nav-link-active' : ''
                ]"
              >
                <Home :size="18" />
                {{ t('nav.dashboard') }}
              </RouterLink>
              <RouterLink 
                to="/projects" 
                :class="[
                  'nav-link flex items-center gap-2',
                  isActiveRoute('/projects') || route.path.startsWith('/projects') ? 'nav-link-active' : ''
                ]"
              >
                <Folder :size="18" />
                {{ t('nav.projects') }}
              </RouterLink>
              <RouterLink 
                to="/teams" 
                :class="[
                  'nav-link flex items-center gap-2',
                  isActiveRoute('/teams') || route.path.startsWith('/teams') ? 'nav-link-active' : ''
                ]"
              >
                <Users :size="18" />
                {{ t('nav.teams') }}
              </RouterLink>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <NotificationDropdown />
            <ThemeToggle />
            <LanguageToggle />
            
            <div class="relative hidden sm:block">
              <button 
                @click="showUserMenu = !showUserMenu" 
                class="flex items-center gap-2 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
              >
                <UserAvatar :user="authStore.user" size="sm" />
                <span class="hidden md:block text-sm font-medium text-gray-700 dark:text-gray-300 max-w-[120px] truncate">
                  {{ authStore.user?.name }}
                </span>
                <ChevronDown :size="16" class="text-gray-500 dark:text-gray-400" />
              </button>
              
              <Transition
                enter-active-class="transition-all duration-200"
                enter-from-class="opacity-0 scale-95 -translate-y-2"
                enter-to-class="opacity-100 scale-100 translate-y-0"
                leave-active-class="transition-all duration-200"
                leave-from-class="opacity-100 scale-100 translate-y-0"
                leave-to-class="opacity-0 scale-95 -translate-y-2"
              >
                <div v-if="showUserMenu" class="dropdown-menu">
                  <RouterLink 
                    to="/profile" 
                    @click="showUserMenu = false" 
                    class="dropdown-item flex items-center gap-2"
                  >
                    <User :size="16" />
                    {{ t('common.profile') }}
                  </RouterLink>
                  <hr class="my-1 border-gray-100 dark:border-dark-700">
                  <button 
                    @click="handleLogout" 
                    class="dropdown-item flex items-center gap-2 text-red-600 dark:text-red-400 w-full"
                  >
                    <LogOut :size="16" />
                    {{ t('common.logout') }}
                  </button>
                </div>
              </Transition>
              
              <div v-if="showUserMenu" class="fixed inset-0 z-40" @click="showUserMenu = false" />
            </div>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 pb-24 lg:pb-8">
      <RouterView v-slot="{ Component }">
        <Transition
          mode="out-in"
          enter-active-class="transition-all duration-200"
          enter-from-class="opacity-0 translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-200"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 -translate-y-2"
        >
          <component :is="Component" />
        </Transition>
      </RouterView>
    </main>

    <BottomNavigation @newTask="handleNewTask" />
    <OnboardingModal />
    <KeyboardShortcuts ref="keyboardShortcutsRef" @newTask="handleNewTask" @search="handleSearch" />
  </div>
</template>
