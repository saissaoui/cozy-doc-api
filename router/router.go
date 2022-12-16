package router

import (
	"time"

	"cozy-doc-api/handlers"
	"cozy-doc-api/services"
	"cozy-doc-api/utils"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(service *services.DocsServiceImpl) *gin.Engine {

	// Set the default gin router
	r := gin.New()

	r.Use(gin.Recovery())

	// Initialize middlewares
	initializeMiddlewares(r)

	// Initialize routes
	initializeRoutes(r, service)

	return r

}

func initializeRoutes(r *gin.Engine, service *services.DocsServiceImpl) {
	untracedGroup := r.Group("/")
	untracedGroup.Use(Ginzap(utils.Logger, time.RFC3339, true, false))

	// health
	untracedGroup.GET("/", handlers.GetHealth)
	untracedGroup.GET("/health", handlers.GetHealth)

	tracedGroup := r.Group("/")
	tracedGroup.Use(Ginzap(utils.Logger, time.RFC3339, true, true))

	//docs
	tracedGroup.POST("/docs/bulk", handlers.BulkInsertDocs(*service))

	// fallback
	r.NoRoute(handlers.NoRoute)
}
