package app

import (
	"os"

	"github.com/Emin-07/TaskManager/internal/adapter/handler"
)

type Config struct {
	Port string
}

type App struct {
	Cfg         *Config
	taskHandler *handler.TaskHandler
	userHandler *handler.UserHandler
}

type Option func(*App)

func NewApp(opts ...Option) *App {
	s := &App{
		Cfg: &Config{
			Port: os.Getenv("TODO_PORT"),
		},
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithTaskHandler(taskHandler *handler.TaskHandler) Option {
	return func(a *App) {
		a.taskHandler = taskHandler
	}
}

func WithUserHandler(userHandler *handler.UserHandler) Option {
	return func(a *App) {
		a.userHandler = userHandler
	}
}
