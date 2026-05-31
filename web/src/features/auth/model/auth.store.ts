import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { setAuthToken } from '@/shared/api'
import type { User } from '@/entities/user'
import { loginRequest, meRequest, registerRequest } from '../api/auth.api'

const TOKEN_KEY = 'flashcards.token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const user = ref<User | null>(null)

  const isAuthenticated = computed(() => token.value !== null)

  // Sync the http client with the token restored from localStorage.
  setAuthToken(token.value)

  function applyToken(value: string | null): void {
    token.value = value
    setAuthToken(value)
    if (value) localStorage.setItem(TOKEN_KEY, value)
    else localStorage.removeItem(TOKEN_KEY)
  }

  async function login(email: string, password: string): Promise<void> {
    const res = await loginRequest(email, password)
    applyToken(res.token)
    user.value = res.user
  }

  async function register(email: string, password: string): Promise<void> {
    const locale = localStorage.getItem('locale') ?? 'en'
    const res = await registerRequest(email, password, locale)
    applyToken(res.token)
    user.value = res.user
  }

  // Called on app start: if a token was restored, fetch the user behind it.
  async function restore(): Promise<void> {
    if (!token.value) return
    try {
      user.value = await meRequest()
    } catch {
      applyToken(null)
    }
  }

  async function loginWithToken(value: string): Promise<void> {
    applyToken(value)
    user.value = await meRequest()
  }

  function logout(): void {
    applyToken(null)
    user.value = null
  }

  function updateUser(updated: User): void {
    user.value = updated
  }

  return { token, user, isAuthenticated, login, register, restore, loginWithToken, logout, updateUser }
})
