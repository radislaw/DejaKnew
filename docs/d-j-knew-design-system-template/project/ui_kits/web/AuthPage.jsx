/* AuthPage — sign in / register */
const { useState: useStateAuth } = React;

function AuthPage({ mode, onAuth, onSwitchMode }) {
  const [email, setEmail] = useStateAuth("");
  const [password, setPassword] = useStateAuth("");
  const [loading, setLoading] = useStateAuth(false);
  const [error, setError] = useStateAuth(null);

  const isLogin = mode === "login";
  const title = isLogin ? "Sign in" : "Create account";
  const submitLabel = isLogin ? "Sign in" : "Register";
  const switchPrompt = isLogin ? "No account?" : "Already have an account?";
  const switchLabel = isLogin ? "Register" : "Sign in";

  function onSubmit(e) {
    e.preventDefault();
    if (!email.trim() || password.length < 8) {
      setError("Sign in failed");
      return;
    }
    setError(null);
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      onAuth(email);
    }, 450);
  }

  return (
    <main data-screen-label={isLogin ? "Sign in" : "Register"} style={{
      minHeight: "100vh", background: "var(--bg-page)",
      display: "flex", alignItems: "center", justifyContent: "center",
      padding: "0 16px",
    }}>
      <div style={{ width: "100%", maxWidth: 384 }}>
        <h1 style={{
          fontSize: 24, fontWeight: 600, color: "var(--text-primary)",
          margin: "0 0 24px 0",
        }}>{title}</h1>

        <form onSubmit={onSubmit} style={{ display: "flex", flexDirection: "column", gap: 16 }}>
          <Input
            type="email"
            placeholder="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoFocus
          />
          <Input
            type="password"
            placeholder="password (min 8 characters)"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          {error && (
            <p style={{ fontSize: 14, color: "var(--error-text-alt)", margin: 0 }}>{error}</p>
          )}
          <Button variant="primaryLg" size="md" disabled={loading} type="submit">
            {loading ? "…" : submitLabel}
          </Button>
        </form>

        <p style={{ marginTop: 16, fontSize: 14, color: "var(--text-muted)" }}>
          {switchPrompt}{" "}
          <a
            href="#"
            onClick={(e) => { e.preventDefault(); onSwitchMode(); }}
            style={{ color: "var(--accent-bright)", textDecoration: "none" }}
            onMouseEnter={(e) => { e.target.style.textDecoration = "underline"; }}
            onMouseLeave={(e) => { e.target.style.textDecoration = "none"; }}
          >{switchLabel}</a>
        </p>
      </div>
    </main>
  );
}

Object.assign(window, { AuthPage });
