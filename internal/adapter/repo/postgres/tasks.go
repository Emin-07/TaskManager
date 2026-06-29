package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id       int       `db:"id"`
	Title    string    `db:"title"`
	Text     string    `db:"text"`
	Priority string    `db:"priority"`
	Created  time.Time `db:"created"`
	Expires  time.Time `db:"expires"`
	UserId   int       `db:"user_id"`
}

type TaskModel struct {
	DB *sqlx.DB
}

func (m *TaskModel) Latest(ctx context.Context) ([]*Task, error) {
	tasks := []*Task{}
	err := m.DB.SelectContext(ctx, &tasks, "SELECT * FROM tasks WHERE expires > UTC_TIMESTAMP() ORDER BY id")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *TaskModel) Get(ctx context.Context, id int) (*Task, error) {
	task := Task{}
	err := m.DB.GetContext(ctx, &task, "SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &task, nil
}

func (m *TaskModel) Insert(ctx context.Context, title, text, priority string, expireDays, userId int) (int64, error) {
	query := `INSERT INTO tasks (title, text, priority, expires, user_id) VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL ? DAY), ?)`
	res, err := m.DB.ExecContext(ctx, query, title, text, priority, expireDays, userId)

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

func (m *TaskModel) Patch(ctx context.Context, title, text, priority string, id, expireDays int) error {
	var query strings.Builder
	var args []any
	var isNotFirst bool
	query.WriteString("UPDATE tasks SET ")
	if title != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(`title = ?`)
		args = append(args, title)
	}
	if text != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(`text = ? `)
		args = append(args, text)
	}
	if priority != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(`priority = ? `)
		args = append(args, priority)
	}
	if expireDays != 0 {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(`expires = DATE_ADD(NOW(), INTERVAL ? DAY) `)
		args = append(args, expireDays)
	}
	if len(args) == 0 {
		return ErrNoData
	}
	args = append(args, id)
	query.WriteString("WHERE id = ?")

	_, err := m.DB.ExecContext(ctx, query.String(), args...)

	if err != nil {
		return err
	}

	return nil
}

func (m *TaskModel) Delete(ctx context.Context, id int) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
