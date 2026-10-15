package rest

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered",
					"error", err,
					"method", r.Method,
					"path", r.URL.Path,
					"stack", string(debug.Stack()),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func MapRoutes(userHandler *UserHandler, env string) http.Handler {
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("GET /users", userHandler.GetAll)
	apiMux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	apiMux.HandleFunc("POST /users", userHandler.Create)
	apiMux.HandleFunc("PUT /users/{id}", userHandler.Update)
	apiMux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	mainMux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))

	if env == "developer" || env == "local" {
		slog.Info("Swagger UI is enabled at /swagger/")
		mainMux.Handle("/swagger/", httpSwagger.Handler())
	} else {
		slog.Info("Swagger UI is disabled for production build")
	}

	handler := loggingMiddleware(mainMux)
	handler = recoveryMiddleware(handler)

	return handler
}
