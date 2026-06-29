package handler

import "github.com/Emin-07/TaskManager/internal/core/port"

type TaskHandler struct {
	service port.TaskService
}

func NewTaskHandler(service port.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

type UserHandler struct {
	service port.UserService
}

func NewUserHandler(service port.UserService) *UserHandler {
	return &UserHandler{service: service}
}
