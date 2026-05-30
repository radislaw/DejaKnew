package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"flashcards/api/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthStateCookie = "oauth_state"

type OAuthHandler struct {
	auth        *service.AuthService
	google      *oauth2.Config
	yandex      *oauth2.Config
	frontendURL string
}

func NewOAuthHandler(auth *service.AuthService, appURL, frontendURL, googleClientID, googleClientSecret, yandexClientID, yandexClientSecret string) *OAuthHandler {
	return &OAuthHandler{
		auth:        auth,
		frontendURL: frontendURL,
		google: &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  appURL + "/api/v1/auth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		yandex: &oauth2.Config{
			ClientID:     yandexClientID,
			ClientSecret: yandexClientSecret,
			RedirectURL:  appURL + "/api/v1/auth/yandex/callback",
			Scopes:       []string{"login:email", "login:info"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://oauth.yandex.com/authorize",
				TokenURL: "https://oauth.yandex.com/token",
			},
		},
	}
}

func (h *OAuthHandler) GoogleRedirect(c *gin.Context) {
	if h.google.ClientID == "" {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Google+OAuth+not+configured")
		return
	}
	state := newState()
	c.SetCookie(oauthStateCookie, state, 300, "/", "", false, true)
	c.Redirect(http.StatusFound, h.google.AuthCodeURL(state))
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	if err := h.verifyState(c); err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Invalid+state")
		return
	}
	token, err := h.google.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Token+exchange+failed")
		return
	}
	info, err := fetchGoogleUser(token.AccessToken)
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Failed+to+fetch+user+info")
		return
	}
	locale := c.GetString("locale")
	jwt, _, err := h.auth.LoginWithOAuth("google", info.ID, info.Email, locale)
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Login+failed")
		return
	}
	c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?token="+jwt)
}

func (h *OAuthHandler) YandexRedirect(c *gin.Context) {
	if h.yandex.ClientID == "" {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Yandex+OAuth+not+configured")
		return
	}
	state := newState()
	c.SetCookie(oauthStateCookie, state, 300, "/", "", false, true)
	c.Redirect(http.StatusFound, h.yandex.AuthCodeURL(state))
}

func (h *OAuthHandler) YandexCallback(c *gin.Context) {
	if err := h.verifyState(c); err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Invalid+state")
		return
	}
	token, err := h.yandex.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Token+exchange+failed")
		return
	}
	info, err := fetchYandexUser(token.AccessToken)
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Failed+to+fetch+user+info")
		return
	}
	locale := c.GetString("locale")
	jwt, _, err := h.auth.LoginWithOAuth("yandex", info.ID, info.Email, locale)
	if err != nil {
		c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?error=Login+failed")
		return
	}
	c.Redirect(http.StatusFound, h.frontendURL+"/auth/callback?token="+jwt)
}

func (h *OAuthHandler) verifyState(c *gin.Context) error {
	cookie, err := c.Cookie(oauthStateCookie)
	if err != nil || cookie != c.Query("state") {
		return fmt.Errorf("state mismatch")
	}
	return nil
}

func newState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type googleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func fetchGoogleUser(accessToken string) (*googleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var info googleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	if info.ID == "" {
		return nil, fmt.Errorf("empty user id from Google")
	}
	return &info, nil
}

type yandexUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"default_email"`
}

func fetchYandexUser(accessToken string) (*yandexUserInfo, error) {
	req, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var info yandexUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	if info.ID == "" {
		return nil, fmt.Errorf("empty user id from Yandex")
	}
	return &info, nil
}
