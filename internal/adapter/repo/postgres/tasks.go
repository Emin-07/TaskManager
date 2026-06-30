package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

func (m *TaskRepo) List(ctx context.Context, limit, offset int) ([]*repo.TaskDb, error) {
	tasks := []*repo.TaskDb{}
	query := `SELECT * FROM tasks WHERE expires > CURRENT_TIMESTAMP  ORDER BY id LIMIT $1 OFFSET $2`
	err := m.DB.SelectContext(ctx, &tasks, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *TaskRepo) Get(ctx context.Context, id int) (*repo.TaskDb, error) {
	task := repo.TaskDb{}
	err := m.DB.GetContext(ctx, &task, "SELECT * FROM tasks WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return &task, nil
}

func (m *TaskRepo) Insert(ctx context.Context, title, text string, priority, expireDays, userId int) (int64, error) {
	query := `INSERT INTO tasks (title, text, priority, expires, user_id) VALUES ($1, $2, $3, CURRENT_TIMESTAMP + MAKE_INTERVAL(days => $4), $5)`
	res, err := m.DB.ExecContext(ctx, query, title, text, expireDays, priority, userId)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func queryOrderTracker(query *strings.Builder, isNotFirst *bool) {
	if *isNotFirst {
		query.WriteString(`, `)
	} else {
		*isNotFirst = true
	}
}

func (m *TaskRepo) Patch(ctx context.Context, title, text string, priority, expireDays, id int) error {
	var query strings.Builder
	var args []any
	var isNotFirst bool
	cnt := 1
	query.WriteString("UPDATE tasks SET ")
	if title != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`title = $%d`, cnt))
		cnt++
		args = append(args, title)
	}
	if text != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`text = $%d `, cnt))
		cnt++
		args = append(args, text)
	}
	if priority != 0 {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`priority = $%d `, cnt))
		cnt++
		args = append(args, priority)
	}
	if expireDays != 0 {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`expires = CURRENT_TIMESTAMP + MAKE_INTERVAL(days => $%d) `, cnt))
		cnt++
		args = append(args, expireDays)
	}
	if len(args) == 0 {
		return domain.ErrNoData
	}
	args = append(args, id)
	query.WriteString(fmt.Sprintf("WHERE id = $%d", cnt))

	_, err := m.DB.ExecContext(ctx, query.String(), args...)

	if err != nil {
		return err
	}

	return nil
}

func (m *TaskRepo) Delete(ctx context.Context, id int) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
