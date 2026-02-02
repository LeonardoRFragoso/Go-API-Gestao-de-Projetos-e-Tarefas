<script setup>
import { onMounted, ref } from 'vue'
import { useProjectsStore } from '@/stores/projects'

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
  if (confirm('Tem certeza que deseja excluir este projeto?')) {
    await projectsStore.deleteProject(id)
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Projetos</h1>
        <p class="text-gray-600 mt-1">Gerencie seus projetos</p>
      </div>
      <button @click="showModal = true" class="btn-primary flex items-center">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Novo Projeto
      </button>
    </div>

    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">Carregando...</div>
    <div v-else-if="projectsStore.projects.length === 0" class="card p-12 text-center">
      <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-gray-900">Nenhum projeto</h3>
      <p class="mt-2 text-gray-500">Crie seu primeiro projeto para começar</p>
      <button @click="showModal = true" class="btn-primary mt-4">Criar Projeto</button>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="project in projectsStore.projects" :key="project.id" class="card hover:shadow-md transition-shadow">
        <RouterLink :to="`/projects/${project.id}`" class="block p-6">
          <div class="flex items-start justify-between">
            <div class="flex items-center">
              <div class="w-4 h-4 rounded-full mr-3" :style="{ backgroundColor: project.color }"></div>
              <h3 class="font-semibold text-gray-900">{{ project.name }}</h3>
            </div>
          </div>
          <p v-if="project.description" class="mt-2 text-sm text-gray-600 line-clamp-2">{{ project.description }}</p>
          <div class="mt-4 flex items-center text-sm text-gray-500">
            <span>{{ project.board_count || 0 }} quadros</span>
            <span class="mx-2">•</span>
            <span>{{ project.member_count || 1 }} membros</span>
          </div>
        </RouterLink>
        <div class="px-6 py-3 bg-gray-50 border-t border-gray-100 flex justify-end">
          <button @click.prevent="handleDelete(project.id)" class="text-sm text-red-600 hover:text-red-700">
            Excluir
          </button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100">
          <h2 class="text-xl font-semibold text-gray-900">Novo Projeto</h2>
        </div>
        <form @submit.prevent="handleCreate" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
            <input v-model="newProject.name" type="text" required class="input" placeholder="Nome do projeto" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
            <textarea v-model="newProject.description" rows="3" class="input" placeholder="Descrição opcional"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Cor</label>
            <div class="flex space-x-2">
              <button v-for="color in colors" :key="color" type="button" @click="newProject.color = color"
                class="w-8 h-8 rounded-full transition-transform hover:scale-110"
                :class="{ 'ring-2 ring-offset-2 ring-gray-400': newProject.color === color }"
                :style="{ backgroundColor: color }">
              </button>
            </div>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showModal = false" class="btn-secondary">Cancelar</button>
            <button type="submit" :disabled="projectsStore.loading" class="btn-primary">Criar</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
