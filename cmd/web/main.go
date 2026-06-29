package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/Emin-07/TaskManager/internal/adapter/repo/postgres"
)

type application struct {
	logger       *zap.Logger
	tasks        *postgres.TaskModel
	readableJSON bool
}

func main() {

}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// TODO: add env files to handle dsn
