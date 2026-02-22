package ai

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joshburnsxyz/lark/api"
)

// LogEntry records a single prompt->response exchange.
type LogEntry struct {
	Timestamp    string         `json:"timestamp"`
	ScenarioID   string         `json:"scenario_id"`
	ScenarioName string         `json:"scenario_name"`
	Language     string         `json:"language"`
	Turn         string         `json:"turn"` // "start" or "input"
	SystemPrompt string         `json:"system_prompt"`
	UserPrompt   string         `json:"user_prompt"`
	Response     *api.GameMessage `json:"response"`
	Correction   *api.Correction  `json:"correction,omitempty"`
}

// ScenarioLog is the top-level structure written to each scenario's JSON file.
type ScenarioLog struct {
	ScenarioID string     `json:"scenario_id"`
	Entries    []LogEntry `json:"entries"`
}

// LoggingClient wraps an ai.Client and logs all prompt/response pairs to JSON files.
type LoggingClient struct {
	inner  Client
	dir    string
	mu     sync.Mutex
}

// NewLoggingClient creates a logging wrapper. Files are written to dir/<scenarioID>.json.
func NewLoggingClient(inner Client, dir string) *LoggingClient {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Printf("[logger] warning: cannot create log dir %s: %v", dir, err)
	}
	return &LoggingClient{inner: inner, dir: dir}
}

func (c *LoggingClient) appendEntry(entry LogEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := filepath.Join(c.dir, entry.ScenarioID+".json")

	var logFile ScenarioLog
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &logFile)
	}

	if logFile.ScenarioID == "" {
		logFile.ScenarioID = entry.ScenarioID
	}
	logFile.Entries = append(logFile.Entries, entry)

	data, err := json.MarshalIndent(logFile, "", "  ")
	if err != nil {
		log.Printf("[logger] marshal error: %v", err)
		return
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Printf("[logger] write error for %s: %v", path, err)
	}
}

func (c *LoggingClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string) (*api.GameMessage, any, float64, error) {
	msg, history, cost, err := c.inner.StartScenario(ctx, scenario, lang, explanationLang)
	if err == nil {
		c.appendEntry(LogEntry{
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			ScenarioID:   scenario.ID,
			ScenarioName: scenario.Name,
			Language:     lang.Code,
			Turn:         "start",
			SystemPrompt: SystemPrompt(scenario, lang, explanationLang),
			UserPrompt:   ScenarioSeed(scenario, lang),
			Response:     msg,
		})
	}
	return msg, history, cost, err
}

func (c *LoggingClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, callback StreamCallback) (*api.GameMessage, any, float64, error) {
	msg, history, cost, err := c.inner.StartScenarioStream(ctx, scenario, lang, explanationLang, callback)
	if err == nil {
		c.appendEntry(LogEntry{
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			ScenarioID:   scenario.ID,
			ScenarioName: scenario.Name,
			Language:     lang.Code,
			Turn:         "start",
			SystemPrompt: SystemPrompt(scenario, lang, explanationLang),
			UserPrompt:   ScenarioSeed(scenario, lang),
			Response:     msg,
		})
	}
	return msg, history, cost, err
}

func (c *LoggingClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string) (*api.GameMessage, *api.Correction, any, float64, error) {
	msg, correction, newHistory, cost, err := c.inner.SendInput(ctx, scenario, lang, explanationLang, history, input)
	if err == nil {
		c.appendEntry(LogEntry{
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			ScenarioID:   scenario.ID,
			ScenarioName: scenario.Name,
			Language:     lang.Code,
			Turn:         "input",
			SystemPrompt: SystemPrompt(scenario, lang, explanationLang),
			UserPrompt:   input,
			Response:     msg,
			Correction:   correction,
		})
	}
	return msg, correction, newHistory, cost, err
}

func (c *LoggingClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, float64, error) {
	msg, correction, newHistory, cost, err := c.inner.SendInputStream(ctx, scenario, lang, explanationLang, history, input, callback)
	if err == nil {
		c.appendEntry(LogEntry{
			Timestamp:    time.Now().UTC().Format(time.RFC3339),
			ScenarioID:   scenario.ID,
			ScenarioName: scenario.Name,
			Language:     lang.Code,
			Turn:         "input",
			SystemPrompt: SystemPrompt(scenario, lang, explanationLang),
			UserPrompt:   input,
			Response:     msg,
			Correction:   correction,
		})
	}
	return msg, correction, newHistory, cost, err
}

func (c *LoggingClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, explanationLang string, responseText string) any {
	return c.inner.BuildStartHistory(scenario, lang, explanationLang, responseText)
}
