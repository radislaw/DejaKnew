import { createRouter, createWebHistory } from 'vue-router'
import { RootPage } from '@/pages/root'
import { LoginPage } from '@/pages/login'
import { RegisterPage } from '@/pages/register'
import { DeckPage } from '@/pages/deck'
import { StudyPage } from '@/pages/study'
import { OnboardingPage } from '@/pages/onboarding'
import { AuthCallbackPage } from '@/pages/auth-callback'
import { SettingsPage } from '@/pages/settings'
import { RoadmapPage, RoadmapDevPage } from '@/pages/roadmap'
import { useAuthStore } from '@/features/auth'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
  }
}

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: RootPage },
    { path: '/decks/:id', name: 'deck', component: DeckPage, meta: { requiresAuth: true } },
    { path: '/study/:id', name: 'study', component: StudyPage, meta: { requiresAuth: true } },
    { path: '/login', name: 'login', component: LoginPage },
    { path: '/register', name: 'register', component: RegisterPage },
    { path: '/onboarding', name: 'onboarding', component: OnboardingPage },
    { path: '/auth/callback', name: 'auth-callback', component: AuthCallbackPage },
    { path: '/settings', name: 'settings', component: SettingsPage, meta: { requiresAuth: true } },
    { path: '/roadmap', name: 'roadmap', component: RoadmapPage },
    { path: '/roadmap/dev', name: 'roadmap-dev', component: RoadmapDevPage },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const onboardingDone = localStorage.getItem('ob_done')

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'home' }
  }
  if ((to.name === 'login' || to.name === 'register') && auth.isAuthenticated) {
    return { name: 'home' }
  }
  if (auth.isAuthenticated && !onboardingDone && to.name !== 'onboarding') {
    return { name: 'onboarding' }
  }
  if (to.name === 'onboarding' && !auth.isAuthenticated) {
    return { name: 'home' }
  }
  if (to.name === 'onboarding' && onboardingDone) {
    return { name: 'home' }
  }
})
