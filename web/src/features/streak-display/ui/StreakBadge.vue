<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { StudyStats } from '@/entities/study-stats'

defineProps<{ stats: StudyStats }>()

const { t } = useI18n()
</script>

<template>
  <div class="flex flex-wrap gap-4 rounded-xl bg-slate-800/50 px-4 py-3 text-sm">
    <div class="flex items-center gap-1.5">
      <span class="text-lg">🔥</span>
      <span v-if="stats.streak.current > 0" class="font-semibold text-orange-400">
        {{ stats.streak.current }}
      </span>
      <span v-else class="text-slate-500">0</span>
      <span class="text-slate-400">{{ t('stats.streak') }}</span>
      <span v-if="stats.streak.longest > 0" class="ml-1 text-xs text-slate-500">
        ({{ t('stats.best') }}: {{ stats.streak.longest }})
      </span>
    </div>

    <div class="h-auto w-px bg-slate-700" />

    <div class="flex items-center gap-1.5">
      <span class="font-semibold text-sky-400">{{ stats.today.reviews }}</span>
      <span class="text-slate-400">{{ t('stats.todayReviews') }}</span>
    </div>

    <div class="h-auto w-px bg-slate-700" />

    <div class="flex items-center gap-1.5">
      <span class="font-semibold text-slate-300">{{ stats.total.reviews }}</span>
      <span class="text-slate-400">{{ t('stats.totalReviews') }}</span>
      <span class="text-slate-600">·</span>
      <span class="font-semibold text-slate-300">{{ stats.total.cardsStudied }}</span>
      <span class="text-slate-400">{{ t('stats.cardsStudied') }}</span>
    </div>
  </div>
</template>
