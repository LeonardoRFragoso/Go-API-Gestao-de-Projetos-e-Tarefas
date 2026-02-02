<script setup>
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'

ChartJS.register(ArcElement, Tooltip, Legend)

const { t } = useI18n()
const settingsStore = useSettingsStore()

const props = defineProps({
  todo: { type: Number, default: 0 },
  inProgress: { type: Number, default: 0 },
  done: { type: Number, default: 0 }
})

const total = computed(() => props.todo + props.inProgress + props.done)

const chartData = computed(() => ({
  labels: [
    t('lists.todo'),
    t('lists.inProgress'),
    t('lists.done')
  ],
  datasets: [{
    data: [props.todo, props.inProgress, props.done],
    backgroundColor: [
      '#94a3b8',
      '#3b82f6',
      '#22c55e'
    ],
    borderWidth: 0,
    hoverOffset: 4
  }]
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  cutout: '70%',
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      backgroundColor: settingsStore.darkMode ? '#1e293b' : '#ffffff',
      titleColor: settingsStore.darkMode ? '#f1f5f9' : '#1f2937',
      bodyColor: settingsStore.darkMode ? '#cbd5e1' : '#4b5563',
      borderColor: settingsStore.darkMode ? '#334155' : '#e5e7eb',
      borderWidth: 1,
      padding: 12,
      cornerRadius: 8,
      displayColors: true,
      boxPadding: 4
    }
  }
}))
</script>

<template>
  <div class="relative">
    <div class="w-full h-48">
      <Doughnut v-if="total > 0" :data="chartData" :options="chartOptions" />
      <div v-else class="w-full h-full flex items-center justify-center text-gray-400 dark:text-gray-500">
        {{ t('tasks.noTasks') }}
      </div>
    </div>
    <div v-if="total > 0" class="absolute inset-0 flex items-center justify-center pointer-events-none">
      <div class="text-center">
        <div class="text-3xl font-bold text-gray-900 dark:text-white">{{ total }}</div>
        <div class="text-xs text-gray-500 dark:text-gray-400">{{ t('tasks.title') }}</div>
      </div>
    </div>
    <div v-if="total > 0" class="flex justify-center gap-4 mt-4">
      <div class="flex items-center gap-2">
        <div class="w-3 h-3 rounded-full bg-gray-400"></div>
        <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('lists.todo') }} ({{ todo }})</span>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-3 h-3 rounded-full bg-blue-500"></div>
        <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('lists.inProgress') }} ({{ inProgress }})</span>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-3 h-3 rounded-full bg-green-500"></div>
        <span class="text-xs text-gray-600 dark:text-gray-400">{{ t('lists.done') }} ({{ done }})</span>
      </div>
    </div>
  </div>
</template>
