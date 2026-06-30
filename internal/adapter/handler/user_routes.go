package handler

import "github.com/gin-gonic/gin"

func (t UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/signup", t.SignUp)

	r.GET("/users/:id", t.GetById)
	r.DELETE("/users/:id", t.Delete)
	//r.PATCH("/users/:id", t.Patch)

	r.GET("/users", t.GetByEmail)
	r.GET("/users", t.GetUserTasks)

}
