<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/stores/settings'
import { FolderKanban, Layout, GripVertical, ChevronRight, ChevronLeft, X } from 'lucide-vue-next'

const { t } = useI18n()
const settingsStore = useSettingsStore()

const currentStep = ref(0)

const steps = [
  {
    icon: FolderKanban,
    titleKey: 'onboarding.step1Title',
    descriptionKey: 'onboarding.step1Description',
    image: '📁'
  },
  {
    icon: Layout,
    titleKey: 'onboarding.step2Title',
    descriptionKey: 'onboarding.step2Description',
    image: '📊'
  },
  {
    icon: GripVertical,
    titleKey: 'onboarding.step3Title',
    descriptionKey: 'onboarding.step3Description',
    image: '✨'
  }
]

const nextStep = () => {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  } else {
    finish()
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const skip = () => {
  settingsStore.completeOnboarding()
}

const finish = () => {
  settingsStore.completeOnboarding()
}
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-300"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="!settingsStore.onboardingCompleted"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
      >
        <div class="relative w-full max-w-lg bg-white dark:bg-dark-800 rounded-2xl shadow-2xl overflow-hidden">
          <button
            @click="skip"
            class="absolute top-4 right-4 p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 rounded-lg hover:bg-gray-100 dark:hover:bg-dark-700 transition-colors z-10"
          >
            <X :size="20" />
          </button>

          <div class="p-8">
            <div class="text-center mb-8">
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
                {{ t('onboarding.welcome') }}
              </h2>
            </div>

            <Transition
              mode="out-in"
              enter-active-class="transition-all duration-300"
              enter-from-class="opacity-0 translate-x-4"
              enter-to-class="opacity-100 translate-x-0"
              leave-active-class="transition-all duration-300"
              leave-from-class="opacity-100 translate-x-0"
              leave-to-class="opacity-0 -translate-x-4"
            >
              <div :key="currentStep" class="text-center">
                <div class="w-20 h-20 mx-auto mb-6 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center">
                  <span class="text-4xl">{{ steps[currentStep].image }}</span>
                </div>
                <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-3">
                  {{ t(steps[currentStep].titleKey) }}
                </h3>
                <p class="text-gray-600 dark:text-gray-400">
                  {{ t(steps[currentStep].descriptionKey) }}
                </p>
              </div>
            </Transition>

            <div class="flex justify-center gap-2 mt-8 mb-6">
              <button
                v-for="(step, index) in steps"
                :key="index"
                @click="currentStep = index"
                :class="[
                  'w-2 h-2 rounded-full transition-all duration-300',
                  currentStep === index
                    ? 'bg-primary-500 w-6'
                    : 'bg-gray-300 dark:bg-dark-600 hover:bg-gray-400'
                ]"
              />
            </div>

            <div class="flex items-center justify-between">
              <button
                @click="skip"
                class="text-sm text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
              >
                {{ t('onboarding.skip') }}
              </button>

              <div class="flex items-center gap-2">
                <button
                  v-if="currentStep > 0"
                  @click="prevStep"
                  class="p-2 text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-dark-700 rounded-lg transition-colors"
                >
                  <ChevronLeft :size="20" />
                </button>
                <button
                  @click="nextStep"
                  class="flex items-center gap-2 px-4 py-2 bg-primary-600 hover:bg-primary-700 text-white rounded-lg transition-colors"
                >
                  <span>{{ currentStep === steps.length - 1 ? t('onboarding.finish') : t('common.next') }}</span>
                  <ChevronRight v-if="currentStep < steps.length - 1" :size="18" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
