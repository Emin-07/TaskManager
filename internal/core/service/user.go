package service

import (
	"context"
	"math/rand/v2"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/core/port"
)

const usernameNAmount = 8

type UserServ struct {
	repo port.UserRepo
}

func NewUserService(repo port.UserRepo) UserServ {
	return UserServ{repo: repo}
}

func (u UserServ) Authenticate(ctx context.Context, email string, password string) (*domain.User, error) {
	err := validateCreds(email, password)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.Authenticate(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, CreatedAt: user.CreatedAt}, nil
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
	return &domain.User{ID: user.Id, Username: user.Username, Role: user.Role, Email: user.Email, CreatedAt: user.CreatedAt}, nil
}

func (u UserServ) List(ctx context.Context) ([]*domain.User, error) {
	usersToConvert, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var users []*domain.User
	for _, user := range usersToConvert {
		newTask := &domain.User{
			ID:        user.Id,
			Username:  user.Username,
			Role:      user.Role,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.Patch(ctx, user.Username, user.Role, user.Email, hashedPassword, idInt)
}

func (u UserServ) Insert(ctx context.Context, user *domain.SignupUser) error {
	err := validateCreds(user.Email, user.Password)
	if err != nil {
		return err
	}
	if user.Role == "" {
		user.Role = "user"
	}
	if user.Username == "" {
		name := strings.Builder{}
		name.WriteString("unknown")
		for range usernameNAmount {
			name.WriteString(strconv.Itoa(rand.N(9)))
		}
		user.Username = name.String()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.Insert(ctx, user.Username, user.Role, user.Email, hashedPassword)
}

func (u UserServ) Delete(ctx context.Context, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, idInt)
}
