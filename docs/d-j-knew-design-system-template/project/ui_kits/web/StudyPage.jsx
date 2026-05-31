/* StudyPage — the signature flashcard flip */
const { useState: useStateStudy, useEffect: useEffectStudy } = React;

function QualityButton({ tone, children, onClick }) {
  const tones = {
    again: { text: "var(--error-text)",   hoverBg: "var(--error-bg)" },
    hard:  { text: "var(--warn-text)",    hoverBg: "var(--warn-bg)" },
    good:  { text: "var(--good-text)",    hoverBg: "var(--good-bg)" },
    easy:  { text: "var(--easy-text)",    hoverBg: "var(--easy-bg)" },
  };
  const t = tones[tone];
  const [hover, setHover] = useStateStudy(false);
  return (
    <button
      onClick={onClick}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        border: "1px solid var(--border-subtle)",
        background: hover ? t.hoverBg : "var(--bg-elevated)",
        color: t.text,
        padding: "12px 0",
        borderRadius: 12,
        fontWeight: 500,
        fontSize: 14,
        cursor: "pointer",
        fontFamily: "inherit",
        transition: "background-color 150ms",
      }}
    >{children}</button>
  );
}

function Flashcard({ card, flipped, onFlip }) {
  return (
    <div
      onClick={() => !flipped && onFlip()}
      style={{
        position: "relative",
        width: "100%",
        cursor: flipped ? "default" : "pointer",
        userSelect: "none",
        background: "var(--bg-surface)",
        border: "1px solid var(--border-default)",
        borderRadius: 16,
        minHeight: "16rem",
        perspective: "1200px",
      }}
    >
      {/* Front */}
      <div style={{
        position: "absolute", inset: 0,
        display: "flex", flexDirection: "column",
        alignItems: "center", justifyContent: "center",
        borderRadius: 16, padding: 32,
        transition: "opacity 500ms, transform 500ms",
        opacity: flipped ? 0 : 1,
        transform: flipped ? "scale(0.95)" : "scale(1)",
        pointerEvents: flipped ? "none" : "auto",
      }}>
        {(card.content.level || card.content.topic) && (
          <div style={{ display: "flex", gap: 8, alignSelf: "flex-start", marginBottom: 12 }}>
            {card.content.level && <Badge variant="accent">{card.content.level}</Badge>}
            {card.content.topic && <Badge variant="soft">{card.content.topic}</Badge>}
          </div>
        )}
        <p style={{
          textAlign: "center", fontSize: 20, fontWeight: 500,
          color: "var(--text-primary)", margin: 0,
        }}>{card.content.front}</p>
        <span style={{ marginTop: 24, fontSize: 12, color: "var(--text-placeholder)" }}>tap to reveal</span>
      </div>

      {/* Back */}
      <div style={{
        position: "absolute", inset: 0,
        display: "flex", flexDirection: "column",
        borderRadius: 16, padding: 32,
        transition: "opacity 500ms, transform 500ms",
        opacity: flipped ? 1 : 0,
        transform: flipped ? "scale(1)" : "scale(0.95)",
        pointerEvents: flipped ? "auto" : "none",
        overflow: "auto",
      }}>
        <div style={{ flex: 1, fontSize: 14, color: "var(--text-secondary)", lineHeight: 1.55 }}
             dangerouslySetInnerHTML={{ __html: card.content.backHtml || `<p>${card.content.back}</p>` }} />
        {card.content.tags?.length > 0 && (
          <div style={{ marginTop: 16, display: "flex", flexWrap: "wrap", gap: 6 }}>
            {card.content.tags.map(t => <Tag key={t}>{t}</Tag>)}
          </div>
        )}
        {card.content.sourceUrl && (
          <a
            href={card.content.sourceUrl}
            target="_blank"
            rel="noopener noreferrer"
            onClick={(e) => e.stopPropagation()}
            style={{
              marginTop: 12, alignSelf: "flex-start",
              fontSize: 12, color: "var(--accent-bright)",
              textDecoration: "none",
            }}
          >Source ↗</a>
        )}
      </div>
    </div>
  );
}

function StudyPage({ deck, queue: initialQueue, onBack }) {
  const [queue, setQueue] = useStateStudy(initialQueue);
  const [flipped, setFlipped] = useStateStudy(false);
  const [done, setDone] = useStateStudy(initialQueue.length === 0);

  const current = queue[0];
  const remaining = queue.length;

  function rate(/* quality */) {
    // Pop the front card. In real product, schedule based on quality.
    const rest = queue.slice(1);
    if (rest.length === 0) {
      setDone(true);
    } else {
      setFlipped(false);
      setQueue(rest);
    }
  }

  return (
    <Page narrow="xl">
      <Header
        left={<IconText onClick={onBack}>← Back</IconText>}
        title="Study session"
        right={<span style={{ fontSize: 14, color: "var(--text-placeholder)" }}>{done ? 0 : remaining} left</span>}
      />

      {done ? (
        <div style={{
          display: "flex", flexDirection: "column", alignItems: "center",
          gap: 16, padding: "80px 0", textAlign: "center",
        }}>
          <p style={{ fontSize: 24, fontWeight: 600, color: "var(--success-text)", margin: 0 }}>Done for today!</p>
          <p style={{ fontSize: 14, color: "var(--text-muted)", margin: 0 }}>All cards reviewed. Come back tomorrow for more.</p>
          <Button onClick={onBack} size="md" style={{ marginTop: 16, padding: "8px 24px" }}>Back to deck</Button>
        </div>
      ) : (
        <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 24 }}>
          <Flashcard card={current} flipped={flipped} onFlip={() => setFlipped(true)} />
          {flipped && (
            <div style={{ width: "100%", display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 8 }}>
              <QualityButton tone="again" onClick={() => rate(0)}>Again</QualityButton>
              <QualityButton tone="hard"  onClick={() => rate(1)}>Hard</QualityButton>
              <QualityButton tone="good"  onClick={() => rate(2)}>Good</QualityButton>
              <QualityButton tone="easy"  onClick={() => rate(3)}>Easy</QualityButton>
            </div>
          )}
        </div>
      )}
    </Page>
  );
}

Object.assign(window, { StudyPage });
