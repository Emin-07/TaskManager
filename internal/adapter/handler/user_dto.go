package handler

import "time"

type UserResponse struct {
	Id        int       `json:"id"`
	Username  string    `json:"title"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
