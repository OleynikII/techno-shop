package main

import (
	_ "GoTest/docs"
	"GoTest/internal/delivery/rest"
	"GoTest/internal/repository/memory"
	"GoTest/internal/service"
	"log/slog"
	"net/http"
	"os"
)

// @title           User App API
// @version         1.0
// @description     API Server for Users Application
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "developer"
	}

	// Repositories
	userRepo := memory.NewUserRepository()

	// Services
	userService := service.NewUserService(userRepo)

	// Handlers
	userHandler := rest.NewUserHandler(userService)

	router := rest.MapRoutes(userHandler, env)
	slog.Info("starting server", "port", "8080", "env", "local")
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("critical error: server failed to start", "error", err)
		os.Exit(1)
	}

}
