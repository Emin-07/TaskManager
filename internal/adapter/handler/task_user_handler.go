package handler

import "github.com/Emin-07/TaskManager/internal/core/port"

type TaskHandler struct {
	service      port.TaskService
	tokenService port.TokenService
}

func NewTaskHandler(service port.TaskService, tokenService port.TokenService) TaskHandler {
	return TaskHandler{service: service, tokenService: tokenService}
}

type UserHandler struct {
	service      port.UserService
	tokenService port.TokenService
}

func NewUserHandler(service port.UserService, tokenService port.TokenService) UserHandler {
	return UserHandler{service: service, tokenService: tokenService}
}
