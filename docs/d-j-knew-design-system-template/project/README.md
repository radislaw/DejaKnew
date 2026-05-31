# DéjàKnew Design System

> Spaced-repetition flashcards for developers and students. Dark-themed,
> minimal, terminal-adjacent. Built on Vue 3 + Tailwind CSS 4.

The product (codenamed **FlashCards** in the source, branded **DéjàKnew**) is a
PWA where users build decks of cards, study them on an SM-2-style spaced
schedule, and rate recall on a four-state scale (Again / Hard / Good / Easy).
It is for people who already live in a dark terminal — every surface leans into
that aesthetic.

---

## Sources

| Source | Path / Link | Notes |
| --- | --- | --- |
| Web app codebase | `web/` (mounted, read-only) | Vue 3 SFCs, Tailwind 4, Pinia, vue-router, FSD architecture |
| Color & component spec | Provided in chat | Authoritative for tokens; mirrored in `colors_and_type.css` |

No Figma file was provided. No marketing site. The web PWA is the only
product surface.

---

## CONTENT FUNDAMENTALS

The voice is **direct, lowercase-ish, and a little terse**. It reads like
someone who would rather just ship the feature than write a tooltip about it.

### Casing
- **Sentence case everywhere.** Page titles, section headings, buttons.
  - `My decks`, `New deck`, `Sign in`, `Add card`, `Back to deck`
- The four review verbs are the only Title-Case-y exception, and they're
  single words: **Again · Hard · Good · Easy**
- Form field placeholders are **all lowercase**: `email`, `password (min 8 characters)`.
- Brand wordmark uses the diacritic — **DéjàKnew** — but the source codename
  "FlashCards" still appears in the running app. Treat **DéjàKnew** as the
  consumer name.

### Pronouns
- **"My"**, not "your". `My decks`, `My cards`. The product speaks *as* the
  user, not *to* them.
- Empty / error states drop pronouns entirely: `No decks yet. Create your
  first one.`, `Loading…`, `Done for today!`

### Punctuation & symbols
- Ellipsis is the single **`…`** glyph, never three dots. Used liberally for
  loading: `Loading…`, `…` (as the button label while a form submits).
- En dashes / em dashes are not really used; the voice favors sentence
  fragments.
- Unicode arrows do real work as inline icons (see Iconography): **`←`** for
  back, **`↗`** for external link, **`✕`** for delete, **`▲ ▼`** for
  collapse / expand. Lucide is in `package.json` but not yet wired up — these
  glyphs are the live icon system.

### Tone examples (lifted verbatim)
- Empty deck list — `No decks yet. Create your first one.`
- Empty card list — `No cards yet. Add your first one.`
- Study session complete — **`Done for today!`** with the subline `All cards
  reviewed. Come back tomorrow for more.`
- Reveal-hint on a flashcard front — `tap to reveal`
- Header counter — `12 left` (no "cards", no parens)
- Add card affordance — `+ New deck` / `+ Add card` — the plus is a literal
  character, not an icon
- "More options" disclosure — `▼ More options` / `▲ Less`

### Emoji & decoration
**No emoji.** None in the codebase, none in copy. The aesthetic is closer to
a `git status` than a Notion page. If you find yourself reaching for 🎉, stop.

### Microcopy rules of thumb
- Prefer **one** noun over a noun phrase. `Front`, `Back`, `Tags`, `Source`.
- Errors are short and accept blame: `Failed to load deck`, `Sign in failed`.
- Never explain the SRS algorithm in the UI. The four rating buttons are
  self-explanatory and that's the design intent.
- Keys, hints, and counts are muted — they sit at `text-slate-500` (`--text-placeholder`).

---

## VISUAL FOUNDATIONS

### Color
A **single dark palette** built on Tailwind's slate ramp, with **sky** as the
sole brand accent. Semantic colors (emerald / red / rose / orange) appear
only where they carry meaning — study CTAs, errors, review-quality ratings.

- **Page background** is a custom near-black navy `#070d1c` — *not* slate-950
  (`#020617`), and *not* black. It's the one bespoke value in the system; the
  rest is the Tailwind palette unmodified.
- **Surfaces are translucent.** Cards use `bg-slate-800/60` over the page,
  which lets the navy bleed through and gives every panel a slight cool tint.
  Solid `slate-900` is reserved for *inputs and dropdowns* — they need to
  read as "things you can type into."
- **Sky is hierarchical**: `sky-600` for primary buttons, `sky-400` for the
  brightest call to action (auth submit, links, focused border). `sky-300`
  is the hover for the bright variant. Never use sky for body text.
- **Emerald is reserved for "Study"** — the only positive action in the app.
  This consistency is load-bearing; don't reuse emerald for generic success.
- **Red vs rose** — `red-400` for inline errors that are about a thing
  (deck failed to delete); `rose-400` for *auth* errors specifically. The
  distinction is subtle but intentional.
- **Imagery vibe**: there is no imagery. No photos, no illustrations, no
  decorative SVG. The design is text + chrome.

### Typography
- **Inter** for everything UI (system-ui as fallback). Default weight 400,
  500 for buttons / titles, 600 for the page title.
- **JetBrains Mono** is recommended for any code inside flashcard backs —
  cards render markdown via `marked`, so inline `code` and ```code blocks```
  should look like code. (Not yet styled in the source.)
- Tracking: only the page title uses `tracking-tight` (`-0.015em`).
- Scale, in actual use:
  - `text-xl` (20px) — page titles, flashcard front
  - `text-lg` (18px) — section headings ("My decks", "Cards")
  - `text-sm` (14px) — body, buttons, card titles
  - `text-xs` (12px) — meta, tags, hints
- **Line-height** is left to Tailwind defaults; nothing is finely tuned.

### Layout
- Page rhythm: `px-4 py-6`, centered, **`max-w-4xl`** (1024px). Auth shrinks
  to **`max-w-sm`** (24rem). Study session is **`max-w-xl`** (36rem) — the
  narrower column makes the flashcard feel weighty.
- Every page is **flat scroll** — no sidebars, no top nav, no breadcrumbs.
  Just a header (title + actions) and a main column.
- Lists use CSS grid with explicit breakpoints: `grid sm:grid-cols-2
  lg:grid-cols-3`. Cards have a fixed gap (`gap-3`).
- Headers use `flex items-center` with a single `← Back` ghost button on the
  left and the primary action (Study, Sign out) on the right.

### Backgrounds, transparency, blur
- **No gradients.** Anywhere.
- **No background images, textures, or patterns.**
- Transparency is the whole vibe. `bg-slate-800/60` on cards over `#070d1c`
  produces a layered, slightly atmospheric depth without any blur. There is
  **no `backdrop-filter`** in the source — depth comes purely from alpha.
- The only "elevation" comes from border color shifting on hover
  (`slate-700 → slate-500`). No box-shadows.

### Borders
- **Hairline 1px borders are the primary structural device.** Every card,
  every input, every flashcard is bounded by a `slate-700` border. Hover
  bumps it to `slate-500`. Focus on inputs jumps to `sky-400`.
- Inside a panel, nested elements step down one notch (`slate-600`). This
  creates a quiet visual hierarchy without color.
- No inner shadows, no inset borders, no double borders.

### Shadows
- **There are none.** The system uses translucency + borders instead. If a
  shadow ever feels needed, it's a signal you're solving the wrong problem.

### Corner radii
A clear hierarchy by surface size:
- **`rounded-2xl`** (1rem) — page-scale cards, flashcard container
- **`rounded-xl`** (0.75rem) — panels, form cards, review-quality buttons
- **`rounded-lg`** (0.5rem) — buttons, inputs, textareas
- **`rounded-full`** — tags, badges, pills

### Hover & press states
- **Buttons (filled)**: shift one step lighter on the same hue.
  `sky-600 → sky-500`, `emerald-600 → emerald-500`, `sky-400 → sky-300`.
- **Buttons (ghost / text)**: text color brightens. `slate-400 → slate-200`.
- **Cards**: border lightens (`slate-700 → slate-500`); no transform, no
  scale, no shadow.
- **Icon-only delete `✕`**: starts at `slate-600`, becomes `red-400` on
  hover. Hidden by default with `opacity-0`, revealed via
  `group-hover:opacity-100`. This "actions appear on hover" pattern is
  consistent across deck list and card list.
- **Review-quality buttons** (Again/Hard/Good/Easy): the text stays its
  semantic color always; hover *fills* the button with the matching
  `…-900/30` wash.
- **Press** states are not styled. No active scale.
- All transitions use Tailwind's default `transition` (150ms, ease) — except
  the flashcard flip, which is **500ms opacity + scale-95**.

### Motion
- **The flashcard flip is the one signature animation**: a 500ms crossfade
  between front and back, with the outgoing face scaling to 95% as it fades.
  No 3D rotation despite `perspective: 1200px` being declared — the design
  uses absolute-positioned faces and opacity.
- Everything else is a Tailwind `transition` (color, border-color). No
  springs, no bounces, no Framer-Motion-style flourishes.

### Visual motifs to remember
1. **Translucent slate panel over near-black-navy page** — the look.
2. **Hairline borders that lighten on hover** — the interaction language.
3. **Sky for "go", emerald for "study", semantic colors for the 4 ratings.**
4. **No emoji, no gradients, no shadows, no images.** Negative space is the
   decoration.
5. **Hover-reveal hidden actions** (`opacity-0 → opacity-100` on group hover).
6. **Pills and tags are full-radius**; everything else has a quiet radius.

---

## ICONOGRAPHY

The current production icon system is **Unicode characters**, used inline as
text. They behave like icons but inherit text color and sizing automatically,
which fits the developer-terminal aesthetic perfectly.

### Live inventory (every "icon" in the app)
| Glyph | Meaning | Context |
| --- | --- | --- |
| `←` | Back | Deck page, study page headers |
| `↗` | External link | "Source ↗" on flashcard back |
| `✕` | Delete | Deck card, card list — hover-revealed |
| `+` | Add | `+ New deck`, `+ Add card` (literal char inside button label) |
| `▼` `▲` | Expand / Collapse | "More options" disclosure in card form |
| `#` | Tag prefix | `#javascript`, rendered as part of the chip text |

### Lucide is installed but unused
`@lucide/vue` is in `package.json` but no component imports it yet. When you
need a real icon (something a Unicode glyph can't express — a chart, a sync
state, a settings cog), reach for Lucide. It's tonally correct: thin
strokes, geometric, neutral — same family of decisions as Inter and slate.

- Stroke width: **1.5** (Lucide default). Don't go heavier.
- Size: 14–16px inline (match `text-sm` / `text-base` line-height).
- Color: inherit from current `text-*` class. Never give icons their own
  color unless they're carrying semantic meaning (red delete, sky link).

If Lucide is unavailable in a given environment, **Heroicons (outline)** is
the closest substitute — same stroke philosophy. Document any substitution
inline.

### SVG / PNG / icon fonts
None used in the codebase. No `assets/` icons, no `<img>` tags, no sprite
sheet. The product is pure CSS + Unicode + Lucide-on-standby.

### Emoji
**Never.** Not in copy, not in icons, not in empty states.

---

## Index

```
.
├── README.md                     ← you are here
├── SKILL.md                      ← skill manifest (cross-compatible w/ Claude Skills)
├── colors_and_type.css           ← CSS variables: colors, type, spacing, motion
├── assets/
│   └── dejaknew-wordmark.svg     ← logo lockup (built from type)
├── preview/                      ← Design System tab cards (700×~)
│   ├── colors-*.html             ← palette, semantic, surfaces
│   ├── type-*.html               ← scale, weights, voice
│   ├── spacing-*.html            ← radii, spacing scale, transparency
│   └── components-*.html         ← buttons, inputs, badges, flashcard, quality
└── ui_kits/
    └── web/                      ← DéjàKnew PWA recreation
        ├── README.md
        ├── index.html            ← click-through prototype: login → home → deck → study
        ├── App.jsx
        ├── HomePage.jsx          ← deck grid + create form
        ├── DeckPage.jsx          ← card list + create / edit
        ├── StudyPage.jsx         ← flashcard + Again/Hard/Good/Easy
        ├── AuthPage.jsx          ← shared sign in / register form
        └── components.jsx        ← Button, Input, Card, Badge, etc.
```

### Caveats
- Inter and JetBrains Mono are loaded from Google Fonts CDN; no local files
  were shipped with the codebase. Swap in licensed copies if going to
  production.
- The codebase brand says "FlashCards", chat brief says **DéjàKnew**.
  The design system uses **DéjàKnew** consistently; the wordmark is
  bespoke (no logo file existed).
- No Figma or marketing site to mine — all recreations are derived from the
  Vue SFCs.
