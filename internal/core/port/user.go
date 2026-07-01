package port

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type UserService interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	//GetUserTasks(ctx context.Context, id string) ([]*domain.Task, error)
	List(ctx context.Context) ([]*domain.User, error)
	Patch(ctx context.Context, user *domain.SignupUser, id string) error
	Insert(ctx context.Context, user *domain.SignupUser) error
	Delete(ctx context.Context, id string) error
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*repo.UserDb, error)
	GetById(ctx context.Context, id int) (*repo.UserDb, error)
	GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error)
	List(ctx context.Context) ([]*repo.UserDb, error)
	Insert(ctx context.Context, username, role, email, passwordHash string) error
	Patch(ctx context.Context, username, role, email, passwordHash string, id int) error
	Delete(ctx context.Context, id int) error
}
