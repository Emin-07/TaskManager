package service

import "github.com/Emin-07/TaskManager/internal/core/port"

type UserServ struct {
	repo port.UserRepo
}

func NewUserService(repo port.UserRepo) *UserServ {
	return &UserServ{repo: repo}
}
