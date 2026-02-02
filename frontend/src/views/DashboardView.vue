<script setup>
import { onMounted, computed } from 'vue'
import { useProjectsStore } from '@/stores/projects'
import { useAuthStore } from '@/stores/auth'
import { useStatsStore } from '@/stores/stats'
import { useI18n } from 'vue-i18n'
import { Folder, CheckCircle, Clock, ChevronRight, Plus, TrendingUp } from 'lucide-vue-next'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'
import TaskStatsChart from '@/components/charts/TaskStatsChart.vue'
import WeeklyProgressChart from '@/components/charts/WeeklyProgressChart.vue'

const { t } = useI18n()
const projectsStore = useProjectsStore()
const authStore = useAuthStore()
const statsStore = useStatsStore()

onMounted(async () => {
  await Promise.all([
    projectsStore.fetchProjects(1, 5),
    statsStore.fetchDashboardStats()
  ])
})

const recentProjects = computed(() => projectsStore.projects.slice(0, 5))

const stats = computed(() => [
  {
    label: t('dashboard.projects'),
    value: statsStore.projectsCount || projectsStore.pagination.total || 0,
    icon: Folder,
    color: 'primary',
    bgClass: 'bg-primary-100 dark:bg-primary-900/30',
    textClass: 'text-primary-600 dark:text-primary-400'
  },
  {
    label: t('dashboard.completedTasks'),
    value: statsStore.taskStats.done || 0,
    icon: CheckCircle,
    color: 'green',
    bgClass: 'bg-green-100 dark:bg-green-900/30',
    textClass: 'text-green-600 dark:text-green-400'
  },
  {
    label: t('dashboard.inProgress'),
    value: statsStore.taskStats.in_progress || 0,
    icon: Clock,
    color: 'yellow',
    bgClass: 'bg-yellow-100 dark:bg-yellow-900/30',
    textClass: 'text-yellow-600 dark:text-yellow-400'
  }
])
</script>

<template>
  <div class="space-y-6">
    <div class="animate-fade-in">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
        {{ t('dashboard.welcome', { name: authStore.user?.name }) }}
      </h1>
      <p class="text-gray-600 dark:text-gray-400 mt-1">{{ t('dashboard.subtitle') }}</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6">
      <div
        v-for="(stat, index) in stats"
        :key="stat.label"
        class="stat-card animate-slide-up"
        :style="{ animationDelay: `${index * 100}ms` }"
      >
        <div :class="[stat.bgClass, 'p-3 rounded-xl']">
          <component :is="stat.icon" :size="24" :class="stat.textClass" />
        </div>
        <div>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ stat.label }}</p>
          <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stat.value }}</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="card p-6 animate-slide-up" style="animation-delay: 300ms">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.taskStats') }}</h2>
          <TrendingUp :size="20" class="text-gray-400" />
        </div>
        <TaskStatsChart
          :todo="statsStore.taskStats.todo || 0"
          :in-progress="statsStore.taskStats.in_progress || 0"
          :done="statsStore.taskStats.done || 0"
        />
      </div>

      <div class="card p-6 animate-slide-up" style="animation-delay: 400ms">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.weeklyProgress') }}</h2>
        </div>
        <WeeklyProgressChart :data="statsStore.weeklyData" />
      </div>
    </div>

    <div class="card animate-slide-up" style="animation-delay: 500ms">
      <div class="p-4 md:p-6 border-b border-gray-100 dark:border-dark-700">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('dashboard.recentProjects') }}</h2>
          <RouterLink 
            to="/projects" 
            class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium flex items-center gap-1"
          >
            {{ t('common.seeAll') }}
            <ChevronRight :size="16" />
          </RouterLink>
        </div>
      </div>

      <div v-if="projectsStore.loading" class="p-6 space-y-4">
        <div v-for="i in 3" :key="i" class="flex items-center gap-4">
          <SkeletonLoader type="circle" height="12px" width="12px" />
          <div class="flex-1">
            <SkeletonLoader type="text" width="60%" />
            <SkeletonLoader type="text" width="30%" height="12px" />
          </div>
        </div>
      </div>

      <div v-else-if="recentProjects.length === 0" class="p-8 text-center">
        <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 dark:bg-dark-700 rounded-full flex items-center justify-center">
          <Folder :size="32" class="text-gray-400" />
        </div>
        <p class="text-gray-500 dark:text-gray-400 mb-4">{{ t('dashboard.noProjects') }}</p>
        <RouterLink to="/projects" class="btn-primary inline-flex items-center gap-2">
          <Plus :size="18" />
          {{ t('projects.newProject') }}
        </RouterLink>
      </div>

      <ul v-else class="divide-y divide-gray-100 dark:divide-dark-700">
        <li 
          v-for="(project, index) in recentProjects" 
          :key="project.id" 
          class="group"
        >
          <RouterLink 
            :to="`/projects/${project.id}`" 
            class="flex items-center p-4 hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors"
          >
            <div 
              class="w-3 h-3 rounded-full mr-3 ring-2 ring-white dark:ring-dark-800" 
              :style="{ backgroundColor: project.color || '#3b82f6' }"
            />
            <div class="flex-1 min-w-0">
              <p class="font-medium text-gray-900 dark:text-white truncate">{{ project.name }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">
                {{ project.board_count || 0 }} {{ t('projects.boards', project.board_count || 0) }}
              </p>
            </div>
            <ChevronRight 
              :size="20" 
              class="text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 group-hover:translate-x-1 transition-all" 
            />
          </RouterLink>
        </li>
      </ul>
    </div>
  </div>
</template>
