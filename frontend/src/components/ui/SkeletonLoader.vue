<script setup>
defineProps({
  type: {
    type: String,
    default: 'text',
    validator: (value) => ['text', 'circle', 'card', 'avatar', 'button', 'title'].includes(value)
  },
  width: {
    type: String,
    default: '100%'
  },
  height: {
    type: String,
    default: null
  },
  count: {
    type: Number,
    default: 1
  }
})

const getHeight = (type) => {
  const heights = {
    text: '16px',
    title: '28px',
    circle: '40px',
    card: '120px',
    avatar: '40px',
    button: '40px'
  }
  return heights[type] || '16px'
}
</script>

<template>
  <div class="space-y-2">
    <div
      v-for="i in count"
      :key="i"
      class="skeleton"
      :class="{
        'rounded-full': type === 'circle' || type === 'avatar',
        'rounded-lg': type === 'card',
        'rounded-md': type === 'button',
        'rounded': type === 'text' || type === 'title'
      }"
      :style="{
        width: type === 'circle' || type === 'avatar' ? (height || getHeight(type)) : width,
        height: height || getHeight(type)
      }"
    />
  </div>
</template>

<style scoped>
.skeleton {
  @apply bg-gray-200 dark:bg-gray-700;
  background: linear-gradient(90deg, #e5e7eb 25%, #f3f4f6 50%, #e5e7eb 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

:global(.dark) .skeleton {
  background: linear-gradient(90deg, #334155 25%, #475569 50%, #334155 75%);
  background-size: 200% 100%;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
