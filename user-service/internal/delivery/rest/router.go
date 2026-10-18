package rest

import (
	"log/slog"
	"net/http"
	"user-service/internal/delivery/rest/handler"
	"user-service/internal/delivery/rest/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func MapRoutes(userHandler *handler.UserHandler, swaggerEnabled bool) http.Handler {
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("GET /users", userHandler.GetPaginated)
	apiMux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	apiMux.HandleFunc("POST /users", userHandler.Create)
	apiMux.HandleFunc("PUT /users/{id}", userHandler.Update)
	apiMux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	if swaggerEnabled {
		slog.Info("Swagger UI is enabled at /swagger/")
		mainMux.Handle("/swagger/", httpSwagger.Handler())
	} else {
		slog.Info("Swagger UI is disabled")
	}

	handler := middleware.LoggingMiddleware(mainMux)
	handler = middleware.RecoveryMiddleware(handler)

	return handler
}
