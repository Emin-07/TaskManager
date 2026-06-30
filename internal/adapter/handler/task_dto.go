package handler

import "time"

type TaskResponse struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Priority int       `json:"priority"`
	Expires  time.Time `json:"expires"`
}

type TaskRequest struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	Priority   int    `json:"priority"`
	ExpireDays int    `json:"expire_days"`
}
