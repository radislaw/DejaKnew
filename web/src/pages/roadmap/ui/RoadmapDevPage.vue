<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LocaleSwitcher from '@/shared/ui/LocaleSwitcher.vue'

const router = useRouter()
const { t, tm } = useI18n()

const visiblePhases = ref<Set<number>>(new Set())
const cleanups: (() => void)[] = []

interface DevPhaseI18n {
  label: string
  title: string
  features: string[]
}

interface PhaseMeta {
  status: 'done' | 'active' | 'planned' | 'vision'
  progress: number
  doneCount: number
}

const PHASE_META: PhaseMeta[] = [
  { status: 'active',  progress: 50, doneCount: 8 },
  { status: 'planned', progress: 0,  doneCount: 0 },
  { status: 'planned', progress: 0,  doneCount: 0 },
  { status: 'planned', progress: 0,  doneCount: 0 },
  { status: 'planned', progress: 0,  doneCount: 0 },
  { status: 'vision',  progress: 0,  doneCount: 0 },
]

const phases = computed(() => {
  const data = tm('roadmapDev.phases') as DevPhaseI18n[]
  return data.map((p, i) => ({ ...p, ...PHASE_META[i], id: i }))
})

const statusMeta = computed(() => ({
  done:    { label: t('roadmap.statusDone'),    color: '#4ade80' },
  active:  { label: t('roadmap.statusActive'),  color: '#0ea5e9' },
  planned: { label: t('roadmap.statusPlanned'), color: '#a78bfa' },
  vision:  { label: t('roadmap.statusVision'),  color: '#64748b' },
}))

onMounted(() => {
  document.querySelectorAll<HTMLElement>('.rd-phase').forEach((el, i) => {
    const obs = new IntersectionObserver((entries) => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          visiblePhases.value = new Set([...visiblePhases.value, i])
          obs.unobserve(el)
        }
      })
    }, { threshold: 0.1 })
    obs.observe(el)
    cleanups.push(() => obs.disconnect())
  })
})

onUnmounted(() => cleanups.forEach(fn => fn()))
</script>

<template>
  <div class="rd-root">
    <div class="rd-page">
      <nav class="rd-nav">
        <a href="/" class="rd-logo" @click.prevent="router.push('/')">
          <span class="rd-logo-mark">D</span>
          <span>DéjàKnew</span>
        </a>
        <div class="rd-nav-right">
          <LocaleSwitcher />
          <span class="rd-dev-badge">{{ t('roadmapDev.devBadge') }}</span>
          <button class="rd-btn rd-btn-ghost" @click="router.push('/roadmap')">{{ t('roadmapDev.userView') }}</button>
        </div>
      </nav>

      <header class="rd-header">
        <div class="rd-eyebrow"><span class="rd-dot"></span> {{ t('roadmapDev.eyebrow') }}</div>
        <h1>{{ t('roadmapDev.title') }} <span class="rd-accent">{{ t('roadmapDev.titleAccent') }}</span></h1>
        <p>{{ t('roadmapDev.subtitle') }}</p>
      </header>

      <div class="rd-legend">
        <span v-for="(meta, key) in statusMeta" :key="key" class="rd-legend-item">
          <span class="rd-legend-dot" :style="{ background: meta.color }"></span>
          {{ meta.label }}
        </span>
      </div>

      <div class="rd-timeline">
        <div class="rd-spine" aria-hidden="true"></div>

        <div
          v-for="(phase, i) in phases"
          :key="i"
          class="rd-phase"
          :class="[`is-${phase.status}`, { 'is-visible': visiblePhases.has(i) }]"
          :style="{ '--delay': i * 60 + 'ms' }"
        >
          <div class="rd-node" aria-hidden="true">
            <div class="rd-node-ring"></div>
            <div class="rd-node-dot"></div>
          </div>

          <article class="rd-card">
            <div class="rd-card-top">
              <span class="rd-phase-label">{{ phase.label }}</span>
              <span
                class="rd-status-badge"
                :style="{
                  color: statusMeta[phase.status].color,
                  borderColor: statusMeta[phase.status].color + '44',
                  background: statusMeta[phase.status].color + '12',
                }"
              >
                <span
                  class="rd-status-dot"
                  :class="{ 'is-pulse': phase.status === 'active' }"
                  :style="{ background: statusMeta[phase.status].color }"
                ></span>
                {{ statusMeta[phase.status].label }}
              </span>
            </div>

            <h2 class="rd-card-title">{{ phase.title }}</h2>

            <div v-if="phase.status !== 'vision'" class="rd-progress-wrap">
              <div class="rd-progress-track">
                <div
                  class="rd-progress-fill"
                  :style="{
                    width: visiblePhases.has(i) ? phase.progress + '%' : '0%',
                    background: statusMeta[phase.status].color,
                    transitionDelay: visiblePhases.has(i) ? '250ms' : '0ms',
                  }"
                ></div>
              </div>
              <span class="rd-progress-label">{{ phase.progress }}%</span>
            </div>

            <ul class="rd-features">
              <li
                v-for="(feat, fi) in phase.features"
                :key="fi"
                class="rd-feature"
              >
                <span
                  class="rd-check"
                  :class="fi < phase.doneCount ? 'is-done' : 'is-todo'"
                >
                  <svg v-if="fi < phase.doneCount" viewBox="0 0 12 12" fill="none">
                    <path d="M2 6.5L5 9.5L10 3" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                  <svg v-else viewBox="0 0 12 12" fill="none">
                    <circle cx="6" cy="6" r="2" fill="currentColor" opacity="0.4"/>
                  </svg>
                </span>
                <span class="rd-feat-text" :class="{ 'is-done': fi < phase.doneCount }">{{ feat }}</span>
              </li>
            </ul>
          </article>
        </div>
      </div>

      <footer class="rd-foot">
        <div>{{ t('roadmapDev.footerCopy') }}</div>
        <div class="rd-foot-links">
          <a href="/" @click.prevent="router.push('/')">{{ t('roadmapDev.footerHome') }}</a>
          <a href="/roadmap" @click.prevent="router.push('/roadmap')">{{ t('roadmapDev.footerUserRoadmap') }}</a>
        </div>
      </footer>
    </div>
  </div>
</template>

<style scoped>
.rd-root {
  --bg: #070d1c;
  --fg: #f1f5f9;
  --fg-dim: #94a3b8;
  --fg-mute: #64748b;
  --sky: #0284c7;
  --card: #0f172a;
  --border: #1e293b;
  --border-strong: #334155;
  font-family: 'Inter', system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  -webkit-font-smoothing: antialiased;
  background: var(--bg);
  color: var(--fg);
  min-height: 100vh;
  overflow-x: hidden;
}

.rd-page {
  max-width: 760px;
  margin: 0 auto;
  padding: 0 32px;
}

.rd-nav {
  display: flex; align-items: center; justify-content: space-between;
  padding: 24px 0; border-bottom: 1px solid var(--border);
}
.rd-logo {
  display: inline-flex; align-items: center; gap: 10px;
  font-weight: 700; font-size: 17px; color: var(--fg); text-decoration: none; cursor: pointer;
}
.rd-logo-mark {
  width: 22px; height: 22px; border-radius: 6px;
  background: var(--card); border: 1px solid var(--border-strong);
  display: grid; place-items: center;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; font-weight: 600; color: #0ea5e9;
}
.rd-nav-right { display: flex; align-items: center; gap: 12px; }
.rd-dev-badge {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: #a78bfa;
  padding: 3px 8px; border-radius: 4px;
  border: 1px solid #a78bfa44; background: #a78bfa12;
}
.rd-btn {
  display: inline-flex; align-items: center; gap: 6px;
  font-family: inherit; font-weight: 500; font-size: 14px;
  padding: 8px 14px; border-radius: 8px; border: 1px solid transparent;
  cursor: pointer; background: transparent;
  transition: background .15s ease;
}
.rd-btn-ghost { color: var(--fg-dim); }
.rd-btn-ghost:hover { color: var(--fg); background: var(--card); }

.rd-header {
  text-align: center; padding: 64px 0 48px;
  border-bottom: 1px solid var(--border);
}
.rd-eyebrow {
  display: inline-flex; align-items: center; gap: 8px;
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-mute);
  padding: 6px 12px; border: 1px solid var(--border); border-radius: 999px;
  margin-bottom: 24px; background: rgba(15,23,42,.5);
}
.rd-dot {
  width: 6px; height: 6px; border-radius: 999px;
  background: #a78bfa;
  animation: rd-pulse 2.4s ease-out infinite;
}
@keyframes rd-pulse {
  0%   { box-shadow: 0 0 0 0   rgba(167,139,250,.5); }
  70%  { box-shadow: 0 0 0 8px rgba(167,139,250,0); }
  100% { box-shadow: 0 0 0 0   rgba(167,139,250,0); }
}
.rd-header h1 {
  font-size: clamp(32px, 5vw, 48px); font-weight: 800;
  letter-spacing: -.035em; margin: 0 0 16px;
}
.rd-accent { color: #a78bfa; }
.rd-header p { font-size: 16px; line-height: 1.6; color: var(--fg-dim); margin: 0; }

.rd-legend {
  display: flex; align-items: center; gap: 20px; flex-wrap: wrap;
  padding: 20px 0; border-bottom: 1px solid var(--border);
  font-size: 13px; color: var(--fg-dim);
}
.rd-legend-item { display: flex; align-items: center; gap: 7px; }
.rd-legend-dot { width: 8px; height: 8px; border-radius: 999px; flex-shrink: 0; }

.rd-timeline { position: relative; padding: 40px 0 32px; }
.rd-spine {
  position: absolute; left: 15px; top: 0; bottom: 0; width: 1px;
  background: linear-gradient(to bottom, transparent, var(--border) 8%, var(--border) 92%, transparent);
}

.rd-phase {
  position: relative; display: flex; gap: 32px; padding: 0 0 40px;
  opacity: 0; transform: translateY(16px);
  transition: opacity .4s ease var(--delay, 0ms), transform .4s ease var(--delay, 0ms);
}
.rd-phase.is-visible { opacity: 1; transform: translateY(0); }

.rd-node {
  position: relative; flex-shrink: 0; width: 30px;
  display: flex; align-items: flex-start; justify-content: center; padding-top: 16px;
}
.rd-node-ring {
  position: absolute; width: 18px; height: 18px; border-radius: 999px;
  border: 1.5px solid var(--border-strong); background: var(--bg);
}
.rd-node-dot {
  position: absolute; width: 8px; height: 8px; border-radius: 999px;
  background: var(--fg-mute); top: 21px;
}
.rd-phase.is-done .rd-node-ring    { border-color: #4ade80; }
.rd-phase.is-done .rd-node-dot     { background: #4ade80; }
.rd-phase.is-active .rd-node-ring  { border-color: #0ea5e9; box-shadow: 0 0 14px -2px rgba(14,165,233,.4); }
.rd-phase.is-active .rd-node-dot   { background: #0ea5e9; }
.rd-phase.is-planned .rd-node-ring { border-color: #a78bfa55; }
.rd-phase.is-planned .rd-node-dot  { background: #a78bfa66; }
.rd-phase.is-vision .rd-node-ring  { border-color: var(--border); }
.rd-phase.is-vision .rd-node-dot   { background: var(--fg-mute); opacity: .4; }

.rd-card {
  flex: 1; background: var(--card); border: 1px solid var(--border);
  border-radius: 14px; padding: 20px 24px;
}
.rd-phase.is-active .rd-card { border-color: rgba(14,165,233,.25); }
.rd-phase.is-vision .rd-card { opacity: .6; }

.rd-card-top {
  display: flex; align-items: center; gap: 10px; margin-bottom: 10px; flex-wrap: wrap;
}
.rd-phase-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute); letter-spacing: .05em;
}
.rd-status-badge {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 2px 9px; border-radius: 999px; border: 1px solid;
  font-size: 11px; font-weight: 500;
}
.rd-status-dot {
  width: 6px; height: 6px; border-radius: 999px; flex-shrink: 0;
}
.rd-status-dot.is-pulse { animation: rd-pulse 2.4s ease-out infinite; }

.rd-card-title {
  font-size: 18px; font-weight: 700; letter-spacing: -.02em; margin: 0 0 14px;
}

.rd-progress-wrap { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.rd-progress-track {
  flex: 1; height: 3px; border-radius: 999px;
  background: rgba(255,255,255,.06); overflow: hidden;
}
.rd-progress-fill {
  height: 100%; border-radius: 999px;
  transition: width .8s cubic-bezier(.4,0,.2,1);
}
.rd-progress-label {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 11px; color: var(--fg-mute); flex-shrink: 0; min-width: 28px; text-align: right;
}

.rd-features { list-style: none; margin: 0; padding: 0; display: flex; flex-direction: column; gap: 7px; }
.rd-feature { display: flex; align-items: baseline; gap: 9px; }
.rd-check {
  width: 16px; height: 16px; border-radius: 4px; flex-shrink: 0;
  display: grid; place-items: center; margin-top: 1px;
}
.rd-check.is-done { background: rgba(74,222,128,.12); color: #4ade80; border: 1px solid rgba(74,222,128,.25); }
.rd-check.is-todo { background: rgba(255,255,255,.03); color: var(--fg-mute); border: 1px solid var(--border); }
.rd-check svg { width: 11px; height: 11px; }
.rd-feat-text {
  font-family: 'JetBrains Mono', ui-monospace, monospace;
  font-size: 12px; color: var(--fg-dim); line-height: 1.5;
}
.rd-feat-text.is-done { color: var(--fg-mute); }

.rd-foot {
  border-top: 1px solid var(--border); padding: 24px 0 40px;
  display: flex; align-items: center; justify-content: space-between;
  font-size: 13px; color: var(--fg-mute);
}
.rd-foot-links { display: flex; gap: 20px; }
.rd-foot-links a {
  color: var(--fg-mute); text-decoration: none;
  transition: color .15s ease; cursor: pointer;
}
.rd-foot-links a:hover { color: var(--fg-dim); }

@media (max-width: 640px) {
  .rd-page { padding: 0 16px; }
  .rd-phase { gap: 18px; }
  .rd-card { padding: 16px 16px; }
  .rd-foot { flex-direction: column; gap: 10px; }
}

@media (prefers-reduced-motion: reduce) {
  .rd-phase { transition: none; opacity: 1; transform: none; }
  .rd-dot, .rd-node-dot, .rd-status-dot { animation: none; }
}
</style>
