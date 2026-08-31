package testutil

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

// MockTaskRepo is a testify mock of port.TaskRepo for unit tests.
type MockTaskRepo struct{ mock.Mock }

func (m *MockTaskRepo) List(ctx context.Context, limit, offset, userId int) ([]*repo.TaskDb, error) {
	args := m.Called(ctx, limit, offset, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repo.TaskDb), args.Error(1)
}

func (m *MockTaskRepo) Get(ctx context.Context, id, userId int, role string) (*repo.TaskDb, error) {
	args := m.Called(ctx, id, userId, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.TaskDb), args.Error(1)
}

func (m *MockTaskRepo) Insert(ctx context.Context, title, text string, priority, expireDays, userId int) error {
	return m.Called(ctx, title, text, priority, expireDays, userId).Error(0)
}

func (m *MockTaskRepo) Patch(ctx context.Context, title, text, userRole string, priority, expireDays, id, userId int) error {
	return m.Called(ctx, title, text, userRole, priority, expireDays, id, userId).Error(0)
}

func (m *MockTaskRepo) Delete(ctx context.Context, id, userId int, role string) error {
	return m.Called(ctx, id, userId, role).Error(0)
}

// MockUserRepo is a testify mock of port.UserRepo for unit tests.
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Authenticate(ctx context.Context, email string) (*repo.UserDb, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.UserDb), args.Error(1)
}

func (m *MockUserRepo) GetById(ctx context.Context, id int) (*repo.UserDb, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repo.UserDb), args.Error(1)
}

func (m *MockUserRepo) GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repo.TaskDb), args.Error(1)
}

func (m *MockUserRepo) List(ctx context.Context) ([]*repo.UserDb, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repo.UserDb), args.Error(1)
}

func (m *MockUserRepo) Insert(ctx context.Context, username, role, email string, passwordHash []byte) error {
	return m.Called(ctx, username, role, email, passwordHash).Error(0)
}

func (m *MockUserRepo) Patch(ctx context.Context, username, role, email string, passwordHash []byte, id int) error {
	return m.Called(ctx, username, role, email, passwordHash, id).Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

// MockTokenService is a testify mock of port.TokenService for unit tests.
type MockTokenService struct{ mock.Mock }

func (m *MockTokenService) CreateToken(id, role string) (string, error) {
	args := m.Called(id, role)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ParseFromRequest(r *http.Request) (map[string]string, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

// MockRateAndCacheService is a testify mock of port.RateAndCacheService.
type MockRateAndCacheService struct{ mock.Mock }

func (m *MockRateAndCacheService) Increment(ctx context.Context, id string) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockRateAndCacheService) Decrement(ctx context.Context, id string) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockRateAndCacheService) Set(ctx context.Context, name, id, userId string, val any, duration time.Duration) error {
	return m.Called(ctx, name, id, userId, val, duration).Error(0)
}

func (m *MockRateAndCacheService) Get(ctx context.Context, name, id, userId string) (string, error) {
	args := m.Called(ctx, name, id, userId)
	return args.String(0), args.Error(1)
}

func (m *MockRateAndCacheService) Del(ctx context.Context, name, id, userId string) error {
	return m.Called(ctx, name, id, userId).Error(0)
}

// MockBroker is a testify mock of port.MessageBrokerOut.
type MockBroker struct{ mock.Mock }

func (m *MockBroker) Publish(data map[string]string, topic string) error {
	return m.Called(data, topic).Error(0)
}

func (m *MockBroker) Consume(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// MockTaskService is a testify mock of port.TaskService.
type MockTaskService struct{ mock.Mock }

func (m *MockTaskService) Get(ctx context.Context, id, userIdStr, role string) (*domain.Task, error) {
	args := m.Called(ctx, id, userIdStr, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *MockTaskService) List(ctx context.Context, limit, offset int, userIdStr string) ([]*domain.Task, error) {
	args := m.Called(ctx, limit, offset, userIdStr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *MockTaskService) Post(ctx context.Context, task *domain.CreateTask, userIdStr string) error {
	return m.Called(ctx, task, userIdStr).Error(0)
}

func (m *MockTaskService) Delete(ctx context.Context, id, userIdStr, role string) error {
	return m.Called(ctx, id, userIdStr, role).Error(0)
}

func (m *MockTaskService) Patch(ctx context.Context, task *domain.CreateTask, id, userIdStr, role string) error {
	return m.Called(ctx, task, id, userIdStr, role).Error(0)
}

// MockUserService is a testify mock of port.UserService.
type MockUserService struct{ mock.Mock }

func (m *MockUserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetById(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserService) Patch(ctx context.Context, user *domain.SignupUser, id string) error {
	return m.Called(ctx, user, id).Error(0)
}

func (m *MockUserService) Insert(ctx context.Context, user *domain.SignupUser) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// Well-known sample values reused across tests.
var (
	SampleTaskDb = &repo.TaskDb{
		Id:       1,
		Title:    "Buy milk",
		Text:     "whole milk",
		Priority: 2,
		UserId:   1,
	}
	SampleUserDb = &repo.UserDb{
		Id:       1,
		Username: "alice",
		Role:     "admin",
		Email:    "alice@example.com",
	}
	SampleDomainTask = &domain.Task{
		ID:       1,
		Title:    "Buy milk",
		Text:     "whole milk",
		Priority: 2,
		UserId:   1,
	}
	SampleDomainUser = &domain.User{
		ID:       1,
		Username: "alice",
		Role:     "admin",
		Email:    "alice@example.com",
	}
)
