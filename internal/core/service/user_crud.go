package service

import (
	"context"
	"strconv"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (u UserServ) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, PasswordHash: user.PasswordHash, CreatedAt: user.CreatedAt}, nil
}

func (u UserServ) GetById(ctx context.Context, id string) (*domain.User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.GetById(ctx, idInt)
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, PasswordHash: user.PasswordHash, CreatedAt: user.CreatedAt}, nil
}

func (u UserServ) List(ctx context.Context) ([]*domain.User, error) {
	usersToConvert, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var users []*domain.User
	for _, user := range usersToConvert {
		newTask := &domain.User{
			ID:           user.Id,
			Username:     user.Username,
			Role:         user.Role,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
			PasswordHash: user.PasswordHash,
		}
		users = append(users, newTask)
	}
	return users, nil

}

//func (u UserServ) GetUserTasks(ctx context.Context, id string) ([]*domain.Task, error) {
//	idInt, err := strconv.Atoi(id)
//	if err != nil {
//		return nil, err
//	}
//	tasksToConvert, err := u.repo.GetUserTasks(ctx, idInt)
//	if err != nil {
//		return nil, err
//	}
//
//	var tasks []*domain.Task
//	for _, task := range tasksToConvert {
//		newTask := &domain.Task{
//			ID:        task.Id,
//			Title:     task.Title,
//			Text:      task.Text,
//			Priority:  task.Priority,
//			CreatedAt: task.CreatedAt,
//			Expires:   task.Expires,
//			UserId:    task.UserId,
//		}
//		tasks = append(tasks, newTask)
//	}
//	return tasks, nil
//}

func (u UserServ) Patch(ctx context.Context, user *domain.SignupUser, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	//TODO: Password hashing
	hashed_password := user.Password

	return u.repo.Patch(ctx, user.Username, user.Role, user.Email, hashed_password, idInt)
}

func (u UserServ) Insert(ctx context.Context, user *domain.SignupUser) error {
	//TODO: Password hashing
	hashed_password := user.Password

	return u.repo.Insert(ctx, user.Username, user.Role, user.Email, hashed_password)
}

func (u UserServ) Delete(ctx context.Context, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, idInt)
}
