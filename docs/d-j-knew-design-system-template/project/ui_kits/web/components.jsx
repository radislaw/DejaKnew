/* DéjàKnew UI Kit — primitives
   Every component reads tokens from ../../colors_and_type.css */
const { useState } = React;

// ── Button ───────────────────────────────────────────────────
function Button({ variant = "primary", size = "sm", disabled, children, ...rest }) {
  const base = {
    fontFamily: "inherit",
    fontWeight: 500,
    borderRadius: 8,
    border: "none",
    cursor: disabled ? "not-allowed" : "pointer",
    transition: "background-color 150ms, color 150ms, border-color 150ms",
    opacity: disabled ? 0.5 : 1,
    lineHeight: 1.25,
  };
  const sizes = {
    sm: { padding: "6px 16px", fontSize: 14 },
    md: { padding: "10px 16px", fontSize: 14 },
    lg: { padding: "10px 24px", fontSize: 14 },
  };
  const variants = {
    primary:    { background: "var(--accent-default)",  color: "#fff" },
    primaryLg:  { background: "var(--accent-bright)",   color: "var(--text-inverse)" },
    success:    { background: "var(--success)",         color: "#fff" },
    ghost:      { background: "transparent",            color: "var(--text-muted)", padding: "6px 12px" },
  };
  const [hover, setHover] = useState(false);
  const hoverStyle = hover && !disabled ? {
    primary:    { background: "var(--accent-hover)" },
    primaryLg:  { background: "var(--accent-bright-hover)" },
    success:    { background: "var(--success-hover)" },
    ghost:      { color: "var(--text-secondary)" },
  }[variant] : null;
  return (
    <button
      {...rest}
      disabled={disabled}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{ ...base, ...sizes[size], ...variants[variant], ...hoverStyle, ...(rest.style || {}) }}
    >
      {children}
    </button>
  );
}

// ── Input / Textarea ────────────────────────────────────────
function Input({ size = "md", style, ...rest }) {
  const [focus, setFocus] = useState(false);
  const sizes = {
    sm: { padding: "8px 12px", fontSize: 14 },
    md: { padding: "10px 16px", fontSize: 14 },
  };
  return (
    <input
      {...rest}
      onFocus={(e) => { setFocus(true); rest.onFocus?.(e); }}
      onBlur={(e) => { setFocus(false); rest.onBlur?.(e); }}
      style={{
        background: "var(--bg-surface-solid)",
        border: `1px solid ${focus ? "var(--border-focus)" : "var(--border-default)"}`,
        borderRadius: 8,
        color: "var(--text-primary)",
        fontFamily: "inherit",
        outline: "none",
        width: "100%",
        ...sizes[size],
        ...style,
      }}
    />
  );
}

function Textarea({ rows = 3, style, ...rest }) {
  const [focus, setFocus] = useState(false);
  return (
    <textarea
      rows={rows}
      {...rest}
      onFocus={(e) => { setFocus(true); rest.onFocus?.(e); }}
      onBlur={(e) => { setFocus(false); rest.onBlur?.(e); }}
      style={{
        background: "var(--bg-surface-solid)",
        border: `1px solid ${focus ? "var(--border-focus-deep)" : "var(--border-subtle)"}`,
        borderRadius: 8,
        color: "var(--text-primary)",
        fontFamily: "inherit",
        fontSize: 14,
        padding: "8px 12px",
        resize: "none",
        outline: "none",
        width: "100%",
        ...style,
      }}
    />
  );
}

// ── Card ────────────────────────────────────────────────────
function Card({ children, hoverable = true, style, ...rest }) {
  const [hover, setHover] = useState(false);
  return (
    <div
      {...rest}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        background: "var(--bg-surface)",
        border: `1px solid ${hover && hoverable ? "var(--border-hover)" : "var(--border-default)"}`,
        borderRadius: 12,
        padding: 16,
        transition: "border-color 150ms",
        ...style,
      }}
    >
      {typeof children === "function" ? children({ hover }) : children}
    </div>
  );
}

// ── Badge / Tag ─────────────────────────────────────────────
function Badge({ variant = "default", children, style }) {
  const styles = {
    default: { background: "#334155", color: "var(--text-secondary)", padding: "2px 10px" },
    accent:  { background: "var(--accent-badge-bg)", color: "var(--accent-badge-text)", padding: "2px 10px", fontWeight: 500 },
    soft:    { background: "rgba(51,65,85,0.8)", color: "var(--text-muted)", padding: "2px 10px" },
  };
  return (
    <span style={{
      display: "inline-flex", alignItems: "center", borderRadius: 9999,
      fontSize: 12, ...styles[variant], ...style,
    }}>{children}</span>
  );
}

function Tag({ children }) {
  return (
    <span style={{
      background: "rgba(51,65,85,0.6)", color: "var(--text-muted)",
      borderRadius: 9999, padding: "2px 8px", fontSize: 12,
    }}>#{children}</span>
  );
}

// ── Page chrome ────────────────────────────────────────────
function Page({ children, narrow }) {
  const max = narrow === "sm" ? 384 : narrow === "xl" ? 576 : 1024;
  return (
    <div data-screen-label={narrow ? "narrow" : "page"} style={{
      minHeight: "100vh", background: "var(--bg-page)",
      color: "var(--text-primary)", padding: "24px 16px",
    }}>
      <div style={{ maxWidth: max, margin: "0 auto" }}>{children}</div>
    </div>
  );
}

function Header({ left, title, subtitle, right }) {
  return (
    <header style={{
      display: "flex", alignItems: "center", gap: 16,
      marginBottom: 32,
    }}>
      {left}
      <div style={{ flex: 1, minWidth: 0 }}>
        <h1 style={{
          fontSize: 20, fontWeight: 600, letterSpacing: "-0.015em",
          color: "var(--text-primary)", margin: 0,
        }}>{title}</h1>
        {subtitle && (
          <p style={{ fontSize: 14, color: "var(--text-muted)", margin: "2px 0 0 0" }}>{subtitle}</p>
        )}
      </div>
      {right}
    </header>
  );
}

// ── Ghost icon-text (back ←, etc) ───────────────────────────
function IconText({ children, onClick, style }) {
  const [hover, setHover] = useState(false);
  return (
    <button
      onClick={onClick}
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        background: "none", border: "none", cursor: "pointer",
        fontFamily: "inherit", fontSize: 14, padding: 0,
        color: hover ? "var(--text-secondary)" : "var(--text-muted)",
        transition: "color 150ms",
        ...style,
      }}
    >{children}</button>
  );
}

// expose
Object.assign(window, { Button, Input, Textarea, Card, Badge, Tag, Page, Header, IconText });
