package postgres

import (
	"github.com/jmoiron/sqlx"
)

type TaskRepo struct {
	DB *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{DB: db}
}

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}
