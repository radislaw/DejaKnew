# BRD: DéjàKnew — Business Requirements Document

| Поле | Значение |
|---|---|
| Версия | v1.0 |
| Дата | 2026-05-29 |
| Статус | Draft |
| Владелец | Product |
| Связанные документы | `spec.md`, `architecture.md`, `docs/product-roadmap.md` |

---

## TL;DR

DéjàKnew — RPG-экосистема обучения с интервальным повторением (SRS) для разработчиков, студентов и ADHD-пользователей. Проблема: инструменты типа Anki технически сильны, но эмоционально мёртвы — нет dopamine-loop, нет прогрессии, нет поддержки ADHD. Решение: игровая механика (персонажи, «Нейроны», Silent Rooms, Creator Economy) поверх проверенного SRS. Монетизация: freemium → подписка → косметика → курсы → creator revenue share. Цель на 3 года: 100-150k WAL, $2-3M ARR.

---

## Оглавление

1. [Executive Summary](#1-executive-summary)
2. [Анализ рынка](#2-анализ-рынка)
3. [Целевые персоны](#3-целевые-персоны)
4. [Бизнес-модель и монетизация](#4-бизнес-модель-и-монетизация)
5. [Unit-экономика](#5-unit-экономика)
6. [Приоритизация Phase 0 (RICE)](#6-приоритизация-phase-0-rice)
7. [Риски и открытые вопросы](#7-риски-и-открытые-вопросы)
8. [Go-to-Market стратегия](#8-go-to-market-стратегия)

---

## 1. Executive Summary

### 1.1 Продукт и видение

**DéjàKnew** — персонализированная RPG-экосистема обучения, построенная на алгоритме интервального повторения (Spaced Repetition System, SRS). Пользователь создаёт или импортирует карточки, проходит сессии повторения, зарабатывает внутреннюю валюту «Нейроны», открывает персонажей и скины, движется по RPG-прогрессии и может совместно работать в Silent Rooms без давления.

**Видение:** стать default-инструментом обучения для технической и ADHD-аудитории — тех, кто нашёл Anki слишком скучным, Duolingo слишком детским, а Notion — слишком беспорядочным. Tagline: «Учи умнее, не дольше».

**Версии:** PWA (Phase 0), нативные iOS/Android (планируются после Phase 2), web-first.

### 1.2 Проблема

Существующие SRS-инструменты страдают тремя фундаментальными дефектами:

1. **Эмоциональная пустота.** Anki — клинически эффективен, визуально мёртв. Нет прогрессии, нет награды, нет смысла возвращаться завтра кроме чувства долга.
2. **ADHD-враждебность.** Большинство EdTech-продуктов игнорируют нейроразнообразие: нет коротких сессий, нет мягкого старта, есть streak-shame и перегруженные интерфейсы.
3. **Изоляция.** Учиться в одиночестве тяжело. Body Doubling — доказанная техника для ADHD — не реализована ни в одном популярном SRS-инструменте.

### 1.3 Уникальное ценностное предложение

> DéjàKnew — единственный SRS-инструмент с RPG-прогрессией, ADHD-first UX (Silent Rooms, Focus Mode, микро-сессии) и Creator Economy поверх курсов, ориентированный на техническую аудиторию с dark/terminal эстетикой.

**Ключевые дифференциаторы:**
- Валюта «Нейроны» создаёт dopamine-loop без реальных денег в базовом флоу.
- Silent Rooms (Body Doubling без видео/аудио) — функция, которой нет ни у одного прямого конкурента.
- Focus Mode: 3 кнопки, 5-10 минутные сессии, никакого shame-фидбека.
- Creator Economy с revenue share — превращает продукт в платформу.
- Позиционирование в dev-нише: Code blocks, импорт из RFC/документации, terminal-эстетика.

### 1.4 Целевой рынок и аудитория

**Рынок:** EdTech / SRS / microlearning / developer upskilling.

**Первичная аудитория (Phase 0-2):**
- Разработчики 22-35 лет, изучающие новые технологии, языки программирования, готовящиеся к интервью.
- Студенты 18-25 лет (технические и гуманитарные специальности, языки, медицина).
- ADHD-пользователи 18-32 лет, которые начинали и бросали десятки инструментов.

**География Phase 0:** RU + EN-говорящие пользователи (СНГ + Европа + США/Канада).

**География Phase 2+:** глобальный выход с EN-контентом и локализациями.

### 1.5 Бизнес-модель (краткая)

Freemium-модель с тремя параллельными потоками дохода: **подписка Pro** (безлимитная AI-генерация карточек, отключение рекламы, полный доступ к курсам), **микротранзакции** (покупка «Нейронов» и косметики за реальные деньги нетерпеливыми пользователями), **rewarded ads** (пользователь сам выбирает посмотреть рекламу ради бонуса). В Phase 6 добавляется **Creator Economy** — платформенный сбор 30% с продаж пользовательских курсов.

---

## 2. Анализ рынка

### 2.1 TAM / SAM / SOM

**TAM — Total Addressable Market**

Глобальный рынок EdTech оценивается в $340-400B к 2027 году (HolonIQ Global EdTech Report 2024). Сегмент self-paced / microlearning, к которому относится SRS, составляет ~$80-90B.

Более узко: рынок language learning apps — $21B к 2027 (Grand View Research, 2024 CAGR 18.7%). Рынок flashcard / spaced repetition software как отдельный сегмент оценить сложнее, но совокупная аудитория Anki (40M+ скачиваний), Quizlet (500M+ зарегистрированных), Brainscape (~8M) и RemNote (~1M) указывает на десятки миллионов активных пользователей.

**Итого TAM (оценка):** $60-80B (SRS + microlearning + developer upskilling вместе).

**SAM — Serviceable Addressable Market**

Целевые сегменты: (1) разработчики, использующие SRS-инструменты для технических знаний, (2) студенты технических/медицинских специальностей, (3) ADHD-пользователи в возрасте 18-35.

По данным Stack Overflow Developer Survey 2024: ~26M профессиональных разработчиков глобально, из которых ~15-20% активно используют инструменты spaced repetition или flashcards. ADHD-аудитория: ~5-7% взрослого населения, из которых значительная часть пересекается с tech/student аудиторией.

Реалистичная оценка SAM с учётом языкового охвата (RU/EN Phase 0-1, +6 языков Phase 3+) и платёжеспособности: **$3-5B** (SRS + dev upskilling подсегмент).

**SOM — Serviceable Obtainable Market (горизонт 3 года)**

| Горизонт | Пользователи (MAU) | WAL (Weekly Active Learners) | ARR (оценка) |
|---|---|---|---|
| 12 месяцев (Phase 1) | 15 000-40 000 | 8 000-20 000 | $80K-250K |
| 24 месяца (Phase 2-3) | 60 000-120 000 | 35 000-70 000 | $700K-1.5M |
| 36 месяцев (Phase 4-5) | 120 000-250 000 | 70 000-150 000 | $1.8M-4M |

Основания: Duolingo вышел на 1M MAU через 14 месяцев после публичного запуска при массовом маркетинге; для нишевого dev-продукта реалистично 10x медленнее при органическом росте. Конверсия free→paid 4-8% при среднем чеке $5-6/мес.

### 2.2 Конкурентная матрица

| Параметр | Anki | Duolingo | Quizlet | RemNote | Brainscape | **DéjàKnew** |
|---|---|---|---|---|---|---|
| **Free tier** | Полный (desktop) | Да (с ограничениями) | Да (с ограничениями) | Да (limited) | Да (limited) | **Да (базовый)** |
| **Мобильное приложение** | iOS платно, Android бесплатно | iOS + Android | iOS + Android | iOS + Android | iOS + Android | **PWA Phase 0, нативно Phase 2+** |
| **Веб-версия** | AnkiWeb (базовая) | Да | Да | Да | Да | **Да (PWA-first)** |
| **SRS-алгоритм** | SM-2 / FSRS | Нет (curriculum) | Упрощённый | SM-2 / FSRS | SuperMemo | **SM-2 Phase 0, апгрейд Phase 2** |
| **ADHD-friendly** | Нет | Частично | Нет | Нет | Нет | **Да (Focus Mode, Silent Rooms)** |
| **Геймификация** | Нет | Да (streak, лиги) | Минимальная | Нет | Нет | **Да (RPG-прогрессия, Нейроны)** |
| **AI-генерация карточек** | Нет | Нет (AI адаптация) | Да (Magic Notes) | Да (AI Flashcards) | Нет | **Да (Phase 2, freemium)** |
| **Социальные функции** | Shared decks | Друзья, лиги | Shared sets, live | Нет | Нет | **Silent Rooms, пары (Phase 4)** |
| **Creator Economy** | Нет | Нет | Нет | Нет | Нет | **Да (Phase 6, revenue share)** |
| **Курсы/уроки** | Нет | Да (встроенные) | Нет | Нет | Нет | **Да (Phase 3, интерактивные)** |
| **Dev/код контент** | Community decks | Нет | Нет | Частично | Нет | **Да (Code blocks, импорт из docs)** |
| **Монетизация** | Платное iOS-приложение | Подписка $6.99/мес | Подписка $7.99/мес | Подписка $6/мес | Подписка $9.99/мес | **Freemium + сubs + косметика + ads** |
| **Desktop-приложение** | Да (нативное) | Нет | Нет | Да (Electron) | Нет | **Нет (PWA)** |

### 2.3 Конкурентные преимущества и уязвимости

**Где DéjàKnew выигрывает:**

- **ADHD-ниша не занята.** Ни один из конкурентов не строит продукт с ADHD-first архитектурой. Это редкая возможность занять категорию.
- **Dev-нища плохо обслужена.** Anki технически гибкий, но UI 90-х. RemNote — сложный. DéjàKnew может стать "Anki for developers, but good".
- **Silent Rooms = уникальный функционал.** Body Doubling доказанно работает для ADHD и концентрации. Ни один SRS-конкурент это не реализовал.
- **RPG-прогрессия vs геймификация Duolingo.** Duolingo — геймификация для детей (лиги, птица-тиран). DéjàKnew — RPG-прогрессия для взрослых (персонажи, косметика, RPG skill tree).
- **Creator Economy как долгосрочный ров.** Platfrom lock-in через контент сложнее всего скопировать.
- **Нет "shame loop".** Duolingo печально известен негативным подкреплением (агрессивная птица). DéjàKnew строится на positive reinforcement-only.

**Где DéjàKnew уязвим:**

- **SRS-алгоритм:** Anki реализовал FSRS (Free Spaced Repetition Scheduler) в 2023, который значительно точнее SM-2. DéjàKnew стартует с SM-2 — нужен план апгрейда до FSRS или собственного алгоритма.
- **Нет desktop-приложения в Phase 0.** Часть серьёзных пользователей (medstudents) привязана к Anki Desktop.
- **Нет established creator base и community.** Anki и Quizlet имеют миллионы shared decks. DéjàKnew стартует с пустой библиотекой.
- **Brand awareness.** Против Duolingo (~600M зарегистрированных) — незначительный игрок, требуется нишевый positioning.
- **AI costs.** При бесплатном лимите AI-генерации costs могут превысить revenue на ранних этапах.
- **PWA-ограничения iOS.** Apple ограничивает push-notifications и background sync в PWA на iOS — влияет на retention механики.

### 2.4 Positioning Statement

> Для разработчиков, студентов и людей с ADHD, которым традиционные SRS-инструменты кажутся выгорающими и скучными, **DéjàKnew** — это персонализированная RPG-экосистема обучения, которая, в отличие от Anki и Quizlet, превращает ежедневное повторение в игровой прогресс с персонажами, «Нейронами» и Silent Rooms — без давления, без негативного подкрепления, с dark/terminal эстетикой, созданной для технической аудитории.

---

## 3. Целевые персоны

### 3.1 Дмитрий — Senior Backend Developer, 28 лет

**Профиль:**
- Работает в продуктовой компании, стек Go/Rust, 5 лет опыта.
- Хочет вырасти до Staff Engineer, активно готовится к интервью в FAANG-like компанию.
- Изучает Kubernetes, system design, читает RFC и технические статьи.
- Использует Anki для технических карточек, но постоянно забрасывает из-за некрасивого UI и отсутствия мотивации продолжать.
- Работает в режиме deep work: 2-3 часа без отвлечений, но нуждается в инструментах, которые помогают фокусироваться.

**Цели:**
- Эффективно запоминать технические концепции: алгоритмы, паттерны, API.
- Готовиться к Staff-интервью без многочасовых boring-сессий.
- Не тратить время на настройку инструментов (Anki требует кастомных шаблонов).

**Боли:**
- Anki UX — "это карточки из 1990-х". Нет visual feedback, нет прогресса, нет смысла возвращаться.
- Нет качественного dev-контента (Code blocks, syntax highlighting, импорт из документации).
- Теряет streak, чувствует вину, бросает на месяц, начинает заново.
- AI-генерация в Quizlet генерирует мусор для технических тем.

**Как использует DéjàKnew:**
- Импортирует существующие Anki-деки с алгоритмами.
- Использует AI-генерацию для создания карточек из статей и RFC.
- Code blocks в карточках и уроках для читабельного кода.
- Focus Mode для 25-минутных session без уведомлений.
- Silent Rooms для совместного учёбы с коллегами (accountability без звонков).

**Готовность платить:** $8-12/мес за Pro (подписка снимает боль AI-лимитов + no ads). Купит 1-2 скина/персонажа за реальные деньги если они "cool, not childish" (~$5-10 разово).

**Ключевые фичи для Дмитрия:**
- AI-генерация карточек из markdown/URL/RFC.
- Code blocks с syntax highlighting в карточках.
- Focus Mode.
- Import из Anki (.apkg).
- Silent Rooms (Phase 4).
- Terminal/dark UI — обязательно.

---

### 3.2 Анна — Студентка медицинского факультета, 21 год

**Профиль:**
- 3-й курс медфака, изучает анатомию, биохимию, фармакологию и латынь.
- Параллельно учит английский для сдачи USMLE в будущем.
- Использует Anki с готовыми медицинскими деками (Anking, Zanki — тысячи карточек).
- Телефон — основное устройство обучения, ноутбук для учёбы в библиотеке.
- Живёт в режиме: сессия (дедлайны и ужас) → каникулы (забрасывает всё).

**Цели:**
- Сдать экзамен по анатомии без выгорания от 2000 карточек.
- Поддерживать streak хотя бы 5 дней в неделю.
- Учиться в паре с подругой (accountability).

**Боли:**
- Anki на телефоне — неудобный, тормозит, синхронизация ломается.
- Тысячи карточек создают anxiety, а не мотивацию.
- Пропускает несколько дней → карточки накапливаются → бросает.
- Нет элемента "социального" обучения в инструментах.

**Как использует DéjàKnew:**
- Импортирует медицинские деки из Anki (.apkg).
- Ежедневные сундуки дают повод открыть приложение.
- Streak freeze в период сессии снимает anxiety ("знаю, что выживу").
- Дружеская пара с подругой — взаимный accountability.
- PWA на iPhone — быстрый доступ с домашнего экрана.

**Готовность платить:** $3-5/мес (студенческий бюджет, но "это моя учёба, а не развлечение"). Предпочтёт earning нейронов через rewarded ads и ежедневные сундуки, чем платить реальными деньгами за косметику. Купит streak freeze за реальные деньги в критичный момент (~$1-2 разово).

**Ключевые фичи для Анны:**
- Imprt из Anki (.apkg) с сохранением форматирования.
- Геймификация (ежедневные сундуки, streak).
- Streak freeze (критично для retention в период экзаменов).
- Push-уведомления в нужное время.
- Дружеские пары (Phase 4).
- Быстрый мобильный PWA.

---

### 3.3 Михаил — Junior Data Analyst, 24 года, ADHD

**Профиль:**
- Диагностирован с ADHD-combined в 19 лет. Принимает терапию, но продуктивность нестабильна.
- Работает junior data analyst в стартапе, хочет изучить Python, SQL, ML-концепции.
- Начинал учиться в Anki, Duolingo, Notion, Obsidian — бросал все через 2-3 недели.
- Основная проблема: не может начать сессию ("paralysis by analysis"), теряет нить через 8 минут, guilt при пропуске нарастает → бросает инструмент.
- Активно присутствует в ADHD-сообществах (Reddit r/ADHD, Telegram-чаты).

**Цели:**
- Удержаться в обучении дольше, чем обычно бросает.
- Изучить Python/SQL для роста в карьере.
- Найти инструмент, который не вызывает shame при пропуске.

**Боли:**
- Традиционный SRS требует дисциплины, которой нет по определению при ADHD.
- "Накопленные карточки" — главный триггер для бросания инструмента.
- Streak-shame: Duolingo и Anki наказывают за пропуск, что усиливает avoidance.
- Сложные интерфейсы перегружают рабочую память.
- Самостоятельная работа невыносима без внешней структуры.

**Как использует DéjàKnew:**
- **Focus Mode** — 3 кнопки, минималистичный экран, никаких отвлечений.
- **Микро-сессии** 5-10 минут: успех = "я сделал сессию", не "я сделал 50 карточек".
- **Silent Rooms** — body doubling с незнакомцами. Присутствие других людей создаёт структуру.
- **No-shame streak freeze** — потерял streak? Freeze куплен/earned заранее, нет guilt-экрана.
- **AI-компаньон** (Phase 5) — виртуальный body doubling companion.
- Нейроны через ежедневные сундуки и короткие rewarded ads.

**Готовность платить:** $5-8/мес если подписка снимает конкретную боль (no-shame UX + full Silent Rooms + AI companion). Купит "персонажа-компаньона" за реальные деньги (~$4-9) если это снижает anxiety от обучения. ADHD-пользователи часто готовы платить за "инструменты, которые наконец работают" — но WTP нестабильна из-за финансовой дисорганизации.

**Ключевые фичи для Михаила:**
- Focus Mode (3 кнопки, 5-10 мин микро-сессии, Phase 5).
- Silent Rooms (Body Doubling без видео/аудио, Phase 4).
- No-shame streak mechanics (freeze, soft restart).
- Positive reinforcement ONLY — никаких "ты пропустил 47 дней".
- AI-компаньон (Phase 5).
- Push-уведомления: мягкие, не guilt-индуцирующие.

---

## 4. Бизнес-модель и монетизация

### 4.1 Revenue Streams по фазам

| Фаза | Revenue Stream | Тип | Модель | Потенциал |
|---|---|---|---|---|
| Phase 0 | Нет монетизации | — | Инвестиции в retention-инфраструктуру | — |
| Phase 1 | Косметика за реальные деньги (скины, персонажи) | One-time / consumption | Прямая покупка | Низкий ($500-5K/мес) |
| Phase 1 | Rewarded ads (видеореклама по выбору) | Consumption / impressions | Revenue share с рекламной сетью | Низкий-средний ($200-2K/мес) |
| Phase 2 | Подписка Pro (AI gen + no ads + полный доступ) | Recurring, monthly/annual | Freemium gating | **Высокий** (основной поток) |
| Phase 3 | Платные курсы за Нейроны / для подписчиков | Freemium gating | Consumption + subscription | Средний |
| Phase 4 | Премиум Silent Rooms (расширенные опции) | Recurring | Upsell для подписчиков | Низкий |
| Phase 5 | AI-компаньон (Body Doubling, Phase 5) | Recurring / add-on | Premium feature | Средний (ADHD-ниша) |
| Phase 6 | Creator Economy: платформенный сбор 30% с продаж курсов | Platform fee | Revenue share | **Высокий долгосрочно** |
| Phase 6 | Promoted courses / creator tools | B2B SaaS-like | Subscription для creators | Средний |

**Примечание:** Rewarded ads работают как валюта "заработка Нейронов", а не как баннерная реклама. Это снижает раздражение пользователей и повышает вовлечённость при просмотре рекламы.

### 4.2 Freemium Воронка

```
Visitor (SEO / ProductHunt / referral / community)
    ↓
Signup (email / Google / Apple)
    ↓
Onboarding Flow (выбор режима: обычный / Focus; первая колода)
    ↓
Activated User (≥1 колода создана, ≥1 сессия завершена)
    — цель: 40% от signups → activated (бенчмарк EdTech: 35-50%)
    ↓
D7-Retained User (вернулся на 7-й день)
    — цель: 25% D7 (бенчмарк SRS apps: 20-30%)
    ↓
Free Habit User (3+ сессии/неделю в течение 4+ недель)
    — зарабатывает Нейроны, открывает сундуки, попадает в Silent Rooms
    ↓
Monetization Trigger — один из:
    • AI-лимит исчерпан (5 генераций / 20 в мес → нужен Pro)
    • Rewarded ad → покупает Нейроны для ускорения
    • Понравился скин → покупает за реальные деньги
    • Хочет полный курс → upgrade или buy за Нейроны
    ↓
Paid Trial / First Purchase
    ↓
Subscriber / Repeat Buyer
```

**Целевые конверсии по воронке:**

| Этап | Conversion | Бенчмарк |
|---|---|---|
| Visitor → Signup | 15-25% | EdTech median 10-20% |
| Signup → Activated | 40% | EdTech median 35-50% |
| Activated → D7 retained | 25% | SRS apps 20-30% |
| D7 → D30 retained | 40% | EdTech 35-45% |
| Free → Paid (M3+) | 4-8% | EdTech freemium 3-8% |
| One-time buyer → Subscriber | 30% | Gaming / edtech 25-40% |

### 4.3 Pricing Strategy

**Подписка Pro:**

| Тип | Цена EN | Цена RU | Заметки |
|---|---|---|---|
| Monthly | $6.99/мес | 499 ₽/мес | Бенчмарк: Duolingo Super $6.99, Quizlet Plus $7.99 |
| Annual | $49.99/год ($4.17/мес) | 3 490 ₽/год | Скидка ~40% vs monthly |
| Lifetime (limited drops) | $149 | 10 990 ₽ | Drops 2-4 раза в год, ограниченный доступ |
| Student discount | $3.99/мес | 299 ₽/мес | Верификация через .edu email / ISIC |

**Нейроны (внутренняя валюта):**

| Пакет | Цена | Нейроны | Цена за 100 Нейронов |
|---|---|---|---|
| Стартовый | $0.99 | 1 000 | $0.099 |
| Стандартный | $4.99 | 5 500 (+10% бонус) | $0.091 |
| Популярный | $9.99 | 12 000 (+20% бонус) | $0.083 |
| Крупный | $19.99 | 26 000 (+30% бонус) | $0.077 |

**Косметика и персонажи:**

| Категория | Цена в Нейронах | Цена за реальные деньги | Заметки |
|---|---|---|---|
| Базовый скин (UI theme) | 500-1 500 Н | — (только за Нейроны) | Зарабатывается через обучение |
| Персонаж (common) | 1 500-3 000 Н | $2.99 | Покупка ускорения |
| Персонаж (rare) | 5 000-8 000 Н | $4.99-6.99 | |
| Персонаж (legendary) | — | $9.99 | Только за реальные деньги |
| Эксклюзив из сундука | — | — | Только через сундуки (не продаётся) |
| Streak Freeze | 200 Н | $0.99 | One-time |
| AI-генерация (доп. лимит) | 1 000 Н за 10 генераций | $0.99 | Для бесплатных пользователей |

**Курсы (Phase 3+):**

| Доступ | Условие |
|---|---|
| 5 уроков из курса | Бесплатно |
| Полный курс (платный) | 5 000-15 000 Нейронов ИЛИ включён в Pro |
| Creator курс | Устанавливает автор (в рамках платформенных лимитов) |

**Ценовая философия:** "Плати деньгами (быстро) или временем и усердием (валюта за обучение)." Нейроны зарабатываются через честное обучение — это защищает от P2W восприятия.

### 4.4 Monetization Risks

| Риск | Описание | Митигация |
|---|---|---|
| Ad-каннибализация подписки | Rewarded ads слишком удобны → нет смысла платить | Лимит rewarded ads/день (3-5 шт), подписка убирает ads полностью |
| Shame-paywall для ADHD | Paywall в середине сессии = ключевой триггер abandonment | Paywall только после завершения сессии; grace period |
| Apple/Google policy | Apple IAP обязателен для in-app purchases на iOS, 30% комиссия | Web-PWA обходит IAP в Phase 0; нативное приложение потребует Review + RevenueCat |
| RU-платёжный рынок | Санкции ограничивают Stripe для RU-пользователей | Параллельно с Stripe — ЮKassa / Robokassa для RU-billing |
| Нейроны как "азартная игра" | Регуляторные риски в некоторых юрисдикциях | Нейроны не конвертируются в реальные деньги; сундуки не содержат loot box за реальные деньги |
| Creator fraud (Phase 6) | Некачественный контент, плагиат, фродовые выплаты | Manual moderation Phase 4, automated content quality scoring Phase 6 |
| AI costs > revenue | При бесплатном лимите 20 AI gen/мес затраты превышают доход | Лимит должен быть экономически обоснован: ~$0.02/генерация × 20 = $0.40/пользователь/мес |

---

## 5. Unit-экономика

> Все цифры — оценки на основе публичных данных из Duolingo investor decks (2023-2024), RevenueCat State of Subscription Apps 2024, Sensor Tower EdTech benchmarks 2024, и анализа аналогичных freemium EdTech продуктов. Требуют валидации на реальных данных.

### 5.1 ARPU (Average Revenue Per User)

**Paying Users (ARPPU — Average Revenue Per Paying User):**

| Когорта | ARPPU Monthly | Обоснование |
|---|---|---|
| Pro Subscriber (monthly) | $6.99 | По цене тарифа |
| Pro Subscriber (annual, amortized) | $4.17 | $49.99 / 12 |
| Pro Subscriber blended | ~$5.50 | 60% monthly / 40% annual |
| One-time cosmetic buyer | $4-7 разово | Medianа in-app purchases EdTech |
| Rewarded ads viewer | $0.05-0.10/мес | $2-4 RPM × 1-2 показа/день × 30 дней / 1000 |

**ARPU All Users (blended):**

| Сценарий | Conversion Paid | ARPU (all users) |
|---|---|---|
| Conservative (Year 1) | 3% | $0.16/мес |
| Base (Year 2) | 5% | $0.28/мес |
| Target (Year 3) | 8% | $0.44/мес |

Для сравнения: Duolingo ARPU всех MAU (Q4 2023) — ~$0.60/мес при ~88M MAU. Для нишевого продукта с более высоким качеством аудитории (tech/students) ARPU может быть выше при меньшей базе.

### 5.2 LTV (Lifetime Value)

Формула: `LTV = ARPPU × (1 / monthly_churn_rate)`

| Когорта | ARPPU | Monthly Churn | Avg Lifetime | LTV |
|---|---|---|---|---|
| Free пользователь | $0.05 (ads) | 15-25% | 4-7 мес | $0.20-0.35 |
| Subscriber (monthly) | $6.99 | 8-10% | 10-12 мес | $70-87 |
| Subscriber (annual) | $4.17 | 25-35% (annual renewal) | 14-18 мес | $58-75 |
| Subscriber blended | $5.50 | ~8% | ~12 мес | **$66** |
| One-time cosmetics buyer | $6 разово | — | — | $6-9 (с учётом 30% repeat) |

**Benchmark:** RevenueCat (2024) — медианный LTV для EdTech subscription apps: $45-80. Duolingo ARPU × avg lifetime ≈ $50-70 для подписчиков. Оценка в $66 находится в реалистичном диапазоне.

### 5.3 CAC (Customer Acquisition Cost)

| Канал | CAC | Conversion Rate | Заметки |
|---|---|---|---|
| SEO / органический поиск | $0.50-2 | Высокая (intent-based) | Инвестиция в контент на 6-12 мес |
| ProductHunt launch | $0.20-1 (по времени команды) | Средняя | Разовый буст при launch |
| GitHub / dev communities | $0.50-3 | Высокая (нишевая) | README badges, dev.to, Habr |
| Reddit ads (r/ADHD, r/learnprogramming) | $4-8 | Средняя | Нишевый таргетинг |
| TikTok edu (organic) | $0.10-0.50 | Низкая-средняя | Требует контент-стратегию |
| TikTok / Instagram paid | $8-15 | Низкая | Broad audience |
| Партнёрства (bootcamps, university) | $1-3 | Средняя | B2B2C |
| ADHD-инфлюенсеры | $2-6 | Высокая (trust-based) | Nano/micro-influencers |

**Blended CAC Target:** <$5 (Phase 0-1), <$8 (Phase 2+, с платными каналами).

**LTV:CAC Ratio Target:** ≥10:1 для organic, ≥5:1 для blended. При LTV $66 и blended CAC $5 → ratio 13:1 — здоровый бизнес.

### 5.4 Payback Period

| Сценарий | CAC | ARPPU | Payback |
|---|---|---|---|
| Organic подписчик | $2 | $5.50/мес | **<1 мес** |
| Paid подписчик | $8 | $5.50/мес | ~1.5 мес |
| Blended (60% organic / 40% paid) | $4.40 | $5.50/мес | **<1 мес** |

Payback на уровне <2 месяцев — отличный показатель для B2C SaaS. Дополнительная маржа — rewarded ads и косметика от non-subscribers.

### 5.5 Churn Rate

| Когорта | Monthly Churn | Обоснование |
|---|---|---|
| Free users (MAU decay) | 15-20%/мес | Норма для бесплатных EdTech пользователей |
| Monthly subscribers | 8-10% | Duolingo ~8%, language apps median 7-12% |
| Annual subscribers | 25-35% при ренью (→ ~2.5-3%/мес amortized) | Commitment effect сильнее |
| **Target (Year 2)** | **<8% monthly paid** | Критично для unit-экономики |

**Чем снижать churn:**
- Streak mechanics + Нейроны создают switching cost.
- Silent Rooms и социальные функции (Phase 4) увеличивают lock-in.
- AI-персонализация (Phase 2+) делает каждый аккаунт уникальным.
- Annual-subscription промо снижает monthly churn за счёт annual commitment.

### 5.6 North Star Metric и Secondary Metrics

**North Star Metric: Weekly Active Learners (WAL)**

Определение: пользователи, завершившие ≥3 SRS-сессии за последние 7 дней.

Обоснование: WAL коррелирует с retention (sticky users возвращаются ≥3х/неделю), paid conversion (активные пользователи исчерпывают лимиты и конвертируются), долгосрочным LTV (привычка → retention → lifetime). Метрика не оптимизируется искусственно (в отличие от DAU/MAU).

**Secondary Metrics:**

| Метрика | Что измеряет | Target (Year 1) | Target (Year 2) |
|---|---|---|---|
| D1 Retention | Onboarding quality | >45% | >55% |
| D7 Retention | Early habit formation | >25% | >30% |
| D30 Retention | Core loop stickiness | >12% | >18% |
| Sessions per WAU per week | Engagement depth | ≥3 | ≥4 |
| Avg session length | Quality of engagement | 8-15 мин | 10-18 мин |
| Free → Paid Conversion | Monetization effectiveness | 3% | 5% |
| ARPPU (monthly) | Revenue per paying user | $5.50 | $5.80 |
| Monthly Paid Churn | Subscription health | <10% | <8% |
| AI Generations / WAU | Feature adoption (Phase 2+) | N/A | ≥2/неделю |
| WAL / MAU ratio | Stickiness index | 35% | 45% |
| Creator GMV (Phase 6) | Platform health | N/A | $10K/мес |

---

## 6. Приоритизация Phase 0 (RICE)

### Методология

**RICE Score = (Reach × Impact × Confidence) / Effort**

- **Reach (0-10):** доля пользователей Phase 0, которых напрямую затрагивает функция (0 = никого, 10 = всех).
- **Impact (0-3):** влияние на North Star Metric или критические бизнес-цели (3 = массивное, 2 = среднее, 1 = низкое, 0.5 = минимальное).
- **Confidence (0-100%):** уверенность в оценках Reach и Impact.
- **Effort (person-weeks):** суммарные трудозатраты (frontend + backend).

### RICE Scoring Table

| # | Задача | Reach | Impact | Confidence | Effort (pw) | RICE Score | Приоритет |
|---|---|---|---|---|---|---|---|
| 1 | Аккаунты (email/Google/Apple) | 10 | 3 | 95% | 3 | **950** | P0 |
| 2 | Режим обучения — SRS (SM-2) | 10 | 3 | 92% | 3 | **920** | P0 |
| 3 | Онбординг flow | 10 | 3 | 85% | 2 | **1 275** | P0 |
| 4 | Базовые карточки (CRUD + import CSV/Anki) | 10 | 3 | 92% | 4 | **690** | P0 |
| 5 | i18n (ru/en) | 9 | 2 | 95% | 2 | **855** | P0 |
| 6 | Аналитика MVP | 10 | 2 | 90% | 1.5 | **1 200** | P0 |
| 7 | GDPR / приватность | 10 | 1 | 95% | 1.5 | **633** | P1 |
| 8 | Push-уведомления | 7 | 2 | 80% | 2 | **560** | P1 |
| 9 | Admin-панель MVP | 6 | 1.5 | 85% | 2 | **382** | P1 |
| 10 | A/B тест фреймворк | 5 | 2 | 70% | 2 | **350** | P1 |
| 11 | Оффлайн PWA (service worker) | 8 | 2 | 72% | 4 | **288** | P2 |
| 12 | Платёжный шлюз (Stripe + RevenueCat) | 4 | 3 | 70% | 3 | **280** | P2 |
| 13 | WebSocket-инфраструктура | 3 | 2 | 60% | 2.5 | **144** | P2 |

### Обоснование оценок

**Онбординг flow (RICE: 1 275 — наивысший):** Достигает 100% пользователей. Impact 3: activation rate — прямой предиктор всей последующей воронки. Confidence 85%: хорошо изученный феномен, но конкретный дизайн требует A/B. Effort 2 pw: UI-инвестиция без сложной backend-логики. Онбординг определяет судьбу 40-60% пользователей в первые 24 часа.

**Аналитика MVP (RICE: 1 200):** Reach 10 — влияет на всех пользователей косвенно (лучшие решения → лучший продукт). Impact 2: без аналитики невозможно оптимизировать воронку. Effort 1.5 pw: MVP-аналитика (Mixpanel/PostHog basic events) — быстро.

**Аккаунты (RICE: 950):** Gate для всего остального. Без аккаунтов нет синхронизации, нет retention, нет монетизации. Confidence высокая — стандартный flow. Effort 3 pw включает OAuth-интеграцию.

**SRS (RICE: 920):** Core product. Impact 3 — прямо влияет на WAL (North Star). Confidence 92% — SM-2 хорошо задокументирован.

**i18n (RICE: 855):** Reach 9 — все пользователи (RU-only без i18n теряет EN-аудиторию). Impact 2 — открывает EN-рынок. Effort 2 pw — i18n-архитектура уже реализована (Phase 0 progress), инкрементальный cost невысок.

**Базовые карточки (RICE: 690):** Несмотря на очевидную важность, RICE ниже из-за Effort 4 pw. Import из Anki (.apkg) — нетривиальная задача. Но это P0 — без карточек нет продукта.

**GDPR (RICE: 633):** Reach 10 (compliance затрагивает всех), Impact 1 (не влияет на WAL напрямую, но блокирует EU-лонч). Effort 1.5 pw: basic data retention policy + delete flow.

**Push-уведомления (RICE: 560):** Reach 7 (зависит от согласия пользователей — ~60-70% дадут), Impact 2 (retention-инструмент), Effort 2 pw. Важно: iOS PWA имеет ограничения — только через APNS с ограничениями.

**Admin-панель (RICE: 382):** Нужна команде, не пользователям. Impact 1.5 (модерация и управление), Effort 2 pw.

**A/B Framework (RICE: 350):** Confidence 70% — неочевидно, насколько Phase 0 dataset будет достаточен для значимых экспериментов. Effort 2 pw. Нужен до Phase 2 монетизации, но может быть упрощён (feature flags vs полноценный A/B).

**Оффлайн PWA (RICE: 288):** Reach 8 (mobile users), Impact 2, но Confidence 72% — offline-UX для SRS сложнее, чем кажется (конфликт синхронизации). Effort 4 pw. Можно отложить на Phase 0 late или Phase 1 early.

**Платёжный шлюз (RICE: 280):** Reach 4 (только paying users Phase 1+), Impact 3 (критично для монетизации), Confidence 70%, Effort 3 pw. Нужен до Phase 1, но не для MVP-лонча.

**WebSocket-инфраструктура (RICE: 144):** Reach 3 (только Phase 4+), Impact 2, Confidence 60%, Effort 2.5 pw. Нужна как архитектурная закладка, но Phase 4 ещё далеко — минимальная инвестиция в Phase 0.

### Финальный приоритизированный список Phase 0

**Группа 1 — Must Launch (до первого публичного релиза):**
1. Аккаунты (email/Google/Apple) — RICE 950
2. Базовые карточки (CRUD + import) — RICE 690
3. Режим обучения SRS (SM-2) — RICE 920
4. Онбординг flow — RICE 1 275
5. i18n (ru/en) — RICE 855

**Группа 2 — Launch + 4 недели:**
6. Аналитика MVP — RICE 1 200
7. GDPR / приватность — RICE 633
8. Push-уведомления — RICE 560

**Группа 3 — Phase 0 завершение (перед Phase 1):**
9. Admin-панель MVP — RICE 382
10. A/B тест фреймворк — RICE 350
11. Оффлайн PWA — RICE 288
12. Платёжный шлюз — RICE 280
13. WebSocket-инфраструктура — RICE 144

---

## 7. Риски и открытые вопросы

### 7.1 Product Risks

| Риск | Описание | Вероятность | Воздействие |
|---|---|---|---|
| Недостаточный retention | SRS без strong gamification теряет пользователей через 2-3 недели | Средняя | Высокое |
| SRS-алгоритм уступает FSRS | SM-2 менее точен, чем FSRS. Пользователи, знакомые с Anki FSRS, заметят разницу | Средняя | Среднее |
| Геймификация воспринимается как "детская" | Dev-аудитория скептична к gamification. Nerf Duolingo-стиля = потеря доверия | Средняя | Среднее |
| Creator quality control | В Phase 6 некачественный или плагиатный контент портит платформу | Низкая (Phase 6 далеко) | Высокое |
| Сложность onboarding для ADHD | Даже упрощённый онбординг может быть слишком длинным для ADHD-аудитории | Средняя | Среднее |
| Feature creep в Phase 0 | Желание "добавить ещё одну фишку" задерживает лонч | Высокая | Высокое |

**Митигация:**
- Провести user interviews с dev/ADHD-аудиторией перед finalizing онбординга.
- Запланировать апгрейд с SM-2 на FSRS в Phase 2 backlog.
- A/B-тест gamification elements с dev-cohort vs student-cohort отдельно.
- Phase 0 scope locked — все добавления через spec.md update.

### 7.2 Business Risks

| Риск | Описание | Вероятность | Воздействие |
|---|---|---|---|
| CAC растёт быстрее LTV | Платные каналы дорожают (iOS CPM вырос на 35% за 2022-2024) | Высокая | Высокое |
| Конкуренция со стороны Quizlet AI / Duolingo | Quizlet запустил Magic Notes (AI gen); Duolingo инвестирует в AI (2023-2025) | Высокая | Среднее |
| RU-платёжная фрагментация | Stripe недоступен для RU-карт, Apple Pay ограничен | Высокая (для RU-рынка) | Высокое |
| Apple IAP 30% комиссия | Нативное приложение теряет 30% revenue на Apple комиссии | Высокая (при выходе на iOS native) | Среднее |
| Нейроны как регуляторный риск | Некоторые регуляторы (EU, AUS) ужесточают требования к virtual currencies | Низкая | Среднее |
| Медленный creator onboarding | Creator Economy требует critical mass creators до того, как станет ценным | Средняя (Phase 6) | Высокое |

**Митигация:**
- Organic-first growth strategy: SEO + dev communities снижает зависимость от paid.
- ЮKassa / Robokassa + Stripe в параллель с Day 1.
- Web PWA в Phase 0 обходит Apple IAP (direct Stripe); нативное приложение — отдельное решение.
- Нейроны закрыты от peer-to-peer transfer и cash-out — снижает regulatory risk.

### 7.3 Technical Risks

| Риск | Описание | Вероятность | Воздействие |
|---|---|---|---|
| WebSocket масштабирование (Phase 4) | Silent Rooms с тысячами одновременных подключений требуют горизонтального scaling | Средняя | Высокое |
| AI generation costs vs revenue | При бесплатном лимите 20 gen/мес: $0.02-0.05/gen × 20 × 50K users = $20-50K/мес | Средняя | Высокое |
| iOS PWA push-ограничения | Apple добавил Web Push для iOS 16.4+, но с ограничениями; background sync нестабилен | Высокая | Среднее |
| GDPR: ADHD-идентификация | ADHD self-identification — sensitive personal data (health category) по GDPR Art. 9 | Средняя | Высокое |
| Offline sync конфликты | Пользователь редактировал карточки offline → конфликт при синхронизации | Средняя | Среднее |
| Go backend — новый стек для команды | Go — новый язык для проекта; learning curve влияет на скорость разработки | Высокая | Среднее |
| Data breach | PostgreSQL + S3 с пользовательскими данными требует security-первого подхода | Низкая | Критическое |

**Митигация:**
- WebSocket: заложить stateless-архитектуру (Redis pub/sub или NATS) с Phase 0, горизонтальный scale Phase 4.
- AI costs: жёсткие rate limits + мониторинг cost/user/day с первого дня. Рассмотреть self-hosted LLM (Ollama) для Phase 3+.
- GDPR: ADHD-идентификация — только opt-in, не хранить как "health data" explicit (хранить как UI preference).
- Offline sync: Last-write-wins + conflict UI для edge cases.

### 7.4 Открытые вопросы

| Вопрос | Опции | Рекомендация | Срок решения |
|---|---|---|---|
| **Web-only PWA vs нативные приложения** | (A) PWA навсегда; (B) PWA Phase 0-1, нативно Phase 2+; (C) Нативно с Phase 0 | (B) — PWA достаточна для Phase 0-1, нативное приложение разблокирует iOS push и App Store discovery | До Phase 2 planning |
| **Целевые рынки Phase 0** | (A) RU-only; (B) RU + EN одновременно; (C) EN-first | (B) — i18n реализован, обе аудитории с первого дня; лонч в RU-сообществах + EN на ProductHunt | До лонча |
| **Push-уведомления: стек** | (A) FCM direct; (B) OneSignal; (C) Firebase с RevenueCat | (B) OneSignal — абстрагирует FCM/APNS, проще для PWA, бесплатный tier до 10K subscribers | Phase 0 |
| **Платёжный шлюз** | (A) Stripe-only; (B) Stripe + ЮKassa; (C) Stripe + Paddle | (B) Stripe для global + ЮKassa для RU; Paddle рассмотреть позже как merchant of record | До Phase 1 |
| **SRS-алгоритм апгрейд** | (A) SM-2 навсегда; (B) SM-2 Phase 0, FSRS Phase 2; (C) Собственный алгоритм | (B) — FSRS open source, доказанно лучше SM-2, апгрейд реалистичен в Phase 2 | Phase 2 planning |
| **LLM для AI-генерации** | (A) OpenAI API; (B) Anthropic Claude API; (C) Self-hosted (Ollama + Llama 3) | (A/B) для Phase 2 — низкие upfront costs; (C) рассмотреть с Phase 3 при росте объёмов | Phase 2 |
| **Creator revenue split** | (A) 70/30 (creator/platform); (B) 80/20; (C) Sliding scale | (A) 70/30 — стандарт для EdTech platforms (Teachable, Udemy 50-97% variable). Чётко и понятно | Phase 6 planning |

### 7.5 Risk Matrix

| | **Низкое воздействие** | **Среднее воздействие** | **Высокое/Критическое воздействие** |
|---|---|---|---|
| **Высокая вероятность** | iOS PWA push-ограничения; Go learning curve | CAC рост; Quizlet AI конкуренция; RU-платёжная фрагментация | Feature creep задерживает лонч |
| **Средняя вероятность** | Нейроны regulatory risk | Geymification восприятие dev-аудиторией; Offline sync конфликты; AI costs > revenue | SRS retention недостаточный; WebSocket scaling (Phase 4); GDPR ADHD data |
| **Низкая вероятность** | — | Creator quality control | Data breach; Creator Economy не набирает критическую массу |

**Приоритет работы с рисками:**
1. Feature creep (высокая вероятность + высокое воздействие) — строгий Phase 0 scope lock.
2. SRS retention (средняя/высокая вероятность) — A/B-тест onboarding и gamification до публичного лонча.
3. AI costs (средняя вероятность + высокое воздействие) — cost-monitoring с первого пользователя AI gen.
4. GDPR ADHD data — юридическая консультация до Phase 0 лонча в EU.

---

## 8. Go-to-Market стратегия

### 8.1 Phase 0 GTM — Закрытая Бета

**Цель:** 300-800 активных beta-пользователей, 25%+ D7 retention, NPS ≥40.

**Аудитория:** Первые пользователи — квалифицированные, дающие фидбек, терпимые к багам.

**Каналы закрытой беты:**
- Личная сеть команды: dev-коллеги, ADHD-сообщества в Telegram/Discord.
- Посты в профессиональных Telegram-каналах: Go/Vue разработчики, медстуденты, ADHD-chats.
- Waitlist лендинг с referral механикой ("пригласи друга → двое попадают раньше").

**KPIs закрытой беты:**
- Signup: 500-1 500 waitlist
- Activated: ≥40% от signups
- D7 retention: ≥25%
- NPS: ≥40
- Критических багов: 0

### 8.2 Каналы привлечения

**Органические каналы (приоритет Phase 0-2):**

| Канал | Тактика | Ожидаемый CAC | Временной горизонт |
|---|---|---|---|
| SEO | Контент-маркетинг: "best anki alternative for developers", "flashcards for ADHD", "learn [tech] with spaced repetition" | $0.50-2 | 3-9 мес для rank |
| GitHub | README badge "Made with DéjàKnew", демо-репозитории с карточками по популярным технологиям | $0.50-3 | Постоянно |
| Reddit | r/Anki, r/ADHD, r/learnprogramming, r/languagelearning — value-first посты, не реклама | $1-4 | Phase 0-1 |
| Habr / dev.to | Технические статьи о SRS, FSRS, ADHD и обучении — с упоминанием DéjàKnew | $0.50-2 | Phase 0-1 |
| ProductHunt | Launch день: цель Top 5 of the Day → email list burst | $0.20-1 | Phase 1 launch |
| TikTok/Instagram | Edu-контент: "learn [tech] in 10 minutes", ADHD-tips, deck showcases | $0.10-0.50 organic | Phase 1+ |
| Twitter/X dev communities | Threads о SRS, обучении, ADHD-productivity. Без hard-sell | $0.50-2 | Постоянно |

**Партнёрства:**

| Партнёр | Механика | Приоритет |
|---|---|---|
| Coding bootcamps (Яндекс Практикум, SkillFactory, etc.) | B2B2C: студенты используют DéjàKnew как инструмент курса | Phase 2 |
| Университетские IT-клубы | Бесплатный Pro для club members в обмен на onboarding sessions | Phase 1 |
| ADHD-инфлюенсеры (YouTube, TikTok, Instagram) | Nano/micro-influencers (10K-100K) — аутентичный обзор | Phase 1-2 |
| Medical student communities (Anking, etc.) | Интеграция / импорт поддержка → organic adoption | Phase 2 |

**Платные каналы (Phase 2+, после достижения product-market fit):**

| Канал | Таргетинг | CAC estimate |
|---|---|---|
| Reddit Ads | r/learnprogramming, r/ADHD, r/Anki | $5-10 |
| Twitter/X Promoted | dev-аудитория, ADHD hashtags | $7-15 |
| Google Ads | "anki alternative", "flashcard app for adhd" | $8-20 |
| Retargeting | Посетители сайта без signup | $2-5 |

### 8.3 Launch Sequence — Milestones

| Milestone | Таймлайн | Событие | KPI Target |
|---|---|---|---|
| **M0** | T+0 | Закрытая бета: 300-800 users, ручной onboarding | Activated ≥40%, D7 ≥25%, NPS ≥40 |
| **M1** | T+6-8 нед | Open beta: публичный лонч + ProductHunt | 2 000-5 000 signups/мес, WAL ≥1 500 |
| **M2** | T+12-16 нед | Phase 1: Нейроны, сундуки, магазин косметики | Sessions/WAU ≥4, D30 ≥15% |
| **M3** | T+18-24 нед | Phase 2: AI-генерация + Pro subscription | Paid conversion ≥3%, ARPPU ≥$5.50 |
| **M4** | T+28-36 нед | Phase 3: Интерактивные курсы | Course completion ≥35%, lessons/paid user/week ≥3 |
| **M5** | T+40-48 нед | Phase 4: Silent Rooms + Дружеские пары | Silent Rooms DAU ≥10% WAU, friend-pair retention +15% |
| **M6** | T+56-64 нед | Phase 5: Full Focus Mode + AI companion | ADHD-cohort D30 retention +20% vs baseline |
| **M7** | T+72-88 нед | Phase 6: Creator Economy beta | ≥50 active creators, GMV ≥$10K/мес |

### 8.4 KPIs по фазам

| Фаза | Primary KPI | Target | Secondary KPIs |
|---|---|---|---|
| **Phase 0 (MVP)** | Activation rate | ≥40% signups → activated | D7 retention ≥25%, NPS ≥40, 0 critical bugs |
| **Phase 1 (Экономика)** | WAL (Weekly Active Learners) | 5 000+ WAL | Sessions/WAU ≥4, D30 ≥15%, нейронная экономика сбалансирована (sink/source) |
| **Phase 2 (AI + Подписка)** | Free → Paid conversion | ≥3% | ARPPU ≥$5.50, Monthly paid churn ≤10%, AI gen adoptions ≥30% MAU |
| **Phase 3 (Курсы)** | Course completion rate | ≥35% | Lessons/paid user/week ≥3, курсовые DAU/MAU ≥25% |
| **Phase 4 (Социальные)** | Silent Rooms adoption | ≥10% WAU weekly | Friend-pair D30 retention +15% vs solo, Social DAU/MAU |
| **Phase 5 (Focus Mode)** | ADHD-cohort retention | D30 +20% vs baseline | Focus Mode adoption ≥40% ADHD-cohort, avg session length ≥8 мин |
| **Phase 6 (Creator Economy)** | Creator GMV | ≥$50K/мес | Active creators ≥200, top creator earnings ≥$1K/мес, content NPS ≥50 |

**Итоговые 3-летние цели:**
- WAL: 100 000-150 000
- ARR: $2-4M
- Paid users: 8 000-15 000
- Active creators (Phase 6): 500+
- NPS: ≥55

---

*Документ создан: 2026-05-29. Следующий review: при входе в Phase 1 или при значимых изменениях в конкурентной среде.*
*Владелец: Product. Для изменений в монетизационной стратегии или TAM/SAM/SOM — требуется обновление секций 2, 4, 5.*
