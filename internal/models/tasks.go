package models

import (
	"database/sql"
	"errors"
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
}

type TaskModel struct {
	DB *sqlx.DB
}

func (m *TaskModel) Latest() ([]*Task, error) {
	tasks := []*Task{}
	err := m.DB.Select(&tasks, "SELECT * FROM tasks WHERE expires > UTC_TIMESTAMP() ORDER BY id")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *TaskModel) Get(id int) (*Task, error) {
	task := Task{}
	err := m.DB.Get(&task, "SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &task, nil
}

func (m *TaskModel) Insert(title, text, priority string, expireDays int) (int64, error) {
	query := `INSERT INTO tasks (title, text, priority, expires) VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL ? DAY))`
	res, err := m.DB.Exec(query, title, text, priority, expireDays)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *TaskModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
func (m *TaskModel) RefreshTasks(limit, days int) error {
	// _, err := m.DB.Exec("UPDATE tasks SET expires = DATE_ADD(NOW(), INTERVAL ? DAY) WHERE id < ? AND expires < UTC_TIMESTAMP()", days, limit)
	_, err := m.DB.Exec("UPDATE tasks SET expires = DATE_ADD(NOW(), INTERVAL ? DAY) WHERE id < ? ", days, limit)
	if err != nil {
		return err
	}
	return nil
}
