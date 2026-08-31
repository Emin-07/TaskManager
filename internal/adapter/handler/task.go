package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	taskID := c.Param("id")
	userID := data["id"]

	taskFromCache, err := t.rateAndCacheService.Get(c.Request.Context(), "task", taskID, userID)
	if errors.Is(err, domain.ErrKeyNotFound) {
		taskToConvert, err := t.service.Get(c.Request.Context(), taskID, userID, data["role"])
		if err != nil {
			if errors.Is(err, domain.ErrNoRecord) {
				c.JSON(http.StatusNotFound, gin.H{"task": fmt.Sprintf("task with %s not found", taskID)})
				return
			}
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		task := TaskResponse{
			Id:       taskToConvert.ID,
			Title:    taskToConvert.Title,
			Text:     taskToConvert.Text,
			Priority: taskToConvert.Priority,
			Expires:  taskToConvert.Expires,
		}

		err = t.rateAndCacheService.Set(c.Request.Context(), "task", strconv.Itoa(task.Id), userID, task, time.Minute*15)
		if err != nil {
			log.Printf("error occurred when setting cache, err: %v\n", err)
		}
		c.JSON(http.StatusOK, gin.H{"task": task})
		return
	} else if err != nil {
		log.Printf("error occurred when getting from cache, err: %v\n", err)
		return
	}
	var res TaskResponse
	err = json.Unmarshal([]byte(taskFromCache), &res)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": res})
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
	if limit == "" {
		limit = "5"
	}
	if offset == "" {
		offset = "0"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	userID := data["id"]

	var cacheTasks = make([]string, limitInt)
	var cacheErr error

	for i := offsetInt; limitInt > i-offsetInt; i++ {
		taskFromCache, err := t.rateAndCacheService.Get(c.Request.Context(), "task", strconv.Itoa(i), userID)
		if err != nil {
			cacheErr = err
			break
		}
		cacheTasks[i-offsetInt] = taskFromCache
	}

	if cacheErr == nil {
		results := make([]TaskResponse, len(cacheTasks))
		for i, cacheTask := range cacheTasks {
			err := json.Unmarshal([]byte(cacheTask), &results[i])
			if err != nil {
				log.Printf("err: %v, task: %v\n", err, cacheTask)
			}
		}
		c.JSON(http.StatusOK, gin.H{"tasks": results})
		return
	} else if !errors.Is(err, domain.ErrKeyNotFound) && err != nil {
		log.Printf("error occurred when getting from cache, err: %v\n", err)
	} else {
		tasksToConvert, err := t.service.List(c.Request.Context(), limitInt, offsetInt, userID)
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var tasks []TaskResponse
		for _, task := range tasksToConvert {
			tasks = append(tasks, TaskResponse{
				Id:       task.ID,
				Title:    task.Title,
				Text:     task.Text,
				Priority: task.Priority,
				Expires:  task.Expires,
			})
		}

		c.JSON(http.StatusOK, gin.H{"tasks": tasks})

		for i := offsetInt; limitInt > i-offsetInt; i++ {
			err := t.rateAndCacheService.Set(c.Request.Context(), "task", strconv.Itoa(i+1), userID, tasks[i], time.Minute*15)
			if err != nil {
				log.Printf("error occurred when setting cache, err: %v\n", err)
			}
		}
	}
}

// @Summary add task
// @Security ApiKeyAuth
// @Description add task into database
// @Tags tasks
// @Accept json
// @Produce json
// @Param tasks body handler.TaskRequest true "task schema for creating a new user"
// @Success 202 {array}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (t *TaskHandler) Post(c *gin.Context) {
	data, err := t.tokenService.ParseFromRequest(c.Request)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	var taskReq TaskRequest
	if err := c.ShouldBindBodyWithJSON(&taskReq); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	serviceData, err := json.Marshal(map[string]any{
		"title":       taskReq.Title,
		"text":        taskReq.Text,
		"priority":    taskReq.Priority,
		"expire_days": taskReq.ExpireDays,
		"user_id":     data["id"],
		"operation":   domain.CreateOperation,
	})

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = t.broker.Publish(map[string]string{"task-0": string(serviceData)}, domain.TopicTasks)

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusAccepted)
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
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	taskID := c.Param("id")
	userID := data["id"]
	serviceData, err := json.Marshal(map[string]any{
		"id":        taskID,
		"role":      data["role"],
		"user_id":   userID,
		"operation": domain.DeleteOperation,
	})

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	key := fmt.Sprintf("task-%s", taskID)
	err = t.broker.Publish(map[string]string{key: string(serviceData)}, domain.TopicTasks)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Status(http.StatusAccepted)

	err = t.rateAndCacheService.Del(c.Request.Context(), "task", taskID, userID)
	if err != nil {
		log.Println(err)
	}
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
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	userID := data["id"]
	taskID := c.Param("id")

	var taskReq TaskRequest
	if err := c.ShouldBindBodyWithJSON(&taskReq); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	serviceData, err := json.Marshal(map[string]any{
		"id":          taskID,
		"title":       taskReq.Title,
		"text":        taskReq.Text,
		"priority":    taskReq.Priority,
		"expire_days": taskReq.ExpireDays,
		"user_id":     userID,
		"role":        data["role"],
		"operation":   domain.PatchOperation,
	})

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	key := fmt.Sprintf("task-%s", taskID)
	err = t.broker.Publish(map[string]string{key: string(serviceData)}, domain.TopicTasks)

	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusAccepted)

	err = t.rateAndCacheService.Del(c.Request.Context(), "task", taskID, userID)
	if err != nil {
		log.Println(err)
	}
}
