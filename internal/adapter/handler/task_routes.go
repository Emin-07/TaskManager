package handler

import "github.com/gin-gonic/gin"

func (t *TaskHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(SecureHeaders())
	authorized := r.Group("/")
	authorized.Use(AuthRequired(t.tokenService))
	{
		authorized.GET("/tasks", t.List)
		authorized.POST("/tasks", t.Post)

		authorized.GET("/tasks/:id", t.Get)
		authorized.DELETE("/tasks/:id", t.Delete)
		authorized.PATCH("/tasks/:id", t.Patch)
	}

}
