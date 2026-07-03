package port

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type TaskService interface {
	Get(ctx context.Context, id, userIdStr, role string) (*domain.Task, error)
	List(ctx context.Context, limit, offset, userIdStr string) ([]*domain.Task, error)
	Post(ctx context.Context, task *domain.CreateTask, userIdStr string) error
	Delete(ctx context.Context, id, userIdStr, role string) error
	Patch(ctx context.Context, task *domain.CreateTask, id, userIdStr, role string) error
}

type TaskRepo interface {
	List(ctx context.Context, limit, offset, userId int) ([]*repo.TaskDb, error)
	Get(ctx context.Context, id, userId int, role string) (*repo.TaskDb, error)
	Insert(ctx context.Context, title, text string, priority, expireDays, userId int) error
	Patch(ctx context.Context, title, text, userRole string, priority, expireDays, id, userId int) error
	Delete(ctx context.Context, id, userId int, role string) error
}
