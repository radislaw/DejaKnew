<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDeckStore } from '@/entities/deck'
import { useCardStore } from '@/entities/card'

const { t } = useI18n()
const router = useRouter()
const deckStore = useDeckStore()
const cardStore = useCardStore()

type Step = 'welcome' | 'deck' | 'card' | 'done'
const step = ref<Step>('welcome')

const deckTitle = ref('')
const deckError = ref('')
const cardFront = ref('')
const cardBack = ref('')
const cardError = ref('')
const cardAdded = ref(false)
const createdDeckId = ref<string | null>(null)
const loading = ref(false)

const stepNumber: Record<Step, number> = { welcome: 0, deck: 1, card: 2, done: 3 }

async function createDeck() {
  if (!deckTitle.value.trim()) return
  loading.value = true
  deckError.value = ''
  try {
    await deckStore.addDeck({ title: deckTitle.value.trim() })
    createdDeckId.value = deckStore.decks[0]?.id ?? null
    step.value = 'card'
  } catch {
    deckError.value = t('onboarding.deckError')
  } finally {
    loading.value = false
  }
}

async function addCard() {
  if (!cardFront.value.trim() || !cardBack.value.trim() || !createdDeckId.value) return
  loading.value = true
  cardError.value = ''
  try {
    await cardStore.addCard(createdDeckId.value, {
      front: cardFront.value.trim(),
      back: cardBack.value.trim(),
    })
    cardAdded.value = true
    finish()
  } catch {
    cardError.value = t('onboarding.cardError')
  } finally {
    loading.value = false
  }
}

function finish() {
  localStorage.setItem('ob_done', '1')
  step.value = 'done'
}

function skip() {
  localStorage.setItem('ob_done', '1')
  router.push({ name: 'home' })
}

function goNext() {
  if (createdDeckId.value && cardAdded.value) {
    router.push({ name: 'study', params: { id: createdDeckId.value } })
  } else if (createdDeckId.value) {
    router.push({ name: 'deck', params: { id: createdDeckId.value } })
  } else {
    router.push({ name: 'home' })
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-950 px-4">
    <div class="w-full max-w-md">

      <!-- progress dots -->
      <div v-if="step !== 'welcome'" class="mb-8 flex justify-center gap-2">
        <div
          v-for="n in 3"
          :key="n"
          class="h-1.5 w-8 rounded-full transition-colors duration-300"
          :class="n <= stepNumber[step] ? 'bg-sky-500' : 'bg-slate-700'"
        />
      </div>

      <!-- step: welcome -->
      <div v-if="step === 'welcome'" class="text-center">
        <div class="mb-6 text-5xl font-black tracking-tight text-white">
          Déjà<span class="text-sky-400">Knew</span>
        </div>
        <p class="mb-2 text-lg font-semibold text-white">{{ t('onboarding.welcomeTitle') }}</p>
        <p class="mb-10 text-sm leading-relaxed text-slate-400">{{ t('onboarding.welcomeSubtitle') }}</p>
        <button
          class="w-full rounded-xl bg-sky-600 py-3 text-sm font-semibold text-white transition hover:bg-sky-500 active:scale-95"
          @click="step = 'deck'"
        >
          {{ t('onboarding.welcomeCta') }}
        </button>
        <button
          class="mt-4 w-full py-2 text-sm text-slate-500 transition hover:text-slate-300"
          @click="skip"
        >
          {{ t('onboarding.welcomeSkip') }}
        </button>
      </div>

      <!-- step: create deck -->
      <div v-else-if="step === 'deck'">
        <p class="mb-1 text-xs font-medium uppercase tracking-widest text-slate-500">
          {{ t('onboarding.step1of3') }}
        </p>
        <h2 class="mb-1 text-xl font-bold text-white">{{ t('onboarding.deckTitle') }}</h2>
        <p class="mb-6 text-sm text-slate-400">{{ t('onboarding.deckSubtitle') }}</p>
        <input
          v-model="deckTitle"
          type="text"
          :placeholder="t('onboarding.deckPlaceholder')"
          class="mb-3 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder-slate-500 outline-none focus:border-sky-500"
          @keydown.enter="createDeck"
        />
        <p v-if="deckError" class="mb-3 text-sm text-red-400">{{ deckError }}</p>
        <button
          :disabled="!deckTitle.trim() || loading"
          class="mb-3 w-full rounded-xl bg-sky-600 py-3 text-sm font-semibold text-white transition hover:bg-sky-500 active:scale-95 disabled:cursor-not-allowed disabled:opacity-40"
          @click="createDeck"
        >
          {{ loading ? t('common.loading') : t('onboarding.deckCta') }}
        </button>
        <button
          class="w-full py-2 text-sm text-slate-500 transition hover:text-slate-300"
          @click="skip"
        >
          {{ t('onboarding.deckSkip') }}
        </button>
      </div>

      <!-- step: add card -->
      <div v-else-if="step === 'card'">
        <p class="mb-1 text-xs font-medium uppercase tracking-widest text-slate-500">
          {{ t('onboarding.step2of3') }}
        </p>
        <h2 class="mb-1 text-xl font-bold text-white">{{ t('onboarding.cardTitle') }}</h2>
        <p class="mb-6 text-sm text-slate-400">{{ t('onboarding.cardSubtitle') }}</p>
        <input
          v-model="cardFront"
          type="text"
          :placeholder="t('onboarding.cardFrontPlaceholder')"
          class="mb-3 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder-slate-500 outline-none focus:border-sky-500"
        />
        <input
          v-model="cardBack"
          type="text"
          :placeholder="t('onboarding.cardBackPlaceholder')"
          class="mb-3 w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder-slate-500 outline-none focus:border-sky-500"
          @keydown.enter="addCard"
        />
        <p v-if="cardError" class="mb-3 text-sm text-red-400">{{ cardError }}</p>
        <button
          :disabled="!cardFront.trim() || !cardBack.trim() || loading"
          class="mb-3 w-full rounded-xl bg-sky-600 py-3 text-sm font-semibold text-white transition hover:bg-sky-500 active:scale-95 disabled:cursor-not-allowed disabled:opacity-40"
          @click="addCard"
        >
          {{ loading ? t('common.loading') : t('onboarding.cardCta') }}
        </button>
        <button
          class="w-full py-2 text-sm text-slate-500 transition hover:text-slate-300"
          @click="finish"
        >
          {{ t('onboarding.cardSkip') }}
        </button>
      </div>

      <!-- step: done -->
      <div v-else-if="step === 'done'" class="text-center">
        <div class="mb-4 text-5xl">🎉</div>
        <h2 class="mb-2 text-xl font-bold text-white">{{ t('onboarding.doneTitle') }}</h2>
        <p class="mb-10 text-sm leading-relaxed text-slate-400">{{ t('onboarding.doneSubtitle') }}</p>
        <button
          class="w-full rounded-xl bg-sky-600 py-3 text-sm font-semibold text-white transition hover:bg-sky-500 active:scale-95"
          @click="goNext"
        >
          {{ cardAdded ? t('onboarding.doneCta') : t('onboarding.doneSkipCta') }}
        </button>
      </div>

    </div>
  </div>
</template>
