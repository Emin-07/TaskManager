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

	"github.com/pressly/goose/v3"

	"github.com/Emin-07/TaskManager/cmd/web/docs"
	"github.com/Emin-07/TaskManager/internal/adapter/handler"
	"github.com/Emin-07/TaskManager/internal/adapter/repo/postgres"
	"github.com/Emin-07/TaskManager/internal/adapter/repo/redis"
	"github.com/Emin-07/TaskManager/internal/app"
	"github.com/Emin-07/TaskManager/internal/core/service"
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

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set goose dialect: %v", err)
	}

	ctxMigration, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := goose.UpContext(ctxMigration, db.DB, "migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations applied successfully!")

	userRepo := postgres.NewUserRepo(db)
	taskRepo := postgres.NewTaskRepo(db)
	redisRepo := redis.NewRedisClientRepo()

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	tokenService := service.NewTokenService(privKeyPath, pubKeyPath)
	redisService := service.NewRateAndCacheService(redisRepo)

	userHandler := handler.NewUserHandler(userService, tokenService, redisService)
	taskHandler := handler.NewTaskHandler(taskService, tokenService, redisService)

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
