package app

import (
	"fmt"

	"cozy-doc-api/router"
	"cozy-doc-api/services"
	"cozy-doc-api/utils"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config utils.AppConfig
	Router *gin.Engine
}

func New() *App {
	app := &App{}
	app.setup()
	return app
}

func (app *App) setup() {

	// Load configuration
	config := utils.LoadConfig()

	// Initialize Services
	servicesWrapper := services.InitService(config)

	// Initialize Router
	r := router.InitializeRouter(servicesWrapper)

	app.Config = config
	app.Router = r

}

func (app *App) Run() {

	// Serving application

	port := app.Config.Port

	app.Router.Run(fmt.Sprintf(":%d", port))

}
