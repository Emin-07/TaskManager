package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

// @Summary get task
// @Security ApiKeyAuth
// @Description get task from database by id
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} handler.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/:id [get]
func (t *TaskHandler) Get(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	taskToConvert, err := t.service.Get(c.Request.Context(), c.Param("id"), data["id"], data["role"])
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			c.JSON(http.StatusNotFound, gin.H{"task": fmt.Sprintf("task with %d not found", data["id"])})
			return
		}
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	task := &TaskResponse{
		Id:       taskToConvert.ID,
		Title:    taskToConvert.Title,
		Text:     taskToConvert.Text,
		Priority: taskToConvert.Priority,
		Expires:  taskToConvert.Expires,
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// @Summary list tasks
// @Security ApiKeyAuth
// @Description list tasks from database using offset and limit
// @Tags tasks
// @Accept json
// @Produce json
// @Param limit query int false "Limit for tasks listing" default(5)
// @Param offset query int false "Offset for tasks listing" default(0)
// @Success 200 {array} handler.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (t *TaskHandler) List(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	tasksToConvert, err := t.service.List(c.Request.Context(), limit, offset, data["id"])
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var tasks []*TaskResponse
	for _, task := range tasksToConvert {
		tasks = append(tasks, &TaskResponse{
			Id:       task.ID,
			Title:    task.Title,
			Text:     task.Text,
			Priority: task.Priority,
			Expires:  task.Expires,
		})
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

// @Summary add task
// @Security ApiKeyAuth
// @Description add task into database
// @Tags tasks
// @Accept json
// @Produce json
// @Param tasks body handler.TaskRequest true "task schema for creating a new user"
// @Success 201 {array} handler.TaskResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (t *TaskHandler) Post(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	var taskReq TaskRequest
	if err := c.ShouldBindBodyWithJSON(&taskReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = t.service.Post(c.Request.Context(), &domain.CreateTask{
		Title:      taskReq.Title,
		Text:       taskReq.Text,
		Priority:   taskReq.Priority,
		ExpireDays: taskReq.ExpireDays,
	}, data["id"])

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusCreated, "/tasks")
}

// @Summary delete task
// @Security ApiKeyAuth
// @Description delete task from database
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "id for deleting a task"
// @Success 200 {object} map[string]string "Returns a confirmation message"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/:id [delete]
func (t *TaskHandler) Delete(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			c.JSON(http.StatusNoContent, gin.H{"no content": fmt.Sprintf("there isn't a task with id = %d", data["id"])})
			return
		}
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	err = t.service.Delete(c.Request.Context(), c.Param("id"), data["id"], data["role"])
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("task with id %d deleted", data["id"])})
}

// @Summary patch task
// @Security ApiKeyAuth
// @Description patch task from database
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "id for patching a task"
// @Param task body handler.TaskRequest true "task schema for patching a task"
// @Success 200 {object} map[string]string "Returns a confirmation message"
// @Success 204 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/:id [patch]
func (t *TaskHandler) Patch(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	id := c.Param("id")
	var taskReq TaskRequest
	if err := c.ShouldBindBodyWithJSON(&taskReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = t.service.Patch(c.Request.Context(), &domain.CreateTask{
		Title:      taskReq.Title,
		Text:       taskReq.Text,
		Priority:   taskReq.Priority,
		ExpireDays: taskReq.ExpireDays,
	}, id, data["id"], data["role"])

	if err != nil {
		if errors.Is(err, domain.ErrNoRecord) {
			c.JSON(http.StatusNoContent, gin.H{"no content": fmt.Sprintf("there isn't a task with id = %d", data["id"])})
			return
		}
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("task with id %d patched", data["id"])})

}
