package service

import (
	"context"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (u *UserServ) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, PasswordHash: user.PasswordHash, CreatedAt: user.CreatedAt}, nil
}

func (u *UserServ) GetById(ctx context.Context, id int) (*domain.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, PasswordHash: user.PasswordHash, CreatedAt: user.CreatedAt}, nil
}

func (u *UserServ) GetUserTasks(ctx context.Context, id int) ([]*domain.Task, error) {
	tasksToConvert, err := u.repo.GetUserTasks(ctx, id)
	if err != nil {
		return nil, err
	}

	var tasks []*domain.Task
	for _, task := range tasksToConvert {
		newTask := &domain.Task{
			ID:        task.Id,
			Title:     task.Title,
			Text:      task.Text,
			Priority:  task.Priority,
			CreatedAt: task.CreatedAt,
			Expires:   task.Expires,
			UserId:    task.UserId,
		}
		tasks = append(tasks, newTask)
	}
	return tasks, nil
}

func (u *UserServ) Insert(ctx context.Context, user *domain.SignupUser) (int, error) {
	//TODO: Password hashing
	hashed_password := user.Password

	return u.repo.Insert(ctx, user.Username, user.Role, user.Email, hashed_password)
}

func (u *UserServ) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}
