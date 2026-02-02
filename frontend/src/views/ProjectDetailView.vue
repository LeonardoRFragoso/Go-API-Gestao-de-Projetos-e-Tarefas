<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectsStore } from '@/stores/projects'
import { useBoardsStore } from '@/stores/boards'

const route = useRoute()
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
  if (confirm('Tem certeza que deseja excluir este quadro?')) {
    await boardsStore.deleteBoard(id)
  }
}
</script>

<template>
  <div>
    <div v-if="projectsStore.loading" class="text-center py-12 text-gray-500">Carregando...</div>
    <template v-else-if="projectsStore.currentProject">
      <div class="mb-8">
        <RouterLink to="/projects" class="text-sm text-gray-500 hover:text-gray-700 flex items-center mb-4">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Voltar para Projetos
        </RouterLink>
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="w-4 h-4 rounded-full mr-3" :style="{ backgroundColor: projectsStore.currentProject.color }"></div>
            <h1 class="text-2xl font-bold text-gray-900">{{ projectsStore.currentProject.name }}</h1>
          </div>
          <button @click="showModal = true" class="btn-primary flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Novo Quadro
          </button>
        </div>
        <p v-if="projectsStore.currentProject.description" class="text-gray-600 mt-2">
          {{ projectsStore.currentProject.description }}
        </p>
      </div>

      <div v-if="boardsStore.loading" class="text-center py-12 text-gray-500">Carregando quadros...</div>
      <div v-else-if="boardsStore.boards.length === 0" class="card p-12 text-center">
        <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2" />
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">Nenhum quadro</h3>
        <p class="mt-2 text-gray-500">Crie um quadro para começar a gerenciar suas tarefas</p>
        <button @click="showModal = true" class="btn-primary mt-4">Criar Quadro</button>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="board in boardsStore.boards" :key="board.id" class="card hover:shadow-md transition-shadow">
          <RouterLink :to="`/boards/${board.id}`" class="block p-6">
            <h3 class="font-semibold text-gray-900">{{ board.name }}</h3>
            <p v-if="board.description" class="mt-2 text-sm text-gray-600 line-clamp-2">{{ board.description }}</p>
            <div class="mt-4 flex items-center text-sm text-gray-500">
              <span>{{ board.list_count || 0 }} listas</span>
              <span class="mx-2">•</span>
              <span>{{ board.task_count || 0 }} tarefas</span>
            </div>
          </RouterLink>
          <div class="px-6 py-3 bg-gray-50 border-t border-gray-100 flex justify-end">
            <button @click.prevent="handleDeleteBoard(board.id)" class="text-sm text-red-600 hover:text-red-700">
              Excluir
            </button>
          </div>
        </div>
      </div>
    </template>

    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl shadow-xl w-full max-w-md mx-4">
        <div class="p-6 border-b border-gray-100">
          <h2 class="text-xl font-semibold text-gray-900">Novo Quadro</h2>
        </div>
        <form @submit.prevent="handleCreateBoard" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
            <input v-model="newBoard.name" type="text" required class="input" placeholder="Nome do quadro" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
            <textarea v-model="newBoard.description" rows="3" class="input" placeholder="Descrição opcional"></textarea>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showModal = false" class="btn-secondary">Cancelar</button>
            <button type="submit" :disabled="boardsStore.loading" class="btn-primary">Criar</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
