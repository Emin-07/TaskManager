package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/Emin-07/TaskManager/internal/adapter/repo"
	"github.com/Emin-07/TaskManager/internal/core/domain"
)

type UserModel struct {
	DB *sqlx.DB
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*repo.UserDb, error) {
	user := &repo.UserDb{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetById(ctx context.Context, id int) (*repo.UserDb, error) {
	user := &repo.UserDb{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error) {
	query := "SELECT * FROM tasks WHERE user_id = ? ORDER BY id DESC"
	tasks := []*repo.TaskDb{}
	err := m.DB.SelectContext(ctx, &tasks, query, id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *UserModel) Insert(ctx context.Context, username, role, email, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, role,  email, password_hash) VALUES (?, ?, ?, ?)"
	res, err := m.DB.ExecContext(ctx, query, username, role, email, passwordHash)
	if err != nil {
		return 0, err
	}
	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(newId), nil
}

//TODO:need to add patch

func (m *UserModel) Delete(ctx context.Context, id int) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
