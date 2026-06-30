package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (t *TaskHandler) Get(c *gin.Context) {
	taskToConvert, err := t.service.Get(c.Request.Context(), c.Param("id"))
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
	tasksToConvert, err := t.service.List(c.Request.Context(), limit, offset)
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
	// TODO: get current user id
	userId := 1
	var reqTask TaskRequest
	if err := c.ShouldBindBodyWithJSON(&reqTask); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := t.service.Post(c.Request.Context(), &domain.CreateTask{
		Title:      reqTask.Title,
		Text:       reqTask.Text,
		Priority:   reqTask.Priority,
		ExpireDays: reqTask.ExpireDays,
	}, userId)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%d", id))

}
func (t *TaskHandler) Delete(c *gin.Context) {
	err := t.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks")
}

func (t *TaskHandler) Patch(c *gin.Context) {
	id := c.Param("id")
	var reqTask TaskRequest
	if err := c.ShouldBindBodyWithJSON(&reqTask); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := t.service.Patch(c.Request.Context(), &domain.CreateTask{
		Title:      reqTask.Title,
		Text:       reqTask.Text,
		Priority:   reqTask.Priority,
		ExpireDays: reqTask.ExpireDays,
	}, id)

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%d", id))

}
