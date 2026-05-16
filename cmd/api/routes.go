package main

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/Emin-07/TaskManager/ui"
	"github.com/gin-gonic/gin"
)

func (app *application) router() *gin.Engine {
	router := gin.Default()

	router.Use(SetSecureHeaders(), app.ErrorHandler())

	router.SetFuncMap(template.FuncMap{
		"humanDate": humanDate,
	})

	router.LoadHTMLFS(http.FS(ui.Files), "**/*.html")
	router.GET("/", app.home)
	router.GET("/json", app.homeJSON)

	router.POST("/signup", app.signup)
	router.POST("/login", app.login)

	//router.Use(Protecter())

	router.GET("/tasks/:id", app.view)
	router.DELETE("/tasks/:id", app.delete)
	router.PATCH("/tasks/:id", app.patch)
	router.POST("/tasks", app.insert)
	router.POST("/tasks/refresh", app.refreshTasks)

	staticFiles, _ := fs.Sub(ui.Files, "static")
	router.StaticFS("static", http.FS(staticFiles))

	return router
}
