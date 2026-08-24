package port

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type UserRepo interface {
	Authenticate(ctx context.Context, email string) (*repo.UserDb, error)
	GetById(ctx context.Context, id int) (*repo.UserDb, error)
	GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error)
	List(ctx context.Context) ([]*repo.UserDb, error)
	Insert(ctx context.Context, username, role, email string, passwordHash []byte) error
	Patch(ctx context.Context, username, role, email string, passwordHash []byte, id int) error
	Delete(ctx context.Context, id int) error
}

type UserService interface {
	Authenticate(ctx context.Context, email string, password string) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	//GetUserTasks(ctx context.Context, id string) ([]*domain.Task, error)
	List(ctx context.Context) ([]*domain.User, error)
	Patch(ctx context.Context, user *domain.SignupUser, id string) error
	Insert(ctx context.Context, user *domain.SignupUser) error
	Delete(ctx context.Context, id string) error
}
