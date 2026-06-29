package domain

import "time"

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Comment   string    `json:"comment"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}

//type User struct {
//	ID      string `json:"id"`
//	Date    string `json:"date"`
//	Title   string `json:"title"`
//	Comment string `json:"comment"`
//	Repeat  string `json:"repeat"`
//}
