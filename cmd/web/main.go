package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/Emin-07/TaskManager/internal/adapter/handler"
	"github.com/Emin-07/TaskManager/internal/adapter/repo/postgres"
	"github.com/Emin-07/TaskManager/internal/app"
	"github.com/Emin-07/TaskManager/internal/core/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

func main() {
	dsn := fmt.Sprintf("postgres://username:password@host:port/dbname?sslmode=disable")
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepo(db)
	taskRepo := postgres.NewTaskRepo(db)

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)

	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)
	application := app.NewApp(app.WithTaskHandler(taskHandler), app.WithUserHandler(userHandler))
	srv := application.NewServer()

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}

	log.Println("Server exiting")
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
