package port

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type TaskService interface {
	Get(ctx context.Context, id string) (*domain.Task, error)
	List(ctx context.Context, limit, offset string) ([]*domain.Task, error)
	Post(ctx context.Context, task *domain.CreateTask, userId int) error
	Delete(ctx context.Context, id string) error
	Patch(ctx context.Context, task *domain.CreateTask, id string) error
}

type TaskRepo interface {
	List(ctx context.Context, limit, offset int) ([]*repo.TaskDb, error)
	Get(ctx context.Context, id int) (*repo.TaskDb, error)
	Insert(ctx context.Context, title, text string, priority, expireDays, userId int) error
	Patch(ctx context.Context, title, text string, priority, expireDays, id int) error
	Delete(ctx context.Context, id int) error
}
