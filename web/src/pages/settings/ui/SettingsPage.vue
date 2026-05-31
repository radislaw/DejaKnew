<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/features/auth'
import { updateProfileRequest, changePasswordRequest, deleteAccountRequest } from '@/features/profile'
import { AppLayout } from '@/widgets/app-layout'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()

// Profile section
const email = ref(auth.user?.email ?? '')
const locale = ref(auth.user?.locale ?? 'en')
const profileLoading = ref(false)
const profileError = ref<string | null>(null)
const profileSuccess = ref(false)

// Password section
const currentPassword = ref('')
const newPassword = ref('')
const passwordLoading = ref(false)
const passwordError = ref<string | null>(null)
const passwordSuccess = ref(false)

// Delete section
const deleteLoading = ref(false)
const deleteError = ref<string | null>(null)
const confirmDelete = ref(false)

const isLocalUser = auth.user?.provider === 'local'

async function saveProfile() {
  profileError.value = null
  profileSuccess.value = false
  profileLoading.value = true
  try {
    const updated = await updateProfileRequest(isLocalUser ? email.value : '', locale.value)
    auth.updateUser(updated)
    profileSuccess.value = true
  } catch (e: unknown) {
    profileError.value = e instanceof Error ? e.message : t('settings.saveFailed')
  } finally {
    profileLoading.value = false
  }
}

async function changePassword() {
  passwordError.value = null
  passwordSuccess.value = false
  passwordLoading.value = true
  try {
    await changePasswordRequest(currentPassword.value, newPassword.value)
    currentPassword.value = ''
    newPassword.value = ''
    passwordSuccess.value = true
  } catch (e: unknown) {
    passwordError.value = e instanceof Error ? e.message : t('settings.passwordFailed')
  } finally {
    passwordLoading.value = false
  }
}

async function deleteAccount() {
  deleteError.value = null
  deleteLoading.value = true
  try {
    await deleteAccountRequest()
    auth.logout()
    router.push({ name: 'home' })
  } catch (e: unknown) {
    deleteError.value = e instanceof Error ? e.message : t('settings.deleteFailed')
    deleteLoading.value = false
  }
}
</script>

<template>
  <AppLayout>
    <main class="mx-auto max-w-2xl space-y-8">
      <h1 class="text-2xl font-semibold">{{ t('settings.title') }}</h1>

      <!-- Profile section -->
      <section class="rounded-xl border border-slate-700 p-6 space-y-4">
        <h2 class="text-lg font-medium">{{ t('settings.profile') }}</h2>

        <div v-if="!isLocalUser" class="text-sm text-slate-400">
          {{ t('settings.oauthProvider', { provider: auth.user?.provider }) }}
        </div>

        <div v-if="isLocalUser" class="space-y-2">
          <label class="block text-sm text-slate-400">{{ t('settings.email') }}</label>
          <input
            v-model="email"
            type="email"
            required
            class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div class="space-y-2">
          <label class="block text-sm text-slate-400">{{ t('settings.language') }}</label>
          <select
            v-model="locale"
            class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="en">English</option>
            <option value="ru">Русский</option>
          </select>
        </div>

        <p v-if="profileError" class="text-sm text-red-400">{{ profileError }}</p>
        <p v-if="profileSuccess" class="text-sm text-green-400">{{ t('settings.saveSuccess') }}</p>

        <button
          :disabled="profileLoading"
          class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:opacity-50"
          @click="saveProfile"
        >
          {{ profileLoading ? t('common.loading') : t('common.save') }}
        </button>
      </section>

      <!-- Change password section (local users only) -->
      <section v-if="isLocalUser" class="rounded-xl border border-slate-700 p-6 space-y-4">
        <h2 class="text-lg font-medium">{{ t('settings.changePassword') }}</h2>

        <div class="space-y-2">
          <label class="block text-sm text-slate-400">{{ t('settings.currentPassword') }}</label>
          <input
            v-model="currentPassword"
            type="password"
            minlength="8"
            required
            class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <div class="space-y-2">
          <label class="block text-sm text-slate-400">{{ t('settings.newPassword') }}</label>
          <input
            v-model="newPassword"
            type="password"
            minlength="8"
            required
            class="w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>

        <p v-if="passwordError" class="text-sm text-red-400">{{ passwordError }}</p>
        <p v-if="passwordSuccess" class="text-sm text-green-400">{{ t('settings.passwordSuccess') }}</p>

        <button
          :disabled="passwordLoading || !currentPassword || !newPassword"
          class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:opacity-50"
          @click="changePassword"
        >
          {{ passwordLoading ? t('common.loading') : t('settings.changePassword') }}
        </button>
      </section>

      <!-- Danger zone -->
      <section class="rounded-xl border border-red-900 p-6 space-y-4">
        <h2 class="text-lg font-medium text-red-400">{{ t('settings.dangerZone') }}</h2>
        <p class="text-sm text-slate-400">{{ t('settings.deleteWarning') }}</p>

        <div v-if="!confirmDelete">
          <button
            class="rounded-lg border border-red-700 px-4 py-2 text-sm font-medium text-red-400 hover:bg-red-900/30"
            @click="confirmDelete = true"
          >
            {{ t('settings.deleteAccount') }}
          </button>
        </div>

        <div v-else class="space-y-3">
          <p class="text-sm font-medium text-red-400">{{ t('settings.deleteConfirm') }}</p>
          <div class="flex gap-3">
            <button
              :disabled="deleteLoading"
              class="rounded-lg bg-red-700 px-4 py-2 text-sm font-medium text-white hover:bg-red-600 disabled:opacity-50"
              @click="deleteAccount"
            >
              {{ deleteLoading ? t('common.loading') : t('settings.deleteConfirmButton') }}
            </button>
            <button
              class="rounded-lg px-4 py-2 text-sm text-slate-400 hover:text-slate-200"
              @click="confirmDelete = false"
            >
              {{ t('common.cancel') }}
            </button>
          </div>
          <p v-if="deleteError" class="text-sm text-red-400">{{ deleteError }}</p>
        </div>
      </section>
    </main>
  </AppLayout>
</template>
