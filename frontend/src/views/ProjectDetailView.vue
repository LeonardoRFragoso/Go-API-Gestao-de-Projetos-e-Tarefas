<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useProjectsStore } from '@/stores/projects'
import { useBoardsStore } from '@/stores/boards'
import { Users, Settings, MoreVertical, Plus, ArrowLeft, Grid3x3, Zap } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const projectsStore = useProjectsStore()
const boardsStore = useBoardsStore()

const showModal = ref(false)
const newBoard = ref({ name: '', description: '' })

onMounted(async () => {
  await projectsStore.fetchProject(route.params.id)
  await boardsStore.fetchBoardsByProject(route.params.id)
})

const handleCreateBoard = async () => {
  const result = await boardsStore.createBoard({
    ...newBoard.value,
    project_id: route.params.id
  })
  if (result) {
    showModal.value = false
    newBoard.value = { name: '', description: '' }
  }
}

const handleDeleteBoard = async (id) => {
  if (confirm(t('boards.deleteConfirm'))) {
    await boardsStore.deleteBoard(id)
  }
}

const handleDeleteProject = async () => {
  if (confirm(t('projects.deleteConfirm'))) {
    await projectsStore.deleteProject(route.params.id)
    router.push('/projects')
  }
}

const totalBoards = computed(() => boardsStore.boards.length)
const totalTasks = computed(() => {
  return boardsStore.boards.reduce((sum, board) => sum + (board.task_count || 0), 0)
})
const totalMembers = computed(() => {
  return projectsStore.currentProject?.member_count || 0
})
</script>

<template>
  <div>
    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">{{ t('common.loading') }}</div>
    <template v-else-if="projectsStore.currentProject">
      <div class="mb-8">
        <RouterLink to="/projects" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 flex items-center mb-4 group">
          <ArrowLeft class="w-4 h-4 mr-1 group-hover:-translate-x-1 transition-transform" />
          {{ t('projects.backToProjects') }}
        </RouterLink>
        
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-3">
            <div class="w-6 h-6 rounded-lg" :style="{ backgroundColor: projectsStore.currentProject.color }"></div>
            <div>
              <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ projectsStore.currentProject.name }}</h1>
              <p v-if="projectsStore.currentProject.description" class="text-gray-600 dark:text-gray-400 mt-1">
                {{ projectsStore.currentProject.description }}
              </p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button @click="showModal = true" class="btn-primary flex items-center gap-2">
              <Plus :size="18" />
              {{ t('boards.newBoard') }}
            </button>
            <RouterLink :to="`/projects/${route.params.id}/settings`" class="p-2 hover:bg-gray-100 dark:hover:bg-dark-700 rounded-lg transition-colors">
              <Settings :size="20" class="text-gray-600 dark:text-gray-400" />
            </RouterLink>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="card p-4 flex items-center gap-4">
            <div class="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
              <Grid3x3 :size="20" class="text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">{{ t('projects.boards') }}</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ totalBoards }}</p>
            </div>
          </div>
          <div class="card p-4 flex items-center gap-4">
            <div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-lg">
              <Zap :size="20" class="text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">{{ t('projects.tasks') }}</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ totalTasks }}</p>
            </div>
          </div>
          <div class="card p-4 flex items-center gap-4">
            <div class="p-3 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
              <Users :size="20" class="text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <p class="text-sm text-gray-600 dark:text-gray-400">{{ t('projects.members') }}</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ totalMembers }}</p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="boardsStore.loading" class="text-center py-12 text-gray-500">{{ t('common.loading') }}</div>
      <div v-else-if="boardsStore.boards.length === 0" class="card p-12 text-center">
        <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">{{ t('boards.noBoards') }}</h3>
        <p class="mt-2 text-gray-500">{{ t('projects.createFirst') }}</p>
        <button @click="showModal = true" class="btn-primary mt-4">{{ t('boards.newBoard') }}</button>
      </div>
      <div v-else class="space-y-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('boards.boards') }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="board in boardsStore.boards" :key="board.id" class="card group hover:shadow-lg transition-all">
            <RouterLink :to="`/boards/${board.id}`" class="block p-6">
              <div class="flex items-start justify-between mb-2">
                <h3 class="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">{{ board.name }}</h3>
              </div>
              <p v-if="board.description" class="text-sm text-gray-600 dark:text-gray-400 line-clamp-2 mb-4">{{ board.description }}</p>
              <div class="flex items-center text-sm text-gray-500 dark:text-gray-400 gap-3">
                <span class="flex items-center gap-1">
                  <span class="font-medium">{{ board.list_count || 0 }}</span>
                  {{ t('projects.lists') }}
                </span>
                <span class="text-gray-300 dark:text-gray-600">•</span>
                <span class="flex items-center gap-1">
                  <span class="font-medium">{{ board.task_count || 0 }}</span>
                  {{ t('projects.tasks') }}
                </span>
              </div>
            </RouterLink>
            <div class="px-6 py-3 bg-gray-50 dark:bg-dark-700 border-t border-gray-100 dark:border-dark-600 flex justify-between items-center">
              <RouterLink :to="`/boards/${board.id}`" class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium">
                {{ t('common.open') }}
              </RouterLink>
              <button @click.prevent="handleDeleteBoard(board.id)" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300">
                {{ t('common.delete') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100 dark:border-dark-600">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('boards.newBoard') }}</h2>
        </div>
        <form @submit.prevent="handleCreateBoard" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('boards.boardName') }}</label>
            <input v-model="newBoard.name" type="text" required class="input" :placeholder="t('boards.boardNamePlaceholder')" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('boards.boardDescription') }}</label>
            <textarea v-model="newBoard.description" rows="3" class="input" :placeholder="t('boards.boardDescriptionPlaceholder')"></textarea>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showModal = false" class="btn-secondary">{{ t('common.cancel') }}</button>
            <button type="submit" :disabled="boardsStore.loading" class="btn-primary">{{ t('common.create') }}</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
