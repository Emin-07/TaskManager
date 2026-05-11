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
		return
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{"tasks": tasks})
}
func (app *application) homeJSON(ctx *gin.Context) {
	tasks, err := app.tasks.Latest()
	if err != nil {
		app.serverError(ctx, err)
	}
	app.JSON(ctx, http.StatusOK, gin.H{"tasks": tasks})

}

func (app *application) view(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		app.badRequestDetailed(ctx, "The 'id' parameter must be an integer", "received: '%v'", ctx.Param("id"))
		return
	}
	task, err := app.tasks.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(ctx, "There is no task with id: '%v'", id)
		} else {

		}
	}
	app.JSON(ctx, http.StatusOK, gin.H{"task": task})
}

type Task struct {
	Title      string `form:"title" json:"title" binding:"omitempty,min=2"`
	Text       string `form:"text" json:"text" binding:"omitempty,min=2"`
	Priority   string `form:"priority" json:"priority" binding:"omitempty,oneof=high middle low"`
	ExpireDays int    `form:"expires" json:"expires" binding:"omitempty,gte=1,lte=365"`
}

func (app *application) patch(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		app.badRequestDetailed(ctx, "The 'id' parameter must be an integer", "received: '%v'", ctx.Param("id"))
		return
	}
	var task = Task{}
	if err := ctx.ShouldBind(&task); err != nil {
		app.badRequest(ctx, err)
		return
	}

	err = app.tasks.Patch(task.Title, task.Text, task.Priority, id, task.ExpireDays)
	if err != nil {
		app.serverError(ctx, err)
		return
	}

	ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/tasks/%v", id))
}

func (app *application) insert(ctx *gin.Context) {
	var task = Task{Priority: "middle", ExpireDays: 7}
	if err := ctx.ShouldBind(&task); err != nil {
		app.badRequest(ctx, err)
		return
	}

	id, err := app.tasks.Insert(task.Title, task.Text, task.Priority, task.ExpireDays)
	if err != nil {
		app.serverError(ctx, err)
		return
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
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/")
}

type TaskQueryParams struct {
	Limit int `form:"limit" binding:"omitempty,gte=1" default:"9999999"`
	Days  int `form:"days" binding:"omitempty,gte=1,lte=365" default:"7"`
}

func (app *application) refreshTasks(ctx *gin.Context) {
	var queryParams TaskQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		app.badRequest(ctx, err)
		return
	}

	app.tasks.RefreshTasks(queryParams.Limit, queryParams.Days)

	ctx.Redirect(http.StatusSeeOther, "/")
}
