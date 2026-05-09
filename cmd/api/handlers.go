package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Emin-07/TaskManager/internal/models"
	"github.com/gin-gonic/gin"
)

func (app *application) home(ctx *gin.Context) {
	tasks, err := app.tasks.Latest()
	if err != nil {
		app.serverError(ctx, err)
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

func (app *application) view(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		app.badRequetDetailed(ctx, "The 'id' parameter must be an integer", "received: '%v'", ctx.Param("id"))
		return
	}
	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(ctx, "There is no task with id: '%v'", id)
		} else {

		}
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

type CreateTask struct {
	Title      string `form:"title" json:"title" binding:"required,min=2"`
	Text       string `form:"text" json:"text" binding:"omitempty,min=1"`
	Priority   string `form:"priority" json:"priority" binding:"omitempty,oneof=high middle low"`
	ExpireDays int    `form:"expires" json:"expires" binding:"omitempty,gte=1,lte=365"`
}

func (app *application) insert(ctx *gin.Context) {
	var task CreateTask
	if err := ctx.ShouldBind(&task); err != nil {
		app.badRequet(ctx, err)
		return
	}
	if task.Priority == "" {
		task.Priority = "middle"
	}
	if task.ExpireDays == 0 {
		task.ExpireDays = 7
	}
	id, err := app.tasks.Insert(task.Title, task.Text, task.Priority, task.ExpireDays)
	if err != nil {
		app.serverError(ctx, err)
	}
	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%v", id))
}

func (app *application) delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		app.clientError(ctx, http.StatusBadRequest)
	}

	err = app.tasks.Delete(id)
	if err != nil {
		app.serverError(ctx, err)
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}

func (app *application) refreshTasks(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "99999999")
	refreshAmountStr := ctx.DefaultQuery("days", "7")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		app.badRequetDetailed(ctx, "The 'limit' query parameter must be an integer", "received: %v'", limitStr)
		return
	}
	refreshAmount, err := strconv.Atoi(refreshAmountStr)
	if err != nil {
		app.badRequetDetailed(ctx, "The 'days' query parameter must be an integer", "received: %v'", refreshAmountStr)
		return
	}

	app.tasks.RefreshTasks(limit, refreshAmount)

	ctx.Redirect(http.StatusSeeOther, "/")
}
