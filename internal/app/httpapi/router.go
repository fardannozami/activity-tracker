package httpapi

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/app/httpapi/handler"
	"github.com/fardannozami/activity-tracker/internal/app/httpapi/middleware"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/gin-gonic/gin"
)

type Dependency struct {
	ClientHandler *handler.ClientHandler
	LogHandler    *handler.LogHandler
	UsageHandler  *handler.UsageHandler
	ClientsRepo   *postgres.ClientRepo
}

func NewRouter(d Dependency) *gin.Engine {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	api := r.Group("/api")
	api.POST("/register", d.ClientHandler.RegisterClient)
	api.POST("/logs", middleware.APIKey(d.ClientsRepo), d.LogHandler.Create)

	usage := api.Group("/usage", middleware.APIKey(d.ClientsRepo))
	usage.GET("/daily", d.UsageHandler.Daily)
	usage.GET("/top", d.UsageHandler.Top)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}
