package main

import (
	"bufio"
	"os"
	"strings"
)

// Config holds server configuration loaded from environment variables.
type Config struct {
	Port string

	// AI_BACKEND: "google" (default) or "openrouter"
	AIBackend string

	// Google AI Studio
	GeminiAPIKey string
	GeminiModel  string

	// OpenRouter
	OpenRouterAPIKey string
	OpenRouterModel  string
}

// LoadConfig loads .env file (if present) then reads environment variables.
func LoadConfig() Config {
	loadDotEnv(".env")

	c := Config{
		Port:             os.Getenv("PORT"),
		AIBackend:        os.Getenv("AI_BACKEND"),
		GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
		GeminiModel:      os.Getenv("GEMINI_MODEL"),
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterModel:  os.Getenv("OPENROUTER_MODEL"),
	}
	if c.Port == "" {
		c.Port = "8080"
	}
	if c.AIBackend == "" {
		c.AIBackend = "google"
	}
	if c.GeminiModel == "" {
		c.GeminiModel = "gemini-3-flash-preview"
	}
	if c.OpenRouterModel == "" {
		c.OpenRouterModel = "google/gemini-3-flash-preview"
	}
	return c
}

// loadDotEnv reads a .env file and sets environment variables that aren't already set.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // .env file is optional
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val_, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val_ = strings.TrimSpace(val_)
		// Don't override existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, val_)
		}
	}
}
