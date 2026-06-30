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
	Priority   string
	ExpireDays string
}

type User struct {
	ID           int
	Username     string
	Role         string
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
}

type SignupUser struct {
	Username string
	Role     string
	Email    string
	Password string
}
