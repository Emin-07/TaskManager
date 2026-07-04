package handler

import "time"

type UserResponse struct {
	Id        int       `json:"id" form:"id" binding:"required,gte=1"`
	Username  string    `json:"username" form:"username" binding:"omitempty"`
	Role      string    `json:"role" form:"role" binding:"oneof=user admin"`
	Email     string    `json:"email" form:"email" binding:"email"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
}

type UserRequest struct {
	Username string `json:"username" form:"username" binding:"omitempty"`
	Role     string `json:"role" form:"role" binding:"oneof=user admin"`
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"min=8,max=32"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Password string `json:"password" form:"password" binding:"min=8,max=32"`
}
