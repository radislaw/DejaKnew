<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/shared/ui/LocaleSwitcher.vue'

const router = useRouter()
const { t, tm } = useI18n()

const visiblePhases = ref<Set<number>>(new Set())
const visibleItems = ref<Set<string>>(new Set())
const cleanups: (() => void)[] = []

interface PhaseMeta {
  status: 'done' | 'active' | 'planned' | 'vision'
  progress: number
  doneCount: number
}

const PHASE_META: PhaseMeta[] = [
  { status: 'active',  progress: 50,  doneCount: 8 },
  { status: 'planned', progress: 0,   doneCount: 0 },
  { status: 'planned', progress: 0,   doneCount: 0 },
  { status: 'planned', progress: 0,   doneCount: 0 },
  { status: 'planned', progress: 0,   doneCount: 0 },
  { status: 'vision',  progress: 0,   doneCount: 0 },
]

interface PhaseI18n {
  label: string
  title: string
  description: string
  features: string[]
}

const phases = computed(() => {
  const data = tm('roadmap.phases') as PhaseI18n[]
  return data.map((p, i) => ({
    ...p,
    ...PHASE_META[i],
    id: i,
  }))
})

const statusMeta = computed(() => ({
  done:    { label: t('roadmap.statusDone'),    color: '#4ade80' },
  active:  { label: t('roadmap.statusActive'),  color: '#0ea5e9' },
  planned: { label: t('roadmap.statusPlanned'), color: '#a78bfa' },
  vision:  { label: t('roadmap.statusVision'),  color: '#64748b' },
}))

onMounted(() => {
  const phaseEls = document.querySelectorAll<HTMLElement>('.rm-phase')
  phaseEls.forEach((el, i) => {
    const obs = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          visiblePhases.value = new Set([...visiblePhases.value, i])
          obs.unobserve(el)
        }
      })
    }, { threshold: 0.12 })
    obs.observe(el)
    cleanups.push(() => obs.disconnect())
  })

  const itemEls = document.querySelectorAll<HTMLElement>('.rm-feature')
  itemEls.forEach((el) => {
    const key = el.dataset.key || ''
    const obs = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          visibleItems.value = new Set([...visibleItems.value, key])
          obs.unobserve(el)
        }
      })
    }, { threshold: 0.5 })
    obs.observe(el)
    cleanups.push(() => obs.disconnect())
  })
})

onUnmounted(() => cleanups.forEach(fn => fn()))
</script>

<template>
  <div class="rm-root">
    <div class="rm-stars" aria-hidden="true"></div>

    <div class="rm-page">
      <!-- Nav -->
      <nav class="rm-nav">
        <a href="/" class="rm-logo" @click.prevent="router.push('/')">
          <span class="rm-logo-mark">D</span>
          <span>DéjàKnew</span>
        </a>
        <div class="rm-nav-right">
          <LocaleSwitcher />
          <button class="rm-btn rm-btn-ghost" @click="router.push('/register')">
            {{ t('roadmap.getStarted') }}
          </button>
        </div>
      </nav>

      <!-- Header -->
      <header class="rm-header">
        <div class="rm-eyebrow"><span class="rm-dot"></span> {{ t('roadmap.eyebrow') }}</div>
        <h1>{{ t('roadmap.title') }}<br><span class="rm-shimmer">{{ t('roadmap.titleHighlight') }}</span></h1>
        <p>{{ t('roadmap.subtitle') }}</p>
      </header>

      <!-- Legend -->
      <div class="rm-legend" aria-label="Status legend">
        <span v-for="(meta, key) in statusMeta" :key="key" class="rm-legend-item">
          <span class="rm-legend-dot" :style="{ background: meta.color }"></span>
          {{ meta.label }}
        </span>
      </div>

      <!-- Timeline -->
      <div class="rm-timeline">
        <div class="rm-spine" aria-hidden="true"></div>

        <div
          v-for="(phase, i) in phases"
          :key="phase.id"
          class="rm-phase"
          :class="[`is-${phase.status}`, { 'is-visible': visiblePhases.has(i) }]"
          :style="{ '--delay': i * 80 + 'ms' }"
        >
          <!-- Node -->
          <div class="rm-node" aria-hidden="true">
            <div class="rm-node-ring"></div>
            <div class="rm-node-dot"></div>
          </div>

          <!-- Card -->
          <article class="rm-card">
            <div class="rm-card-top">
              <div class="rm-card-labels">
                <span class="rm-phase-label">{{ phase.label }}</span>
                <span
                  class="rm-status-badge"
                  :style="{
                    color: statusMeta[phase.status].color,
                    borderColor: statusMeta[phase.status].color + '44',
                    background: statusMeta[phase.status].color + '12',
                  }"
                >
                  <span
                    class="rm-status-dot"
                    :class="{ 'is-pulse': phase.status === 'active' }"
                    :style="{ background: statusMeta[phase.status].color }"
                  ></span>
                  {{ statusMeta[phase.status].label }}
                </span>
              </div>
            </div>

            <h2 class="rm-card-title">{{ phase.title }}</h2>
            <p class="rm-card-desc">{{ phase.description }}</p>

            <!-- Progress bar -->
            <div v-if="phase.status !== 'vision'" class="rm-progress-wrap">
              <div class="rm-progress-track">
                <div
                  class="rm-progress-fill"
                  :style="{
                    width: visiblePhases.has(i) ? phase.progress + '%' : '0%',
                    background: statusMeta[phase.status].color,
                    transitionDelay: visiblePhases.has(i) ? '300ms' : '0ms',
                  }"
                ></div>
              </div>
              <span class="rm-progress-label">{{ phase.progress }}%</span>
            </div>

            <!-- Features -->
            <ul class="rm-features">
              <li
                v-for="(feat, fi) in phase.features"
                :key="fi"
                class="rm-feature"
                :data-key="`${i}-${fi}`"
                :class="{ 'is-visible': visibleItems.has(`${i}-${fi}`) }"
                :style="{ '--fi-delay': fi * 40 + 'ms' }"
              >
                <span
                  class="rm-check"
                  :class="fi < phase.doneCount ? 'is-done' : 'is-todo'"
                  :aria-label="fi < phase.doneCount ? 'Done' : 'Planned'"
                >
                  <svg v-if="fi < phase.doneCount" viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <path d="M2 6.5L5 9.5L10 3" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <svg v-else viewBox="0 0 12 12" fill="none" aria-hidden="true">
                    <circle cx="6" cy="6" r="2" fill="currentColor" opacity="0.4"/>
                  </svg>
                </span>
                <span class="rm-feat-text" :class="{ 'is-done': fi < phase.doneCount }">{{ feat }}</span>
              </li>
            </ul>
          </article>
        </div>
      </div>

      <!-- CTA -->
      <section class="rm-cta">
        <div class="rm-cta-inner">
          <div class="rm-cta-label">{{ t('roadmap.ctaLabel') }}</div>
          <h2>{{ t('roadmap.ctaTitle') }}</h2>
          <p>{{ t('roadmap.ctaSubtitle') }}</p>
          <button class="rm-btn rm-btn-primary rm-btn-lg" @click="router.push('/register')">
            <span>{{ t('roadmap.ctaButton') }}</span>
          </button>
        </div>
      </section>

      <!-- Footer -->
      <footer class="rm-foot">
        <div>© 2026 DéjàKnew</div>
        <div class="rm-foot-links">
          <a href="/" @click.prevent="router.push('/')">{{ t('roadmap.footerHome') }}</a>
        </div>
      </footer>
    </div>
  </div>
</template>

<style scoped>
/* ── Tokens ──────────────────────────────────────── */
.rm-root {
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
  font-feature-settings: 'cv11', 'ss01';
  line-height: 1.5;
  background: radial-gradient(ellipse 120% 60% at 50% 0%, #1a1545 0%, #0e1a3a 40%, var(--bg) 70%);
  color: var(--fg);
  min-height: 100vh;
  overflow-x: hidden;
}

::selection { background: rgba(14, 165, 233, 0.3); color: var(--fg); }

/* ── Starfield ───────────────────────────────────── */
.rm-stars {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background-image:
    radial-gradient(1px 1px at 11% 17%, rgba(186,230,253,.30), transparent 70%),
    radial-gradient(1px 1px at 79% 9%,  rgba(196,181,253,.25), transparent 70%),
    radial-gradient(1px 1px at 33% 68%, rgba(186,230,253,.22), transparent 70%),
    radial-gradient(1px 1px at 92% 54%, rgba(186,230,253,.18), transparent 70%),
    radial-gradient(1px 1px at 55% 29%, rgba(196,181,253,.20), transparent 70%),
    radial-gradient(1px 1px at 4%  81%, rgba(186,230,253,.18), transparent 70%),
    radial-gradient(1.5px 1.5px at 66% 42%, rgba(125,211,252,.35), transparent 70%),
    radial-gradient(1px 1px at 21% 93%, rgba(186,230,253,.18), transparent 70%),
    radial-gradient(1px 1px at 87% 76%, rgba(196,181,253,.18), transparent 70%),
    radial-gradient(1px 1px at 44% 5%,  rgba(125,211,252,.28), transparent 70%),
    radial-gradient(1px 1px at 73% 88%, rgba(186,230,253,.20), transparent 70%);
}

/* ── Page ────────────────────────────────────────── */
.rm-page {
  position: relative;
  z-index: 1;
  max-width: 760px;
  margin: 0 auto;
  padding: 0 32px 0;
}

/* ── Nav ─────────────────────────────────────────── */
.rm-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
}
.rm-logo {
  display: inline-flex; align-items: center; gap: 10px;
  font-weight: 700; font-size: 17px; letter-spacing: -.01em;
  color: var(--fg); text-decoration: none; cursor: pointer;
}
.rm-logo-mark {
  width: 22px; height: 22px; border-radius: 6px;
  background: var(--card); border: 1px solid var(--border-strong);
  display: grid; place-items: center;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; font-weight: 600; color: var(--sky-hover);
}

/* ── Buttons ─────────────────────────────────────── */
.rm-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 6px;
  font-family: inherit; font-weight: 500; font-size: 14px; line-height: 1;
  padding: 10px 16px; border-radius: 8px; border: 1px solid transparent;
  cursor: pointer; color: var(--fg); background: transparent;
  transition: background .15s ease, border-color .15s ease, color .15s ease;
  white-space: nowrap;
}
.rm-nav-right { display: flex; align-items: center; gap: 8px; }
.rm-btn-ghost { color: var(--fg-dim); }
.rm-btn-ghost:hover { color: var(--fg); background: var(--card); }
.rm-btn-primary {
  position: relative; background: var(--sky); color: #f0f9ff;
  font-weight: 600; overflow: hidden; isolation: isolate;
}
.rm-btn-primary > * { position: relative; z-index: 1; }
.rm-btn-primary::after {
  content: ""; position: absolute; top: 0; left: -120%; width: 60%; height: 100%;
  background: linear-gradient(110deg, transparent 30%, rgba(255,255,255,.22) 50%, transparent 70%);
  transition: left .7s ease; pointer-events: none;
}
.rm-btn-primary:hover { background: var(--sky-hover); }
.rm-btn-primary:hover::after { left: 130%; }
.rm-btn-lg { font-size: 15px; padding: 14px 22px; border-radius: 10px; }

/* ── Header ──────────────────────────────────────── */
.rm-header {
  text-align: center;
  padding: 80px 0 56px;
  border-bottom: 1px solid var(--border);
}
.rm-eyebrow {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-mute);
  padding: 6px 12px; border: 1px solid var(--border); border-radius: 999px;
  margin-bottom: 28px; background: rgba(15,23,42,.5);
}
.rm-dot {
  width: 6px; height: 6px; border-radius: 999px;
  background: var(--sky-hover);
  box-shadow: 0 0 0 0 rgba(14,165,233,.6);
  animation: rm-pulse 2.4s ease-out infinite;
}
@keyframes rm-pulse {
  0%   { box-shadow: 0 0 0 0   rgba(14,165,233,.5); }
  70%  { box-shadow: 0 0 0 8px rgba(14,165,233,0);  }
  100% { box-shadow: 0 0 0 0   rgba(14,165,233,0);  }
}
.rm-header h1 {
  font-size: clamp(36px, 5.5vw, 54px); font-weight: 800;
  line-height: 1.05; letter-spacing: -.035em;
  margin: 0 0 20px; color: var(--fg);
}
.rm-shimmer {
  background-image: linear-gradient(100deg, #f1f5f9 0%, #f1f5f9 25%, #7dd3fc 40%, #c4b5fd 55%, #f1f5f9 70%, #f1f5f9 100%);
  background-size: 250% 100%;
  -webkit-background-clip: text; background-clip: text;
  -webkit-text-fill-color: transparent; color: transparent;
  animation: shimmer 6s ease-in-out infinite;
}
@keyframes shimmer {
  0%, 100% { background-position: 130% 0; }
  50%       { background-position: -30% 0; }
}
.rm-header p {
  font-size: 17px; line-height: 1.6; color: var(--fg-dim);
  max-width: 480px; margin: 0 auto;
}

/* ── Legend ──────────────────────────────────────── */
.rm-legend {
  display: flex; align-items: center; gap: 20px; flex-wrap: wrap;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
  font-size: 13px; color: var(--fg-dim);
}
.rm-legend-item { display: flex; align-items: center; gap: 7px; }
.rm-legend-dot { width: 8px; height: 8px; border-radius: 999px; flex-shrink: 0; }

/* ── Timeline ────────────────────────────────────── */
.rm-timeline {
  position: relative;
  padding: 48px 0 32px;
}
.rm-spine {
  position: absolute;
  left: 15px;
  top: 0; bottom: 0;
  width: 1px;
  background: linear-gradient(to bottom, transparent, var(--border) 8%, var(--border) 92%, transparent);
}

/* ── Phase ───────────────────────────────────────── */
.rm-phase {
  position: relative;
  display: flex;
  gap: 32px;
  padding: 0 0 48px;
  opacity: 0;
  transform: translateY(20px);
  transition: opacity .5s ease var(--delay, 0ms), transform .5s ease var(--delay, 0ms);
}
.rm-phase.is-visible { opacity: 1; transform: translateY(0); }

/* Node */
.rm-node {
  position: relative; flex-shrink: 0; width: 30px;
  display: flex; align-items: flex-start; justify-content: center; padding-top: 20px;
}
.rm-node-ring {
  position: absolute; width: 18px; height: 18px; border-radius: 999px;
  border: 1.5px solid var(--border-strong); background: var(--bg);
  transition: border-color .3s ease, box-shadow .3s ease;
}
.rm-node-dot {
  position: absolute; width: 8px; height: 8px; border-radius: 999px;
  background: var(--fg-mute); top: 25px;
  transition: background .3s ease;
}
.rm-phase.is-done .rm-node-ring    { border-color: #4ade80; box-shadow: 0 0 12px -2px rgba(74,222,128,.35); }
.rm-phase.is-done .rm-node-dot     { background: #4ade80; }
.rm-phase.is-active .rm-node-ring  { border-color: #0ea5e9; box-shadow: 0 0 16px -2px rgba(14,165,233,.45); }
.rm-phase.is-active .rm-node-dot   { background: #0ea5e9; animation: node-pulse 2s ease-out infinite; }
@keyframes node-pulse {
  0%   { box-shadow: 0 0 0 0   rgba(14,165,233,.6); }
  70%  { box-shadow: 0 0 0 8px rgba(14,165,233,0); }
  100% { box-shadow: 0 0 0 0   rgba(14,165,233,0); }
}
.rm-phase.is-planned .rm-node-ring { border-color: #a78bfa55; }
.rm-phase.is-planned .rm-node-dot  { background: #a78bfa66; }
.rm-phase.is-vision .rm-node-ring  { border-color: var(--border); }
.rm-phase.is-vision .rm-node-dot   { background: var(--fg-mute); opacity: .4; }

/* Card */
.rm-card {
  flex: 1; background: var(--card); border: 1px solid var(--border);
  border-radius: 16px; padding: 24px 28px;
  position: relative; overflow: hidden;
  transition: border-color .2s ease;
}
.rm-card::before {
  content: ""; position: absolute; top: 0; left: 28px; right: 28px; height: 1px;
}
.rm-phase.is-done .rm-card         { border-color: rgba(74,222,128,.18); }
.rm-phase.is-done .rm-card::before { background: linear-gradient(90deg, transparent, rgba(74,222,128,.5), transparent); }
.rm-phase.is-active .rm-card       { border-color: rgba(14,165,233,.28); }
.rm-phase.is-active .rm-card::before { background: linear-gradient(90deg, transparent, rgba(14,165,233,.6), transparent); }
.rm-phase.is-planned .rm-card::before { background: linear-gradient(90deg, transparent, rgba(167,139,250,.35), transparent); }
.rm-phase.is-vision .rm-card       { opacity: .65; }
.rm-phase.is-vision .rm-card::before { background: linear-gradient(90deg, transparent, rgba(100,116,139,.3), transparent); }

.rm-card-top {
  display: flex; align-items: center; justify-content: space-between; gap: 12px;
  margin-bottom: 14px; flex-wrap: wrap;
}
.rm-card-labels { display: flex; align-items: center; gap: 8px; }
.rm-phase-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute); letter-spacing: .05em;
}
.rm-status-badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 3px 10px; border-radius: 999px; border: 1px solid;
  font-size: 12px; font-weight: 500; letter-spacing: .01em;
}
.rm-status-dot { width: 6px; height: 6px; border-radius: 999px; flex-shrink: 0; }
.rm-status-dot.is-pulse { animation: rm-pulse 2.4s ease-out infinite; }

.rm-card-title {
  font-size: 20px; font-weight: 700; letter-spacing: -.02em;
  margin: 0 0 8px; color: var(--fg);
}
.rm-card-desc {
  font-size: 14px; line-height: 1.6; color: var(--fg-dim); margin: 0 0 20px;
}

/* Progress */
.rm-progress-wrap { display: flex; align-items: center; gap: 10px; margin-bottom: 20px; }
.rm-progress-track {
  flex: 1; height: 4px; border-radius: 999px;
  background: rgba(255,255,255,.07); overflow: hidden;
}
.rm-progress-fill {
  height: 100%; border-radius: 999px;
  transition: width .8s cubic-bezier(.4,0,.2,1);
}
.rm-progress-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute); flex-shrink: 0; min-width: 28px; text-align: right;
}

/* Features */
.rm-features {
  list-style: none; margin: 0; padding: 0;
  display: grid; grid-template-columns: 1fr 1fr; gap: 8px 16px;
}
.rm-feature {
  display: flex; align-items: center; gap: 9px;
  opacity: 0; transform: translateX(-8px);
  transition: opacity .35s ease var(--fi-delay, 0ms), transform .35s ease var(--fi-delay, 0ms);
}
.rm-feature.is-visible { opacity: 1; transform: translateX(0); }
.rm-check {
  width: 18px; height: 18px; border-radius: 5px; flex-shrink: 0;
  display: grid; place-items: center;
}
.rm-check.is-done {
  background: rgba(74,222,128,.15); color: #4ade80;
  border: 1px solid rgba(74,222,128,.3);
}
.rm-check.is-todo {
  background: rgba(255,255,255,.04); color: var(--fg-mute);
  border: 1px solid var(--border);
}
.rm-check svg { width: 12px; height: 12px; }
.rm-feat-text { font-size: 13px; color: var(--fg-dim); line-height: 1.3; }
.rm-feat-text.is-done { color: var(--fg-mute); }

/* ── CTA ─────────────────────────────────────────── */
.rm-cta { padding: 72px 0; border-top: 1px solid var(--border); }
.rm-cta-inner {
  background: var(--card); border: 1px solid var(--border);
  border-radius: 20px; padding: 48px; text-align: center;
  position: relative; overflow: hidden;
}
.rm-cta-inner::before {
  content: ""; position: absolute; top: 0; left: 10%; right: 10%; height: 1px;
  background: linear-gradient(90deg, transparent, rgba(14,165,233,.5), rgba(167,139,250,.4), transparent);
}
.rm-cta-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-mute); margin-bottom: 16px;
}
.rm-cta-inner h2 {
  font-size: 28px; font-weight: 700; letter-spacing: -.02em; margin: 0 0 10px; color: var(--fg);
}
.rm-cta-inner p { font-size: 15px; color: var(--fg-dim); margin: 0 0 28px; }

/* ── Footer ──────────────────────────────────────── */
.rm-foot {
  border-top: 1px solid var(--border); padding: 28px 0 40px;
  display: flex; align-items: center; justify-content: space-between;
  font-size: 13px; color: var(--fg-mute);
}
.rm-foot-links { display: flex; gap: 20px; }
.rm-foot-links a {
  color: var(--fg-mute); text-decoration: none;
  transition: color .15s ease; cursor: pointer;
}
.rm-foot-links a:hover { color: var(--fg-dim); }

/* ── Responsive ──────────────────────────────────── */
@media (max-width: 640px) {
  .rm-page { padding: 0 16px; }
  .rm-header { padding: 56px 0 40px; }
  .rm-features { grid-template-columns: 1fr; }
  .rm-phase { gap: 20px; }
  .rm-card { padding: 18px 18px; }
  .rm-cta-inner { padding: 32px 20px; }
  .rm-foot { flex-direction: column; gap: 12px; }
}

/* ── Reduced motion ──────────────────────────────── */
@media (prefers-reduced-motion: reduce) {
  .rm-phase { transition: none; opacity: 1; transform: none; }
  .rm-feature { transition: none; opacity: 1; transform: none; }
  .rm-progress-fill { transition: none; }
  .rm-dot, .rm-node-dot, .rm-status-dot { animation: none; }
  .rm-shimmer { animation: none; }
}
</style>
