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
func Router(tokens *auth.TokenManager, authHandler *AuthHandler, deckHandler *DeckHandler, cardHandler *CardHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

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

			protected.POST("/decks/:id/cards", cardHandler.Create)
			protected.PATCH("/cards/:id", cardHandler.Update)
			protected.DELETE("/cards/:id", cardHandler.Delete)
		}
	}

	return r
}
