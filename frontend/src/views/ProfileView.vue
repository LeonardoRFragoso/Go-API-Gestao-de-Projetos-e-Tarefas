<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const name = ref(authStore.user?.name || '')
const success = ref('')

const handleUpdate = async () => {
  success.value = ''
  const result = await authStore.updateProfile({ name: name.value })
  if (result) {
    success.value = 'Perfil atualizado com sucesso!'
  }
}
</script>

<template>
  <div class="max-w-2xl">
    <h1 class="text-2xl font-bold text-gray-900 mb-8">Meu Perfil</h1>

    <div class="card p-6">
      <form @submit.prevent="handleUpdate" class="space-y-6">
        <div v-if="success" class="bg-green-50 text-green-600 p-3 rounded-lg text-sm">
          {{ success }}
        </div>
        <div v-if="authStore.error" class="bg-red-50 text-red-600 p-3 rounded-lg text-sm">
          {{ authStore.error }}
        </div>

        <div class="flex items-center space-x-4">
          <div class="w-20 h-20 bg-primary-600 rounded-full flex items-center justify-center text-white text-2xl font-bold">
            {{ authStore.user?.name?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          <div>
            <p class="font-medium text-gray-900">{{ authStore.user?.name }}</p>
            <p class="text-sm text-gray-500">{{ authStore.user?.email }}</p>
          </div>
        </div>

        <hr class="border-gray-200">

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
          <input v-model="name" type="text" required class="input" />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input :value="authStore.user?.email" type="email" disabled class="input bg-gray-50 cursor-not-allowed" />
          <p class="text-xs text-gray-500 mt-1">O email não pode ser alterado</p>
        </div>

        <div class="flex justify-end">
          <button type="submit" :disabled="authStore.loading" class="btn-primary">
            <span v-if="authStore.loading">Salvando...</span>
            <span v-else>Salvar Alterações</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
