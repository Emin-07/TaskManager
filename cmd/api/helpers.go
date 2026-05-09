package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func (app *application) serverError(ctx *gin.Context, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": trace})
}

func (app *application) clientError(ctx *gin.Context, status int) {
	ctx.AbortWithStatusJSON(status, gin.H{"Error": http.StatusText(status)})
}

func (app *application) badRequetDetailed(ctx *gin.Context, message, detailFormat string, a ...any) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": message, "detailFormat": fmt.Sprintf(detailFormat, a...)})
}

func (app *application) badRequet(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
}

func (app *application) notFound(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": http.StatusText(http.StatusNotFound)})
}
