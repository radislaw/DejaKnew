import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  RotateCw, ChevronLeft, ChevronRight, Check, X,
  BarChart3, Brain, Sparkles, Zap, Target,
  TrendingUp, Clock, Flame, BookOpen, Filter, Trash2
} from 'lucide-react';

const CARDS = [
  {id: 1, q: "Что такое L1 автономии ИИ?", a: "Autocomplete — tab-подсказки, inline-suggestions. Ты пишешь, ИИ дополняет токены.", tag: "autonomy"},
  {id: 2, q: "Что такое L2 автономии ИИ?", a: "Chat / парный диалог — ИИ предлагает код по запросу, ты ревьюишь каждую строку перед мерджем.", tag: "autonomy"},
  {id: 3, q: "Что такое L3 автономии ИИ?", a: "Supervised agent — ИИ выполняет многошаговую задачу под надзором (Claude Code, Cursor agent). Ты проверяешь перед коммитом.", tag: "autonomy"},
  {id: 4, q: "Что такое L4 автономии ИИ?", a: "Autonomous agent — работает в фоне, открывает PR. Ты ревьюишь PR, но не следишь за выполнением.", tag: "autonomy"},
  {id: 5, q: "Какой уровень автономии для багфикса?", a: "L2. ИИ — хороший co-detective, плохой solo-surgeon.", tag: "scenario"},
  {id: 6, q: "Какой уровень автономии для работы с legacy?", a: "L1 — L2 строго. Это сценарий с самыми драматичными провалами автономных агентов.", tag: "scenario"},
  {id: 7, q: "Какой уровень автономии для рефакторинга?", a: "L2 — L3, но только при наличии тестов. Без тестов — не рефактори.", tag: "scenario"},
  {id: 8, q: "Какой уровень автономии для обучения?", a: "L1 — L2 строго. Низкая автономия — это фича, не баг: иначе получишь приложение без навыка.", tag: "scenario"},
  {id: 9, q: "Какой уровень автономии для прототипа/MVP?", a: "L3 — L4. Единственный сценарий, где vibe coding уместен.", tag: "scenario"},
  {id: 10, q: "Какой уровень автономии для нового проекта?", a: "L2 для фундамента, L3 после того, как заложен constitution.", tag: "scenario"},
  {id: 11, q: "Какой уровень автономии для новой фичи?", a: "L2 на старте (пока загружаешь контекст), L3 после.", tag: "scenario"},
  {id: 12, q: "Какой уровень автономии для code review?", a: "L2 — L3. AI как first-pass filter, решение о merge — за человеком.", tag: "scenario"},
  {id: 13, q: "Принцип 1: Specs persist…", a: "…prompts are transient. Промпт живёт сессию, спека — пока живёт проект. Любая задача >30 минут начинается со спеки.", tag: "principles"},
  {id: 14, q: "Принцип 2: Context is the bottleneck", a: "65% разработчиков жалуются на потерю контекста. Агенты тратят 60–80% токенов на ориентацию. Давай ровно нужный контекст — не весь репо.", tag: "principles"},
  {id: 15, q: "Принцип 3: Trust gradient", a: "Выбирай уровень автономии под риск. Цена ошибки↑ + понимание области↓ = автономия↓.", tag: "principles"},
  {id: 16, q: "Принцип 4: Verify everything", a: "Never merge what you don't understand. ИИ пишет убедительно неправильный код. Тесты + линтер + read every line.", tag: "principles"},
  {id: 17, q: "Принцип 5: Small diffs, atomic commits", a: "Коммиты — это save points. Реальный кейс: 2M-токенов → 31 файл → 260 ошибок. Малыми шагами тот же рефакторинг — 11 минут, 0 ошибок.", tag: "principles"},
  {id: 18, q: "Принцип 6: AI amplifies, doesn't replace", a: "Senior × ИИ = умножение скорости. Junior × ИИ = умножение ошибок. ИИ усиливает то, что уже есть.", tag: "principles"},
  {id: 19, q: "Сценарий A", a: "Новый проект (greenfield). Spec → Architecture → Constitution → Scaffold → Reference implementation.", tag: "scenarios"},
  {id: 20, q: "Сценарий B", a: "Прототип / MVP / hackathon. Throwaway-код, vibe coding допустим, но 1 e2e-тест обязателен. README с warning.", tag: "scenarios"},
  {id: 21, q: "Сценарий C", a: "Новая фича в существующий проект. Spec → Context (CLAUDE.md + похожие фичи) → Decompose в 3–7 шагов → Step-by-step.", tag: "scenarios"},
  {id: 22, q: "Сценарий D", a: "Багфикс. Reproduce → Failing test → Hypotheses → Narrow down → Minimal fix → Green test → Regression-class thinking.", tag: "scenarios"},
  {id: 23, q: "Сценарий E", a: "Legacy / unfamiliar code. Археология → Dependency map → Characterization tests ДО изменений → Strangler Fig, не rewrite.", tag: "scenarios"},
  {id: 24, q: "Сценарий F", a: "Рефакторинг. Тесты обязательны. Scope → Green baseline → One type of change → Behavior diff = 0.", tag: "scenarios"},
  {id: 25, q: "Сценарий G", a: "Обучение через проект. Концепт → попытка сам → сравни с ИИ → объясни своими словами → песочница для экспериментов.", tag: "scenarios"},
  {id: 26, q: "Сценарий H", a: "Code review. Review sandwich: AI первым проходом → Human вторым → AI никогда не authority. Merge решение — человек.", tag: "scenarios"},
  {id: 27, q: "Что делать, если задача попадает в два сценария?", a: "Побеждает более осторожный сценарий. Обычно это E (legacy) или G (обучение). Снижай автономию.", tag: "rules"},
  {id: 28, q: "Что делать, если нужен рефакторинг, а тестов нет?", a: "СТОП. Не рефактори. Сначала переходишь в сценарий E и добавляешь characterization tests. Потом возвращаешься.", tag: "rules"},
  {id: 29, q: "Что такое characterization tests?", a: "Тесты, которые фиксируют текущее поведение кода (даже если странное). Не валидируют корректность — только consistency. Сафети-сеть для изменений в legacy.", tag: "concepts"},
  {id: 30, q: "Что такое Strangler Fig?", a: "Стратегия замены legacy: новый код рядом со старым, постепенная миграция трафика, удаление старого только когда никто не вызывает. Альтернатива опасным rewrite.", tag: "concepts"},
  {id: 31, q: "Что такое constitution проекта?", a: "Файл правил для ИИ (CLAUDE.md, .cursorrules, AGENTS.md) в корне репо: стек, стиль, паттерны, что не делать. Самый высокорычажный артефакт проекта.", tag: "concepts"},
  {id: 32, q: "Что такое review sandwich?", a: "AI review первым проходом → Human review вторым → AI никогда не принимает merge-решение. AI — фильтр, не authority.", tag: "concepts"},
  {id: 33, q: "Что такое vibe coding и когда оно допустимо?", a: "Открытые промпты без спеки. Допустимо ТОЛЬКО в сценарии B (прототип). На продакшене — антипаттерн.", tag: "concepts"},
  {id: 34, q: "Сколько шагов в цикле задачи?", a: "8: Spec → Context → Decompose → Generate → Verify → Commit → Loop → Automate.", tag: "cycle"},
  {id: 35, q: "Что такое hooks в Claude Code?", a: "User-defined shell-команды, которые автоматически запускаются в определённые моменты жизни ИИ. Превращают рекомендации в принудительные правила.", tag: "hooks"},
  {id: 36, q: "Какой hook единственный может блокировать действие?", a: "PreToolUse. Exit code 2 останавливает tool call и отправляет stderr Claude как новый input — не просто блокирует, а учит.", tag: "hooks"},
  {id: 37, q: "Какой hook закрепляет принцип Verify everything?", a: "PostToolUse на Edit/Write. Запускает линтер, type checker, тесты после каждого изменения файла.", tag: "hooks"},
  {id: 38, q: "Какой hook закрепляет принцип Specs persist?", a: "SessionStart. В начале сессии читает CLAUDE.md, spec.md, current-plan.md — Claude автоматически грузит контекст.", tag: "hooks"},
  {id: 39, q: "4 слоя контроля ИИ — какие?", a: "CLAUDE.md (рекомендация) → Skills (умения) → MCP (доступ) → Hooks (закон на уровне shell). Серьёзные правила — в hooks.", tag: "concepts"},
  {id: 40, q: "Что такое progressive disclosure в skills?", a: "Архитектура: метаданные skills всегда в контексте (~100 токенов), полное содержимое (<5k) подгружается только когда skill срабатывает.", tag: "concepts"},
  {id: 41, q: "Антипаттерн: «просто перепиши этот модуль»", a: "Гарантированная катастрофа на legacy. Канонический кейс: 2M токенов, 31 файл, 260 ошибок типов, сломанный билд. Декомпозируй ruthlessly.", tag: "antipatterns"},
  {id: 42, q: "Антипаттерн: «мерджу, потом разберусь»", a: "Не работает никогда. Либо разбираешься сейчас, либо через месяц переписываешь.", tag: "antipatterns"},
  {id: 43, q: "Антипаттерн: один монструозный промпт", a: "«Сделай мне всю систему X» = мусор без архитектуры. Декомпозиция важнее всего.", tag: "antipatterns"},
  {id: 44, q: "Антипаттерн: прототип, который ушёл в прод", a: "Главный источник legacy-ада. Если случилось — немедленно переключайся на сценарий E (legacy) для любых изменений.", tag: "antipatterns"},
  {id: 45, q: "Антипаттерн: trust ИИ в незнакомой области", a: "ИИ уверенно объясняет несуществующие features. Кросс-проверяй через документацию для всех нетривиальных claim'ов.", tag: "antipatterns"},
  {id: 46, q: "Антипаттерн: L4 при обучении", a: "Получаешь приложение, не получаешь навык. Через месяц всё забыто. В обучении низкая автономия — это фича.", tag: "antipatterns"},
  {id: 47, q: "Антипаттерн: рефакторинг без тестов", a: "Это не рефакторинг, а слепые изменения. Без тестов нет способа узнать, сохранил ли ты поведение.", tag: "antipatterns"},
  {id: 48, q: "Сколько коммитов на одну задачу?", a: "Один задача = один атомарный коммит. Если задача декомпозирована на N шагов — N коммитов. Каждый коммит — save point.", tag: "rules"},
  {id: 49, q: "Чек-лист перед началом задачи", a: "1) Знаю сценарий A-H? 2) Выбран уровень L1-L4? 3) Есть спека? 4) Загружен правильный контекст? 5) Декомпозирована? 6) Тесты зелёные? 7) Есть CLAUDE.md? Если хоть одно «нет» — стоп.", tag: "checklist"},
  {id: 50, q: "Что делать на старте багфикса?", a: "До любого кода: воспроизведи баг локально. Затем напиши failing test. Без воспроизведения ИИ-помощь — спекуляция.", tag: "cycle"},
  {id: 51, q: "Что делать на старте работы с legacy?", a: "Определи цель узко (understand? extend? replace?). Затем археология модуль за модулем, dependency map, characterization tests — только потом изменения.", tag: "cycle"},
  {id: 52, q: "Что делать при обнаружении бага во время рефакторинга?", a: "Записать, исправить в отдельном PR. НЕ совмещай рефакторинг и багфикс в одном PR.", tag: "rules"},
  {id: 53, q: "Что значит «behavior diff = 0» в рефакторинге?", a: "Те же тесты проходят, тот же external API, те же логи/метрики/UI. Если что-то изменилось — это не рефакторинг, а изменение поведения.", tag: "concepts"}
];

// === DEEP NAVY PALETTE ===
// Tags use light cool/warm tints that read well on navy background
const TAGS = {
  all: { label: "Все", color: "#94a3b8", icon: BookOpen },
  principles: { label: "Принципы", color: "#7dd3fc", icon: Brain },
  autonomy: { label: "Уровни", color: "#67e8f9", icon: TrendingUp },
  scenarios: { label: "Сценарии", color: "#fcd34d", icon: Target },
  scenario: { label: "Выбор", color: "#fdba74", icon: Filter },
  concepts: { label: "Концепты", color: "#a5b4fc", icon: Sparkles },
  antipatterns: { label: "Антипаттерны", color: "#fda4af", icon: X },
  hooks: { label: "Hooks", color: "#c4b5fd", icon: Zap },
  rules: { label: "Правила", color: "#86efac", icon: Check },
  cycle: { label: "Цикл", color: "#f9a8d4", icon: RotateCw },
  checklist: { label: "Чек-лист", color: "#93c5fd", icon: Check }
};

// Design tokens — single source of truth
const C = {
  bg: '#070d1c',          // page bg, deepest midnight
  bgGradTop: '#0d1530',   // subtle navy lift at top
  cardFront: '#0f172a',   // slate-900-ish, primary card
  cardFrontHi: '#1e293b', // slate-800 highlight
  cardBack: '#131c34',    // back of card — dark but distinct from front
  cardBackHi: '#1a2547',  // back gradient highlight
  border: 'rgba(148, 163, 184, 0.10)',
  borderHi: 'rgba(148, 163, 184, 0.18)',
  text: '#e2e8f0',        // slate-200 main text on dark
  textDim: '#94a3b8',     // slate-400 secondary
  textFaint: '#64748b',   // slate-500 tertiary
  textGhost: '#475569',   // slate-600 hint
  accent: '#7dd3fc',      // sky-300, primary accent
  accentHot: '#38bdf8',   // sky-400, hover/active
  accentWarm: '#fbbf24',  // amber-400, secondary warmth
  good: '#86efac',
  warn: '#fdba74',
  bad: '#fda4af',
};

const FONT_STACK = "'Manrope', system-ui, -apple-system, BlinkMacSystemFont, sans-serif";
const MONO_STACK = "'JetBrains Mono', ui-monospace, monospace";

function calculateNextReview(card, quality) {
  const now = Date.now();
  let { interval = 0, ease = 2.5, reps = 0, lapses = 0 } = card.srs || {};

  if (quality === 0) {
    reps = 0;
    interval = 0;
    lapses += 1;
    ease = Math.max(1.3, ease - 0.2);
  } else {
    reps += 1;
    ease = Math.max(1.3, ease + (0.1 - (3 - quality) * (0.08 + (3 - quality) * 0.02)));
    if (reps === 1) interval = 1;
    else if (reps === 2) interval = quality === 1 ? 1 : 3;
    else interval = Math.round(interval * ease * (quality === 1 ? 0.6 : quality === 3 ? 1.3 : 1));
  }

  return {
    interval, ease, reps, lapses,
    nextReview: now + interval * 24 * 60 * 60 * 1000,
    lastReview: now
  };
}

function getCardStatus(card, now = Date.now()) {
  const srs = card.srs;
  if (!srs || srs.reps === 0) return 'new';
  if (srs.nextReview <= now) return 'due';
  if (srs.lapses >= 2) return 'difficult';
  if (srs.reps >= 5) return 'learned';
  return 'learning';
}

function useStorage(key, initial) {
  const [value, setValue] = useState(initial);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const r = await window.storage.get(key);
        if (!cancelled && r?.value) {
          try { setValue(JSON.parse(r.value)); }
          catch { setValue(r.value); }
        }
      } catch { /* not found is ok */ }
    })();
    return () => { cancelled = true; };
  }, [key]);

  const save = useCallback(async (newValue) => {
    setValue(newValue);
    try {
      await window.storage.set(key, JSON.stringify(newValue));
    } catch (e) { console.error('save failed', e); }
  }, [key]);

  return [value, save];
}

function StatRing({ value, max, color, size = 80, label, sublabel }) {
  const radius = (size - 8) / 2;
  const circumference = 2 * Math.PI * radius;
  const progress = max > 0 ? value / max : 0;
  const offset = circumference * (1 - progress);

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative" style={{ width: size, height: size }}>
        <svg width={size} height={size} className="-rotate-90">
          <circle cx={size/2} cy={size/2} r={radius} stroke="rgba(148, 163, 184, 0.10)" strokeWidth="4" fill="none" />
          <motion.circle
            cx={size/2} cy={size/2} r={radius}
            stroke={color} strokeWidth="4" fill="none" strokeLinecap="round"
            strokeDasharray={circumference}
            initial={{ strokeDashoffset: circumference }}
            animate={{ strokeDashoffset: offset }}
            transition={{ duration: 0.8, ease: "easeOut" }}
            style={{ filter: `drop-shadow(0 0 8px ${color}40)` }}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <div className="text-xl font-semibold tabular-nums" style={{ color, fontFamily: FONT_STACK }}>{value}</div>
          {sublabel && <div className="text-[9px] uppercase tracking-wider" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>{sublabel}</div>}
        </div>
      </div>
      {label && <div className="text-xs uppercase tracking-wider" style={{ color: C.textDim, fontFamily: MONO_STACK }}>{label}</div>}
    </div>
  );
}

function TagPill({ tag, active, count, onClick }) {
  const info = TAGS[tag] || { label: tag, color: '#888', icon: BookOpen };
  const Icon = info.icon;
  return (
    <motion.button
      whileTap={{ scale: 0.96 }}
      onClick={onClick}
      className="inline-flex items-center gap-1.5 px-3 py-1.5 text-[11px] rounded-full uppercase tracking-wider transition-all border"
      style={{
        fontFamily: MONO_STACK,
        ...(active
          ? { backgroundColor: info.color, color: C.bg, borderColor: 'transparent', boxShadow: `0 0 20px ${info.color}30` }
          : { color: C.textDim, borderColor: C.border, background: 'transparent' }
        )
      }}
      onMouseEnter={(e) => {
        if (!active) {
          e.currentTarget.style.color = C.text;
          e.currentTarget.style.borderColor = C.borderHi;
        }
      }}
      onMouseLeave={(e) => {
        if (!active) {
          e.currentTarget.style.color = C.textDim;
          e.currentTarget.style.borderColor = C.border;
        }
      }}
    >
      <Icon size={11} strokeWidth={2.2} />
      <span>{info.label}</span>
      {count !== undefined && (
        <span className="ml-0.5 tabular-nums" style={{ opacity: active ? 0.9 : 0.6 }}>{count}</span>
      )}
    </motion.button>
  );
}

function FlashCard({ card, flipped, onFlip, status }) {
  const tagInfo = TAGS[card.tag] || TAGS.all;
  const statusLabels = { new: 'новая', learning: 'учим', due: 'повторить', learned: 'выучено', difficult: 'сложная' };
  const statusColors = {
    due: C.warn,
    difficult: C.bad,
    learned: C.good,
    learning: C.accent
  };

  return (
    <div className="w-full" style={{ perspective: '2000px' }}>
      <motion.div
        onClick={onFlip}
        className="relative w-full cursor-pointer"
        style={{ transformStyle: 'preserve-3d', minHeight: '360px' }}
        animate={{ rotateY: flipped ? 180 : 0 }}
        transition={{ duration: 0.6, ease: [0.4, 0, 0.2, 1] }}
      >
        {/* FRONT — deepest navy */}
        <div
          className="absolute inset-0 rounded-2xl p-8 md:p-10 flex flex-col"
          style={{
            backfaceVisibility: 'hidden',
            WebkitBackfaceVisibility: 'hidden',
            background: `linear-gradient(155deg, ${C.cardFrontHi} 0%, ${C.cardFront} 60%, ${C.bg} 100%)`,
            border: `1px solid ${C.border}`,
            boxShadow: `0 30px 60px -20px rgba(0, 0, 0, 0.8), inset 0 1px 0 rgba(255, 255, 255, 0.04), 0 0 0 1px rgba(125, 211, 252, 0.04)`
          }}
        >
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-2">
              <div
                className="w-2 h-2 rounded-full"
                style={{
                  backgroundColor: tagInfo.color,
                  boxShadow: `0 0 8px ${tagInfo.color}80`
                }}
              />
              <span className="text-[10px] uppercase tracking-[0.2em]" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>
                {tagInfo.label}
              </span>
            </div>
            {status && status !== 'new' && (
              <span
                className="text-[9px] uppercase tracking-wider px-2 py-1 rounded-full"
                style={{
                  fontFamily: MONO_STACK,
                  color: statusColors[status] || C.textDim,
                  backgroundColor: `${statusColors[status] || C.textDim}15`,
                  border: `1px solid ${statusColors[status] || C.textDim}25`
                }}
              >
                {statusLabels[status]}
              </span>
            )}
          </div>

          <div className="flex-1 flex items-center">
            <h2 className="text-2xl md:text-3xl leading-tight tracking-tight" style={{
              color: C.text,
              fontFamily: FONT_STACK,
              fontWeight: 500,
              letterSpacing: '-0.02em'
            }}>
              {card.q}
            </h2>
          </div>

          <div className="flex items-center justify-center gap-2 text-[10px] uppercase tracking-[0.15em] mt-6"
            style={{ color: C.textGhost, fontFamily: MONO_STACK }}>
            <RotateCw size={11} />
            <span>Тап чтобы перевернуть</span>
          </div>
        </div>

        {/* BACK — slightly lighter navy, NOT white. Tag-tinted glow */}
        <div
          className="absolute inset-0 rounded-2xl p-8 md:p-10 flex flex-col"
          style={{
            backfaceVisibility: 'hidden',
            WebkitBackfaceVisibility: 'hidden',
            transform: 'rotateY(180deg)',
            background: `linear-gradient(155deg, ${C.cardBackHi} 0%, ${C.cardBack} 60%, ${C.cardFront} 100%)`,
            border: `1px solid ${tagInfo.color}25`,
            boxShadow: `0 30px 60px -20px rgba(0, 0, 0, 0.8), inset 0 1px 0 rgba(255, 255, 255, 0.04), 0 0 40px -20px ${tagInfo.color}40`
          }}
        >
          <div className="flex items-center gap-2 mb-6">
            <div
              className="w-2 h-2 rounded-full"
              style={{
                backgroundColor: tagInfo.color,
                boxShadow: `0 0 8px ${tagInfo.color}80`
              }}
            />
            <span className="text-[10px] uppercase tracking-[0.2em]" style={{ color: tagInfo.color, fontFamily: MONO_STACK, opacity: 0.85 }}>
              ответ
            </span>
          </div>

          <div className="flex-1 flex items-center">
            <p className="text-lg md:text-xl leading-relaxed" style={{
              color: C.text,
              fontFamily: FONT_STACK,
              fontWeight: 400,
              letterSpacing: '-0.005em'
            }}>
              {card.a}
            </p>
          </div>

          <div className="text-[10px] uppercase tracking-[0.15em] text-center mt-6" style={{ color: C.textGhost, fontFamily: MONO_STACK }}>
            Оцени свой ответ
          </div>
        </div>
      </motion.div>
    </div>
  );
}

function ReviewButtons({ onReview, disabled }) {
  const buttons = [
    { quality: 0, label: 'Забыл', sub: '<1d', color: C.bad, icon: X, key: '1' },
    { quality: 1, label: 'Сложно', sub: '1d', color: C.warn, icon: Clock, key: '2' },
    { quality: 2, label: 'Хорошо', sub: '3d', color: C.good, icon: Check, key: '3' },
    { quality: 3, label: 'Легко', sub: '7d+', color: C.accent, icon: Sparkles, key: '4' }
  ];

  return (
    <div className="grid grid-cols-4 gap-2 md:gap-3">
      {buttons.map(b => {
        const Icon = b.icon;
        return (
          <motion.button
            key={b.quality}
            whileHover={!disabled ? { y: -2 } : {}}
            whileTap={!disabled ? { scale: 0.97 } : {}}
            onClick={() => !disabled && onReview(b.quality)}
            disabled={disabled}
            className="relative flex flex-col items-center gap-1.5 p-3 md:p-4 rounded-xl border transition-all disabled:opacity-40 disabled:cursor-not-allowed"
            style={{
              borderColor: disabled ? C.border : `${b.color}40`,
              backgroundColor: disabled ? 'transparent' : `${b.color}10`,
              boxShadow: disabled ? 'none' : `0 0 24px -10px ${b.color}40`
            }}
          >
            <Icon size={18} style={{ color: b.color }} strokeWidth={2.2} />
            <span className="text-xs font-medium" style={{ color: C.text, fontFamily: FONT_STACK }}>{b.label}</span>
            <span className="text-[9px] uppercase tracking-wider" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>{b.sub}</span>
            <span className="absolute top-1.5 right-1.5 text-[9px]" style={{ color: C.textGhost, fontFamily: MONO_STACK }}>{b.key}</span>
          </motion.button>
        );
      })}
    </div>
  );
}

function StatsView({ cards, srsData, sessionStats, onBack }) {
  const stats = useMemo(() => {
    const now = Date.now();
    let nNew = 0, nLearning = 0, nDue = 0, nLearned = 0, nDifficult = 0;
    cards.forEach(c => {
      const status = getCardStatus({ ...c, srs: srsData[c.id] }, now);
      if (status === 'new') nNew++;
      else if (status === 'learning') nLearning++;
      else if (status === 'due') nDue++;
      else if (status === 'learned') nLearned++;
      else if (status === 'difficult') nDifficult++;
    });
    return { nNew, nLearning, nDue, nLearned, nDifficult };
  }, [cards, srsData]);

  const byTag = useMemo(() => {
    const result = {};
    Object.keys(TAGS).filter(t => t !== 'all').forEach(t => {
      const tagCards = cards.filter(c => c.tag === t);
      const learned = tagCards.filter(c => getCardStatus({ ...c, srs: srsData[c.id] }) === 'learned').length;
      result[t] = { total: tagCards.length, learned, percent: tagCards.length ? Math.round(learned / tagCards.length * 100) : 0 };
    });
    return result;
  }, [cards, srsData]);

  return (
    <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0, y: -20 }} className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl tracking-tight" style={{ color: C.text, fontFamily: FONT_STACK, fontWeight: 600, letterSpacing: '-0.02em' }}>
          Статистика
        </h2>
        <button onClick={onBack} className="text-xs uppercase tracking-wider transition-colors"
          style={{ color: C.textDim, fontFamily: MONO_STACK }}
          onMouseEnter={e => e.currentTarget.style.color = C.text}
          onMouseLeave={e => e.currentTarget.style.color = C.textDim}
        >
          ← Назад
        </button>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-5 gap-4 p-6 rounded-2xl"
        style={{
          background: `linear-gradient(155deg, ${C.cardFrontHi} 0%, ${C.cardFront} 100%)`,
          border: `1px solid ${C.border}`
        }}
      >
        <StatRing value={stats.nLearned} max={cards.length} color={C.good} label="Выучено" sublabel="из 53" size={88} />
        <StatRing value={stats.nDue} max={cards.length} color={C.warn} label="Повторить" sublabel="cейчас" size={88} />
        <StatRing value={stats.nLearning} max={cards.length} color={C.accent} label="Учим" sublabel="" size={88} />
        <StatRing value={stats.nDifficult} max={cards.length} color={C.bad} label="Сложные" sublabel="" size={88} />
        <StatRing value={stats.nNew} max={cards.length} color={C.textDim} label="Новые" sublabel="" size={88} />
      </div>

      {sessionStats.reviewed > 0 && (
        <div className="p-5 rounded-xl flex items-center justify-between"
          style={{
            background: `linear-gradient(155deg, ${C.accent}10 0%, ${C.accent}05 100%)`,
            border: `1px solid ${C.accent}25`
          }}
        >
          <div className="flex items-center gap-3">
            <Flame size={20} style={{ color: C.accent }} />
            <div>
              <div className="text-xs uppercase tracking-wider" style={{ color: C.textDim, fontFamily: MONO_STACK }}>Эта сессия</div>
              <div className="text-lg mt-0.5" style={{ color: C.text, fontFamily: FONT_STACK, fontWeight: 500 }}>
                {sessionStats.reviewed} карточек · {sessionStats.correct} правильно
              </div>
            </div>
          </div>
          <div className="text-3xl tabular-nums" style={{ fontFamily: FONT_STACK, fontWeight: 600, color: C.accent, letterSpacing: '-0.03em' }}>
            {sessionStats.reviewed ? Math.round(sessionStats.correct / sessionStats.reviewed * 100) : 0}%
          </div>
        </div>
      )}

      <div>
        <h3 className="text-xs uppercase tracking-wider mb-3" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>По темам</h3>
        <div className="space-y-2">
          {Object.entries(byTag).map(([tag, data]) => {
            const info = TAGS[tag];
            return (
              <div key={tag} className="p-3 rounded-lg flex items-center gap-3"
                style={{ background: 'rgba(148, 163, 184, 0.03)', border: `1px solid ${C.border}` }}
              >
                <div className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: info.color, boxShadow: `0 0 6px ${info.color}80` }} />
                <div className="text-xs uppercase tracking-wider w-28" style={{ color: C.text, fontFamily: MONO_STACK }}>{info.label}</div>
                <div className="flex-1 h-1.5 rounded-full overflow-hidden" style={{ backgroundColor: 'rgba(148, 163, 184, 0.08)' }}>
                  <motion.div
                    initial={{ width: 0 }}
                    animate={{ width: `${data.percent}%` }}
                    transition={{ duration: 0.8, ease: 'easeOut' }}
                    className="h-full rounded-full"
                    style={{ backgroundColor: info.color, boxShadow: `0 0 8px ${info.color}60` }}
                  />
                </div>
                <div className="text-xs tabular-nums w-16 text-right" style={{ color: C.textDim, fontFamily: MONO_STACK }}>
                  {data.learned}/{data.total}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </motion.div>
  );
}

export default function App() {
  const [mode, setMode] = useState('study');
  const [studyMode, setStudyMode] = useState('smart');
  const [filter, setFilter] = useState('all');
  const [srsData, saveSrsData] = useStorage('srs-data-v1', {});
  const [sessionStats, setSessionStats] = useState({ reviewed: 0, correct: 0 });
  const [currentIdx, setCurrentIdx] = useState(0);
  const [flipped, setFlipped] = useState(false);
  const [queue, setQueue] = useState([]);

  useEffect(() => {
    const now = Date.now();
    let cards = CARDS.filter(c => filter === 'all' || c.tag === filter);

    if (studyMode === 'smart') {
      const due = cards.filter(c => srsData[c.id] && srsData[c.id].nextReview <= now);
      const newCards = cards.filter(c => !srsData[c.id] || srsData[c.id].reps === 0);
      const learning = cards.filter(c => srsData[c.id] && srsData[c.id].reps > 0 && srsData[c.id].nextReview > now);
      cards = [...due, ...newCards.slice(0, 5), ...learning];
    } else if (studyMode === 'new') {
      cards = cards.filter(c => !srsData[c.id] || srsData[c.id].reps === 0);
    } else if (studyMode === 'due') {
      cards = cards.filter(c => srsData[c.id] && srsData[c.id].nextReview <= now);
    } else if (studyMode === 'difficult') {
      cards = cards.filter(c => srsData[c.id] && (srsData[c.id].lapses >= 2 || srsData[c.id].ease < 2.0));
    }

    setQueue(cards);
    setCurrentIdx(0);
    setFlipped(false);
  }, [studyMode, filter, srsData]);

  const currentCard = queue[currentIdx];
  const cardWithSrs = currentCard ? { ...currentCard, srs: srsData[currentCard.id] } : null;
  const status = cardWithSrs ? getCardStatus(cardWithSrs) : null;

  const handleReview = useCallback((quality) => {
    if (!currentCard) return;
    const newSrs = calculateNextReview({ srs: srsData[currentCard.id] }, quality);
    const updated = { ...srsData, [currentCard.id]: newSrs };
    saveSrsData(updated);

    setSessionStats(s => ({
      reviewed: s.reviewed + 1,
      correct: s.correct + (quality >= 2 ? 1 : 0)
    }));

    setTimeout(() => {
      setFlipped(false);
      setTimeout(() => {
        if (currentIdx < queue.length - 1) {
          setCurrentIdx(currentIdx + 1);
        } else {
          setCurrentIdx(queue.length);
        }
      }, 100);
    }, 200);
  }, [currentCard, srsData, saveSrsData, currentIdx, queue.length]);

  useEffect(() => {
    const handler = (e) => {
      if (e.target.tagName === 'INPUT') return;
      if (mode !== 'study' || !currentCard) return;

      if (e.key === ' ' || e.key === 'ArrowUp' || e.key === 'ArrowDown') {
        e.preventDefault();
        setFlipped(f => !f);
      } else if (e.key === 'ArrowRight') {
        e.preventDefault();
        if (currentIdx < queue.length - 1) {
          setFlipped(false);
          setCurrentIdx(currentIdx + 1);
        }
      } else if (e.key === 'ArrowLeft') {
        e.preventDefault();
        if (currentIdx > 0) {
          setFlipped(false);
          setCurrentIdx(currentIdx - 1);
        }
      } else if (flipped && ['1','2','3','4'].includes(e.key)) {
        e.preventDefault();
        handleReview(parseInt(e.key) - 1);
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, [mode, flipped, currentCard, currentIdx, queue.length, handleReview]);

  const tagCounts = useMemo(() => {
    const counts = { all: CARDS.length };
    Object.keys(TAGS).filter(t => t !== 'all').forEach(t => {
      counts[t] = CARDS.filter(c => c.tag === t).length;
    });
    return counts;
  }, []);

  const learnedCount = Object.values(srsData).filter(s => s.reps >= 5).length;
  const dueCount = Object.values(srsData).filter(s => s.nextReview <= Date.now() && s.reps > 0).length;
  const sessionDone = currentIdx >= queue.length && queue.length > 0;

  const resetProgress = async () => {
    if (!confirm('Сбросить весь прогресс? Это удалит данные о повторениях.')) return;
    await saveSrsData({});
    setSessionStats({ reviewed: 0, correct: 0 });
  };

  return (
    <div className="min-h-screen flex flex-col" style={{
      background: `radial-gradient(ellipse 80% 60% at 50% 0%, ${C.bgGradTop} 0%, ${C.bg} 70%)`,
      fontFamily: FONT_STACK,
      color: C.text
    }}>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap');
        ::selection { background: ${C.accent}30; color: ${C.text}; }
        html { color-scheme: dark; }
      `}</style>

      {/* Subtle accent bar */}
      <div className="h-px w-full" style={{
        background: `linear-gradient(90deg, transparent, ${C.accent}80 30%, ${C.accentWarm}40 70%, transparent)`
      }} />

      <div className="max-w-3xl w-full mx-auto px-5 py-6 md:py-10 flex-1 flex flex-col">
        <header className="mb-8">
          <div className="flex items-start justify-between gap-4">
            <div>
              <div className="text-[10px] uppercase tracking-[0.25em] mb-2" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>
                AI-augmented разработка
              </div>
              <h1 className="text-3xl md:text-4xl tracking-tight leading-none" style={{
                fontFamily: FONT_STACK,
                fontWeight: 600,
                letterSpacing: '-0.03em',
                color: C.text
              }}>
                Карточки <span style={{ fontWeight: 400, color: C.accent }}>для запоминания</span>
              </h1>
            </div>
            <div className="flex gap-2">
              <button onClick={() => setMode(mode === 'stats' ? 'study' : 'stats')}
                className="p-2.5 rounded-lg border transition-all"
                style={{ borderColor: C.border, color: C.textDim }}
                onMouseEnter={e => { e.currentTarget.style.borderColor = C.borderHi; e.currentTarget.style.color = C.text; }}
                onMouseLeave={e => { e.currentTarget.style.borderColor = C.border; e.currentTarget.style.color = C.textDim; }}
                title="Статистика">
                <BarChart3 size={16} />
              </button>
              <button onClick={resetProgress}
                className="p-2.5 rounded-lg border transition-all"
                style={{ borderColor: C.border, color: C.textDim }}
                onMouseEnter={e => { e.currentTarget.style.borderColor = `${C.bad}40`; e.currentTarget.style.color = C.bad; }}
                onMouseLeave={e => { e.currentTarget.style.borderColor = C.border; e.currentTarget.style.color = C.textDim; }}
                title="Сбросить прогресс">
                <Trash2 size={16} />
              </button>
            </div>
          </div>

          {mode === 'study' && (
            <div className="flex items-center gap-4 mt-4 text-xs" style={{ fontFamily: MONO_STACK }}>
              <div className="flex items-center gap-1.5" style={{ color: C.textDim }}>
                <span style={{ color: C.good, textShadow: `0 0 6px ${C.good}80` }}>●</span>
                <span className="tabular-nums">{learnedCount}</span>
                <span className="opacity-60 uppercase tracking-wider">выучено</span>
              </div>
              {dueCount > 0 && (
                <div className="flex items-center gap-1.5" style={{ color: C.textDim }}>
                  <span style={{ color: C.warn, textShadow: `0 0 6px ${C.warn}80` }}>●</span>
                  <span className="tabular-nums">{dueCount}</span>
                  <span className="opacity-60 uppercase tracking-wider">к повторению</span>
                </div>
              )}
              {sessionStats.reviewed > 0 && (
                <div className="flex items-center gap-1.5 ml-auto" style={{ color: C.textDim }}>
                  <Flame size={11} style={{ color: C.accentWarm }} />
                  <span className="tabular-nums">{sessionStats.reviewed}</span>
                  <span className="opacity-60 uppercase tracking-wider">сессия</span>
                </div>
              )}
            </div>
          )}
        </header>

        <AnimatePresence mode="wait">
          {mode === 'stats' ? (
            <StatsView key="stats" cards={CARDS} srsData={srsData} sessionStats={sessionStats} onBack={() => setMode('study')} />
          ) : (
            <motion.div key="study" initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="flex-1 flex flex-col">
              <div className="flex gap-2 mb-3 flex-wrap">
                {[
                  { id: 'smart', label: 'Умный режим', icon: Brain },
                  { id: 'all', label: 'Все', icon: BookOpen },
                  { id: 'new', label: 'Новые', icon: Sparkles },
                  { id: 'due', label: 'Повторить', icon: Clock },
                  { id: 'difficult', label: 'Сложные', icon: Flame }
                ].map(m => {
                  const Icon = m.icon;
                  const active = studyMode === m.id;
                  return (
                    <button key={m.id} onClick={() => setStudyMode(m.id)}
                      className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs uppercase tracking-wider transition-all"
                      style={{
                        fontFamily: MONO_STACK,
                        ...(active
                          ? { backgroundColor: C.accent, color: C.bg, boxShadow: `0 0 20px ${C.accent}30` }
                          : { color: C.textDim, background: 'transparent' }
                        )
                      }}
                      onMouseEnter={e => {
                        if (!active) { e.currentTarget.style.color = C.text; e.currentTarget.style.background = 'rgba(148, 163, 184, 0.05)'; }
                      }}
                      onMouseLeave={e => {
                        if (!active) { e.currentTarget.style.color = C.textDim; e.currentTarget.style.background = 'transparent'; }
                      }}
                    >
                      <Icon size={11} strokeWidth={2.2} />
                      {m.label}
                    </button>
                  );
                })}
              </div>

              <div className="flex gap-1.5 mb-6 flex-wrap">
                {Object.keys(TAGS).map(t => (
                  <TagPill key={t} tag={t} active={filter === t} count={tagCounts[t]} onClick={() => setFilter(t)} />
                ))}
              </div>

              <div className="flex-1 flex flex-col justify-center">
                {sessionDone ? (
                  <motion.div initial={{ scale: 0.9, opacity: 0 }} animate={{ scale: 1, opacity: 1 }}
                    className="text-center py-16 px-8 rounded-2xl"
                    style={{
                      background: `linear-gradient(155deg, ${C.cardFrontHi} 0%, ${C.cardFront} 100%)`,
                      border: `1px solid ${C.accent}25`,
                      boxShadow: `0 0 60px -20px ${C.accent}40`
                    }}
                  >
                    <Flame size={48} style={{ color: C.accentWarm }} className="mx-auto mb-4" strokeWidth={1.5} />
                    <h2 className="text-2xl tracking-tight mb-2" style={{ fontFamily: FONT_STACK, fontWeight: 600, color: C.text, letterSpacing: '-0.02em' }}>
                      Сессия завершена
                    </h2>
                    <p className="mb-6" style={{ color: C.textDim }}>
                      Прошёл {sessionStats.reviewed} {sessionStats.reviewed === 1 ? 'карточку' : 'карточек'} ·
                      {' '}{Math.round(sessionStats.correct / sessionStats.reviewed * 100)}% правильно
                    </p>
                    <button onClick={() => { setCurrentIdx(0); setSessionStats({ reviewed: 0, correct: 0 }); }}
                      className="px-6 py-3 rounded-lg text-xs uppercase tracking-wider transition-all"
                      style={{
                        backgroundColor: C.accent,
                        color: C.bg,
                        fontFamily: MONO_STACK,
                        boxShadow: `0 0 24px ${C.accent}40`
                      }}
                      onMouseEnter={e => e.currentTarget.style.backgroundColor = C.accentHot}
                      onMouseLeave={e => e.currentTarget.style.backgroundColor = C.accent}
                    >
                      Ещё одну сессию
                    </button>
                  </motion.div>
                ) : queue.length === 0 ? (
                  <div className="text-center py-16 px-8 rounded-2xl"
                    style={{
                      background: `linear-gradient(155deg, ${C.cardFrontHi} 0%, ${C.cardFront} 100%)`,
                      border: `1px solid ${C.border}`
                    }}
                  >
                    <Check size={48} style={{ color: C.good }} className="mx-auto mb-4" strokeWidth={1.5} />
                    <h2 className="text-2xl tracking-tight mb-2" style={{ fontFamily: FONT_STACK, fontWeight: 600, color: C.text, letterSpacing: '-0.02em' }}>
                      Всё чисто
                    </h2>
                    <p style={{ color: C.textDim }}>
                      В этой категории нет карточек для этого режима. Попробуй другой фильтр.
                    </p>
                  </div>
                ) : (
                  <AnimatePresence mode="wait">
                    <motion.div key={currentCard?.id} initial={{ opacity: 0, x: 20 }} animate={{ opacity: 1, x: 0 }} exit={{ opacity: 0, x: -20 }} transition={{ duration: 0.25 }}>
                      <FlashCard card={currentCard} flipped={flipped} onFlip={() => setFlipped(!flipped)} status={status} />
                    </motion.div>
                  </AnimatePresence>
                )}
              </div>

              {!sessionDone && queue.length > 0 && (
                <div className="mt-6 space-y-4">
                  <div className="flex items-center gap-3">
                    <button onClick={() => { if (currentIdx > 0) { setFlipped(false); setCurrentIdx(currentIdx - 1); } }}
                      disabled={currentIdx === 0}
                      className="p-2 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                      style={{ color: C.textFaint }}
                      onMouseEnter={e => { if (currentIdx > 0) e.currentTarget.style.color = C.text; }}
                      onMouseLeave={e => e.currentTarget.style.color = C.textFaint}
                    >
                      <ChevronLeft size={18} />
                    </button>
                    <div className="flex-1">
                      <div className="flex items-center justify-between text-[10px] uppercase tracking-wider mb-1" style={{ color: C.textFaint, fontFamily: MONO_STACK }}>
                        <span className="tabular-nums">{currentIdx + 1} / {queue.length}</span>
                      </div>
                      <div className="h-1 rounded-full overflow-hidden" style={{ backgroundColor: 'rgba(148, 163, 184, 0.08)' }}>
                        <motion.div className="h-full rounded-full"
                          style={{ backgroundColor: C.accent, boxShadow: `0 0 8px ${C.accent}80` }}
                          initial={{ width: 0 }}
                          animate={{ width: `${((currentIdx + 1) / queue.length) * 100}%` }}
                          transition={{ duration: 0.3 }}
                        />
                      </div>
                    </div>
                    <button onClick={() => { if (currentIdx < queue.length - 1) { setFlipped(false); setCurrentIdx(currentIdx + 1); } }}
                      disabled={currentIdx >= queue.length - 1}
                      className="p-2 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                      style={{ color: C.textFaint }}
                      onMouseEnter={e => { if (currentIdx < queue.length - 1) e.currentTarget.style.color = C.text; }}
                      onMouseLeave={e => e.currentTarget.style.color = C.textFaint}
                    >
                      <ChevronRight size={18} />
                    </button>
                  </div>

                  <motion.div initial={false} animate={{ opacity: flipped ? 1 : 0.3, pointerEvents: flipped ? 'auto' : 'none' }} transition={{ duration: 0.3 }}>
                    <ReviewButtons onReview={handleReview} disabled={!flipped} />
                  </motion.div>
                </div>
              )}
            </motion.div>
          )}
        </AnimatePresence>

        <footer className="mt-8 pt-4 text-center" style={{ borderTop: `1px solid ${C.border}` }}>
          <div className="text-[10px] uppercase tracking-wider space-x-3" style={{ color: C.textGhost, fontFamily: MONO_STACK }}>
            <span>← →  переключение</span>
            <span>пробел  перевернуть</span>
            <span>1—4  оценка</span>
          </div>
        </footer>
      </div>
    </div>
  );
}
