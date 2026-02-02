<script setup>
import { Search, X } from 'lucide-vue-next'
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  autofocus: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'search', 'clear'])
const inputRef = ref(null)

const clearSearch = () => {
  emit('update:modelValue', '')
  emit('clear')
  inputRef.value?.focus()
}

const handleKeydown = (e) => {
  if (e.key === 'Escape') {
    clearSearch()
  }
}

defineExpose({ focus: () => inputRef.value?.focus() })
</script>

<template>
  <div class="relative">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <Search class="w-5 h-5 text-gray-400" />
    </div>
    <input
      ref="inputRef"
      type="text"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
      @keydown="handleKeydown"
      @keydown.enter="$emit('search', modelValue)"
      :placeholder="placeholder || t('common.search')"
      :autofocus="autofocus"
      class="w-full pl-10 pr-10 py-2 bg-gray-50 dark:bg-dark-700 border border-gray-200 dark:border-dark-600 rounded-lg text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-all"
    />
    <button
      v-if="modelValue"
      @click="clearSearch"
      class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
    >
      <X class="w-5 h-5" />
    </button>
  </div>
</template>
