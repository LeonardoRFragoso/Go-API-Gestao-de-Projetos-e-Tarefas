<script setup>
import { onMounted, ref } from 'vue'
import { useTeamsStore } from '@/stores/teams'
import { Users, Plus, Trash2, Settings } from 'lucide-vue-next'

const teamsStore = useTeamsStore()
const showModal = ref(false)
const newTeam = ref({ name: '', description: '', color: '#8B5CF6' })

const colors = ['#8B5CF6', '#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#EC4899', '#06B6D4', '#84CC16']

onMounted(() => {
  teamsStore.fetchTeams()
})

const handleCreate = async () => {
  const result = await teamsStore.createTeam(newTeam.value)
  if (result) {
    showModal.value = false
    newTeam.value = { name: '', description: '', color: '#8B5CF6' }
  }
}

const handleDelete = async (id) => {
  if (confirm('Tem certeza que deseja excluir esta equipe?')) {
    await teamsStore.deleteTeam(id)
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Equipes</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">Gerencie suas equipes e squads</p>
      </div>
      <button @click="showModal = true" class="btn-primary flex items-center">
        <Plus :size="20" class="mr-2" />
        Nova Equipe
      </button>
    </div>

    <div v-if="teamsStore.loading" class="text-center py-12 text-gray-500">Carregando...</div>
    <div v-else-if="teamsStore.teams.length === 0" class="card p-12 text-center">
      <Users :size="64" class="mx-auto text-gray-300 dark:text-gray-600" />
      <h3 class="mt-4 text-lg font-medium text-gray-900 dark:text-white">Nenhuma equipe</h3>
      <p class="mt-2 text-gray-500 dark:text-gray-400">Crie uma equipe para colaborar com outros membros</p>
      <button @click="showModal = true" class="btn-primary mt-4">Criar Equipe</button>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="team in teamsStore.teams" :key="team.id" class="card hover:shadow-md transition-shadow">
        <RouterLink :to="`/teams/${team.id}`" class="block p-6">
          <div class="flex items-start justify-between">
            <div class="flex items-center">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center mr-3" :style="{ backgroundColor: team.color + '20' }">
                <Users :size="20" :style="{ color: team.color }" />
              </div>
              <div>
                <h3 class="font-semibold text-gray-900 dark:text-white">{{ team.name }}</h3>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ team.member_count }} membros</p>
              </div>
            </div>
          </div>
          <p v-if="team.description" class="mt-3 text-sm text-gray-600 dark:text-gray-400 line-clamp-2">{{ team.description }}</p>
          <div class="mt-4 flex items-center text-sm text-gray-500 dark:text-gray-400">
            <span>{{ team.project_count || 0 }} projetos</span>
          </div>
        </RouterLink>
        <div class="px-6 py-3 bg-gray-50 dark:bg-dark-700 border-t border-gray-100 dark:border-dark-600 flex justify-end gap-2">
          <RouterLink :to="`/teams/${team.id}/settings`" class="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white flex items-center">
            <Settings :size="16" class="mr-1" />
            Configurar
          </RouterLink>
          <button @click.prevent="handleDelete(team.id)" class="text-sm text-red-600 hover:text-red-700 flex items-center">
            <Trash2 :size="16" class="mr-1" />
            Excluir
          </button>
        </div>
      </div>
    </div>

    <!-- Modal de criação -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-dark-800 rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100 dark:border-dark-600">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Nova Equipe</h2>
        </div>
        <form @submit.prevent="handleCreate" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nome</label>
            <input v-model="newTeam.name" type="text" required class="input" placeholder="Nome da equipe" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Descrição</label>
            <textarea v-model="newTeam.description" rows="3" class="input" placeholder="Descrição opcional"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Cor</label>
            <div class="flex space-x-2">
              <button v-for="color in colors" :key="color" type="button" @click="newTeam.color = color"
                class="w-8 h-8 rounded-full transition-transform hover:scale-110"
                :class="{ 'ring-2 ring-offset-2 ring-gray-400': newTeam.color === color }"
                :style="{ backgroundColor: color }">
              </button>
            </div>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showModal = false" class="btn-secondary">Cancelar</button>
            <button type="submit" :disabled="teamsStore.loading" class="btn-primary">Criar</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
