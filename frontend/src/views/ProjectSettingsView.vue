<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useProjectsStore } from '@/stores/projects'
import { useToast } from 'vue-toastification'
import { ArrowLeft, Save, Trash2, Users, Lock, AlertCircle } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const projectsStore = useProjectsStore()

const loading = ref(false)
const formData = ref({
  name: '',
  description: '',
  color: '#3B82F6'
})

const showDeleteConfirm = ref(false)

onMounted(async () => {
  await projectsStore.fetchProject(route.params.id)
  if (projectsStore.currentProject) {
    formData.value = {
      name: projectsStore.currentProject.name,
      description: projectsStore.currentProject.description || '',
      color: projectsStore.currentProject.color || '#3B82F6'
    }
  }
})

const handleSave = async () => {
  loading.value = true
  try {
    await projectsStore.updateProject(route.params.id, formData.value)
    toast.success(t('projects.updateSuccess'))
  } catch (error) {
    toast.error(t('projects.updateError'))
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  loading.value = true
  try {
    await projectsStore.deleteProject(route.params.id)
    toast.success(t('projects.deleteSuccess'))
    router.push('/projects')
  } catch (error) {
    toast.error(t('projects.deleteError'))
  } finally {
    loading.value = false
    showDeleteConfirm.value = false
  }
}

const isModified = computed(() => {
  return formData.value.name !== projectsStore.currentProject?.name ||
         formData.value.description !== (projectsStore.currentProject?.description || '') ||
         formData.value.color !== (projectsStore.currentProject?.color || '#3B82F6')
})
</script>

<template>
  <div>
    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">{{ t('common.loading') }}</div>
    <template v-else-if="projectsStore.currentProject">
      <div class="mb-8">
        <RouterLink :to="`/projects/${route.params.id}`" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 flex items-center mb-4 group">
          <ArrowLeft class="w-4 h-4 mr-1 group-hover:-translate-x-1 transition-transform" />
          {{ t('projects.backToProject') }}
        </RouterLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('projects.settings') }}</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-2">{{ t('projects.settingsDescription') }}</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <div class="card p-6">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <Lock :size="20" />
              {{ t('projects.basicInfo') }}
            </h2>
            <form @submit.prevent="handleSave" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.name') }}</label>
                <input v-model="formData.name" type="text" required class="input" />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.description') }}</label>
                <textarea v-model="formData.description" rows="4" class="input" :placeholder="t('projects.descriptionPlaceholder')"></textarea>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.color') }}</label>
                <div class="flex items-center gap-3">
                  <input v-model="formData.color" type="color" class="w-12 h-12 rounded cursor-pointer" />
                  <span class="text-sm text-gray-600 dark:text-gray-400">{{ formData.color }}</span>
                </div>
              </div>

              <div class="flex justify-end pt-4">
                <button 
                  type="submit" 
                  :disabled="!isModified || loading"
                  class="btn-primary flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Save :size="18" />
                  {{ t('common.save') }}
                </button>
              </div>
            </form>
          </div>

          <div class="card p-6">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
              <Users :size="20" />
              {{ t('projects.members') }}
            </h2>
            <p class="text-gray-600 dark:text-gray-400 mb-4">{{ t('projects.membersDescription') }}</p>
            <RouterLink :to="`/projects/${route.params.id}/members`" class="btn-secondary">
              {{ t('projects.managemembers') }}
            </RouterLink>
          </div>

          <div class="card p-6 border-l-4 border-red-500">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-2 flex items-center gap-2">
              <AlertCircle :size="20" class="text-red-500" />
              {{ t('projects.dangerZone') }}
            </h2>
            <p class="text-gray-600 dark:text-gray-400 mb-4">{{ t('projects.dangerZoneDescription') }}</p>
            <button 
              @click="showDeleteConfirm = true"
              class="btn-danger flex items-center gap-2"
            >
              <Trash2 :size="18" />
              {{ t('projects.deleteProject') }}
            </button>
          </div>
        </div>

        <div class="space-y-6">
          <div class="card p-6">
            <h3 class="font-semibold text-gray-900 dark:text-white mb-4">{{ t('projects.preview') }}</h3>
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg" :style="{ backgroundColor: formData.color }"></div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ formData.name || t('projects.untitled') }}</p>
                  <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-1">{{ formData.description || t('projects.noDescription') }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="card p-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
            <h3 class="font-semibold text-blue-900 dark:text-blue-100 mb-2">{{ t('projects.tip') }}</h3>
            <p class="text-sm text-blue-800 dark:text-blue-200">{{ t('projects.settingsTip') }}</p>
          </div>
        </div>
      </div>

      <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md">
          <div class="p-6 border-b border-gray-100 dark:border-dark-600">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white flex items-center gap-2">
              <AlertCircle :size="20" class="text-red-500" />
              {{ t('projects.confirmDelete') }}
            </h2>
          </div>
          <div class="p-6 space-y-4">
            <p class="text-gray-600 dark:text-gray-400">{{ t('projects.deleteWarning') }}</p>
            <p class="text-sm text-gray-500 dark:text-gray-500">{{ t('projects.deleteWarningDetails') }}</p>
          </div>
          <div class="p-6 border-t border-gray-100 dark:border-dark-600 flex justify-end gap-3">
            <button 
              @click="showDeleteConfirm = false" 
              class="btn-secondary"
              :disabled="loading"
            >
              {{ t('common.cancel') }}
            </button>
            <button 
              @click="handleDelete" 
              class="btn-danger"
              :disabled="loading"
            >
              {{ loading ? t('common.loading') : t('projects.deleteProject') }}
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
