package main

import (
	"flag"

	"github.com/Emin-07/TaskManager/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type application struct {
	logger       *zap.Logger
	tasks        *models.TaskModel
	readableJSON bool
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	readableJSON := flag.Bool("readable", true, "Makes JSON in API's better structured for human to read")
	port := flag.String("addr", ":8080", "Specify address of the api, like :8080")
	db, err := openDB("todo:mysql@/todoApp?parseTime=true")

	flag.Parse()

	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
	}

	app := application{
		logger:       logger,
		tasks:        &models.TaskModel{DB: db},
		readableJSON: *readableJSON,
	}

	app.router().Run(*port)
}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// / home
// view users account
// view, delete, update, post your tasks
// log in log out sign in

// id, text, created, expires 1day
// TODO: add env files to handle dsn
