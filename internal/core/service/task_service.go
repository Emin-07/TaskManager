package service

import "github.com/Emin-07/TaskManager/internal/core/port"

type TaskServ struct {
	repo port.TaskRepo
}

func NewTaskService(repo port.TaskRepo) TaskServ {
	return TaskServ{repo: repo}
}
