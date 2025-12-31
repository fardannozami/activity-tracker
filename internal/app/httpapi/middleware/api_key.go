package middleware

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/gin-gonic/gin"
)

func APIKey(clients *postgres.ClientRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing X-API-Key"})
			return
		}

		prefix := service.APIKeyPrefix(key)
		row, ok, err := clients.GetByAPIKeyPrefix(c.Request.Context(), prefix)
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		if !service.ComparAPIKey(row.APIKeyHash, key) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}

		c.Set("client_id", row.ID)
		c.Next()
	}
}
