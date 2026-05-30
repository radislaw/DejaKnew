import { onMounted, ref } from 'vue'

const CIS = new Set(['RU', 'BY', 'KZ', 'UZ', 'AZ', 'AM', 'GE', 'KG', 'MD', 'TJ', 'TM'])
const CACHE_KEY = 'geo_country'
const CACHE_TTL = 24 * 60 * 60 * 1000

async function detectCountry(): Promise<string> {
  const raw = localStorage.getItem(CACHE_KEY)
  if (raw) {
    const { code, ts } = JSON.parse(raw) as { code: string; ts: number }
    if (Date.now() - ts < CACHE_TTL) return code
  }
  try {
    const res = await fetch('http://ip-api.com/json?fields=countryCode')
    const { countryCode } = (await res.json()) as { countryCode: string }
    localStorage.setItem(CACHE_KEY, JSON.stringify({ code: countryCode, ts: Date.now() }))
    return countryCode
  } catch {
    return ''
  }
}

export function useGeoOAuth() {
  const showYandex = ref(false)

  onMounted(async () => {
    const country = await detectCountry()
    showYandex.value = CIS.has(country)
  })

  return { showYandex }
}
