package domain

import "time"

type User struct {
	ID        int
	Username  string
	Role      string
	Email     string
	CreatedAt time.Time
}

type SignupUser struct {
	Username string
	Role     string
	Email    string
	Password string
}
