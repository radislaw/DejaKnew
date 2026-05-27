<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const emit = defineEmits<{
  submit: [title: string, description: string]
  cancel: []
}>()

const title = ref('')
const description = ref('')

function onSubmit() {
  const trimmed = title.value.trim()
  if (!trimmed) return
  emit('submit', trimmed, description.value.trim())
  title.value = ''
  description.value = ''
}
</script>

<template>
  <form
    class="flex flex-col gap-3 rounded-xl border border-slate-700 bg-slate-800/60 p-4"
    @submit.prevent="onSubmit"
  >
    <h2 class="text-sm font-medium text-slate-300">{{ t('deckCreate.title') }}</h2>
    <input
      v-model="title"
      type="text"
      :placeholder="t('deckCreate.titlePlaceholder')"
      maxlength="200"
      required
      class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
    />
    <input
      v-model="description"
      type="text"
      :placeholder="t('deckCreate.descriptionPlaceholder')"
      class="rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 placeholder-slate-500 outline-none focus:border-sky-500"
    />
    <div class="flex gap-2">
      <button
        type="submit"
        class="rounded-lg bg-sky-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-sky-500"
      >
        {{ t('deckCreate.create') }}
      </button>
      <button
        type="button"
        class="rounded-lg px-4 py-1.5 text-sm text-slate-400 hover:text-slate-200"
        @click="emit('cancel')"
      >
        {{ t('common.cancel') }}
      </button>
    </div>
  </form>
</template>
