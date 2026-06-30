package handler

import "github.com/gin-gonic/gin"

func (t TaskHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/tasks/:id", t.Get)
	r.DELETE("/tasks/:id", t.Delete)
	r.PATCH("/tasks/:id", t.Patch)

	r.GET("/tasks", t.List)
	r.POST("/tasks", t.Post)

}
