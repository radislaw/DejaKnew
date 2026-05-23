<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDeck } from '@/entities/deck'
import { useCardStore } from '@/entities/card'
import { CardCreateForm } from '@/features/card-create'

const route = useRoute()
const router = useRouter()
const deckId = route.params.id as string

const cardStore = useCardStore()
const deckTitle = ref('')
const deckDescription = ref('')
const loading = ref(true)
const error = ref<string | null>(null)

const showForm = ref(false)
const createError = ref<string | null>(null)

const editingId = ref<string | null>(null)
const editFront = ref('')
const editBack = ref('')
const editError = ref<string | null>(null)

onMounted(async () => {
  try {
    const detail = await getDeck(deckId)
    deckTitle.value = detail.title
    deckDescription.value = detail.description
    cardStore.setCards(detail.cards)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load deck'
  } finally {
    loading.value = false
  }
})

async function onCreate(front: string, back: string) {
  createError.value = null
  try {
    await cardStore.addCard(deckId, { front, back })
    showForm.value = false
  } catch (e) {
    createError.value = e instanceof Error ? e.message : 'Failed to add card'
  }
}

function startEdit(cardId: string, front: string, back: string) {
  editingId.value = cardId
  editFront.value = front
  editBack.value = back
  editError.value = null
}

function cancelEdit() {
  editingId.value = null
}

async function saveEdit(cardId: string) {
  editError.value = null
  try {
    await cardStore.updateCardItem(cardId, { front: editFront.value, back: editBack.value })
    editingId.value = null
  } catch (e) {
    editError.value = e instanceof Error ? e.message : 'Failed to save'
  }
}

async function onDelete(cardId: string) {
  try {
    await cardStore.removeCard(cardId)
  } catch {
    // best-effort
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#070d1c] px-4 py-6 text-slate-100">
    <header class="mx-auto mb-8 flex max-w-4xl items-center gap-4">
      <button
        class="text-slate-400 hover:text-slate-200 text-sm"
        @click="router.push({ name: 'home' })"
      >
        ← Back
      </button>
      <div class="flex-1">
        <h1 class="text-xl font-semibold tracking-tight">{{ deckTitle }}</h1>
        <p v-if="deckDescription" class="text-sm text-slate-400 mt-0.5">{{ deckDescription }}</p>
      </div>
      <button
        class="rounded-lg bg-emerald-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-500"
        @click="router.push({ name: 'study', params: { id: deckId } })"
      >
        Study
      </button>
    </header>

    <main class="mx-auto max-w-4xl">
      <p v-if="loading" class="text-sm text-slate-400">Loading…</p>
      <p v-else-if="error" class="text-sm text-red-400">{{ error }}</p>

      <template v-else>
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-medium">
            Cards <span class="text-slate-500 text-sm font-normal">({{ cardStore.cards.length }})</span>
          </h2>
          <button
            v-if="!showForm"
            class="rounded-lg bg-sky-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-sky-500"
            @click="showForm = true"
          >
            + Add card
          </button>
        </div>

        <CardCreateForm
          v-if="showForm"
          class="mb-4"
          @submit="onCreate"
          @cancel="showForm = false"
        />

        <p v-if="createError" class="mb-3 text-sm text-red-400">{{ createError }}</p>

        <p v-if="cardStore.cards.length === 0" class="text-sm text-slate-500">
          No cards yet. Add your first one.
        </p>

        <ul v-else class="flex flex-col gap-3">
          <li
            v-for="card in cardStore.cards"
            :key="card.id"
            class="group rounded-xl border border-slate-700 bg-slate-800/60 p-4"
          >
            <template v-if="editingId === card.id">
              <div class="flex flex-col gap-2">
                <textarea
                  v-model="editFront"
                  rows="2"
                  class="resize-none rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-sky-500"
                />
                <textarea
                  v-model="editBack"
                  rows="2"
                  class="resize-none rounded-lg border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-slate-100 outline-none focus:border-sky-500"
                />
                <p v-if="editError" class="text-xs text-red-400">{{ editError }}</p>
                <div class="flex gap-2">
                  <button
                    class="rounded-lg bg-sky-600 px-3 py-1 text-sm font-medium text-white hover:bg-sky-500"
                    @click="saveEdit(card.id)"
                  >
                    Save
                  </button>
                  <button
                    class="rounded-lg px-3 py-1 text-sm text-slate-400 hover:text-slate-200"
                    @click="cancelEdit"
                  >
                    Cancel
                  </button>
                </div>
              </div>
            </template>

            <template v-else>
              <div class="flex items-start justify-between gap-2">
                <div class="flex-1 grid sm:grid-cols-2 gap-x-6 gap-y-1">
                  <div>
                    <p class="text-xs text-slate-500 mb-0.5">Front</p>
                    <p class="text-sm text-slate-100 whitespace-pre-wrap">{{ card.content.front }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-slate-500 mb-0.5">Back</p>
                    <p class="text-sm text-slate-300 whitespace-pre-wrap">{{ card.content.back }}</p>
                  </div>
                </div>
                <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition shrink-0">
                  <button
                    class="text-slate-500 hover:text-sky-400 text-xs"
                    @click="startEdit(card.id, card.content.front, card.content.back)"
                  >
                    Edit
                  </button>
                  <button
                    class="text-slate-500 hover:text-red-400 text-xs"
                    @click="onDelete(card.id)"
                  >
                    ✕
                  </button>
                </div>
              </div>
            </template>
          </li>
        </ul>
      </template>
    </main>
  </div>
</template>
