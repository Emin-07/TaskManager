package handler

import "github.com/Emin-07/TaskManager/internal/core/port"

type TaskHandler struct {
	service             port.TaskService
	tokenService        port.TokenService
	rateAndCacheService port.RateAndCacheService
	broker              port.MessageBrokerOut
}

func NewTaskHandler(service port.TaskService, tokenService port.TokenService, rateAndCacheService port.RateAndCacheService, broker port.MessageBrokerOut) TaskHandler {
	return TaskHandler{service: service, tokenService: tokenService, rateAndCacheService: rateAndCacheService, broker: broker}
}

type UserHandler struct {
	service             port.UserService
	tokenService        port.TokenService
	rateAndCacheService port.RateAndCacheService
	broker              port.MessageBrokerOut
}

func NewUserHandler(service port.UserService, tokenService port.TokenService, rateAndCacheService port.RateAndCacheService, broker port.MessageBrokerOut) UserHandler {
	return UserHandler{service: service, tokenService: tokenService, rateAndCacheService: rateAndCacheService, broker: broker}
}
