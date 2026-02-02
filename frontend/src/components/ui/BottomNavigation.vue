<script setup>
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { Home, FolderKanban, PlusCircle, User } from 'lucide-vue-next'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const emit = defineEmits(['newTask'])

const navItems = [
  { path: '/', icon: Home, label: 'nav.dashboard' },
  { path: '/projects', icon: FolderKanban, label: 'nav.projects' },
  { action: 'newTask', icon: PlusCircle, label: 'tasks.newTask', primary: true },
  { path: '/profile', icon: User, label: 'common.profile' },
]

const handleClick = (item) => {
  if (item.action) {
    emit(item.action)
  } else {
    router.push(item.path)
  }
}

const isActive = (item) => {
  if (item.action) return false
  return route.path === item.path
}
</script>

<template>
  <nav class="fixed bottom-0 left-0 right-0 z-40 bg-white dark:bg-dark-800 border-t border-gray-200 dark:border-dark-700 lg:hidden safe-area-bottom">
    <div class="flex items-center justify-around h-16">
      <button
        v-for="item in navItems"
        :key="item.path || item.action"
        @click="handleClick(item)"
        :class="[
          'flex flex-col items-center justify-center w-full h-full transition-colors',
          item.primary 
            ? 'text-white' 
            : isActive(item) 
              ? 'text-primary-600 dark:text-primary-400' 
              : 'text-gray-500 dark:text-gray-400'
        ]"
      >
        <div
          v-if="item.primary"
          class="w-12 h-12 -mt-6 bg-primary-600 rounded-full flex items-center justify-center shadow-lg hover:bg-primary-700 transition-colors"
        >
          <component :is="item.icon" :size="24" />
        </div>
        <template v-else>
          <component :is="item.icon" :size="22" />
          <span class="text-xs mt-1">{{ t(item.label) }}</span>
        </template>
      </button>
    </div>
  </nav>
</template>

<style scoped>
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
