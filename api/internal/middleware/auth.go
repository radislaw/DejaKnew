package middleware

import (
	"net/http"
	"strings"

	"flashcards/api/internal/auth"

	"github.com/gin-gonic/gin"
)

// ContextUserID is the gin context key holding the authenticated user ID.
const ContextUserID = "userID"

// RequireAuth validates the Bearer token and stores the user ID in context.
func RequireAuth(tokens *auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or malformed Authorization header"})
			return
		}
		userID, err := tokens.Verify(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set(ContextUserID, userID)
		c.Next()
	}
}
