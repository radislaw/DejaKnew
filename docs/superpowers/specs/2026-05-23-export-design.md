# Export Feature Design

Date: 2026-05-23

## Overview

Two scopes of export (single deck, all decks), three formats (CSV, JSON, Anki-CSV).
Plus a UX fix: entire card row clickable, not just title.

## Backend

### Endpoints

Both require JWT auth.

```
GET /api/v1/decks/:id/export?format=csv|json|anki-csv   — single deck
GET /api/v1/export?format=csv|json|anki-csv             — all user's decks
```

Response headers:
- `Content-Disposition: attachment; filename=flashcards-export.{ext}`
- `Content-Type: text/csv` | `application/json` | `text/plain`

### Formats

**csv** (`deck_name,front,back,topic,level,source_url,tags`):
```
deck_name,front,back,topic,level,source_url,tags
Go Basics,What is a goroutine?,A lightweight thread,Go,beginner,,concurrency;go
```
Tags joined with `;`.

**json** (array of decks with nested cards):
```json
[
  {
    "name": "Go Basics",
    "description": "...",
    "cards": [
      { "front": "...", "back": "...", "topic": "...", "level": "...", "sourceUrl": "...", "tags": [] }
    ]
  }
]
```

**anki-csv** (tab-separated, no header, `.txt` extension — Anki requirement):
```
What is a goroutine?	A lightweight thread	concurrency go
```
Tags joined with space (Anki tag format).

### Handler methods

In `api/internal/handler/card_handler.go`:
- `ExportDeck(c *gin.Context)` — fetches cards for one deck, serializes, writes response
- `ExportAll(c *gin.Context)` — fetches all decks+cards for user, serializes, writes response

Shared serialization helpers: `serializeCSV`, `serializeJSON`, `serializeAnkiCSV`.

### Router

In `api/internal/handler/router.go`:
```go
auth.GET("/export", cardHandler.ExportAll)
auth.GET("/decks/:id/export", cardHandler.ExportDeck)
```

### Service / Repository

No new service methods needed if `ListByDeck` and a new `ListAllForUser` are added to card repository. `CardService` needs `ListAllForUser(userID string) ([]domain.Card, error)` wrapping the repo call.

## Frontend

### New files

**`web/src/entities/card/api/card.api.ts`** — add:
```ts
export async function exportCards(format: 'csv' | 'json' | 'anki-csv', deckId?: string): Promise<Blob>
```
Uses `fetch` with auth header, returns `Blob`.

**`web/src/features/card-export/ui/CardExportButton.vue`**:
- Props: `deckId?: string`
- UI: "Export" button with dropdown — CSV / JSON / Anki CSV
- On select: calls `exportCards(format, deckId?)`, creates object URL, triggers `<a>` download, revokes URL

**`web/src/features/card-export/index.ts`** — barrel export.

### Integration

- `DeckPage.vue` — add `<CardExportButton :deckId="deckId" />` next to import button
- `HomePage.vue` — add `<CardExportButton />` (no deckId = all decks)

## UX Fix: Clickable Card Row

Currently: click on card title to navigate.
Fix: wrap entire card row in `<RouterLink>` or add `@click` on the row container.

File: `web/src/features/card-create/ui/` or wherever card list items render (check `DeckPage`).

## Error Handling

- Unknown format → 400 `{"error": "unsupported format"}`
- No cards in deck → return empty file (not 404)
- Frontend fetch error → show inline error message
