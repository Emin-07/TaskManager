package handler

import "github.com/Emin-07/TaskManager/internal/core/port"

type TaskHandler struct {
	service             port.TaskService
	tokenService        port.TokenService
	rateAndCacheService port.RateAndCacheService
}

func NewTaskHandler(service port.TaskService, tokenService port.TokenService, rateAndCacheService port.RateAndCacheService) TaskHandler {
	return TaskHandler{service: service, tokenService: tokenService, rateAndCacheService: rateAndCacheService}
}

type UserHandler struct {
	service             port.UserService
	tokenService        port.TokenService
	rateAndCacheService port.RateAndCacheService
}

func NewUserHandler(service port.UserService, tokenService port.TokenService, rateAndCacheService port.RateAndCacheService) UserHandler {
	return UserHandler{service: service, tokenService: tokenService, rateAndCacheService: rateAndCacheService}
}
