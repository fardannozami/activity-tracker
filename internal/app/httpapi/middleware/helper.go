package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortInvalidAPIKey(c *gin.Context) {
	c.AbortWithStatusJSON(
		http.StatusUnauthorized,
		gin.H{"error": "invalid api key"},
	)
}
