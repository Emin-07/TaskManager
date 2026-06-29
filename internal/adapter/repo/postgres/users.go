package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id           int       `db:"id"`
	Username     string    `db:"title"`
	Role         string    `db:"role"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserModel struct {
	DB *sqlx.DB
}

func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetById(ctx context.Context, id int) (*User, error) {
	user := &User{}
	err := m.DB.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetUserTasks(ctx context.Context, id int) ([]*Task, error) {
	query := "SELECT * FROM tasks WHERE user_id = ? ORDER BY id DESC"
	tasks := []*Task{}
	err := m.DB.SelectContext(ctx, &tasks, query, id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *UserModel) Insert(ctx context.Context, username, email, passwordHash string) (int, error) {
	query := "INSERT INTO users (title, email, password_hash) VALUES (?, ?, ?)"
	res, err := m.DB.ExecContext(ctx, query, username, email, passwordHash)
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
