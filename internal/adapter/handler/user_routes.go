package handler

import "github.com/gin-gonic/gin"

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/users", u.SignUp)

	r.GET("/users/:id", u.GetById)
	r.DELETE("/users/:id", u.Delete)
	r.PATCH("/users/:id", u.Patch)

	//r.GET("/users/:email", u.GetByEmail)

	r.GET("/users", u.ListUsers)

}
