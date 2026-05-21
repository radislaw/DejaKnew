# FlashCards — Architecture

Версия: 0.1 (draft)
Дата: 2026-05-22
Сопровождает: spec.md v0.3

Лёгкий ADR — решения по стеку и структуре с кратким обоснованием.

## Форма репозитория

Монорепо, два пакета:

```
flash-cards/
  web/                 # Vue 3 + Vite фронтенд
  api/                 # Go бэкенд
  docker-compose.yml   # локальная разработка: postgres + minio + api
  spec.md
  architecture.md
  CLAUDE.md            # конституция проекта
  README.md
  docs/
    flashcards-app.jsx # прототип — UX-референс, не код проекта
```

Решение: монорепо — один разработчик, общий цикл релиза, проще.

## Frontend — Feature-Sliced Design v2.1

```
web/src/
  app/         # вход, router, провайдеры, глобальные стили (Tailwind), App.vue
  pages/       # страницы-роуты
  widgets/     # крупные самодостаточные UI-блоки
  features/    # действия с бизнес-ценностью (глаголы)
  entities/    # бизнес-сущности (существительные)
  shared/      # переиспользуемое, бизнес-агностичное
```

Слои v1 (наполнение):

- `app` — router, Pinia, Tailwind base, App.vue
- `pages` — Login, Register, DeckList, DeckEdit, Study, Progress, SharedDeck
- `widgets` — DeckList, DeckEditor, StudySession, ProgressCharts
- `features` — auth (login/register формы), card-review, deck-create,
  card-edit, deck-import-export
- `entities` — deck, card, user, study-progress (сегменты ui/model/api)
- `shared` — ui (kit), api (http-клиент), lib (srs.ts, форматтеры), config

Правила FSD:
- Импорт строго вниз по слоям. Слой импортит только из слоёв ниже.
- Слайсы одного слоя не импортят друг друга напрямую.
  Кросс-импорт сущностей — через `@x`-нотацию (`entities/card/@x/deck`).
- Каждый слайс — публичный API через `index.ts`.

SRS-алгоритм (`calculateNextReview`, `getCardStatus`) — чистая функция,
переносится из прототипа в `shared/lib/srs.ts`.

Стек фронта: Vue 3 (Composition API) + TypeScript + Vite + Tailwind CSS,
Pinia (состояние), Vue Router, vite-plugin-pwa (PWA),
@lucide/vue (иконки), vue-chartjs + Chart.js (графики прогресса),
CSS-transitions для анимаций (без framer-motion).

## Backend — Go, слоистая архитектура

```
api/
  cmd/server/main.go
  internal/
    config/        # чтение env
    domain/        # GORM-модели, sentinel-ошибки
    handler/       # HTTP-обработчики (Gin), файл на ресурс
    service/       # бизнес-логика
    repository/    # доступ к БД (GORM)
    middleware/    # JWT-auth, логирование
    auth/          # JWT + bcrypt
    storage/       # S3-клиент
  go.mod
```

Поток зависимостей: `handler → service → repository`.
Зависимости внедряются через конструкторы.

Решения по Go (обоснование — см. CLAUDE.md):
- Роутинг: **Gin** — самый популярный Go-фреймворк 2026 (~48% использо-
  вания, 75k+ звёзд). Огромное комьюнити = максимум обучающих материалов,
  лучший выбор для изучения Go. net/http-совместим
- БД-доступ: **GORM** — ORM, миграции через AutoMigrate
- Драйвер: GORM Postgres driver (pgx под капотом)
- Пароли: **bcrypt** (golang.org/x/crypto)
- JWT: **golang-jwt/jwt/v5**
- Конфиг — только из env-переменных

Миграции v1 — GORM AutoMigrate (схема в коде моделей). Перед масштаби-
рованием перейти на версионированные миграции (gormigrate) — v2-задача
упрочнения.

## Модель данных (PostgreSQL)

```
users
  id            uuid pk
  email         text unique not null
  password_hash text not null
  created_at    timestamptz not null default now()

decks
  id          uuid pk
  owner_id    uuid fk -> users
  title       text not null
  description text not null default ''
  tags        text[] not null default '{}'
  is_public   boolean not null default false
  share_slug  text unique            -- nullable; есть = доступна по ссылке
  created_at  timestamptz not null default now()
  updated_at  timestamptz not null default now()

cards
  id         uuid pk
  deck_id    uuid fk -> decks (on delete cascade)
  type       text not null default 'basic'   -- v1: только 'basic'
  content    jsonb not null                  -- полиморфно по type
  position   integer not null default 0      -- порядок в колоде
  created_at timestamptz not null default now()
  updated_at timestamptz not null default now()

study_progress                              -- SRS-состояние, на юзера+карту
  id          uuid pk
  user_id     uuid fk -> users
  card_id     uuid fk -> cards (on delete cascade)
  interval    integer not null default 0
  ease        real    not null default 2.5
  reps        integer not null default 0
  lapses      integer not null default 0
  next_review timestamptz
  last_review timestamptz
  unique (user_id, card_id)

review_log                                  -- лог повторений, для графиков
  id          uuid pk
  user_id     uuid fk -> users
  card_id     uuid fk -> cards (on delete cascade)
  quality     smallint not null               -- 0..3
  reviewed_at timestamptz not null default now()

media                                       -- загруженные картинки
  id         uuid pk
  owner_id   uuid fk -> users
  s3_key     text not null
  created_at timestamptz not null default now()
```

Ключевое решение: **`cards.content` — JSONB**. v1 для type='basic':
`{ "front": "...", "back": "...", "front_image": "<media-key>|null",
"back_image": "<media-key>|null" }`. Типы v2 (Cloze, MCQ) добавляются
новым `type` + своей JSONB-схемой — без миграции таблицы.

## REST API (`/api/v1`)

```
POST   /auth/register          email + password -> tokens
POST   /auth/login             -> access + refresh tokens
POST   /auth/refresh           refresh -> новый access
GET    /me                     текущий юзер

GET    /decks                  свои колоды
POST   /decks                  создать
GET    /decks/:id              одна колода + карты
PATCH  /decks/:id              изменить
DELETE /decks/:id              удалить
POST   /decks/:id/share        включить/выключить share_slug
POST   /decks/:id/import       тело: csv|json
GET    /decks/:id/export       ?format=csv|json

POST   /decks/:id/cards        добавить карту
PATCH  /cards/:id              изменить
DELETE /cards/:id              удалить

GET    /study/:deckId/queue    SRS-очередь на сегодня
POST   /study/:cardId/review   тело: quality 0..3 -> обновить SRS + review_log

GET    /stats                  агрегаты для страницы прогресса
                               (статусы карт, по тегам, активность по дням)

POST   /media                  multipart upload -> { key }

GET    /shared/:slug           публичное чтение колоды (без auth)
```

Auth: JWT access (короткоживущий, ~15 мин) + refresh-токен.
Middleware читает `Authorization: Bearer <token>`, кладёт user в контекст.
Эндпоинт `/shared/:slug` — единственный публичный, без middleware.

## Локальная разработка

`docker-compose.yml` поднимает: PostgreSQL, MinIO (S3-совместимое
локально), api-контейнер. Фронт — `npm run dev` (Vite) отдельно.

## Хостинг

Архитектура провайдер-нейтральна (всё в Docker, S3-совместимое
хранилище). Аудитория СНГ/международная — провайдер свободный, выбор
отложен до деплоя. Фронт-SPA — статика на любой CDN; api — контейнер;
БД — managed либо self-hosted.

## Что НЕ делаем в архитектуре v1 (YAGNI)

- Нет очередей, кэша (Redis), воркеров — не нужны
- Нет GraphQL — REST достаточно
- Нет микросервисов — один Go-бинарник
- Нет SSR/Nuxt — SPA; SEO публичных колод решим в v2
