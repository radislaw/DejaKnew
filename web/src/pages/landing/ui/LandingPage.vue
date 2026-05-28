<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { LocaleSwitcher } from '@/shared/ui'

const router = useRouter()
const { t, tm, locale } = useI18n()
const cycleWords = computed(() => tm('landing.cycleWords') as string[])

watch(locale, () => {
  if (cycleTimer) { clearInterval(cycleTimer); cycleTimer = null }
  flipTimers.splice(0).forEach(clearTimeout)
  startCycleWord()
}, { flush: 'post' })

const memoLayerRef = ref<HTMLDivElement | null>(null)
const demoCardRef = ref<HTMLDivElement | null>(null)
const cycleRef = ref<HTMLSpanElement | null>(null)
const counterRef = ref<HTMLSpanElement | null>(null)

const demoFrontText = ref('serendipity')
const demoIpaText = ref('/ˌsɛr.ənˈdɪp.ɪ.ti/')
const demoExampleEn = ref('It was pure <em>serendipity</em> that we met that day.')
const demoBackText = ref('счастливая случайность')
const demoExampleRu = ref('Наша встреча в тот день — чистая <em>случайность</em>.')
const isFlipped = ref(false)

const DECK = [
  {
    en: 'serendipity', ipa: '/ˌsɛr.ənˈdɪp.ɪ.ti/',
    ru: 'счастливая случайность',
    exEn: 'It was pure <em>serendipity</em> that we met that day.',
    exRu: 'Наша встреча в тот день — чистая <em>случайность</em>.'
  },
  {
    en: 'ephemeral', ipa: '/ɪˈfɛm.ər.əl/',
    ru: 'мимолётный',
    exEn: 'Cherry blossoms are famously <em>ephemeral</em>.',
    exRu: 'Цветение сакуры знаменито своей <em>мимолётностью</em>.'
  },
  {
    en: 'petrichor', ipa: '/ˈpɛt.rɪ.kɔːr/',
    ru: 'запах земли после дождя',
    exEn: 'The <em>petrichor</em> rose from the warm pavement.',
    exRu: '<em>Запах земли</em> поднимался от тёплого асфальта.'
  },
  {
    en: 'sonder', ipa: '/ˈsɒn.dər/',
    ru: 'осознание чужих жизней',
    exEn: 'In the crowded train I felt a wave of <em>sonder</em>.',
    exRu: 'В переполненном поезде накрыло <em>ощущением чужих жизней</em>.'
  },
  {
    en: 'limerence', ipa: '/ˈlɪm.ər.əns/',
    ru: 'влюблённая одержимость',
    exEn: 'She mistook <em>limerence</em> for lasting love.',
    exRu: 'Она приняла <em>одержимость</em> за настоящую любовь.'
  },
]

let demoIdx = 0
let demoTimer: ReturnType<typeof setInterval> | null = null
let cycleTimer: ReturnType<typeof setInterval> | null = null
const flipTimers: ReturnType<typeof setTimeout>[] = []
const cleanups: (() => void)[] = []

function setDeckContent(d: typeof DECK[0]) {
  demoFrontText.value = d.en
  demoIpaText.value = d.ipa
  demoExampleEn.value = d.exEn
  demoBackText.value = d.ru
  demoExampleRu.value = d.exRu
}

function flip() {
  isFlipped.value = !isFlipped.value
  if (!isFlipped.value) {
    demoIdx = (demoIdx + 1) % DECK.length
    setDeckContent(DECK[demoIdx])
  }
}

function handleDemoClick() {
  flip()
  if (demoTimer) clearInterval(demoTimer)
  demoTimer = setInterval(flip, 3200)
}

function handleKeydown(e: KeyboardEvent) {
  if (e.code === 'Space' && document.activeElement === document.body) {
    e.preventDefault()
    handleDemoClick()
  }
}

onMounted(() => {
  setDeckContent(DECK[0])

  demoTimer = setInterval(flip, 3200)

  document.addEventListener('keydown', handleKeydown)
  cleanups.push(() => document.removeEventListener('keydown', handleKeydown))

  // Cursor-follow glow on feature cards
  const cards = document.querySelectorAll<HTMLElement>('.lp-card')
  cards.forEach(card => {
    const fn = (e: PointerEvent) => {
      const r = card.getBoundingClientRect()
      card.style.setProperty('--mx', (e.clientX - r.left) + 'px')
      card.style.setProperty('--my', (e.clientY - r.top) + 'px')
    }
    card.addEventListener('pointermove', fn)
    cleanups.push(() => card.removeEventListener('pointermove', fn))
  })

  // Scroll parallax on ghost field
  let ticking = false
  const onScroll = () => {
    if (ticking) return
    ticking = true
    requestAnimationFrame(() => {
      document.documentElement.style.setProperty('--net-y', (window.scrollY * -0.12) + 'px')
      ticking = false
    })
  }
  window.addEventListener('scroll', onScroll, { passive: true })
  cleanups.push(() => window.removeEventListener('scroll', onScroll))

  buildMemo()
  startCycleWord()

  // Forgetting curve draw on scroll into view
  const curveEl = document.querySelector<SVGElement>('.lp-curve')
  if (curveEl) {
    const obs = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) { curveEl.classList.add('is-in'); obs.unobserve(curveEl) }
      })
    }, { threshold: 0.4 })
    obs.observe(curveEl)
    cleanups.push(() => obs.disconnect())
  }

  // Counter count-up animation
  const counterEl = counterRef.value
  if (counterEl) {
    const target = 47823
    const fmt = (n: number) => n.toLocaleString('ru-RU').replace(/,/g, ' ')
    const countObs = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (!e.isIntersecting) return
        countObs.unobserve(e.target)
        const start = performance.now()
        const dur = 1600
        function tick(now: number) {
          const t = Math.min((now - start) / dur, 1)
          const eased = 1 - Math.pow(1 - t, 3)
          if (counterEl) counterEl.textContent = fmt(Math.round(target * eased))
          if (t < 1) requestAnimationFrame(tick)
        }
        requestAnimationFrame(tick)
      })
    }, { threshold: 0.5 })
    countObs.observe(counterEl)
    cleanups.push(() => countObs.disconnect())
  }
})

function startCycleWord() {
  document.fonts.ready.then(() => {
    const cycleEl = cycleRef.value
    if (!cycleEl) return
    const wordEls = Array.from(cycleEl.querySelectorAll<HTMLElement>('.cycle-word'))
    const widths = wordEls.map(w => Math.ceil(w.getBoundingClientRect().width) + 4)
    cycleEl.style.width = widths[0] + 'px'

    let i = 0
    cycleTimer = setInterval(() => {
      const cur = wordEls[i]
      i = (i + 1) % wordEls.length
      const next = wordEls[i]
      cur.classList.remove('is-active')
      cur.classList.add('is-exit')
      next.classList.add('is-active')
      cycleEl.style.width = widths[i] + 'px'
      const t = setTimeout(() => cur.classList.remove('is-exit'), 600)
      flipTimers.push(t)
    }, 2400)
  })
}

function buildMemo() {
  const layer = memoLayerRef.value
  if (!layer) return

  const VOCAB = [
    { tag: 'en→ru', w: 'serendipity', t: 'случайность', wc: '', tc: 'ru' },
    { tag: 'en→ru', w: 'ephemeral', t: 'мимолётный', wc: '', tc: 'ru' },
    { tag: 'en→ru', w: 'petrichor', t: 'запах дождя', wc: '', tc: 'ru' },
    { tag: 'ru→en', w: 'тоска', t: 'longing', wc: 'ru', tc: '' },
    { tag: 'pt→ru', w: 'saudade', t: 'тоска', wc: '', tc: 'ru' },
    { tag: 'da→ru', w: 'hygge', t: 'уют', wc: '', tc: 'ru' },
    { tag: 'jp→en', w: '木漏れ日', t: 'sunlight·leaves', wc: 'cjk', tc: '' },
    { tag: 'jp→ru', w: '記憶', t: 'память', wc: 'cjk', tc: 'ru' },
    { tag: 'jp→en', w: '学', t: 'learn', wc: 'cjk', tc: '' },
    { tag: 'fr→ru', w: "l'esprit", t: 'ум', wc: '', tc: 'ru' },
    { tag: 'fr→en', w: 'flâneur', t: 'wanderer', wc: '', tc: '' },
    { tag: 'de→ru', w: 'Erinnerung', t: 'воспоминание', wc: '', tc: 'ru' },
    { tag: 'de→en', w: 'Sehnsucht', t: 'yearning', wc: '', tc: '' },
    { tag: 'es→ru', w: 'mañana', t: 'завтра', wc: '', tc: 'ru' },
    { tag: 'gr→ru', w: 'mneme', t: 'память', wc: '', tc: 'ru' },
    { tag: 'en→ru', w: 'limerence', t: 'одержимость', wc: '', tc: 'ru' },
    { tag: 'en→ru', w: 'sonder', t: 'чужие жизни', wc: '', tc: 'ru' },
    { tag: 'jp→ru', w: '夢', t: 'сон', wc: 'cjk', tc: 'ru' },
    { tag: 'la→en', w: 'memento', t: 'remember', wc: '', tc: '' },
    { tag: 'fr→ru', w: 'déjà vu', t: 'дежавю', wc: '', tc: 'ru' },
  ]

  const LAYOUT = [
    { x: 4, y: 6, rot: -7, z: 0 }, { x: 18, y: 14, rot: 4, z: 1 },
    { x: 82, y: 4, rot: 6, z: 0 }, { x: 72, y: 13, rot: -3, z: 1 },
    { x: 2, y: 28, rot: 3, z: 1 }, { x: 88, y: 26, rot: -5, z: 0 },
    { x: 12, y: 44, rot: -4, z: 1 }, { x: 84, y: 42, rot: 5, z: 1 },
    { x: 3, y: 58, rot: 6, z: 0 }, { x: 78, y: 56, rot: -6, z: 0 },
    { x: 16, y: 72, rot: -3, z: 1 }, { x: 86, y: 70, rot: 4, z: 1 },
    { x: 5, y: 86, rot: 5, z: 0 }, { x: 70, y: 84, rot: -4, z: 0 },
    { x: 30, y: 92, rot: 3, z: 1 }, { x: 55, y: 96, rot: -5, z: 1 },
  ]

  const pool = VOCAB.slice().sort(() => Math.random() - 0.5)
  const ghosts: HTMLElement[] = []

  LAYOUT.forEach((pos, i) => {
    const v = pool[i % pool.length]
    const g = document.createElement('div')
    g.className = 'ghost'
    g.style.left = pos.x + '%'
    g.style.top = pos.y + '%'
    g.style.setProperty('--base', `translate(-50%, -50%) rotate(${pos.rot}deg) scale(${pos.z === 1 ? 1 : 0.78})`)
    g.style.opacity = String(pos.z === 1 ? 0.55 : 0.32)
    g.style.setProperty('--float-dur', (14 + Math.random() * 10).toFixed(2) + 's')
    g.style.setProperty('--float-del', (-Math.random() * 14).toFixed(2) + 's')
    g.innerHTML =
      `<div class="ghost-face front"><div class="ghost-meta">${v.tag}</div><div class="ghost-word ${v.wc}">${v.w}</div></div>` +
      `<div class="ghost-face back"><div class="ghost-meta">${t('landing.demoTranslation')}</div><div class="ghost-word ${v.tc}">${v.t}</div></div>`
    layer.appendChild(g)
    ghosts.push(g)
  })

  const initFlips = Math.floor(ghosts.length * 0.35)
  ghosts.map((_, i) => i).sort(() => Math.random() - 0.5).slice(0, initFlips)
    .forEach(i => ghosts[i].classList.add('is-flipped'))

  let lastFlipped = -1
  function flipOne() {
    let i = Math.floor(Math.random() * ghosts.length)
    if (i === lastFlipped) i = (i + 1) % ghosts.length
    lastFlipped = i
    const g = ghosts[i]
    g.classList.add('is-active')
    const t1 = setTimeout(() => {
      g.classList.toggle('is-flipped')
      const t2 = setTimeout(() => g.classList.remove('is-active'), 1400)
      flipTimers.push(t2)
    }, 280)
    flipTimers.push(t1)
    const t3 = setTimeout(flipOne, 1400 + Math.random() * 1600)
    flipTimers.push(t3)
  }
  const t0 = setTimeout(flipOne, 1200)
  flipTimers.push(t0)
}

onUnmounted(() => {
  if (demoTimer) clearInterval(demoTimer)
  if (cycleTimer) clearInterval(cycleTimer)
  flipTimers.forEach(clearTimeout)
  cleanups.forEach(fn => fn())
  document.documentElement.style.removeProperty('--net-y')
})
</script>

<template>
  <div class="lp-root">
    <div class="memo" aria-hidden="true">
      <div class="memo-layer" ref="memoLayerRef"></div>
    </div>

    <div class="lp-page">
      <!-- Nav -->
      <nav class="lp-nav">
        <a href="/" class="lp-logo">
          <span class="lp-logo-mark">D</span>
          <span>DéjàKnew</span>
        </a>
        <div class="lp-nav-actions">
          <button class="lp-btn lp-btn-ghost" @click="router.push('/login')">{{ t('landing.signIn') }}</button>
          <button class="lp-btn lp-btn-primary" @click="router.push('/register')">{{ t('landing.start') }}</button>
          <LocaleSwitcher />
        </div>
      </nav>

      <!-- Hero -->
      <header class="lp-hero">
        <span class="lp-eyebrow"><span class="lp-dot"></span> {{ t('landing.eyebrow') }}</span>
        <h1>{{ t('landing.heroTitlePre') }} <span class="lp-cycle" ref="cycleRef"
          ><span
            v-for="(word, i) in cycleWords"
            :key="word"
            class="cycle-word"
            :class="{ 'is-active': i === 0 }"
          >{{ word }}</span
        ></span>{{ t('landing.heroTitlePost') }}</h1>
        <p>{{ t('landing.heroSubtitle') }}</p>
        <div class="lp-hero-actions">
          <button class="lp-btn lp-btn-primary lp-btn-lg" @click="router.push('/register')">
            <span>{{ t('landing.getStarted') }}</span>
          </button>
          <button class="lp-btn lp-btn-outline lp-btn-lg" @click="router.push('/login')">{{ t('landing.signIn') }}</button>
        </div>

        <div class="lp-demo" :aria-label="t('landing.demoAriaLabel')">
          <div class="lp-demo-stack">
            <div
              class="lp-demo-card"
              :class="{ 'is-flipped': isFlipped }"
              ref="demoCardRef"
              @click="handleDemoClick"
            >
              <div class="lp-demo-face front">
                <div class="lp-demo-top">
                  <span>card&nbsp;12&nbsp;/&nbsp;47</span>
                  <span class="lp-pill"><span class="lp-swatch"></span>english&nbsp;·&nbsp;b2</span>
                </div>
                <div class="lp-demo-body">
                  <div class="lp-demo-word">{{ demoFrontText }}</div>
                  <div class="lp-demo-ipa">{{ demoIpaText }}</div>
                  <div class="lp-demo-example" v-html="demoExampleEn"></div>
                </div>
                <div class="lp-demo-reveal">
                  <span>{{ t('landing.demoReveal') }}</span>
                  <kbd>Space</kbd>
                </div>
              </div>
              <div class="lp-demo-face back">
                <div class="lp-demo-top">
                  <span>card&nbsp;12&nbsp;/&nbsp;47</span>
                  <span class="lp-pill"><span class="lp-swatch"></span>{{ t('landing.demoTranslation') }}</span>
                </div>
                <div class="lp-demo-body">
                  <div class="lp-demo-word">{{ demoBackText }}</div>
                  <div class="lp-demo-example" v-html="demoExampleRu"></div>
                </div>
                <div class="lp-demo-bottom">
                  <button class="lp-rate" data-tone="again" type="button">
                    <span class="rate-label">Again</span><span class="rate-interval">10м</span>
                  </button>
                  <button class="lp-rate" data-tone="hard" type="button">
                    <span class="rate-label">Hard</span><span class="rate-interval">1д</span>
                  </button>
                  <button class="lp-rate" data-tone="good" type="button">
                    <span class="rate-label">Good</span><span class="rate-interval">4д</span>
                  </button>
                  <button class="lp-rate" data-tone="easy" type="button">
                    <span class="rate-label">Easy</span><span class="rate-interval">9д</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div class="lp-demo-hint" v-html="t('landing.demoHint')"></div>
        </div>
      </header>

      <!-- Features -->
      <section class="lp-section">
        <div class="lp-section-head">
          <div class="lp-section-label">// features</div>
          <h2 class="lp-section-title">{{ t('landing.featuresTitle') }}</h2>
        </div>
        <div class="lp-features">
          <article class="lp-card">
            <svg class="lp-curve" viewBox="0 0 240 36" preserveAspectRatio="none" aria-hidden="true">
              <path d="M 4 8 Q 30 24 56 30 L 56 12 Q 90 22 124 28 L 124 10 Q 160 18 196 24 L 196 8 Q 218 14 236 18" />
              <circle cx="4"   cy="8"  r="2.5" />
              <circle cx="56"  cy="12" r="2.5" />
              <circle cx="124" cy="10" r="2.5" />
              <circle cx="196" cy="8"  r="2.5" />
            </svg>
            <div class="lp-card-icon">🧠</div>
            <h3 class="lp-card-title">{{ t('landing.feature1Title') }}</h3>
            <p class="lp-card-body">{{ t('landing.feature1Body') }}</p>
          </article>
          <article class="lp-card">
            <div class="lp-card-icon">📱</div>
            <h3 class="lp-card-title">{{ t('landing.feature2Title') }}</h3>
            <p class="lp-card-body">{{ t('landing.feature2Body') }}</p>
          </article>
          <article class="lp-card">
            <div class="lp-card-icon">✈️</div>
            <h3 class="lp-card-title">{{ t('landing.feature3Title') }}</h3>
            <p class="lp-card-body">{{ t('landing.feature3Body') }}</p>
          </article>
        </div>
      </section>

      <!-- Footer CTA -->
      <section class="lp-cta">
        <h2>{{ t('landing.ctaTitlePre') }}<span class="lp-counter" ref="counterRef">47&nbsp;823</span>{{ t('landing.ctaTitlePost') }}</h2>
        <button class="lp-btn lp-btn-primary lp-btn-xl" @click="router.push('/register')">
          <span>{{ t('landing.createAccountCta') }}</span>
        </button>
      </section>

      <!-- Footer -->
      <footer class="lp-foot">
        <div>© 2026 DéjàKnew</div>
        <div class="lp-foot-links">
          <a href="#">{{ t('landing.footerPolicy') }}</a>
          <a href="#">{{ t('landing.footerTerms') }}</a>
          <a href="#">{{ t('landing.footerContact') }}</a>
        </div>
      </footer>
    </div>
  </div>
</template>

<style scoped>
/* ── Reset & tokens ──────────────────────────────── */
.lp-root {
  --bg: #070d1c;
  --fg: #f1f5f9;
  --fg-dim: #94a3b8;
  --fg-mute: #64748b;
  --sky: #0284c7;
  --sky-hover: #0ea5e9;
  --card: #0f172a;
  --border: #1e293b;
  --border-strong: #334155;
  font-family: 'Inter', system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-feature-settings: 'cv11', 'ss01';
  line-height: 1.5;
  background: var(--bg);
  color: var(--fg);
  min-height: 100vh;
  overflow-x: hidden;
}

::selection { background: rgba(14, 165, 233, 0.3); color: var(--fg); }

/* ── Ghost memo field ────────────────────────────── */
.memo {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  perspective: 1400px;
  mask-image: radial-gradient(ellipse 110% 85% at 50% 45%, #000 25%, transparent 95%);
  -webkit-mask-image: radial-gradient(ellipse 110% 85% at 50% 45%, #000 25%, transparent 95%);
}
.memo::before {
  content: "";
  position: absolute;
  inset: -10%;
  background-image:
    radial-gradient(1px 1px at 12% 18%, rgba(186,230,253,.35), transparent 60%),
    radial-gradient(1px 1px at 78% 32%, rgba(196,181,253,.30), transparent 60%),
    radial-gradient(1px 1px at 34% 72%, rgba(186,230,253,.28), transparent 60%),
    radial-gradient(1px 1px at 91% 81%, rgba(186,230,253,.22), transparent 60%),
    radial-gradient(1px 1px at 56% 11%, rgba(196,181,253,.25), transparent 60%),
    radial-gradient(1px 1px at 5%  64%, rgba(186,230,253,.22), transparent 60%),
    radial-gradient(1.5px 1.5px at 67% 58%, rgba(125,211,252,.40), transparent 60%),
    radial-gradient(1px 1px at 22% 45%, rgba(186,230,253,.22), transparent 60%),
    radial-gradient(1px 1px at 88% 14%, rgba(196,181,253,.22), transparent 60%);
  transform: translateY(var(--net-y, 0px));
  will-change: transform;
}
.memo-layer {
  position: absolute;
  inset: 0;
  transform: translateY(var(--net-y, 0px));
  will-change: transform;
}

/* ghost cards created via JS — :deep() reaches them */
.memo :deep(.ghost) {
  position: absolute;
  width: 140px;
  height: 88px;
  transform-style: preserve-3d;
  transform: var(--base) rotateY(var(--tilt, 0deg));
  transition: transform 1.4s cubic-bezier(.5,.05,.2,1), opacity 1.4s ease;
  animation: ghost-float var(--float-dur, 18s) ease-in-out infinite;
  animation-delay: var(--float-del, 0s);
}
.memo :deep(.ghost.is-flipped) { --tilt: 180deg; }
@keyframes ghost-float {
  0%, 100% { translate: 0 0; }
  50%       { translate: 0 -6px; }
}
.memo :deep(.ghost-face) {
  position: absolute;
  inset: 0;
  border-radius: 9px;
  border: 1px solid rgba(125,211,252,.18);
  background: linear-gradient(180deg, rgba(15,23,42,.55), rgba(7,13,28,.35));
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
  box-shadow: 0 8px 30px -12px rgba(2,6,23,.8), inset 0 1px 0 rgba(186,230,253,.06);
}
.memo :deep(.ghost-face.back) { transform: rotateY(180deg); }
.memo :deep(.ghost-face::before) {
  content: "";
  position: absolute;
  top: 0; left: 12px; right: 12px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(125,211,252,.45), transparent);
  opacity: 0.7;
}
.memo :deep(.ghost-meta) {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 8.5px;
  color: rgba(148,163,184,.55);
  letter-spacing: .04em;
  text-transform: lowercase;
}
.memo :deep(.ghost-word) {
  font-size: 15px; font-weight: 600; letter-spacing: -.01em;
  color: rgba(241,245,249,.78); line-height: 1.1;
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.memo :deep(.ghost-word.ru)  { font-size: 13px; color: rgba(186,230,253,.80); }
.memo :deep(.ghost-word.cjk) { font-size: 22px; font-weight: 500; }
.memo :deep(.ghost.is-active .ghost-face) {
  border-color: rgba(125,211,252,.45);
  box-shadow: 0 0 24px -4px rgba(14,165,233,.35), 0 8px 30px -12px rgba(2,6,23,.8), inset 0 1px 0 rgba(186,230,253,.10);
}

/* ── Page layout ─────────────────────────────────── */
.lp-page {
  position: relative;
  z-index: 1;
  max-width: 1024px;
  margin: 0 auto;
  padding: 0 32px;
}

/* ── Nav ─────────────────────────────────────────── */
.lp-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
}
.lp-logo {
  display: inline-flex; align-items: center; gap: 10px;
  font-weight: 700; font-size: 17px; letter-spacing: -.01em;
  color: var(--fg); text-decoration: none;
}
.lp-logo-mark {
  width: 22px; height: 22px; border-radius: 6px;
  background: var(--card); border: 1px solid var(--border-strong);
  display: grid; place-items: center;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; font-weight: 600; color: var(--sky-hover);
}
.lp-nav-actions { display: flex; align-items: center; gap: 8px; }

/* ── Buttons ─────────────────────────────────────── */
.lp-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  font-family: inherit; font-weight: 500; font-size: 14px; line-height: 1;
  padding: 10px 16px; border-radius: 8px; border: 1px solid transparent;
  cursor: pointer; text-decoration: none; color: var(--fg); background: transparent;
  transition: background .15s ease, border-color .15s ease, color .15s ease, transform .05s ease;
  white-space: nowrap;
}
.lp-btn:active { transform: translateY(1px); }
.lp-btn-ghost { color: var(--fg-dim); }
.lp-btn-ghost:hover { color: var(--fg); background: var(--card); }
.lp-btn-outline { color: var(--fg); border-color: var(--border-strong); }
.lp-btn-outline:hover { border-color: #475569; background: var(--card); }
.lp-btn-primary {
  position: relative; background: var(--sky); color: #f0f9ff;
  font-weight: 600; overflow: hidden; isolation: isolate;
}
.lp-btn-primary > * { position: relative; z-index: 1; }
.lp-btn-primary::after {
  content: ""; position: absolute; top: 0; left: -120%; width: 60%; height: 100%;
  background: linear-gradient(110deg, transparent 30%, rgba(255,255,255,.22) 50%, transparent 70%);
  transition: left .7s ease; pointer-events: none;
}
.lp-btn-primary:hover { background: var(--sky-hover); }
.lp-btn-primary:hover::after { left: 130%; }
.lp-btn-lg { font-size: 15px; padding: 14px 22px; border-radius: 10px; }
.lp-btn-xl { font-size: 16px; padding: 16px 28px; border-radius: 10px; font-weight: 600; }

/* ── Hero ────────────────────────────────────────── */
.lp-hero {
  text-align: center;
  padding: 120px 0 96px;
  border-bottom: 1px solid var(--border);
}
.lp-eyebrow {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-mute);
  padding: 6px 12px; border: 1px solid var(--border); border-radius: 999px;
  margin-bottom: 32px; background: rgba(15,23,42,.5);
}
.lp-dot {
  width: 6px; height: 6px; border-radius: 999px;
  background: var(--sky-hover); position: relative;
  box-shadow: 0 0 0 0 rgba(14,165,233,.6);
  animation: pulse 2.4s ease-out infinite;
}
@keyframes pulse {
  0%   { box-shadow: 0 0 0 0   rgba(14,165,233,.5); }
  70%  { box-shadow: 0 0 0 8px rgba(14,165,233,0);  }
  100% { box-shadow: 0 0 0 0   rgba(14,165,233,0);  }
}
.lp-hero h1 {
  font-size: clamp(40px, 6vw, 60px); font-weight: 800;
  line-height: 1.05; letter-spacing: -.03em;
  margin: 0 0 20px; color: var(--fg);
}
.lp-cycle {
  display: inline-grid; vertical-align: bottom; overflow: hidden;
  line-height: 1.05; transition: width .55s cubic-bezier(.4,0,.2,1);
}
.cycle-word {
  grid-area: 1 / 1; justify-self: start; width: max-content;
  opacity: 0; transform: translateY(0.6em);
  transition: opacity .35s ease, transform .5s cubic-bezier(.4,0,.2,1);
  pointer-events: none; white-space: nowrap;
  background-image: linear-gradient(100deg, #f1f5f9 0%, #f1f5f9 35%, #7dd3fc 45%, #c4b5fd 55%, #f1f5f9 65%, #f1f5f9 100%);
  background-size: 250% 100%;
  -webkit-background-clip: text; background-clip: text;
  -webkit-text-fill-color: transparent; color: transparent;
  animation: shimmer 6s ease-in-out infinite;
}
.cycle-word.is-active  { opacity: 1; transform: translateY(0); pointer-events: auto; }
.cycle-word.is-exit    { opacity: 0; transform: translateY(-0.6em); }
@keyframes shimmer {
  0%, 100% { background-position: 130% 0; }
  50%       { background-position: -30% 0; }
}
.lp-hero p {
  font-size: 18px; line-height: 1.55; color: var(--fg-dim);
  max-width: 600px; margin: 0 auto 40px;
}
.lp-hero-actions {
  display: flex; align-items: center; justify-content: center;
  gap: 12px; flex-wrap: wrap;
}

/* ── Demo card ───────────────────────────────────── */
.lp-demo { margin: 72px auto 0; width: 380px; max-width: 100%; }
.lp-demo-stack {
  position: relative; width: 100%; aspect-ratio: 4 / 3; perspective: 1400px;
}
.lp-demo-stack::before,
.lp-demo-stack::after {
  content: ""; position: absolute; inset: 0;
  background: var(--card); border: 1px solid var(--border);
  border-radius: 16px; pointer-events: none; transform-origin: center bottom;
}
.lp-demo-stack::before { transform: translateY(10px) scale(0.96); opacity: 0.55; }
.lp-demo-stack::after  { transform: translateY(20px) scale(0.92); opacity: 0.30; }
.lp-demo-card {
  position: absolute; inset: 0; transform-style: preserve-3d;
  transition: transform .9s cubic-bezier(.5,.05,.2,1); cursor: pointer;
}
.lp-demo-card.is-flipped { transform: rotateY(180deg); }
.lp-demo-face {
  position: absolute; inset: 0; background: var(--card);
  border: 1px solid var(--border); border-radius: 16px;
  padding: 22px 24px 20px; display: flex; flex-direction: column;
  backface-visibility: hidden; -webkit-backface-visibility: hidden;
  text-align: left; overflow: hidden;
}
.lp-demo-face.back { transform: rotateY(180deg); }
.lp-demo-face::before {
  content: ""; position: absolute; top: 0; left: 22px; right: 22px;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--sky-hover), transparent);
  opacity: 0.65;
}
.lp-demo-top {
  display: flex; justify-content: space-between; align-items: center; gap: 12px;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute); white-space: nowrap;
}
.lp-pill {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 3px 9px; border: 1px solid var(--border-strong); border-radius: 999px;
  color: var(--fg-dim); white-space: nowrap;
}
.lp-swatch {
  width: 6px; height: 6px; border-radius: 999px;
  background: var(--sky-hover); flex-shrink: 0;
}
.lp-demo-body {
  flex: 1; display: flex; flex-direction: column; justify-content: center; padding: 8px 0;
}
.lp-demo-word { font-size: 30px; font-weight: 700; letter-spacing: -.02em; color: var(--fg); line-height: 1.1; }
.lp-demo-ipa {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; font-weight: 500; color: var(--fg-mute); margin-top: 8px;
}
.lp-demo-example {
  font-size: 12.5px; color: var(--fg-mute); font-style: italic;
  line-height: 1.5; margin-top: 10px;
}
.lp-demo-example :deep(em) { color: var(--fg-dim); font-style: italic; font-weight: 500; }
.lp-demo-bottom { display: flex; gap: 6px; margin-top: auto; }
.lp-rate {
  flex: 1; display: flex; flex-direction: column; align-items: center; gap: 2px;
  padding: 8px 4px; border: 1px solid var(--border); border-radius: 8px;
  background: rgba(7,13,28,.4);
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  transition: border-color .15s ease, background .15s ease; cursor: pointer;
}
.lp-rate:hover { border-color: var(--border-strong); background: rgba(7,13,28,.7); }
.rate-label   { font-size: 10.5px; color: var(--fg-dim); font-weight: 500; }
.rate-interval { font-size: 9.5px;  color: var(--fg-mute); }
.lp-rate[data-tone="again"] .rate-label { color: #f87171; }
.lp-rate[data-tone="hard"]  .rate-label { color: #fbbf24; }
.lp-rate[data-tone="good"]  .rate-label { color: var(--sky-hover); }
.lp-rate[data-tone="easy"]  .rate-label { color: #4ade80; }
.lp-demo-reveal {
  display: flex; justify-content: space-between; align-items: center;
  margin-top: auto; padding-top: 12px; border-top: 1px dashed var(--border);
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute);
}
.lp-demo-reveal kbd, .lp-demo-hint kbd {
  display: inline-block; padding: 1px 6px;
  border: 1px solid var(--border-strong); border-bottom-width: 2px;
  border-radius: 4px; background: var(--bg); color: var(--fg-dim);
  font-family: 'JetBrains Mono', ui-monospace, monospace; font-size: 10.5px;
}
.lp-demo-hint {
  font-size: 12px; color: var(--fg-mute); text-align: center; margin-top: 36px;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
}

/* ── Features section ────────────────────────────── */
.lp-section { padding: 96px 0; border-bottom: 1px solid var(--border); }
.lp-section-head { margin-bottom: 48px; }
.lp-section-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-mute); text-transform: lowercase;
  letter-spacing: .02em; margin-bottom: 12px;
}
.lp-section-title { font-size: 28px; font-weight: 700; letter-spacing: -.02em; margin: 0; color: var(--fg); }
.lp-features { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.lp-card {
  position: relative; background: var(--card); border: 1px solid var(--border);
  border-radius: 12px; padding: 28px 24px;
  transition: border-color .2s ease, transform .2s ease;
  overflow: hidden; isolation: isolate;
}
.lp-card::before {
  content: ""; position: absolute; inset: 0; border-radius: inherit;
  background: radial-gradient(300px circle at var(--mx,50%) var(--my,50%), rgba(14,165,233,.10), transparent 60%);
  opacity: 0; transition: opacity .25s ease; pointer-events: none; z-index: -1;
}
.lp-card:hover { border-color: var(--border-strong); }
.lp-card:hover::before { opacity: 1; }
.lp-card-icon {
  display: inline-block; font-size: 24px; line-height: 1; margin-bottom: 20px;
  transition: transform .35s cubic-bezier(.34,1.56,.64,1);
}
.lp-card:hover .lp-card-icon { transform: translateY(-2px) rotate(-4deg); }
.lp-card-title { font-size: 16px; font-weight: 600; color: var(--fg); margin: 0 0 8px; letter-spacing: -.01em; }
.lp-card-body  { font-size: 14px; line-height: 1.6; color: var(--fg-dim); margin: 0; }

/* ── Forgetting curve sparkline ──────────────────── */
.lp-curve {
  display: block; width: 100%; height: 36px; margin-bottom: 22px; overflow: visible;
}
.lp-curve path {
  fill: none; stroke: var(--sky-hover); stroke-width: 1.5; stroke-linecap: round;
  stroke-dasharray: 200; stroke-dashoffset: 200;
}
.lp-curve circle {
  fill: var(--bg); stroke: var(--sky-hover); stroke-width: 1.5; opacity: 0;
}
.lp-curve.is-in path { animation: draw 1.6s ease forwards; }
.lp-curve.is-in circle { animation: pop .3s ease forwards; }
.lp-curve.is-in circle:nth-of-type(1) { animation-delay: .2s; }
.lp-curve.is-in circle:nth-of-type(2) { animation-delay: .7s; }
.lp-curve.is-in circle:nth-of-type(3) { animation-delay: 1.1s; }
.lp-curve.is-in circle:nth-of-type(4) { animation-delay: 1.4s; }
@keyframes draw { to { stroke-dashoffset: 0; } }
@keyframes pop  { to { opacity: 1; } }

/* ── Footer CTA ──────────────────────────────────── */
.lp-cta { text-align: center; padding: 120px 0; }
.lp-cta h2 {
  font-size: clamp(28px, 4vw, 40px); font-weight: 700; letter-spacing: -.02em;
  margin: 0 0 32px; color: var(--fg);
}
.lp-counter { font-variant-numeric: tabular-nums; color: var(--sky-hover); }

/* ── Footer ──────────────────────────────────────── */
.lp-foot {
  border-top: 1px solid var(--border); padding: 32px 0;
  display: flex; align-items: center; justify-content: space-between;
  font-size: 13px; color: var(--fg-mute);
}
.lp-foot-links { display: flex; gap: 20px; }
.lp-foot-links a { color: var(--fg-mute); text-decoration: none; transition: color .15s ease; }
.lp-foot-links a:hover { color: var(--fg-dim); }

/* ── Responsive ──────────────────────────────────── */
@media (max-width: 720px) {
  .lp-page { padding: 0 20px; }
  .lp-features { grid-template-columns: 1fr; }
  .lp-hero { padding: 72px 0 64px; }
  .lp-section, .lp-cta { padding: 72px 0; }
  .lp-foot { flex-direction: column; gap: 12px; }
}

/* ── Reduced motion ──────────────────────────────── */
@media (prefers-reduced-motion: reduce) {
  .cycle-word { animation: none; }
  .lp-dot { animation: none; }
  .lp-demo-card { transition: none; }
  .lp-curve path  { stroke-dashoffset: 0; }
  .lp-curve circle { opacity: 1; }
}
</style>
