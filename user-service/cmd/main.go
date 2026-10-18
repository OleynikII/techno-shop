package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"
	"user-service/config"
	_ "user-service/docs"
	"user-service/internal/delivery/rest"
	"user-service/internal/delivery/rest/handler"
	"user-service/internal/repository/psql"
	"user-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// @title user-service
// @version 1.0
// @description Веб апи для работы с пользователями
// @BasePath /api/v1

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", cfg.PG.URL)
	if err != nil {
		slog.Error("failed to open db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		slog.Error("База недоступна", "error", err)
		os.Exit(1)
	}
	// Repositories
	userRepo := psql.NewUserRepository(db)

	// Services
	userService := service.NewUserService(userRepo)

	// Handlers
	baseHandler := handler.NewBaseHandler()
	userHandler := handler.NewUserHandler(baseHandler, userService)

	router := rest.MapRoutes(userHandler, cfg.Swagger.Enabled)
	slog.Info("starting server", "port", cfg.HTTP.Port)
	if err := http.ListenAndServe(":"+cfg.HTTP.Port, router); err != nil {
		slog.Error("critical error: server failed to start", "error", err)
		os.Exit(1)
	}

}
