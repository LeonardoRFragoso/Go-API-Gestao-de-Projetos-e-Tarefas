<script setup>
import { X } from 'lucide-vue-next'
import { onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl', 'full'].includes(value)
  },
  closable: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close'])

const sizeClasses = {
  sm: 'max-w-sm',
  md: 'max-w-md',
  lg: 'max-w-lg',
  xl: 'max-w-xl',
  full: 'max-w-4xl'
}

const close = () => {
  if (props.closable) {
    emit('close')
  }
}

const handleKeydown = (e) => {
  if (e.key === 'Escape' && props.show && props.closable) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
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
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="close"
      >
        <div class="fixed inset-0 bg-black/50 dark:bg-black/70" />
        <Transition
          enter-active-class="transition-all duration-200"
          enter-from-class="opacity-0 scale-95"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition-all duration-200"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-95"
        >
          <div
            v-if="show"
            :class="[
              sizeClasses[size],
              'relative w-full bg-white dark:bg-dark-800 rounded-xl shadow-2xl animate-scale-in max-h-[90vh] overflow-hidden flex flex-col'
            ]"
          >
            <div v-if="title || closable" class="flex items-center justify-between p-4 border-b border-gray-100 dark:border-dark-700">
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ title }}</h2>
              <button
                v-if="closable"
                @click="close"
                class="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors"
              >
                <X :size="20" />
              </button>
            </div>
            <div class="flex-1 overflow-y-auto">
              <slot />
            </div>
          </div>
        </Transition>
      </div>
    </Transition>
  </Teleport>
</template>
