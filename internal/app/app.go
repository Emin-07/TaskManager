package app

import (
	"os"

	"github.com/Emin-07/TaskManager/internal/adapter/handler"
)

type Config struct {
	Addr     string
	Port     string
	Host     string
	Name     string
	Password string
	User     string
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
			Addr:     os.Getenv("WEB_ADDR"),
			Port:     os.Getenv("DB_PORT"),
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_NAME"),
			Password: os.Getenv("DB_PASSWORD"),
			User:     os.Getenv("DB_USER"),
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
