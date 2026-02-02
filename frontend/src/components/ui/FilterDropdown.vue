<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Filter, ChevronDown, X } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  },
  priorities: {
    type: Boolean,
    default: true
  },
  assignees: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)

const priorityOptions = [
  { value: '', label: 'common.all' },
  { value: 'low', label: 'tasks.priorities.low' },
  { value: 'medium', label: 'tasks.priorities.medium' },
  { value: 'high', label: 'tasks.priorities.high' },
  { value: 'urgent', label: 'tasks.priorities.urgent' }
]

const updateFilter = (key, value) => {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}

const clearFilters = () => {
  emit('update:modelValue', {})
  isOpen.value = false
}

const hasActiveFilters = () => {
  return Object.values(props.modelValue).some(v => v)
}
</script>

<template>
  <div class="relative">
    <button
      @click="isOpen = !isOpen"
      :class="[
        'flex items-center gap-2 px-3 py-2 rounded-lg border transition-colors',
        hasActiveFilters()
          ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
          : 'border-gray-200 dark:border-dark-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-dark-700'
      ]"
    >
      <Filter :size="18" />
      <span class="text-sm font-medium">{{ t('common.filter') }}</span>
      <ChevronDown :size="16" :class="['transition-transform', isOpen ? 'rotate-180' : '']" />
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
        v-if="isOpen"
        class="absolute right-0 mt-2 w-64 bg-white dark:bg-dark-800 rounded-xl shadow-lg border border-gray-100 dark:border-dark-700 p-4 z-50"
      >
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-medium text-gray-900 dark:text-white">{{ t('common.filter') }}</h4>
          <button
            v-if="hasActiveFilters()"
            @click="clearFilters"
            class="text-xs text-primary-600 dark:text-primary-400 hover:underline"
          >
            {{ t('common.clear') || 'Limpar' }}
          </button>
        </div>

        <div v-if="priorities" class="mb-4">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {{ t('tasks.priority') }}
          </label>
          <select
            :value="modelValue.priority || ''"
            @change="updateFilter('priority', $event.target.value)"
            class="w-full px-3 py-2 bg-gray-50 dark:bg-dark-700 border border-gray-200 dark:border-dark-600 rounded-lg text-sm text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option v-for="option in priorityOptions" :key="option.value" :value="option.value">
              {{ t(option.label) }}
            </option>
          </select>
        </div>

        <div v-if="assignees.length > 0">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            {{ t('tasks.assignees') }}
          </label>
          <select
            :value="modelValue.assignee || ''"
            @change="updateFilter('assignee', $event.target.value)"
            class="w-full px-3 py-2 bg-gray-50 dark:bg-dark-700 border border-gray-200 dark:border-dark-600 rounded-lg text-sm text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="">{{ t('common.all') }}</option>
            <option v-for="user in assignees" :key="user.id" :value="user.id">
              {{ user.name }}
            </option>
          </select>
        </div>
      </div>
    </Transition>

    <div v-if="isOpen" class="fixed inset-0 z-40" @click="isOpen = false" />
  </div>
</template>
