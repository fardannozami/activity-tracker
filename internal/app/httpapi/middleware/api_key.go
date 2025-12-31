package middleware

import (
	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/gin-gonic/gin"
)

func APIKey(clients *postgres.ClientRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			abortInvalidAPIKey(c)
			return
		}

		prefix := service.APIKeyPrefix(key)
		row, ok, err := clients.GetByAPIKeyPrefix(c.Request.Context(), prefix)
		if err != nil || !ok {
			abortInvalidAPIKey(c)
			return
		}
		if !service.ComparAPIKey(row.APIKeyHash, key) {
			abortInvalidAPIKey(c)
			return
		}

		c.Set("client_id", row.ID)
		c.Next()
	}
}
