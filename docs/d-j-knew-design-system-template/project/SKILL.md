---
name: dejaknew-design
description: Use this skill to generate well-branded interfaces and assets for DéjàKnew (a dark-themed, developer-flavored spaced-repetition flashcard app), either for production or throwaway prototypes/mocks/etc. Contains essential design guidelines, colors, type, fonts, assets, and UI kit components for prototyping.
user-invocable: true
---

Read the `README.md` file within this skill, and explore the other available
files (`colors_and_type.css`, `assets/`, `preview/`, and `ui_kits/web/`).

Anchors of the brand, restated tersely:
- Page background is **#070d1c** (not slate-950, not black). Surfaces are
  **`bg-slate-800/60` translucent** over the page. Inputs are solid `slate-900`.
- Sky is the only accent (`sky-400` bright / `sky-600` primary). Emerald
  belongs to "Study". Semantic colors map to the four review qualities only.
- No emoji, no gradients, no shadows, no imagery. Borders are hairline
  slate, hover lifts to `slate-500`. Hover-revealed actions are a signature.
- Inter for UI, JetBrains Mono for code inside flashcard backs. Sentence
  case everywhere. Microcopy is terse — "No decks yet. Create your first one."
- The four-state review quality buttons (Again · Hard · Good · Easy) are the
  product's iconic component. Recreate them precisely; the colors carry meaning.
- Icons are Unicode glyphs (`← ↗ ✕ + ▼ ▲`). Lucide is the fallback if real
  icons are needed.

If creating visual artifacts (slides, mocks, throwaway prototypes, etc),
copy assets out and create static HTML files for the user to view. The
`ui_kits/web/` folder contains React JSX recreations of every screen — read
those to understand exactly how the patterns compose, and lift components
directly when you can.

If working on production code (Vue 3 + Tailwind 4), the `colors_and_type.css`
file mirrors the project's Tailwind utility classes; map them back to
`bg-slate-800/60`, `border-slate-700`, etc.

If the user invokes this skill without any other guidance, ask them what
they want to build or design, ask some questions, and act as an expert
designer who outputs HTML artifacts _or_ production code, depending on the
need.
