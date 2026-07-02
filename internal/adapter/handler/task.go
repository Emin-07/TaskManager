package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

// @Summary get task
// @Schemes
// @Description get task from db
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {string} {"task": {Id, Title, Text, Priority, Expires}}
// @Router /tasks/:id [get]
func (t *TaskHandler) Get(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	taskToConvert, err := t.service.Get(c.Request.Context(), c.Param("id"), data["id"])
	if err != nil {
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

	c.Redirect(http.StatusSeeOther, "/tasks")

}
func (t *TaskHandler) Delete(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	err = t.service.Delete(c.Request.Context(), c.Param("id"), data["id"])
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks")
}

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
	}, id, data["id"])

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%d", id))

}
