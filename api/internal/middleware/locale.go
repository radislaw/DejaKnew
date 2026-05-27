package middleware

import "github.com/gin-gonic/gin"

const ContextLocale = "locale"

// LocaleFromRequest reads the Accept-Language header and stores a normalised
// locale ("en" or "ru") in the Gin context.
func LocaleFromRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		if len(lang) >= 2 && lang[:2] == "ru" {
			c.Set(ContextLocale, "ru")
		} else {
			c.Set(ContextLocale, "en")
		}
		c.Next()
	}
}
