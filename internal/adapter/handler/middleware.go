package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 && c.Writer.Written() {
			status := c.Writer.Status()
			if status == http.StatusOK || status == 0 {
				status = http.StatusInternalServerError
			}
			c.JSON(status, gin.H{"error": c.Errors[0].Error()})
		}
	}
}

func AuthRequiredTasks(t *TaskHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := t.tokenService.ParseFromRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("to access this endpoint you need to authorize first. err : %v", err)})
			return
		}
		c.Next()
	}
}

func AuthRequiredUsers(u *UserHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := u.tokenService.ParseFromRequest(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("to access this endpoint you need to authorize first. err : %v", err)})
			return
		}
		if data["role"] != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("to access this endpoint you need to be admin. : %v", err)})
			return
		}
		c.Next()
	}
}
