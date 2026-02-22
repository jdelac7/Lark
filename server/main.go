package main

import (
	"context"
	"log"
	"net/http"

	"github.com/joshburnsxyz/lark/server/ai"
	"github.com/joshburnsxyz/lark/server/cost"
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

	// Wrap AI client with response cache for premade scenarios
	client = ai.NewCachedClient(client)

	// Wrap with prompt logger for testing/prompt improvement
	client = ai.NewLoggingClient(client, "prompt_logs")

	sessionStore := session.NewMemoryStore()
	progressStore := progress.NewMemoryStore()
	costStore, err := cost.NewSQLiteStore("cost.db")
	if err != nil {
		log.Fatalf("Failed to open cost database: %v", err)
	}
	defer costStore.Close()

	eventRecorder, err := cost.NewEventRecorder(costStore.DB())
	if err != nil {
		log.Fatalf("Failed to create event recorder: %v", err)
	}

	scenarioHandler := &handlers.ScenarioHandler{
		AI:       client,
		Sessions: sessionStore,
		Cost:     costStore,
		Events:   eventRecorder,
	}
	gameHandler := &handlers.GameHandler{
		AI:       client,
		Sessions: sessionStore,
		Progress: progressStore,
		Cost:     costStore,
		Events:   eventRecorder,
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
