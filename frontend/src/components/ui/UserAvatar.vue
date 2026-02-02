<script setup>
import { computed } from 'vue'
import { User } from 'lucide-vue-next'

const props = defineProps({
  user: {
    type: Object,
    default: null
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
  },
  showName: {
    type: Boolean,
    default: false
  }
})

const sizeClasses = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-12 h-12 text-lg',
  xl: 'w-16 h-16 text-xl'
}

const iconSizes = {
  xs: 12,
  sm: 14,
  md: 18,
  lg: 22,
  xl: 28
}

const initials = computed(() => {
  if (!props.user?.name) return ''
  const names = props.user.name.split(' ')
  if (names.length >= 2) {
    return (names[0][0] + names[names.length - 1][0]).toUpperCase()
  }
  return names[0].substring(0, 2).toUpperCase()
})

const bgColor = computed(() => {
  if (!props.user?.id) return 'bg-gray-400'
  const colors = [
    'bg-blue-500',
    'bg-green-500',
    'bg-yellow-500',
    'bg-red-500',
    'bg-purple-500',
    'bg-pink-500',
    'bg-indigo-500',
    'bg-teal-500'
  ]
  const index = props.user.id.charCodeAt(0) % colors.length
  return colors[index]
})
</script>

<template>
  <div class="flex items-center gap-2">
    <div
      :class="[
        sizeClasses[size],
        bgColor,
        'rounded-full flex items-center justify-center text-white font-medium ring-2 ring-white dark:ring-dark-800 transition-transform hover:scale-105'
      ]"
    >
      <template v-if="user?.avatar_url">
        <img
          :src="user.avatar_url"
          :alt="user.name"
          class="w-full h-full rounded-full object-cover"
        />
      </template>
      <template v-else-if="initials">
        {{ initials }}
      </template>
      <template v-else>
        <User :size="iconSizes[size]" />
      </template>
    </div>
    <span v-if="showName && user?.name" class="text-sm font-medium text-gray-700 dark:text-gray-300">
      {{ user.name }}
    </span>
  </div>
</template>
