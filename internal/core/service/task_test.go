package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	dbmodel "github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
	"github.com/Emin-07/TaskManager/internal/testutil"
)

func newTaskService(repo *testutil.MockTaskRepo) TaskServ {
	return NewTaskService(repo)
}

func TestTaskServ_Get(t *testing.T) {
	t.Run("returns mapped task", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		dbTask := &dbmodel.TaskDb{
			Id:        7,
			Title:     "Write docs",
			Text:      "more docs",
			Priority:  1,
			CreatedAt: time.Now().Add(-time.Hour),
			Expires:   time.Now().Add(time.Hour),
			UserId:    3,
		}
		repo.On("Get", mock.Anything, 7, 3, "user").Return(dbTask, nil)

		task, err := svc.Get(context.Background(), "7", "3", "user")
		require.NoError(t, err)
		assert.Equal(t, 7, task.ID)
		assert.Equal(t, "Write docs", task.Title)
		assert.Equal(t, 3, task.UserId)
	})

	t.Run("returns ErrNoRecord when not found", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		repo.On("Get", mock.Anything, 99, 1, "user").Return(nil, domain.ErrNoRecord)

		_, err := svc.Get(context.Background(), "99", "1", "user")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})

	t.Run("rejects non-numeric id", func(t *testing.T) {
		svc := newTaskService(&testutil.MockTaskRepo{})
		_, err := svc.Get(context.Background(), "abc", "1", "user")
		assert.Error(t, err)
	})

	t.Run("rejects negative id", func(t *testing.T) {
		svc := newTaskService(&testutil.MockTaskRepo{})
		_, err := svc.Get(context.Background(), "-1", "1", "user")
		assert.Error(t, err)
	})
}

func TestTaskServ_List(t *testing.T) {
	t.Run("returns mapped tasks", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		dbTasks := []*dbmodel.TaskDb{
			{Id: 1, Title: "a", UserId: 2},
			{Id: 2, Title: "b", UserId: 2},
		}
		repo.On("List", mock.Anything, 5, 0, 2).Return(dbTasks, nil)

		tasks, err := svc.List(context.Background(), 5, 0, "2")
		require.NoError(t, err)
		require.Len(t, tasks, 2)
		assert.Equal(t, "a", tasks[0].Title)
		assert.Equal(t, "b", tasks[1].Title)
	})

	t.Run("defaults limit to 5 when zero", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		repo.On("List", mock.Anything, 5, 0, 1).Return([]*dbmodel.TaskDb{}, nil)

		tasks, err := svc.List(context.Background(), 0, 0, "1")
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("propagates repo error", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		want := errors.New("db down")
		repo.On("List", mock.Anything, 5, 0, 1).Return(nil, want)

		_, err := svc.List(context.Background(), 5, 0, "1")
		assert.ErrorIs(t, err, want)
	})
}

func TestTaskServ_Post(t *testing.T) {
	t.Run("inserts valid task", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		repo.On("Insert", mock.Anything, "Title", "Text", 2, 1, 5).Return(nil)

		err := svc.Post(context.Background(), &domain.CreateTask{Title: "Title", Text: "Text", Priority: 2, ExpireDays: 1}, "5")
		assert.NoError(t, err)
	})

	t.Run("rejects empty title", func(t *testing.T) {
		svc := newTaskService(&testutil.MockTaskRepo{})
		err := svc.Post(context.Background(), &domain.CreateTask{Title: ""}, "5")
		assert.Error(t, err)
	})

	t.Run("rejects non-positive user id", func(t *testing.T) {
		svc := newTaskService(&testutil.MockTaskRepo{})
		err := svc.Post(context.Background(), &domain.CreateTask{Title: "x"}, "0")
		assert.Error(t, err)
	})
}

func TestTaskServ_Delete(t *testing.T) {
	t.Run("deletes owned task", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		repo.On("Delete", mock.Anything, 3, 1, "user").Return(nil)

		err := svc.Delete(context.Background(), "3", "1", "user")
		assert.NoError(t, err)
	})

	t.Run("propagates ErrNoRecord", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		repo.On("Delete", mock.Anything, 3, 1, "user").Return(domain.ErrNoRecord)

		err := svc.Delete(context.Background(), "3", "1", "user")
		assert.ErrorIs(t, err, domain.ErrNoRecord)
	})
}

func TestTaskServ_Patch(t *testing.T) {
	t.Run("patches valid task", func(t *testing.T) {
		repo := &testutil.MockTaskRepo{}
		svc := newTaskService(repo)
		task := &domain.CreateTask{Title: "New", Text: "Text", Priority: 1, ExpireDays: 2}
		repo.On("Patch", mock.Anything, "New", "Text", "user", 1, 2, 4, 1).Return(nil)

		err := svc.Patch(context.Background(), task, "4", "1", "user")
		assert.NoError(t, err)
	})

	t.Run("rejects negative id", func(t *testing.T) {
		svc := newTaskService(&testutil.MockTaskRepo{})
		err := svc.Patch(context.Background(), &domain.CreateTask{}, "-2", "1", "user")
		assert.Error(t, err)
	})
}
