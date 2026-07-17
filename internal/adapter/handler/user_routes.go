package handler

import (
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(SecureHeaders())
	r.POST("/users", u.SignUp)
	r.POST("/login", u.Authenticate)

	authorized := r.Group("/")
	authorized.Use(AuthRequiredUsers(u))
	{
		authorized.GET("/users/:id", u.GetById)
		authorized.DELETE("/users/:id", u.Delete)
		authorized.PATCH("/users/:id", u.Patch)

		authorized.GET("/users", u.ListUsers)
	}

}
