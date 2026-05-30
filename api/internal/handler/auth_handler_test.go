package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/domain"
	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
)

// fakeUserRepo is an in-memory service.UserRepository for handler tests —
// the same pattern as in service tests, scoped to this package.
type fakeUserRepo struct {
	byEmail    map[string]*domain.User
	byID       map[string]*domain.User
	byProvider map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail:    map[string]*domain.User{},
		byID:       map[string]*domain.User{},
		byProvider: map[string]*domain.User{},
	}
}

func (f *fakeUserRepo) Create(u *domain.User) error {
	f.byEmail[u.Email] = u
	f.byID[u.ID] = u
	if u.Provider != "" && u.ProviderID != "" {
		f.byProvider[u.Provider+":"+u.ProviderID] = u
	}
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

func (f *fakeUserRepo) FindByProvider(provider, providerID string) (*domain.User, error) {
	if u, ok := f.byProvider[provider+":"+providerID]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func testRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	tokens := auth.NewTokenManager("test-secret")
	authSvc := service.NewAuthService(newFakeUserRepo(), tokens)
	deckRepo := newFakeDeckRepo()
	deckSvc := service.NewDeckService(deckRepo)
	cardRepo := newFakeCardRepo()
	cardSvc := service.NewCardService(cardRepo, deckRepo)
	studySvc := service.NewStudyService(newFakeStudyProgressRepo(), cardRepo, deckRepo)
	oauthH := NewOAuthHandler(authSvc, "http://localhost:8080", "http://localhost:5175", "", "", "", "")
	return Router(tokens, NewAuthHandler(authSvc), oauthH, NewDeckHandler(deckSvc, cardSvc), NewCardHandler(cardSvc), NewStudyHandler(studySvc))
}

func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
}

func registerUser(t *testing.T, r http.Handler, email string) string {
	t.Helper()
	w := doJSON(t, r, "POST", "/api/v1/auth/register",
		`{"email":"`+email+`","password":"password123"}`, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register %s failed: %d %s", email, w.Code, w.Body.String())
	}
	return tokenFrom(t, w)
}

func doJSON(t *testing.T, r http.Handler, method, path, body, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func tokenFrom(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp.Token
}

func TestRegisterReturns201WithToken(t *testing.T) {
	r := testRouter()
	w := doJSON(t, r, "POST", "/api/v1/auth/register",
		`{"email":"a@b.com","password":"password123"}`, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if tokenFrom(t, w) == "" {
		t.Fatal("expected a token in the response")
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	r := testRouter()
	w := doJSON(t, r, "POST", "/api/v1/auth/register",
		`{"email":"a@b.com","password":"short"}`, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for short password, got %d", w.Code)
	}
}

func TestRegisterDuplicateReturns409(t *testing.T) {
	r := testRouter()
	body := `{"email":"a@b.com","password":"password123"}`
	doJSON(t, r, "POST", "/api/v1/auth/register", body, "")
	w := doJSON(t, r, "POST", "/api/v1/auth/register", body, "")
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate email, got %d", w.Code)
	}
}

func TestLoginThenMeReturnsUser(t *testing.T) {
	r := testRouter()
	reg := doJSON(t, r, "POST", "/api/v1/auth/register",
		`{"email":"a@b.com","password":"password123"}`, "")
	token := tokenFrom(t, reg)

	w := doJSON(t, r, "GET", "/api/v1/me", "", token)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from /me, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMeWithoutTokenReturns401(t *testing.T) {
	r := testRouter()
	w := doJSON(t, r, "GET", "/api/v1/me", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}
