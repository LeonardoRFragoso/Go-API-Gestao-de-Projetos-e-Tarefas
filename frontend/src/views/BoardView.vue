<script setup>
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useBoardsStore } from '@/stores/boards'
import { useTasksStore } from '@/stores/tasks'
import { useI18n } from 'vue-i18n'
import { useToast } from 'vue-toastification'
import draggable from 'vuedraggable'
import { ChevronLeft, Plus, X, Search } from 'lucide-vue-next'
import BaseModal from '@/components/ui/BaseModal.vue'
import TaskPriorityBadge from '@/components/ui/TaskPriorityBadge.vue'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import SearchBar from '@/components/ui/SearchBar.vue'
import FilterDropdown from '@/components/ui/FilterDropdown.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const { t } = useI18n()
const toast = useToast()
const route = useRoute()
const boardsStore = useBoardsStore()
const tasksStore = useTasksStore()

const showTaskModal = ref(false)
const selectedTask = ref(null)
const newTask = ref({ title: '', description: '', priority: 'medium', list_id: '' })
const searchQuery = ref('')
const filters = ref({})

onMounted(async () => {
  await boardsStore.fetchBoard(route.params.id)
})

const getColumnClass = (index) => {
  const classes = ['kanban-column-todo', 'kanban-column-progress', 'kanban-column-done']
  return classes[index] || 'kanban-column'
}

const filteredLists = computed(() => {
  if (!boardsStore.currentBoard?.lists) return []
  
  return boardsStore.currentBoard.lists.map(list => ({
    ...list,
    tasks: (list.tasks || []).filter(task => {
      const matchesSearch = !searchQuery.value || 
        task.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        task.description?.toLowerCase().includes(searchQuery.value.toLowerCase())
      
      const matchesPriority = !filters.value.priority || task.priority === filters.value.priority
      
      return matchesSearch && matchesPriority
    })
  }))
})

const openNewTaskModal = (listId) => {
  newTask.value = { title: '', description: '', priority: 'medium', list_id: listId }
  showTaskModal.value = true
}

const handleCreateTask = async () => {
  const result = await tasksStore.createTask(newTask.value)
  if (result) {
    showTaskModal.value = false
    toast.success(t('toast.taskCreated'))
    await boardsStore.fetchBoard(route.params.id)
  }
}

const handleTaskClick = async (task) => {
  selectedTask.value = await tasksStore.fetchTask(task.id)
}

const closeTaskDetail = () => {
  selectedTask.value = null
}

const handleUpdateTask = async () => {
  if (selectedTask.value) {
    await tasksStore.updateTask(selectedTask.value.id, {
      title: selectedTask.value.title,
      description: selectedTask.value.description,
      priority: selectedTask.value.priority
    })
    toast.success(t('toast.taskUpdated'))
    closeTaskDetail()
    await boardsStore.fetchBoard(route.params.id)
  }
}

const handleDeleteTask = async (taskId) => {
  if (confirm(t('tasks.deleteConfirm'))) {
    await tasksStore.deleteTask(taskId)
    toast.success(t('toast.taskDeleted'))
    selectedTask.value = null
    await boardsStore.fetchBoard(route.params.id)
  }
}

const onTaskMove = async (evt, listId) => {
  if (evt.added || evt.moved) {
    const task = evt.added?.element || evt.moved?.element
    const newIndex = evt.added?.newIndex ?? evt.moved?.newIndex
    await tasksStore.moveTask(task.id, listId, newIndex)
    toast.success(t('toast.taskMoved'))
    await boardsStore.fetchBoard(route.params.id)
  }
}
</script>

<template>
  <div class="h-[calc(100vh-8rem)]">
    <div v-if="boardsStore.loading" class="flex items-center justify-center py-12">
      <div class="flex space-x-4">
        <div v-for="i in 3" :key="i" class="w-80">
          <SkeletonLoader type="title" class="mb-4" />
          <SkeletonLoader type="card" />
          <SkeletonLoader type="card" class="mt-3" />
        </div>
      </div>
    </div>

    <template v-else-if="boardsStore.currentBoard">
      <div class="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <RouterLink 
            :to="`/projects/${boardsStore.currentBoard.project_id}`" 
            class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 flex items-center mb-2 transition-colors"
          >
            <ChevronLeft :size="16" class="mr-1" />
            {{ t('common.back') }}
          </RouterLink>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ boardsStore.currentBoard.name }}</h1>
        </div>
        
        <div class="flex items-center gap-3">
          <SearchBar v-model="searchQuery" class="w-64" />
          <FilterDropdown v-model="filters" />
        </div>
      </div>

      <div class="flex space-x-4 overflow-x-auto pb-4 h-full scrollbar-thin">
        <div 
          v-for="(list, index) in filteredLists" 
          :key="list.id" 
          :class="getColumnClass(index)"
        >
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-gray-700 dark:text-gray-300">{{ list.name }}</h3>
            <span class="text-sm text-gray-500 dark:text-gray-400 bg-gray-200 dark:bg-dark-600 px-2 py-0.5 rounded-full">
              {{ list.tasks?.length || 0 }}
            </span>
          </div>
          
          <draggable
            :list="list.tasks || []"
            group="tasks"
            item-key="id"
            class="space-y-3 min-h-[200px]"
            ghost-class="opacity-50"
            drag-class="shadow-2xl rotate-2"
            @change="(evt) => onTaskMove(evt, list.id)"
          >
            <template #item="{ element: task }">
              <div 
                @click="handleTaskClick(task)" 
                class="task-card group animate-fade-in"
              >
                <div class="flex items-start justify-between">
                  <h4 class="font-medium text-gray-900 dark:text-white text-sm group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
                    {{ task.title }}
                  </h4>
                </div>
                <p v-if="task.description" class="text-xs text-gray-500 dark:text-gray-400 mt-1 line-clamp-2">
                  {{ task.description }}
                </p>
                <div class="flex items-center justify-between mt-3">
                  <TaskPriorityBadge :priority="task.priority" size="xs" />
                  <div v-if="task.assignees?.length" class="flex -space-x-1">
                    <UserAvatar 
                      v-for="assignee in task.assignees.slice(0, 3)" 
                      :key="assignee.id"
                      :user="assignee"
                      size="xs"
                    />
                  </div>
                </div>
              </div>
            </template>
          </draggable>

          <button 
            @click="openNewTaskModal(list.id)" 
            class="w-full mt-3 py-2 text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-200 dark:hover:bg-dark-600 rounded-lg transition-colors flex items-center justify-center group"
          >
            <Plus :size="16" class="mr-1 group-hover:scale-110 transition-transform" />
            {{ t('tasks.addTask') }}
          </button>
        </div>
      </div>
    </template>

    <BaseModal :show="showTaskModal" :title="t('tasks.newTask')" @close="showTaskModal = false">
      <form @submit.prevent="handleCreateTask" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.taskTitle') }}</label>
          <input v-model="newTask.title" type="text" required class="input" :placeholder="t('tasks.taskTitle')" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.taskDescription') }}</label>
          <textarea v-model="newTask.description" rows="3" class="input" :placeholder="t('tasks.taskDescription')"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.priority') }}</label>
          <select v-model="newTask.priority" class="input">
            <option value="low">{{ t('tasks.priorities.low') }}</option>
            <option value="medium">{{ t('tasks.priorities.medium') }}</option>
            <option value="high">{{ t('tasks.priorities.high') }}</option>
            <option value="urgent">{{ t('tasks.priorities.urgent') }}</option>
          </select>
        </div>
        <div class="flex justify-end space-x-3 pt-4">
          <button type="button" @click="showTaskModal = false" class="btn-secondary">{{ t('common.cancel') }}</button>
          <button type="submit" :disabled="tasksStore.loading" class="btn-primary">{{ t('common.create') }}</button>
        </div>
      </form>
    </BaseModal>

    <BaseModal :show="!!selectedTask" :title="t('tasks.taskDetails')" @close="closeTaskDetail">
      <div v-if="selectedTask" class="p-6 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.taskTitle') }}</label>
          <input v-model="selectedTask.title" type="text" class="input" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.taskDescription') }}</label>
          <textarea v-model="selectedTask.description" rows="4" class="input"></textarea>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('tasks.priority') }}</label>
          <select v-model="selectedTask.priority" class="input">
            <option value="low">{{ t('tasks.priorities.low') }}</option>
            <option value="medium">{{ t('tasks.priorities.medium') }}</option>
            <option value="high">{{ t('tasks.priorities.high') }}</option>
            <option value="urgent">{{ t('tasks.priorities.urgent') }}</option>
          </select>
        </div>
        <div class="flex justify-between pt-4">
          <button @click="handleDeleteTask(selectedTask.id)" class="btn-danger">{{ t('common.delete') }}</button>
          <div class="space-x-3">
            <button @click="closeTaskDetail" class="btn-secondary">{{ t('common.cancel') }}</button>
            <button @click="handleUpdateTask" class="btn-primary">{{ t('common.save') }}</button>
          </div>
        </div>
      </div>
    </BaseModal>
  </div>
</template>
