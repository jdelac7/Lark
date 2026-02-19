package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/joshburnsxyz/lark/api"
)

// Client is an HTTP client for the Lark server API.
type Client struct {
	baseURL  string
	playerID string
	http     *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, playerID string) *Client {
	return &Client{
		baseURL:  baseURL,
		playerID: playerID,
		http:     &http.Client{},
	}
}

func (c *Client) do(method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Player-ID", c.playerID)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp api.ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// GetScenarios fetches available scenarios.
func (c *Client) GetScenarios() ([]api.Scenario, error) {
	var scenarios []api.Scenario
	err := c.do("GET", "/api/v1/scenarios", nil, &scenarios)
	return scenarios, err
}

// GetLanguages fetches supported languages.
func (c *Client) GetLanguages() ([]api.Language, error) {
	var languages []api.Language
	err := c.do("GET", "/api/v1/languages", nil, &languages)
	return languages, err
}

// StartScenario begins a new scenario session.
func (c *Client) StartScenario(scenarioID, language string) (*api.StartResponse, error) {
	return c.StartScenarioCustom(scenarioID, language, "")
}

// StartScenarioCustom begins a scenario session, optionally with a custom prompt.
func (c *Client) StartScenarioCustom(scenarioID, language, customPrompt string) (*api.StartResponse, error) {
	var resp api.StartResponse
	err := c.do("POST", "/api/v1/scenarios/start", api.StartRequest{
		ScenarioID:   scenarioID,
		Language:     language,
		CustomPrompt: customPrompt,
	}, &resp)
	return &resp, err
}

// SendChoice sends a choice selection.
func (c *Client) SendChoice(sessionID string, choiceIndex int) (*api.PlayerInputResponse, error) {
	var resp api.PlayerInputResponse
	err := c.do("POST", "/api/v1/game/input", api.PlayerInputRequest{
		SessionID:   sessionID,
		Mode:        api.InputModeChoice,
		ChoiceIndex: choiceIndex,
	}, &resp)
	return &resp, err
}

// SendFreeText sends free text input.
func (c *Client) SendFreeText(sessionID, text string) (*api.PlayerInputResponse, error) {
	var resp api.PlayerInputResponse
	err := c.do("POST", "/api/v1/game/input", api.PlayerInputRequest{
		SessionID: sessionID,
		Mode:      api.InputModeFreeText,
		Text:      text,
	}, &resp)
	return &resp, err
}

// GetProgress fetches player progress.
func (c *Client) GetProgress() (*api.ProgressResponse, error) {
	var resp api.ProgressResponse
	err := c.do("GET", "/api/v1/progress?playerId="+c.playerID, nil, &resp)
	return &resp, err
}

// sseEvent represents a parsed SSE data line.
type sseEvent struct {
	Token      string          `json:"token,omitempty"`
	Done       bool            `json:"done,omitempty"`
	Error      string          `json:"error,omitempty"`
	SessionID  string          `json:"sessionId,omitempty"`
	Message    *api.GameMessage `json:"message,omitempty"`
	Correction *api.Correction `json:"correction,omitempty"`
}

// doStream performs a streaming HTTP request, calling onToken for each token,
// and unmarshaling the final "done" event into result.
func (c *Client) doStream(method, path string, body any, onToken func(string), result any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Player-ID", c.playerID)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp api.ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var evt sseEvent
		if err := json.Unmarshal([]byte(data), &evt); err != nil {
			continue
		}

		if evt.Error != "" {
			return fmt.Errorf("stream error: %s", evt.Error)
		}

		if evt.Token != "" && onToken != nil {
			onToken(evt.Token)
		}

		if evt.Done {
			// Re-unmarshal the full data into the result type
			if result != nil {
				if err := json.Unmarshal([]byte(data), result); err != nil {
					return fmt.Errorf("decoding final event: %w", err)
				}
			}
			return nil
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading stream: %w", err)
	}

	return fmt.Errorf("stream ended without done event")
}

// streamDoneStartResponse is the final SSE event for start/stream.
type streamDoneStartResponse struct {
	Done      bool            `json:"done"`
	SessionID string          `json:"sessionId"`
	Message   api.GameMessage `json:"message"`
}

// streamDoneInputResponse is the final SSE event for input/stream.
type streamDoneInputResponse struct {
	Done       bool            `json:"done"`
	Message    api.GameMessage `json:"message"`
	Correction *api.Correction `json:"correction,omitempty"`
}

// StreamStartScenario begins a new scenario with SSE streaming.
func (c *Client) StreamStartScenario(scenarioID, language string, onToken func(string)) (*api.StartResponse, error) {
	return c.StreamStartScenarioCustom(scenarioID, language, "", onToken)
}

// StreamStartScenarioCustom begins a scenario with SSE streaming, optionally with a custom prompt.
func (c *Client) StreamStartScenarioCustom(scenarioID, language, customPrompt string, onToken func(string)) (*api.StartResponse, error) {
	var done streamDoneStartResponse
	err := c.doStream("POST", "/api/v1/scenarios/start/stream", api.StartRequest{
		ScenarioID:   scenarioID,
		Language:     language,
		CustomPrompt: customPrompt,
	}, onToken, &done)
	if err != nil {
		return nil, err
	}
	return &api.StartResponse{
		SessionID: done.SessionID,
		Message:   done.Message,
	}, nil
}

// StreamSendChoice sends a choice selection with SSE streaming.
func (c *Client) StreamSendChoice(sessionID string, choiceIndex int, onToken func(string)) (*api.PlayerInputResponse, error) {
	var done streamDoneInputResponse
	err := c.doStream("POST", "/api/v1/game/input/stream", api.PlayerInputRequest{
		SessionID:   sessionID,
		Mode:        api.InputModeChoice,
		ChoiceIndex: choiceIndex,
	}, onToken, &done)
	if err != nil {
		return nil, err
	}
	return &api.PlayerInputResponse{
		Message:    done.Message,
		Correction: done.Correction,
	}, nil
}

// StreamSendFreeText sends free text input with SSE streaming.
func (c *Client) StreamSendFreeText(sessionID, text string, onToken func(string)) (*api.PlayerInputResponse, error) {
	var done streamDoneInputResponse
	err := c.doStream("POST", "/api/v1/game/input/stream", api.PlayerInputRequest{
		SessionID: sessionID,
		Mode:      api.InputModeFreeText,
		Text:      text,
	}, onToken, &done)
	if err != nil {
		return nil, err
	}
	return &api.PlayerInputResponse{
		Message:    done.Message,
		Correction: done.Correction,
	}, nil
}
