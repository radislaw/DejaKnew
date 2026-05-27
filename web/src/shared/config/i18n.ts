import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'

export type MessageSchema = typeof en
export type Locale = 'en' | 'ru'

const SUPPORTED: Locale[] = ['en', 'ru']

function detectLocale(): Locale {
  const saved = localStorage.getItem('locale')
  if (saved === 'ru' || saved === 'en') return saved
  const browser = navigator.language.split('-')[0]
  return (SUPPORTED as string[]).includes(browser) ? (browser as Locale) : 'en'
}

export const i18n = createI18n<[MessageSchema], Locale>({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: { en, ru },
  pluralRules: {
    ru: (choice: number): number => {
      if (choice === 0) return 2
      const teen = choice > 10 && choice < 20
      const end = choice % 10
      if (!teen && end === 1) return 0
      if (!teen && end >= 2 && end <= 4) return 1
      return 2
    },
  },
})
