<script setup>
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useProjectsStore } from '@/stores/projects'

const { t } = useI18n()
const projectsStore = useProjectsStore()
const showModal = ref(false)
const newProject = ref({ name: '', description: '', color: '#3B82F6' })

const colors = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#06B6D4', '#84CC16']

onMounted(() => {
  projectsStore.fetchProjects()
})

const handleCreate = async () => {
  const result = await projectsStore.createProject(newProject.value)
  if (result) {
    showModal.value = false
    newProject.value = { name: '', description: '', color: '#3B82F6' }
  }
}

const handleDelete = async (id) => {
  if (confirm(t('projects.deleteConfirm'))) {
    await projectsStore.deleteProject(id)
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ t('projects.title') }}</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">{{ t('projects.subtitle') }}</p>
      </div>
      <button @click="showModal = true" class="btn-primary flex items-center">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('projects.newProject') }}
      </button>
    </div>

    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">{{ t('common.loading') }}</div>
    <div v-else-if="projectsStore.projects.length === 0" class="card p-12 text-center">
      <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">{{ t('projects.noProjects') }}</h3>
      <p class="mt-2 text-gray-500">{{ t('projects.createFirst') }}</p>
      <button @click="showModal = true" class="btn-primary mt-4">{{ t('projects.newProject') }}</button>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="project in projectsStore.projects" :key="project.id" class="card hover:shadow-md transition-shadow">
        <RouterLink :to="`/projects/${project.id}`" class="block p-6">
          <div class="flex items-start justify-between">
            <div class="flex items-center">
              <div class="w-4 h-4 rounded-full mr-3" :style="{ backgroundColor: project.color }"></div>
              <h3 class="font-semibold text-gray-900 dark:text-white">{{ project.name }}</h3>
            </div>
          </div>
          <p v-if="project.description" class="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">{{ project.description }}</p>
          <div class="mt-4 flex items-center text-sm text-gray-500 dark:text-gray-400">
            <span>{{ project.board_count || 0 }} {{ t('projects.boards') }}</span>
            <span class="mx-2">•</span>
            <span>{{ project.member_count || 1 }} {{ t('projects.members') }}</span>
          </div>
        </RouterLink>
        <div class="px-6 py-3 bg-gray-50 dark:bg-dark-700 border-t border-gray-100 dark:border-dark-600 flex justify-end">
          <button @click.prevent="handleDelete(project.id)" class="text-sm text-red-600 hover:text-red-700">
            {{ t('common.delete') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100 dark:border-dark-600">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('projects.newProject') }}</h2>
        </div>
        <form @submit.prevent="handleCreate" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.projectName') }}</label>
            <input v-model="newProject.name" type="text" required class="input" :placeholder="t('projects.projectNamePlaceholder')" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.projectDescription') }}</label>
            <textarea v-model="newProject.description" rows="3" class="input" :placeholder="t('projects.projectDescriptionPlaceholder')"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">{{ t('projects.projectColor') }}</label>
            <div class="flex space-x-2">
              <button v-for="color in colors" :key="color" type="button" @click="newProject.color = color"
                class="w-8 h-8 rounded-full transition-transform hover:scale-110"
                :class="{ 'ring-2 ring-offset-2 ring-gray-400': newProject.color === color }"
                :style="{ backgroundColor: color }">
              </button>
            </div>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showModal = false" class="btn-secondary">{{ t('common.cancel') }}</button>
            <button type="submit" :disabled="projectsStore.loading" class="btn-primary">{{ t('common.create') }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
