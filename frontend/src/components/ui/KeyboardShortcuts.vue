<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { Keyboard, X } from 'lucide-vue-next'

const { t } = useI18n()
const router = useRouter()

const showModal = ref(false)
const emit = defineEmits(['newTask', 'search'])

const shortcuts = [
  { keys: ['N'], action: 'shortcuts.newTask' },
  { keys: ['/', 'Ctrl', 'K'], action: 'shortcuts.search' },
  { keys: ['Ctrl', 'S'], action: 'shortcuts.save' },
  { keys: ['Esc'], action: 'shortcuts.escape' },
  { keys: ['G', 'D'], action: 'nav.dashboard' },
  { keys: ['G', 'P'], action: 'nav.projects' },
]

const handleKeydown = (e) => {
  if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
    return
  }

  if (e.key === '?' && e.shiftKey) {
    e.preventDefault()
    showModal.value = !showModal.value
    return
  }

  if (e.key === 'n' && !e.ctrlKey && !e.metaKey) {
    e.preventDefault()
    emit('newTask')
    return
  }

  if (e.key === '/' || (e.key === 'k' && (e.ctrlKey || e.metaKey))) {
    e.preventDefault()
    emit('search')
    return
  }

  if (e.key === 'Escape') {
    showModal.value = false
    return
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})

defineExpose({ showModal })
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="showModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="showModal = false"
      >
        <div class="fixed inset-0 bg-black/50 dark:bg-black/70" />
        <div class="relative w-full max-w-md bg-white dark:bg-dark-800 rounded-xl shadow-2xl">
          <div class="flex items-center justify-between p-4 border-b border-gray-100 dark:border-dark-700">
            <div class="flex items-center gap-2">
              <Keyboard :size="20" class="text-primary-500" />
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
                {{ t('shortcuts.title') }}
              </h2>
            </div>
            <button
              @click="showModal = false"
              class="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-700"
            >
              <X :size="18" />
            </button>
          </div>
          <div class="p-4 space-y-3">
            <div
              v-for="shortcut in shortcuts"
              :key="shortcut.action"
              class="flex items-center justify-between py-2"
            >
              <span class="text-gray-600 dark:text-gray-400">{{ t(shortcut.action) }}</span>
              <div class="flex items-center gap-1">
                <kbd
                  v-for="(key, i) in shortcut.keys"
                  :key="i"
                  class="px-2 py-1 text-xs font-medium text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-dark-700 rounded border border-gray-200 dark:border-dark-600"
                >
                  {{ key }}
                </kbd>
              </div>
            </div>
          </div>
          <div class="p-4 bg-gray-50 dark:bg-dark-900 rounded-b-xl">
            <p class="text-xs text-center text-gray-500 dark:text-gray-400">
              Press <kbd class="px-1.5 py-0.5 text-xs bg-gray-100 dark:bg-dark-700 rounded">?</kbd> to toggle this menu
            </p>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
