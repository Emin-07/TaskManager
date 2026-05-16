package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id           int       `db:"id"`
	Username     string    `db:"title"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserModel struct {
	DB *sqlx.DB
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := m.DB.Get(user, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	user := &User{}
	err := m.DB.Get(user, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *UserModel) GetUserTasks(id int) ([]*Task, error) {
	query := "SELECT * FROM tasks WHERE user_id = ? ORDER BY id DESC"
	tasks := []*Task{}
	err := m.DB.Select(&tasks, query, id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *UserModel) Insert(username, email, passwordHash string) (int, error) {
	query := "INSERT INTO users (title, email, password_hash) VALUES (?, ?, ?)"
	res, err := m.DB.Exec(query, username, email, passwordHash)
	if err != nil {
		return 0, err
	}
	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(newId), nil
}

//need to add patch

func (m *UserModel) Delete(id int) error {
	_, err := m.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
