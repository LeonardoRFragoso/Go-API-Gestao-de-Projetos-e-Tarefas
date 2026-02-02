<script setup>
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend)

const { t } = useI18n()
const settingsStore = useSettingsStore()

const props = defineProps({
  data: {
    type: Array,
    default: () => [3, 5, 2, 8, 4, 6, 1]
  }
})

const days = computed(() => {
  const locale = settingsStore.locale === 'pt-BR' ? 'pt-BR' : 'en-US'
  const today = new Date()
  const labels = []
  for (let i = 6; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    labels.push(date.toLocaleDateString(locale, { weekday: 'short' }))
  }
  return labels
})

const chartData = computed(() => ({
  labels: days.value,
  datasets: [{
    label: t('tasks.completedTasks') || 'Completed',
    data: props.data,
    backgroundColor: '#3b82f6',
    borderRadius: 6,
    borderSkipped: false,
    maxBarThickness: 32
  }]
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
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
      cornerRadius: 8
    }
  },
  scales: {
    x: {
      grid: {
        display: false
      },
      ticks: {
        color: settingsStore.darkMode ? '#94a3b8' : '#6b7280'
      }
    },
    y: {
      beginAtZero: true,
      grid: {
        color: settingsStore.darkMode ? '#334155' : '#e5e7eb'
      },
      ticks: {
        color: settingsStore.darkMode ? '#94a3b8' : '#6b7280',
        stepSize: 1
      }
    }
  }
}))
</script>

<template>
  <div class="w-full h-48">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
