package handler

import "time"

type TaskResponse struct {
	Id       int       `json:"id" form:"id" binding:"required,gte=1"`
	Title    string    `json:"title" form:"title" binding:"required,min=1"`
	Text     string    `json:"text" form:"text" binding:"omitempty"`
	Priority int       `json:"priority" form:"priority" binding:"required,gte=0,lte=4"`
	Expires  time.Time `json:"expires" form:"expires"`
}

type TaskRequest struct {
	Title      string `json:"title" form:"title" binding:"required,min=1"`
	Text       string `json:"text" form:"text" binding:"omitempty"`
	Priority   int    `json:"priority" form:"priority" binding:"required,gte=0,lte=4"`
	ExpireDays int    `json:"expire_days"  form:"expire_days"`
}
