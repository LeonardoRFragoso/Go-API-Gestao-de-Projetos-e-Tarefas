<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useProjectsStore } from '@/stores/projects'
import { useToast } from 'vue-toastification'
import { ArrowLeft, Plus, Trash2, Shield, Users } from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const projectsStore = useProjectsStore()

const loading = ref(false)
const members = ref([])
const showAddModal = ref(false)
const newMemberEmail = ref('')
const newMemberRole = ref('member')

const roles = [
  { value: 'owner', label: 'Owner', description: 'Controle total do projeto' },
  { value: 'admin', label: 'Admin', description: 'Gerenciar membros e conteúdo' },
  { value: 'member', label: 'Member', description: 'Criar e editar tarefas' },
  { value: 'viewer', label: 'Viewer', description: 'Apenas visualizar' }
]

onMounted(async () => {
  await projectsStore.fetchProject(route.params.id)
  await fetchMembers()
})

const fetchMembers = async () => {
  loading.value = true
  try {
    const response = await projectsStore.fetchProjectMembers(route.params.id)
    members.value = response || []
  } catch (error) {
    toast.error(t('projects.fetchMembersError'))
  } finally {
    loading.value = false
  }
}

const handleAddMember = async () => {
  if (!newMemberEmail.value) {
    toast.error(t('projects.emailRequired'))
    return
  }

  loading.value = true
  try {
    await projectsStore.addProjectMember(route.params.id, {
      email: newMemberEmail.value,
      role: newMemberRole.value
    })
    toast.success(t('projects.memberAdded'))
    showAddModal.value = false
    newMemberEmail.value = ''
    newMemberRole.value = 'member'
    await fetchMembers()
  } catch (error) {
    toast.error(t('projects.addMemberError'))
  } finally {
    loading.value = false
  }
}

const handleRemoveMember = async (memberId) => {
  if (confirm(t('projects.removeMemberConfirm'))) {
    loading.value = true
    try {
      await projectsStore.removeProjectMember(route.params.id, memberId)
      toast.success(t('projects.memberRemoved'))
      await fetchMembers()
    } catch (error) {
      toast.error(t('projects.removeMemberError'))
    } finally {
      loading.value = false
    }
  }
}

const getRoleLabel = (role) => {
  return roles.find(r => r.value === role)?.label || role
}

const getRoleDescription = (role) => {
  return roles.find(r => r.value === role)?.description || ''
}
</script>

<template>
  <div>
    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">{{ t('common.loading') }}</div>
    <template v-else-if="projectsStore.currentProject">
      <div class="mb-8">
        <RouterLink :to="`/projects/${route.params.id}/settings`" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 flex items-center mb-4 group">
          <ArrowLeft class="w-4 h-4 mr-1 group-hover:-translate-x-1 transition-transform" />
          {{ t('projects.backToSettings') }}
        </RouterLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t('projects.manageMembers') }}</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-2">{{ t('projects.manageMembersDescription') }}</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2">
          <div class="card">
            <div class="p-6 border-b border-gray-100 dark:border-dark-700 flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
                <Users :size="20" />
                {{ t('projects.members') }} ({{ members.length }})
              </h2>
              <button @click="showAddModal = true" class="btn-primary flex items-center gap-2">
                <Plus :size="18" />
                {{ t('projects.addMember') }}
              </button>
            </div>

            <div v-if="loading" class="p-6 text-center text-gray-500">{{ t('common.loading') }}</div>
            <div v-else-if="members.length === 0" class="p-12 text-center">
              <Users :size="32" class="mx-auto text-gray-300 mb-4" />
              <p class="text-gray-500 dark:text-gray-400">{{ t('projects.noMembers') }}</p>
            </div>
            <div v-else class="divide-y divide-gray-100 dark:divide-dark-700">
              <div v-for="member in members" :key="member.id" class="p-6 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-dark-700/50 transition-colors">
                <div class="flex-1">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-full bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center text-white font-semibold">
                      {{ member.name?.charAt(0).toUpperCase() || 'U' }}
                    </div>
                    <div>
                      <p class="font-medium text-gray-900 dark:text-white">{{ member.name }}</p>
                      <p class="text-sm text-gray-600 dark:text-gray-400">{{ member.email }}</p>
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-4">
                  <div class="flex items-center gap-2">
                    <Shield :size="16" class="text-gray-400" />
                    <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ getRoleLabel(member.role) }}</span>
                  </div>
                  <button 
                    @click="handleRemoveMember(member.id)"
                    class="p-2 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors"
                    :disabled="loading"
                  >
                    <Trash2 :size="18" class="text-red-600 dark:text-red-400" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="card p-6">
            <h3 class="font-semibold text-gray-900 dark:text-white mb-4">{{ t('projects.roles') }}</h3>
            <div class="space-y-3">
              <div v-for="role in roles" :key="role.value" class="p-3 bg-gray-50 dark:bg-dark-700 rounded-lg">
                <p class="font-medium text-gray-900 dark:text-white text-sm">{{ role.label }}</p>
                <p class="text-xs text-gray-600 dark:text-gray-400 mt-1">{{ role.description }}</p>
              </div>
            </div>
          </div>

          <div class="card p-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800">
            <h3 class="font-semibold text-blue-900 dark:text-blue-100 mb-2">{{ t('projects.tip') }}</h3>
            <p class="text-sm text-blue-800 dark:text-blue-200">{{ t('projects.membersTip') }}</p>
          </div>
        </div>
      </div>

      <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md">
          <div class="p-6 border-b border-gray-100 dark:border-dark-600">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">{{ t('projects.addMember') }}</h2>
          </div>
          <form @submit.prevent="handleAddMember" class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.email') }}</label>
              <input v-model="newMemberEmail" type="email" required class="input" :placeholder="t('projects.emailPlaceholder')" />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('projects.role') }}</label>
              <select v-model="newMemberRole" class="input">
                <option v-for="role in roles" :key="role.value" :value="role.value">
                  {{ role.label }}
                </option>
              </select>
            </div>

            <div class="flex justify-end gap-3 pt-4">
              <button 
                type="button" 
                @click="showAddModal = false" 
                class="btn-secondary"
                :disabled="loading"
              >
                {{ t('common.cancel') }}
              </button>
              <button 
                type="submit" 
                class="btn-primary"
                :disabled="loading"
              >
                {{ loading ? t('common.loading') : t('common.add') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </template>
  </div>
</template>
