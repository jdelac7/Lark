package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joshburnsxyz/lark/server/ai"
	"github.com/joshburnsxyz/lark/server/handlers"
	"github.com/joshburnsxyz/lark/server/progress"
	"github.com/joshburnsxyz/lark/server/session"
)

func main() {
	cfg := LoadConfig()

	var client ai.Client
	var backendLabel string

	switch cfg.AIBackend {
	case "openrouter":
		if cfg.OpenRouterAPIKey == "" {
			log.Fatal("OPENROUTER_API_KEY is required when AI_BACKEND=openrouter")
		}
		client = ai.NewOpenRouterClient(cfg.OpenRouterAPIKey, cfg.OpenRouterModel)
		backendLabel = "openrouter/" + cfg.OpenRouterModel

	default: // "google"
		if cfg.GeminiAPIKey == "" {
			log.Fatal("GEMINI_API_KEY is required when AI_BACKEND=google")
		}
		ctx := context.Background()
		gc, err := ai.NewGoogleClient(ctx, cfg.GeminiAPIKey, cfg.GeminiModel)
		if err != nil {
			log.Fatalf("Failed to create Google AI client: %v", err)
		}
		client = gc
		backendLabel = "google/" + cfg.GeminiModel
	}

	sessionStore := session.NewMemoryStore()
	progressStore := progress.NewMemoryStore()

	scenarioHandler := &handlers.ScenarioHandler{
		AI:       client,
		Sessions: sessionStore,
	}
	gameHandler := &handlers.GameHandler{
		AI:       client,
		Sessions: sessionStore,
		Progress: progressStore,
	}
	progressHandler := &handlers.ProgressHandler{
		Progress: progressStore,
	}

	router := NewRouter(scenarioHandler, gameHandler, progressHandler)

	log.Printf("Lark server starting on :%s (backend: %s)", cfg.Port, backendLabel)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
