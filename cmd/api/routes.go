package main

import (
	"net/http"
	"text/template"

	"github.com/Emin-07/TaskManager/ui"
	"github.com/gin-gonic/gin"
)

func (app *application) router() *gin.Engine {
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"humanDate": humanDate,
	})

	router.LoadHTMLFS(http.FS(ui.Files), "**/*.html")

	router.GET("/", app.home)
	router.GET("/json", app.homeJSON)
	router.GET("/tasks/:id", app.view)
	router.DELETE("/tasks/:id", app.delete)
	router.PATCH("/tasks/:id", app.patch)
	router.POST("/tasks", app.insert)
	router.POST("/tasks/refresh", app.refreshTasks)

	// cssFS, _ := fs.Sub(ui.Files, "css")
	router.Static("/static", "./ui/static")

	return router
}
