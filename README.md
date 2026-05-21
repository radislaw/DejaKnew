# FlashCards

Web app for creating, studying, and sharing spaced-repetition flashcards.

See `spec.md` (scope), `architecture.md` (design), `CLAUDE.md` (conventions).

## Stack

- Frontend: Vue 3 + TypeScript + Vite + Tailwind CSS (FSD v2.1), PWA
- Backend: Go + Gin + GORM
- PostgreSQL, S3-compatible storage

## Local development

```sh
docker compose up -d          # postgres + minio

cd api && go run ./cmd/server # backend on :8080
cd web && npm install && npm run dev  # frontend on :5173
```

## Commands

Frontend (`web/`): `npm run dev` · `build` · `lint` · `typecheck` · `test`
Backend (`api/`): `go run ./cmd/server` · `go test ./...` · `golangci-lint run`
