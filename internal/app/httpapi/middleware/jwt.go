package middleware

import (
	"net/http"
	"strings"

	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func JWT(ts service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		clientID, err := ts.Verify(token)
		if err != nil || clientID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("client_id", clientID)
		c.Next()
	}
}
