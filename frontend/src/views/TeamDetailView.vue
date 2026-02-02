<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTeamsStore } from '@/stores/teams'
import { useProjectsStore } from '@/stores/projects'
import { Users, UserPlus, Folder, Trash2, Crown, Shield, User } from 'lucide-vue-next'

const route = useRoute()
const teamsStore = useTeamsStore()
const projectsStore = useProjectsStore()

const showAddMemberModal = ref(false)
const showAddProjectModal = ref(false)
const newMember = ref({ email: '', role: 'member' })
const selectedProjectId = ref('')

const roleIcons = {
  lead: Crown,
  admin: Shield,
  member: User
}

const roleLabels = {
  lead: 'Líder',
  admin: 'Admin',
  member: 'Membro'
}

onMounted(async () => {
  await teamsStore.fetchTeam(route.params.id)
  await projectsStore.fetchProjects(1, 100)
})

const availableProjects = computed(() => {
  if (!teamsStore.currentTeam?.projects) return projectsStore.projects
  const teamProjectIds = teamsStore.currentTeam.projects.map(p => p.id)
  return projectsStore.projects.filter(p => !teamProjectIds.includes(p.id))
})

const handleRemoveMember = async (memberId) => {
  if (confirm('Tem certeza que deseja remover este membro?')) {
    await teamsStore.removeMember(route.params.id, memberId)
  }
}

const handleUpdateRole = async (memberId, role) => {
  await teamsStore.updateMemberRole(route.params.id, memberId, role)
}

const handleAddProject = async () => {
  if (selectedProjectId.value) {
    await teamsStore.addProject(route.params.id, selectedProjectId.value)
    showAddProjectModal.value = false
    selectedProjectId.value = ''
  }
}

const handleRemoveProject = async (projectId) => {
  if (confirm('Tem certeza que deseja remover este projeto da equipe?')) {
    await teamsStore.removeProject(route.params.id, projectId)
  }
}
</script>

<template>
  <div>
    <div v-if="teamsStore.loading" class="text-center py-12 text-gray-500">Carregando...</div>
    <template v-else-if="teamsStore.currentTeam">
      <div class="mb-8">
        <RouterLink to="/teams" class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 flex items-center mb-4">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Voltar para Equipes
        </RouterLink>
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="w-12 h-12 rounded-lg flex items-center justify-center mr-4" :style="{ backgroundColor: teamsStore.currentTeam.color + '20' }">
              <Users :size="24" :style="{ color: teamsStore.currentTeam.color }" />
            </div>
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ teamsStore.currentTeam.name }}</h1>
              <p v-if="teamsStore.currentTeam.description" class="text-gray-600 dark:text-gray-400 mt-1">
                {{ teamsStore.currentTeam.description }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Membros -->
        <div class="card">
          <div class="p-4 border-b border-gray-100 dark:border-dark-600 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
              <Users :size="20" class="mr-2" />
              Membros ({{ teamsStore.currentTeam.members?.length || 0 }})
            </h2>
            <button @click="showAddMemberModal = true" class="btn-primary text-sm py-1 px-3 flex items-center">
              <UserPlus :size="16" class="mr-1" />
              Adicionar
            </button>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-600">
            <div v-for="member in teamsStore.currentTeam.members" :key="member.user_id" 
              class="p-4 flex items-center justify-between">
              <div class="flex items-center">
                <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center text-primary-600 dark:text-primary-400 font-medium mr-3">
                  {{ member.user?.name?.charAt(0)?.toUpperCase() || '?' }}
                </div>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">{{ member.user?.name || 'Usuário' }}</p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">{{ member.user?.email }}</p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <select 
                  :value="member.role" 
                  @change="handleUpdateRole(member.user_id, $event.target.value)"
                  class="input text-sm py-1 px-2 w-auto"
                  :disabled="member.role === 'lead'"
                >
                  <option value="lead">Líder</option>
                  <option value="admin">Admin</option>
                  <option value="member">Membro</option>
                </select>
                <button 
                  v-if="member.role !== 'lead'"
                  @click="handleRemoveMember(member.user_id)" 
                  class="text-red-600 hover:text-red-700 p-1"
                >
                  <Trash2 :size="16" />
                </button>
              </div>
            </div>
            <div v-if="!teamsStore.currentTeam.members?.length" class="p-8 text-center text-gray-500 dark:text-gray-400">
              Nenhum membro na equipe
            </div>
          </div>
        </div>

        <!-- Projetos -->
        <div class="card">
          <div class="p-4 border-b border-gray-100 dark:border-dark-600 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
              <Folder :size="20" class="mr-2" />
              Projetos ({{ teamsStore.currentTeam.projects?.length || 0 }})
            </h2>
            <button @click="showAddProjectModal = true" class="btn-primary text-sm py-1 px-3 flex items-center">
              <Folder :size="16" class="mr-1" />
              Adicionar
            </button>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-600">
            <div v-for="project in teamsStore.currentTeam.projects" :key="project.id" 
              class="p-4 flex items-center justify-between">
              <RouterLink :to="`/projects/${project.id}`" class="flex items-center hover:text-primary-600">
                <div class="w-4 h-4 rounded-full mr-3" :style="{ backgroundColor: project.color }"></div>
                <span class="font-medium text-gray-900 dark:text-white">{{ project.name }}</span>
              </RouterLink>
              <button @click="handleRemoveProject(project.id)" class="text-red-600 hover:text-red-700 p-1">
                <Trash2 :size="16" />
              </button>
            </div>
            <div v-if="!teamsStore.currentTeam.projects?.length" class="p-8 text-center text-gray-500 dark:text-gray-400">
              Nenhum projeto vinculado
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Modal adicionar projeto -->
    <div v-if="showAddProjectModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100 dark:border-dark-600">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Adicionar Projeto</h2>
        </div>
        <form @submit.prevent="handleAddProject" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Selecione um projeto</label>
            <select v-model="selectedProjectId" class="input" required>
              <option value="">Selecione...</option>
              <option v-for="project in availableProjects" :key="project.id" :value="project.id">
                {{ project.name }}
              </option>
            </select>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showAddProjectModal = false" class="btn-secondary">Cancelar</button>
            <button type="submit" class="btn-primary">Adicionar</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
