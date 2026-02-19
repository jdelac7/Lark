package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/joshburnsxyz/lark/server/handlers"
)

// NewRouter creates the HTTP router with all routes and middleware.
func NewRouter(scenario *handlers.ScenarioHandler, game *handlers.GameHandler, progress *handlers.ProgressHandler) http.Handler {
	mux := http.NewServeMux()

	// API v1 routes
	mux.HandleFunc("GET /api/v1/health", handlers.Health)
	mux.HandleFunc("GET /api/v1/languages", scenario.Languages)
	mux.HandleFunc("GET /api/v1/scenarios", scenario.List)
	mux.HandleFunc("POST /api/v1/scenarios/start", scenario.Start)
	mux.HandleFunc("POST /api/v1/scenarios/start/stream", scenario.StartStream)
	mux.HandleFunc("POST /api/v1/game/input", game.Input)
	mux.HandleFunc("POST /api/v1/game/input/stream", game.InputStream)
	mux.HandleFunc("GET /api/v1/game/state", game.State)
	mux.HandleFunc("GET /api/v1/progress", progress.Get)

	// Apply middleware
	var handler http.Handler = mux
	handler = corsMiddleware(handler)
	handler = loggingMiddleware(handler)
	handler = recoveryMiddleware(handler)

	return handler
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Player-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v\n%s", err, debug.Stack())
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
