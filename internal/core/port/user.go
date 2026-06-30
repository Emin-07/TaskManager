package port

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id int) (*domain.User, error)
	GetUserTasks(ctx context.Context, id int) ([]*domain.Task, error)
	Insert(ctx context.Context, user *domain.SignupUser) (int, error)
	Delete(ctx context.Context, id int) error
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*repo.UserDb, error)
	GetById(ctx context.Context, id int) (*repo.UserDb, error)
	GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error)
	Insert(ctx context.Context, username, role, email, passwordHash string) (int, error)
	Delete(ctx context.Context, id int) error
}
