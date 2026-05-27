package service

import (
	"errors"
	"testing"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/domain"
)

// fakeUserRepo is an in-memory UserRepository for unit tests.
type fakeUserRepo struct {
	byEmail map[string]*domain.User
	byID    map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail: map[string]*domain.User{},
		byID:    map[string]*domain.User{},
	}
}

func (f *fakeUserRepo) Create(u *domain.User) error {
	f.byEmail[u.Email] = u
	f.byID[u.ID] = u
	return nil
}

func (f *fakeUserRepo) FindByEmail(email string) (*domain.User, error) {
	if u, ok := f.byEmail[email]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (f *fakeUserRepo) FindByID(id string) (*domain.User, error) {
	if u, ok := f.byID[id]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func newService() *AuthService {
	return NewAuthService(newFakeUserRepo(), auth.NewTokenManager("test-secret"))
}

func TestRegisterIssuesTokenAndHashesPassword(t *testing.T) {
	s := newService()
	token, user, err := s.Register("a@b.com", "password123", "en")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token")
	}
	if user.ID == "" || user.Email != "a@b.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
	if user.PasswordHash == "password123" {
		t.Fatal("password stored in plaintext")
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	s := newService()
	if _, _, err := s.Register("a@b.com", "password123", "en"); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	_, _, err := s.Register("a@b.com", "other-password", "en")
	if !errors.Is(err, domain.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestLoginSucceedsWithCorrectPassword(t *testing.T) {
	s := newService()
	if _, _, err := s.Register("a@b.com", "password123", "en"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	token, _, err := s.Login("a@b.com", "password123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token")
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	s := newService()
	if _, _, err := s.Register("a@b.com", "password123", "en"); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	_, _, err := s.Login("a@b.com", "wrong")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginRejectsUnknownEmail(t *testing.T) {
	s := newService()
	_, _, err := s.Login("nobody@b.com", "password123")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
