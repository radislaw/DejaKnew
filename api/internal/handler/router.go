package handler

import (
	"net/http"

	"flashcards/api/internal/auth"
	"flashcards/api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router wires all routes. New resource handlers are registered here.
func Router(tokens *auth.TokenManager, authHandler *AuthHandler, deckHandler *DeckHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(tokens))
		{
			protected.GET("/me", authHandler.Me)

			protected.POST("/decks", deckHandler.Create)
			protected.GET("/decks", deckHandler.List)
			protected.GET("/decks/:id", deckHandler.Get)
			protected.PATCH("/decks/:id", deckHandler.Update)
			protected.DELETE("/decks/:id", deckHandler.Delete)
		}
	}

	return r
}
