package httpapi

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/app/httpapi/handler"
	"github.com/fardannozami/activity-tracker/internal/app/httpapi/middleware"
	_ "github.com/fardannozami/activity-tracker/internal/docs"
	"github.com/fardannozami/activity-tracker/internal/domain/service"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Dependency struct {
	ClientHandler *handler.ClientHandler
	LogHandler    *handler.LogHandler
	UsageHandler  *handler.UsageHandler
	AuthHandler   *handler.AuthHandler
	ClientsRepo   *postgres.ClientRepo
	Token         service.TokenService
}

func NewRouter(d Dependency) *gin.Engine {
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	api := r.Group("/api")
	api.POST("/register", d.ClientHandler.RegisterClient)
	api.POST("/auth/token", d.AuthHandler.Token)

	api.POST("/logs", middleware.APIKey(d.ClientsRepo), d.LogHandler.Create)

	usage := api.Group("/usage", middleware.JWT(d.Token))
	usage.GET("/daily", d.UsageHandler.Daily)
	usage.GET("/top", d.UsageHandler.Top)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}
