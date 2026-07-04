package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Emin-07/TaskManager/internal/adapter/handler"
	"github.com/Emin-07/TaskManager/internal/adapter/repo/postgres"
	"github.com/Emin-07/TaskManager/internal/app"
	"github.com/Emin-07/TaskManager/internal/core/service"

	docs "github.com/Emin-07/TaskManager/cmd/web/docs"
)

const (
	privKeyPath = "certs/private.pem"
	pubKeyPath  = "certs/public.pem"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.Title = "TaskManager"
	docs.SwaggerInfo.Description = "Task managing api with users and tasks"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	cfg := app.NewConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepo(db)
	taskRepo := postgres.NewTaskRepo(db)

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	tokenService := service.NewTokenService(privKeyPath, pubKeyPath)

	userHandler := handler.NewUserHandler(userService, tokenService)
	taskHandler := handler.NewTaskHandler(taskService, tokenService)

	application := app.NewApp(app.WithTaskHandler(taskHandler), app.WithUserHandler(userHandler))
	srv := application.NewServer()

	go func() {
		log.Printf("Staring Server at http://localhost%s \n...", application.Addr)
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
