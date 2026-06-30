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

func (m *UserRepo) GetByEmail(ctx context.Context, email string) (*repo.UserDb, error) {
	user := &repo.UserDb{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (m *UserRepo) GetById(ctx context.Context, id int) (*repo.UserDb, error) {
	user := &repo.UserDb{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNoRecord
		}
		return nil, err
	}
	return user, nil
}

func (m *UserRepo) List(ctx context.Context) ([]*repo.UserDb, error) {
	query := "SELECT * FROM users ORDER BY id DESC"
	users := []*repo.UserDb{}
	err := m.DB.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (m *UserRepo) GetUserTasks(ctx context.Context, id int) ([]*repo.TaskDb, error) {
	query := "SELECT * FROM tasks WHERE user_id = $1 ORDER BY id DESC"
	tasks := []*repo.TaskDb{}
	err := m.DB.SelectContext(ctx, &tasks, query, id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *UserRepo) Insert(ctx context.Context, username, role, email, passwordHash string) (int, error) {
	query := "INSERT INTO users (username, role,  email, password_hash) VALUES ($1, $2, $3, $4)"
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

func (m *UserRepo) Patch(ctx context.Context, username, role, email, passwordHash string, id int) error {
	var query strings.Builder
	var args []any
	var isNotFirst bool
	cnt := 1
	query.WriteString("UPDATE users SET ")
	if username != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`username = $%d `, cnt))
		cnt++
		args = append(args, username)
	}
	if role != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`role = $%d `, cnt))
		cnt++
		args = append(args, role)
	}
	if email != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`email = $%d `, cnt))
		cnt++
		args = append(args, email)
	}
	if passwordHash != "" {
		queryOrderTracker(&query, &isNotFirst)
		query.WriteString(fmt.Sprintf(`password_hash = $%d `, cnt))
		cnt++
		args = append(args, passwordHash)
	}
	if len(args) == 0 {
		return domain.ErrNoData
	}
	args = append(args, id)
	query.WriteString(fmt.Sprintf("WHERE id = $%d"))

	_, err := m.DB.ExecContext(ctx, query.String(), args...)

	if err != nil {
		return err
	}

	return nil
}

func (m *UserRepo) Delete(ctx context.Context, id int) error {
	_, err := m.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
