package domain

import "time"

type Task struct {
	ID        int
	Title     string
	Text      string
	Priority  int
	CreatedAt time.Time
	Expires   time.Time
	UserId    int
}

type CreateTask struct {
	Title      string
	Text       string
	Priority   int
	ExpireDays int
}
