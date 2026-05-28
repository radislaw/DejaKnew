<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDeckStore } from '@/entities/deck'
import { DeckCreateForm } from '@/features/deck-create'
import { DeckList } from '@/widgets/deck-list'
import { CardExportButton } from '@/features/card-export'
import { AppLayout } from '@/widgets/app-layout'

const { t } = useI18n()
const decks = useDeckStore()
const showForm = ref(false)
const createError = ref<string | null>(null)

onMounted(() => decks.fetchDecks())

async function onCreate(title: string, description: string) {
  createError.value = null
  try {
    await decks.addDeck({ title, description })
    showForm.value = false
  } catch (e) {
    createError.value = e instanceof Error ? e.message : t('home.failedToCreateDeck')
  }
}

async function onDelete(id: string) {
  try {
    await decks.removeDeck(id)
  } catch {
    // best-effort; leave deck in list on failure
  }
}
</script>

<template>
  <AppLayout>
    <main class="mx-auto max-w-4xl">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-medium">{{ t('home.myDecks') }}</h2>
        <div class="flex gap-2">
          <CardExportButton />
          <button
            v-if="!showForm"
            class="rounded-lg bg-sky-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-sky-500"
            @click="showForm = true"
          >
            {{ t('home.newDeck') }}
          </button>
        </div>
      </div>

      <DeckCreateForm
        v-if="showForm"
        class="mb-4"
        @submit="onCreate"
        @cancel="showForm = false"
      />

      <p v-if="createError" class="mb-3 text-sm text-red-400">{{ createError }}</p>

      <p v-if="decks.loading" class="text-sm text-slate-400">{{ t('common.loading') }}</p>
      <p v-else-if="decks.error" class="text-sm text-red-400">{{ decks.error }}</p>
      <p v-else-if="decks.decks.length === 0" class="text-sm text-slate-500">
        {{ t('home.noDecks') }}
      </p>
      <DeckList v-else :decks="decks.decks" @delete="onDelete" />
    </main>
  </AppLayout>
</template>
