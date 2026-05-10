package main

import (
	"flag"
	"log"
	"os"

	"github.com/Emin-07/TaskManager/internal/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	tasks        *models.TaskModel
	readableJSON bool
}

func main() {
	readableJSON := flag.Bool("readable", true, "Makes JSON in API's better structured for human to read")
	infoLog := log.New(os.Stdin, "INFO: \t", log.LUTC|log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR: \t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB("todo:mysql@/todoApp?parseTime=true")

	flag.Parse()

	if err != nil {
		errorLog.Fatal(err)
	}

	app := application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		tasks:        &models.TaskModel{DB: db},
		readableJSON: *readableJSON,
	}

	app.router().Run()
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
