/* App — tiny in-memory router + state */
const { useState: useStateApp } = React;

const SEED_DECKS = [
  {
    id: "d1",
    title: "JavaScript fundamentals",
    description: "Hoisting, closures, the event loop — daily review.",
    tags: ["js", "basics"],
    cards: [
      {
        id: "c1",
        content: {
          front: "What does `var x = 1` do during hoisting?",
          back: "The declaration is hoisted to the top of the function scope and initialized to `undefined`.",
          backHtml: `<p>The declaration is hoisted to the top of the function scope and initialized to <code style="font-family:var(--font-mono);color:var(--accent-bright);background:var(--bg-surface-solid);padding:1px 6px;border-radius:4px;font-size:13px">undefined</code>.</p>
<pre style="font-family:var(--font-mono);font-size:12px;background:var(--bg-surface-solid);border:1px solid var(--border-subtle);border-radius:8px;padding:10px;color:var(--text-secondary);overflow:auto;margin-top:8px">console.log(x); <span style="color:var(--text-placeholder)">// undefined</span>
var x = 1;</pre>`,
          level: "Уровень 1",
          topic: "hoisting",
          tags: ["js", "hoisting", "var"],
          sourceUrl: "https://developer.mozilla.org",
        },
      },
      {
        id: "c2",
        content: {
          front: "What's the difference between `let` and `var`?",
          back: "`let` is block-scoped and not hoisted to the top of the function. Accessing it before declaration throws a ReferenceError (the temporal dead zone).",
          backHtml: `<p><code style="font-family:var(--font-mono);color:var(--accent-bright);background:var(--bg-surface-solid);padding:1px 6px;border-radius:4px;font-size:13px">let</code> is block-scoped and not hoisted to the top of the function. Accessing it before declaration throws a ReferenceError (the temporal dead zone).</p>`,
          level: "Уровень 2",
          topic: "scoping",
          tags: ["js", "let", "scoping"],
        },
      },
      {
        id: "c3",
        content: {
          front: "What is a closure?",
          back: "A function bundled with its lexical environment — the variables in scope when it was created.",
          backHtml: `<p>A function bundled with its lexical environment — the variables in scope when it was created.</p>`,
          tags: ["js", "closures"],
        },
      },
    ],
  },
  {
    id: "d2",
    title: "Rust ownership",
    description: "Borrowing, lifetimes, Send/Sync.",
    tags: ["rust"],
    cards: [
      { id: "c4", content: { front: "What does `&mut T` mean?", back: "An exclusive, mutable reference to a value of type T.", backHtml: "<p>An exclusive, mutable reference to a value of type T.</p>", tags: ["rust", "borrowing"] } },
    ],
  },
  {
    id: "d3",
    title: "Tailwind tokens",
    description: "Slate ramp, sky accent, semantic colors.",
    tags: ["css", "design"],
    cards: [],
  },
];

function App() {
  const [user, setUser] = useStateApp(null);
  const [authMode, setAuthMode] = useStateApp("login");
  const [route, setRoute] = useStateApp({ name: "home" });
  const [decks, setDecks] = useStateApp(SEED_DECKS);

  function login(email) {
    setUser(email);
    setRoute({ name: "home" });
  }
  function logout() {
    setUser(null);
    setAuthMode("login");
  }

  function createDeck(title, description) {
    const id = "d" + (Date.now());
    setDecks(d => [{ id, title, description, tags: [], cards: [] }, ...d]);
  }
  function deleteDeck(id) {
    setDecks(d => d.filter(x => x.id !== id));
  }
  function addCardToDeck(deckId, input) {
    setDecks(d => d.map(deck => deck.id !== deckId ? deck : {
      ...deck,
      cards: [{ id: "c" + Date.now(), content: { ...input, backHtml: `<p>${input.back}</p>` } }, ...deck.cards],
    }));
  }
  function deleteCardFromDeck(deckId, cardId) {
    setDecks(d => d.map(deck => deck.id !== deckId ? deck : {
      ...deck, cards: deck.cards.filter(c => c.id !== cardId),
    }));
  }

  if (!user) {
    return (
      <AuthPage
        mode={authMode}
        onAuth={login}
        onSwitchMode={() => setAuthMode(authMode === "login" ? "register" : "login")}
      />
    );
  }

  if (route.name === "deck") {
    const deck = decks.find(d => d.id === route.id);
    if (!deck) { setRoute({ name: "home" }); return null; }
    return (
      <DeckPage
        deck={deck}
        cards={deck.cards}
        onBack={() => setRoute({ name: "home" })}
        onStudy={() => setRoute({ name: "study", id: deck.id })}
        onAddCard={(input) => addCardToDeck(deck.id, input)}
        onDeleteCard={(cardId) => deleteCardFromDeck(deck.id, cardId)}
      />
    );
  }

  if (route.name === "study") {
    const deck = decks.find(d => d.id === route.id);
    if (!deck) { setRoute({ name: "home" }); return null; }
    return (
      <StudyPage
        deck={deck}
        queue={deck.cards}
        onBack={() => setRoute({ name: "deck", id: deck.id })}
      />
    );
  }

  return (
    <HomePage
      user={user}
      decks={decks}
      onCreateDeck={createDeck}
      onDeleteDeck={deleteDeck}
      onOpenDeck={(id) => setRoute({ name: "deck", id })}
      onLogout={logout}
    />
  );
}

Object.assign(window, { App });
