<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/features/auth'
import { useDeckStore } from '@/entities/deck'
import { DeckCreateForm } from '@/features/deck-create'
import { DeckList } from '@/widgets/deck-list'

const router = useRouter()
const auth = useAuthStore()
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
    createError.value = e instanceof Error ? e.message : 'Failed to create deck'
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
  <div class="min-h-screen bg-[#070d1c] px-4 py-6 text-slate-100">
    <header class="mx-auto mb-8 flex max-w-4xl items-center justify-between">
      <h1 class="text-xl font-semibold tracking-tight">FlashCards</h1>
      <div class="flex items-center gap-4">
        <span class="text-sm text-slate-400">{{ auth.user?.email }}</span>
        <button
          class="rounded-lg px-3 py-1.5 text-sm text-slate-400 hover:text-slate-200"
          @click="auth.logout(); router.push('/login')"
        >
          Sign out
        </button>
      </div>
    </header>

    <main class="mx-auto max-w-4xl">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-medium">My decks</h2>
        <button
          v-if="!showForm"
          class="rounded-lg bg-sky-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-sky-500"
          @click="showForm = true"
        >
          + New deck
        </button>
      </div>

      <DeckCreateForm
        v-if="showForm"
        class="mb-4"
        @submit="onCreate"
        @cancel="showForm = false"
      />

      <p v-if="createError" class="mb-3 text-sm text-red-400">{{ createError }}</p>

      <p v-if="decks.loading" class="text-sm text-slate-400">Loading…</p>
      <p v-else-if="decks.error" class="text-sm text-red-400">{{ decks.error }}</p>
      <p v-else-if="decks.decks.length === 0" class="text-sm text-slate-500">
        No decks yet. Create your first one.
      </p>
      <DeckList v-else :decks="decks.decks" @delete="onDelete" />
    </main>
  </div>
</template>
