package auth

import "testing"

func TestIssueAndVerify(t *testing.T) {
	m := NewTokenManager("test-secret")
	token, err := m.Issue("user-123")
	if err != nil {
		t.Fatalf("issue failed: %v", err)
	}
	sub, err := m.Verify(token)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if sub != "user-123" {
		t.Fatalf("expected subject user-123, got %q", sub)
	}
}

func TestVerifyRejectsGarbage(t *testing.T) {
	m := NewTokenManager("test-secret")
	if _, err := m.Verify("garbage.token.value"); err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestVerifyRejectsWrongSecret(t *testing.T) {
	token, _ := NewTokenManager("secret-a").Issue("user-1")
	if _, err := NewTokenManager("secret-b").Verify(token); err == nil {
		t.Fatal("expected error for token signed with a different secret")
	}
}
