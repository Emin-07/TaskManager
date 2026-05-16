package models

import (
	"database/sql"
	"errors"
	"fmt"
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
		}
		return nil, err
	}
	return &task, nil
}

func (m *TaskModel) Insert(title, text, priority string, expireDays, userId int) (int64, error) {
	query := `INSERT INTO tasks (title, text, priority, expires, user_id) VALUES (?, ?, ?, DATE_ADD(NOW(), INTERVAL ? DAY), ?)`
	res, err := m.DB.Exec(query, title, text, priority, expireDays, userId)

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

func funcName(title, title_sql string, isNotFirst bool, query strings.Builder, args []any) (bool, []any) {
	if title != "" {
		isNotFirst = true
		query.WriteString(fmt.Sprintf("%s = ?", title_sql))
		args = append(args, title)
	}
	return isNotFirst, args
}

func (m *TaskModel) Patch(title, text, priority string, id, userId, expireDays int) error {
	var query strings.Builder
	var args []any
	var isNotFirst bool
	query.WriteString("UPDATE tasks SET ")
	isNotFirst, args = funcName(title, "title", isNotFirst, query, args)
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
	if userId != 0 {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(`user_id = ? `)
		args = append(args, userId)
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

	_, err := m.DB.Exec(query.String(), args...)

	if err != nil {
		return err
	}

	return nil
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
