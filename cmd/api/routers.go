package main

import "github.com/gin-gonic/gin"

func (app *application) router() *gin.Engine {
	router := gin.Default()

	router.GET("/", app.home)
	router.GET("/tasks/:id", app.view)
	router.DELETE("/tasks/:id", app.delete)
	router.PATCH("/tasks/:id", app.patch)
	router.POST("/tasks", app.insert)
	router.POST("/tasks/refresh", app.refreshTasks)
	return router
}
