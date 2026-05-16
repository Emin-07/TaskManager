package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) serverError(ctx *gin.Context, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logger.Error("server error", zap.String("trace", trace))
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": trace})
}

func (app *application) clientError(ctx *gin.Context, status int) {
	ctx.AbortWithStatusJSON(status, gin.H{"Error": http.StatusText(status)})
}

func (app *application) badRequestDetailed(ctx *gin.Context, message, detailFormat string, a ...any) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": http.StatusText(http.StatusBadRequest), "message": message, "detailFormat": fmt.Sprintf(detailFormat, a...)})
}

func (app *application) badRequest(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
}

func (app *application) notFound(ctx *gin.Context, detailFormat string, a ...any) {
	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Error": http.StatusText(http.StatusNotFound), "detail": fmt.Sprintf(detailFormat, a...)})
}

func (app *application) JSON(ctx *gin.Context, code int, obj any) {
	if app.readableJSON {
		ctx.IndentedJSON(code, obj)
	} else {
		ctx.JSON(code, obj)
	}
}

func (app *application) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (app *application) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}
