# FlashCards — Project Constitution

Rules for this repo. Highest priority. When unsure, ask — don't guess.

Read `spec.md` (scope + non-goals) and `architecture.md` (stack, data
model, API) before any non-trivial task.

## Stack

- Frontend: Vue 3 + TypeScript + Vite + Tailwind, FSD v2.1, PWA
- Backend: Go + Gin + GORM, layered handler/service/repository
- PostgreSQL, S3-compatible storage
- Monorepo: `web/` + `api/`

## Commands

Frontend (`web/`): `npm run dev` · `build` · `lint` · `typecheck` · `test`
Backend (`api/`): `go run ./cmd/server` · `go test ./...` · `golangci-lint run`
Local infra: `docker compose up` (postgres + minio)

## Frontend — FSD

Layers, top to bottom: `app → pages → widgets → features → entities → shared`.

- Imports go strictly downward. A layer imports only from layers below it.
- Slices in the same layer must not import each other directly.
  Cross-import entities via `@x` (`entities/card/@x/deck`).
- Each slice exposes a public API through `index.ts`. Import slices only
  through `index.ts`, never their internal files.
- Slice segments: `ui/`, `model/`, `api/`, `lib/`, `config/`.
- `<script setup lang="ts">` only. Pinia stores live in a slice's `model/`.
- Custom CSS only for design tokens; everything else is Tailwind utilities.

## Backend — Go

- Flow: `handler → service → repository`. Never skip a layer.
- Inject dependencies via constructors. No package-level globals.
- Handlers stay thin: parse request → call service → write response.
  No business logic in handlers.
- Wrap errors with `%w`. Sentinel errors live in `internal/domain`.
- GORM models live in `internal/domain`. Config from env only.
- Backend runs at L2 — review every line (Go is new to this project).

## Boundaries

- Do not use `@synergy/front-libs-vue3` or ELKA-specific skills/components.
- Do not build anything in `spec.md` non-goals. Want it? Update spec first.
- No Redis/queues/microservices/GraphQL/SSR in v1.

## Canonical feature pattern

After the reference implementation (spec.md, step 5) is built, a link to
the canonical file goes here — copy its structure for new features.
