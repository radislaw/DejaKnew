<script setup lang="ts">
import { useAuthStore } from '@/features/auth'
import { HomePage } from '@/pages/home'
import { LandingPage } from '@/pages/landing'

const auth = useAuthStore()
</script>

<template>
  <LandingPage v-if="!auth.isAuthenticated" />
  <HomePage v-else />
</template>
