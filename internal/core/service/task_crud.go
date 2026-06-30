package service

import (
	"context"
	"strconv"

	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (t *TaskServ) Get(ctx context.Context, id string) (*domain.Task, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	task, err := t.repo.Get(ctx, idInt)
	if err != nil {
		return nil, err
	}
	return &domain.Task{
		ID:        task.Id,
		Title:     task.Title,
		Text:      task.Text,
		Priority:  task.Priority,
		CreatedAt: task.CreatedAt,
		Expires:   task.Expires,
		UserId:    task.UserId,
	}, nil
}
func (t *TaskServ) List(ctx context.Context, limit, offset string) ([]*domain.Task, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}

	tasksToConvert, err := t.repo.List(ctx, limitInt, offsetInt)
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
func (t *TaskServ) Post(ctx context.Context, task *domain.CreateTask, userId int) (int64, error) {
	return t.repo.Insert(ctx, task.Title, task.Text, task.Priority, task.ExpireDays, userId)
}

func (t *TaskServ) Delete(ctx context.Context, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return t.repo.Delete(ctx, idInt)
}

func (t *TaskServ) Patch(ctx context.Context, task *domain.CreateTask, id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return t.repo.Patch(ctx, task.Title, task.Text, task.Priority, task.ExpireDays, idInt)
}
