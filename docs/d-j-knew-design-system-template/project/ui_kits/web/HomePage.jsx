/* HomePage — deck grid */
const { useState: useStateHome } = React;

function DeckCreateForm({ onSubmit, onCancel }) {
  const [title, setTitle] = useStateHome("");
  const [description, setDescription] = useStateHome("");
  function submit(e) {
    e.preventDefault();
    if (!title.trim()) return;
    onSubmit(title.trim(), description.trim());
    setTitle(""); setDescription("");
  }
  return (
    <form onSubmit={submit} style={{
      display: "flex", flexDirection: "column", gap: 12,
      background: "var(--bg-surface)", border: "1px solid var(--border-default)",
      borderRadius: 12, padding: 16, marginBottom: 16,
    }}>
      <h2 style={{ fontSize: 14, fontWeight: 500, color: "var(--text-secondary)", margin: 0 }}>New deck</h2>
      <Input size="sm" placeholder="Title" maxLength={200} value={title} onChange={(e) => setTitle(e.target.value)} autoFocus required style={{ borderColor: "var(--border-subtle)" }} />
      <Input size="sm" placeholder="Description (optional)" value={description} onChange={(e) => setDescription(e.target.value)} style={{ borderColor: "var(--border-subtle)" }} />
      <div style={{ display: "flex", gap: 8 }}>
        <Button type="submit">Create</Button>
        <Button variant="ghost" type="button" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}

function DeckCard({ deck, onOpen, onDelete }) {
  const [hover, setHover] = useStateHome(false);
  const [delHover, setDelHover] = useStateHome(false);
  return (
    <li
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        listStyle: "none",
        background: "var(--bg-surface)",
        border: `1px solid ${hover ? "var(--border-hover)" : "var(--border-default)"}`,
        borderRadius: 12, padding: 16,
        transition: "border-color 150ms",
        display: "flex", flexDirection: "column", gap: 8,
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: 8 }}>
        <a
          href="#"
          onClick={(e) => { e.preventDefault(); onOpen(deck.id); }}
          style={{
            fontSize: 14, fontWeight: 500, color: "var(--text-primary)",
            lineHeight: 1.375, textDecoration: "none", transition: "color 150ms",
          }}
          onMouseEnter={(e) => { e.target.style.color = "var(--accent-bright)"; }}
          onMouseLeave={(e) => { e.target.style.color = "var(--text-primary)"; }}
        >{deck.title}</a>
        <button
          title="Delete deck"
          onClick={() => onDelete(deck.id)}
          onMouseEnter={() => setDelHover(true)}
          onMouseLeave={() => setDelHover(false)}
          style={{
            background: "none", border: "none", cursor: "pointer",
            color: delHover ? "var(--error-text)" : "var(--text-disabled)",
            opacity: hover ? 1 : 0,
            transition: "opacity 150ms, color 150ms",
            fontSize: 14, padding: 0, flexShrink: 0,
          }}
        >✕</button>
      </div>
      {deck.description && (
        <p style={{
          fontSize: 12, color: "var(--text-muted)", margin: 0,
          display: "-webkit-box", WebkitLineClamp: 2, WebkitBoxOrient: "vertical",
          overflow: "hidden",
        }}>{deck.description}</p>
      )}
      {deck.tags?.length > 0 && (
        <div style={{ display: "flex", flexWrap: "wrap", gap: 4 }}>
          {deck.tags.map(t => <Badge key={t}>{t}</Badge>)}
        </div>
      )}
    </li>
  );
}

function HomePage({ user, decks, onCreateDeck, onDeleteDeck, onOpenDeck, onLogout }) {
  const [showForm, setShowForm] = useStateHome(false);

  return (
    <Page>
      <Header
        title="DéjàKnew"
        right={
          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <span style={{ fontSize: 14, color: "var(--text-muted)" }}>{user}</span>
            <IconText onClick={onLogout}>Sign out</IconText>
          </div>
        }
      />

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h2 style={{ fontSize: 18, fontWeight: 500, color: "var(--text-primary)", margin: 0 }}>My decks</h2>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>+ New deck</Button>
        )}
      </div>

      {showForm && (
        <DeckCreateForm
          onSubmit={(t, d) => { onCreateDeck(t, d); setShowForm(false); }}
          onCancel={() => setShowForm(false)}
        />
      )}

      {decks.length === 0 ? (
        <p style={{ fontSize: 14, color: "var(--text-placeholder)" }}>No decks yet. Create your first one.</p>
      ) : (
        <ul style={{
          display: "grid", gap: 12, padding: 0, margin: 0,
          gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))",
        }}>
          {decks.map(d => (
            <DeckCard key={d.id} deck={d} onOpen={onOpenDeck} onDelete={onDeleteDeck} />
          ))}
        </ul>
      )}
    </Page>
  );
}

Object.assign(window, { HomePage });
