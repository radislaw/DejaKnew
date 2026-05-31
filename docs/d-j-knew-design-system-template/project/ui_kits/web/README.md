# DéjàKnew Web — UI Kit

A pixel-faithful recreation of the DéjàKnew Vue 3 / Tailwind PWA, rebuilt as
React + Babel JSX so it's drop-in for prototyping.

## Files
- `index.html` — interactive click-through. Boots at Sign-in; lets you fake-log-in, see your decks, open a deck, add cards, and run a study session with the four-state quality buttons.
- `components.jsx` — primitives: `Button`, `Input`, `Textarea`, `Card`, `Badge`, `Tag`, `Header`, `Page`.
- `AuthPage.jsx` — shared sign-in / register layout.
- `HomePage.jsx` — deck grid + new-deck form.
- `DeckPage.jsx` — card list + new-card form (with the "▼ More options" disclosure).
- `StudyPage.jsx` — the flashcard flip + Again/Hard/Good/Easy.
- `App.jsx` — micro router glue.

## What's faithful
- Token values come straight from `colors_and_type.css`.
- `bg-slate-800/60` translucency over `#070d1c`, hairline borders, hover-reveal `✕`, 500ms flip — all preserved.
- Empty / loading / error copy lifted verbatim from the Vue source.

## What's cut
- No real API. Decks and cards live in component state and seed from a small fixture.
- No SRS scheduling math. Quality buttons advance the queue, that's it.
- No PWA / service-worker.
