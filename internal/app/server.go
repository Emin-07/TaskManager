package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *App) NewServer() *http.Server {
	return &http.Server{
		Addr:         app.Cfg.Addr,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func (app *App) routes() http.Handler {
	router := gin.Default()
	app.taskHandler.RegisterRoutes(router)
	app.userHandler.RegisterRoutes(router)
	return router
}
