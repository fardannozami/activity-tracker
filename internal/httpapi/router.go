package httpapi

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/httpapi/handler"
	"github.com/gin-gonic/gin"
)

type Dependency struct {
	ClientHandler *handler.ClientHandler
}

func NewRouter(d Dependency) *gin.Engine {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	api := r.Group("/api")
	api.POST("/register", d.ClientHandler.RegisterClient)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}
