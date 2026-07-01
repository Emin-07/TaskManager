package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

// TODO: Add getting tasks by users
// TODO: Handle get email
// TODO: Add limit and offset for listing users
func (u *UserHandler) GetByEmail(c *gin.Context) {
	userToConvert, err := u.service.GetByEmail(c.Request.Context(), c.Param("email"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": UserResponse{
		Id:        userToConvert.ID,
		Username:  userToConvert.Username,
		Role:      userToConvert.Role,
		Email:     userToConvert.Email,
		CreatedAt: userToConvert.CreatedAt}})
}

func (u *UserHandler) GetById(c *gin.Context) {
	userToConvert, err := u.service.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": UserResponse{
		Id:        userToConvert.ID,
		Username:  userToConvert.Username,
		Role:      userToConvert.Role,
		Email:     userToConvert.Email,
		CreatedAt: userToConvert.CreatedAt}})
}

func (u *UserHandler) ListUsers(c *gin.Context) {
	usersToConvert, err := u.service.List(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var users []*UserResponse
	for _, user := range usersToConvert {
		users = append(users, &UserResponse{
			Id:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			Email:     user.Email,
			CreatedAt: user.CreatedAt})
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (u *UserHandler) SignUp(c *gin.Context) {
	var userReq UserRequest

	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := u.service.Insert(c.Request.Context(), &domain.SignupUser{
		Username: userReq.Username,
		Role:     userReq.Role,
		Email:    userReq.Email,
		Password: userReq.Password,
	})

	if err = c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/users")
}

func (u *UserHandler) Patch(c *gin.Context) {
	var userReq UserRequest

	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := c.Param("id")
	err := u.service.Patch(c.Request.Context(), &domain.SignupUser{
		Username: userReq.Username,
		Role:     userReq.Role,
		Email:    userReq.Email,
		Password: userReq.Password,
	}, id)

	if err = c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/users/%d", id))
}

func (u *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := u.service.Delete(c.Request.Context(), id); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/users/%v", id))
}
