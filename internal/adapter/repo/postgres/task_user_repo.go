package postgres

import "github.com/jmoiron/sqlx"

type TaskRepo struct {
	DB *sqlx.DB
}

func newTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{DB: db}
}

type UserRepo struct {
	DB *sqlx.DB
}

func newUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{DB: db}
}
