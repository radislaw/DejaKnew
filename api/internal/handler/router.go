package handler

import (
	"net/http"
	"time"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Router wires all routes. New resource handlers are registered here.
func Router(tokens *auth.TokenManager, authHandler *AuthHandler, oauthHandler *OAuthHandler, deckHandler *DeckHandler, cardHandler *CardHandler, studyHandler *StudyHandler, profileHandler *ProfileHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:5176"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept-Language"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	api.Use(middleware.LocaleFromRequest())
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/google", oauthHandler.GoogleRedirect)
		api.GET("/auth/google/callback", oauthHandler.GoogleCallback)
		api.GET("/auth/yandex", oauthHandler.YandexRedirect)
		api.GET("/auth/yandex/callback", oauthHandler.YandexCallback)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(tokens))
		{
			protected.GET("/me", authHandler.Me)

			protected.POST("/decks", deckHandler.Create)
			protected.GET("/decks", deckHandler.List)
			protected.GET("/decks/:id", deckHandler.Get)
			protected.PATCH("/decks/:id", deckHandler.Update)
			protected.DELETE("/decks/:id", deckHandler.Delete)

			protected.POST("/decks/:id/cards", cardHandler.Create)
				protected.POST("/decks/:id/import", cardHandler.Import)
				protected.GET("/decks/:id/export", cardHandler.ExportDeck)
				protected.GET("/export", cardHandler.ExportAll)
			protected.PATCH("/cards/:id", cardHandler.Update)
			protected.DELETE("/cards/:id", cardHandler.Delete)

			protected.PATCH("/profile", profileHandler.Update)
			protected.PUT("/profile/password", profileHandler.ChangePassword)
			protected.DELETE("/profile", profileHandler.Delete)

			protected.GET("/stats", studyHandler.GetStats)
			protected.GET("/study/:deckId/queue", studyHandler.GetQueue)
			protected.POST("/study/:cardId/review", studyHandler.SubmitReview)
		}
	}

	return r
}
