/* DeckPage — list cards, add new */
const { useState: useStateDeck } = React;

function CardCreateForm({ onSubmit, onCancel }) {
  const [front, setFront] = useStateDeck("");
  const [back, setBack] = useStateDeck("");
  const [topic, setTopic] = useStateDeck("");
  const [level, setLevel] = useStateDeck("");
  const [sourceUrl, setSourceUrl] = useStateDeck("");
  const [tagsRaw, setTagsRaw] = useStateDeck("");
  const [showOptional, setShowOptional] = useStateDeck(false);

  function submit(e) {
    e.preventDefault();
    if (!front.trim() || !back.trim()) return;
    const tags = tagsRaw.split(",").map(s => s.trim()).filter(Boolean);
    onSubmit({
      front: front.trim(), back: back.trim(),
      topic: topic.trim() || undefined,
      level: level.trim() || undefined,
      sourceUrl: sourceUrl.trim() || undefined,
      tags: tags.length ? tags : undefined,
    });
    setFront(""); setBack(""); setTopic(""); setLevel(""); setSourceUrl(""); setTagsRaw("");
    setShowOptional(false);
  }

  return (
    <form onSubmit={submit} style={{
      display: "flex", flexDirection: "column", gap: 12,
      background: "var(--bg-surface)", border: "1px solid var(--border-default)",
      borderRadius: 12, padding: 16, marginBottom: 16,
    }}>
      <h2 style={{ fontSize: 14, fontWeight: 500, color: "var(--text-secondary)", margin: 0 }}>New card</h2>
      <Textarea rows={2} placeholder="Front" value={front} onChange={(e) => setFront(e.target.value)} required />
      <Textarea rows={3} placeholder="Back (markdown supported)" value={back} onChange={(e) => setBack(e.target.value)} required />
      <button
        type="button"
        onClick={() => setShowOptional(!showOptional)}
        style={{
          alignSelf: "flex-start", background: "none", border: "none",
          fontSize: 12, color: "var(--text-placeholder)", cursor: "pointer",
          padding: 0, fontFamily: "inherit",
          transition: "color 150ms",
        }}
        onMouseEnter={(e) => { e.target.style.color = "var(--text-secondary)"; }}
        onMouseLeave={(e) => { e.target.style.color = "var(--text-placeholder)"; }}
      >{showOptional ? "▲ Less" : "▼ More options"}</button>

      {showOptional && (
        <React.Fragment>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
            <Input size="sm" placeholder="Topic (e.g. Hoisting)" value={topic} onChange={(e) => setTopic(e.target.value)} style={{ borderColor: "var(--border-subtle)" }} />
            <Input size="sm" placeholder="Level (e.g. Уровень 1)" value={level} onChange={(e) => setLevel(e.target.value)} style={{ borderColor: "var(--border-subtle)" }} />
          </div>
          <Input size="sm" placeholder="Source URL (optional)" type="url" value={sourceUrl} onChange={(e) => setSourceUrl(e.target.value)} style={{ borderColor: "var(--border-subtle)" }} />
          <Input size="sm" placeholder="Tags, comma-separated (e.g. js, hoisting, let)" value={tagsRaw} onChange={(e) => setTagsRaw(e.target.value)} style={{ borderColor: "var(--border-subtle)" }} />
        </React.Fragment>
      )}

      <div style={{ display: "flex", gap: 8 }}>
        <Button type="submit">Add card</Button>
        <Button variant="ghost" type="button" onClick={onCancel}>Cancel</Button>
      </div>
    </form>
  );
}

function CardRow({ card, onDelete, onEdit }) {
  const [hover, setHover] = useStateDeck(false);
  const [editHover, setEditHover] = useStateDeck(false);
  const [delHover, setDelHover] = useStateDeck(false);
  return (
    <li
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        listStyle: "none",
        background: "var(--bg-surface)",
        border: "1px solid var(--border-default)",
        borderRadius: 12, padding: 16,
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", gap: 8 }}>
        <div style={{ flex: 1, display: "grid", gridTemplateColumns: "1fr 1fr", columnGap: 24, rowGap: 4 }}>
          <div>
            <p style={{ fontSize: 12, color: "var(--text-placeholder)", margin: "0 0 2px 0" }}>Front</p>
            <p style={{ fontSize: 14, color: "var(--text-primary)", margin: 0, whiteSpace: "pre-wrap" }}>{card.content.front}</p>
          </div>
          <div>
            <p style={{ fontSize: 12, color: "var(--text-placeholder)", margin: "0 0 2px 0" }}>Back</p>
            <p style={{ fontSize: 14, color: "var(--text-secondary)", margin: 0, whiteSpace: "pre-wrap" }}>{card.content.back}</p>
          </div>
        </div>
        <div style={{
          display: "flex", gap: 8, flexShrink: 0,
          opacity: hover ? 1 : 0, transition: "opacity 150ms",
        }}>
          <button
            onClick={() => onEdit(card.id)}
            onMouseEnter={() => setEditHover(true)}
            onMouseLeave={() => setEditHover(false)}
            style={{
              background: "none", border: "none", cursor: "pointer",
              fontSize: 12, color: editHover ? "var(--accent-bright)" : "var(--text-placeholder)",
              padding: 0, fontFamily: "inherit",
            }}
          >Edit</button>
          <button
            onClick={() => onDelete(card.id)}
            onMouseEnter={() => setDelHover(true)}
            onMouseLeave={() => setDelHover(false)}
            style={{
              background: "none", border: "none", cursor: "pointer",
              fontSize: 12, color: delHover ? "var(--error-text)" : "var(--text-placeholder)",
              padding: 0, fontFamily: "inherit",
            }}
          >✕</button>
        </div>
      </div>
    </li>
  );
}

function DeckPage({ deck, cards, onBack, onStudy, onAddCard, onDeleteCard }) {
  const [showForm, setShowForm] = useStateDeck(false);

  return (
    <Page>
      <Header
        left={<IconText onClick={onBack}>← Back</IconText>}
        title={deck.title}
        subtitle={deck.description}
        right={<Button variant="success" onClick={onStudy}>Study</Button>}
      />

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
        <h2 style={{ fontSize: 18, fontWeight: 500, color: "var(--text-primary)", margin: 0 }}>
          Cards <span style={{ color: "var(--text-placeholder)", fontSize: 14, fontWeight: 400 }}>({cards.length})</span>
        </h2>
        {!showForm && (
          <Button onClick={() => setShowForm(true)}>+ Add card</Button>
        )}
      </div>

      {showForm && (
        <CardCreateForm
          onSubmit={(input) => { onAddCard(input); setShowForm(false); }}
          onCancel={() => setShowForm(false)}
        />
      )}

      {cards.length === 0 ? (
        <p style={{ fontSize: 14, color: "var(--text-placeholder)" }}>No cards yet. Add your first one.</p>
      ) : (
        <ul style={{ display: "flex", flexDirection: "column", gap: 12, padding: 0, margin: 0 }}>
          {cards.map(c => (
            <CardRow key={c.id} card={c} onDelete={onDeleteCard} onEdit={() => {}} />
          ))}
        </ul>
      )}
    </Page>
  );
}

Object.assign(window, { DeckPage });
