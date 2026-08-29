package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

// @Summary login
// @Description login to get token for accessing other endpoints
// @Tags auth
// @Accept x-www-form-urlencoded,multipart/form-data
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password" minlength(8)
// @Success 200 {object} map[string]string "Token for accessing endpoints"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (u *UserHandler) Authenticate(c *gin.Context) {
	userToConvert, err := u.service.Authenticate(c.Request.Context(), c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token, err := u.tokenService.CreateToken(strconv.Itoa(userToConvert.ID), userToConvert.Role)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary signup
// @Description create a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body handler.UserRequest true "Credentials to signup"
// @Success 200 {object} map[string]string "Token for accessing endpoints"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (u *UserHandler) SignUp(c *gin.Context) {
	var userReq UserRequest

	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	serviceData, err := json.Marshal(map[string]any{
		"username":  userReq.Username,
		"role":      userReq.Role,
		"email":     userReq.Email,
		"password":  userReq.Password,
		"operation": domain.CreateOperation,
	})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = u.broker.Publish(map[string]string{"user-0": string(serviceData)}, domain.TopicUsers)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data, err := u.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	token, err := u.tokenService.CreateToken(data["id"], data["role"])
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary get user
// @Security ApiKeyAuth
// @Description get user from database by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} handler.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/:id [get]
func (u *UserHandler) GetById(c *gin.Context) {
	userToConvert, err := u.service.GetById(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			c.JSON(http.StatusNotFound, gin.H{"task": fmt.Sprintf("task with %s not found", c.Param("id"))})
			return
		}
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

// @Summary get users
// @Security ApiKeyAuth
// @Description get users from database by id
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} handler.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [get]
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

// @Summary patch user
// @Security ApiKeyAuth
// @Description patch user from database by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body handler.UserRequest true "user schema for patching a user"
// @Success 200 {object} map[string]string
// @Success 200 {object} map[string]string
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/:id [patch]
func (u *UserHandler) Patch(c *gin.Context) {
	var userReq UserRequest

	if err := c.ShouldBindBodyWithJSON(&userReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id := c.Param("id")
	serviceData, err := json.Marshal(map[string]any{
		"username":  userReq.Username,
		"role":      userReq.Role,
		"email":     userReq.Email,
		"password":  userReq.Password,
		"id":        id,
		"operation": domain.PatchOperation,
	})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	key := fmt.Sprintf("user-%s", id)
	err = u.broker.Publish(map[string]string{key: string(serviceData)}, domain.TopicUsers)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusAccepted)
}

// @Summary delete user
// @Security ApiKeyAuth
// @Description delete user from database by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/:id [delete]
func (u *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	serviceData, err := json.Marshal(map[string]any{
		"id":        id,
		"operation": domain.DeleteOperation,
	})

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	key := fmt.Sprintf("user-%s", id)
	err = u.broker.Publish(map[string]string{key: string(serviceData)}, domain.TopicUsers)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//if errors.Is(err, domain.ErrNoRecord) {
	//	c.JSON(http.StatusNoContent, gin.H{"no content": fmt.Sprintf("there isn't a task with id = %s", id)})
	//	return
	//}
	//c.AbortWithError(http.StatusBadRequest, err)
	//return

	c.Status(http.StatusAccepted)
}
