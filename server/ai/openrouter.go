package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/joshburnsxyz/lark/api"
)

const (
	openRouterBaseURL   = "https://openrouter.ai/api/v1/chat/completions"
	openRouterTimeout   = 15 * time.Second
	openRouterMaxRetries = 2
	openRouterRetryDelay = 2 * time.Second
)

// OpenRouterClient uses the OpenRouter API (OpenAI-compatible).
type OpenRouterClient struct {
	apiKey string
	model  string
	http   *http.Client
}

// NewOpenRouterClient creates an OpenRouter client.
func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey: apiKey,
		model:  model,
		http:   &http.Client{},
	}
}

// orMessage is an OpenAI-compatible chat message.
type orMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// orRequest is the request body for OpenRouter.
type orRequest struct {
	Model          string            `json:"model"`
	Messages       []orMessage       `json:"messages"`
	ResponseFormat *orResponseFormat `json:"response_format,omitempty"`
	Temperature    float32           `json:"temperature"`
	MaxTokens      int               `json:"max_tokens"`
	Stream         bool              `json:"stream"`
}

type orResponseFormat struct {
	Type string `json:"type"`
}

// orResponse is the response from OpenRouter.
type orResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    any    `json:"code"`
	} `json:"error,omitempty"`
}

func (c *OpenRouterClient) send(ctx context.Context, messages []orMessage) (string, error) {
	reqBody := orRequest{
		Model:    c.model,
		Messages: messages,
		ResponseFormat: &orResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.8,
		MaxTokens:   2048,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= openRouterMaxRetries; attempt++ {
		if attempt > 0 {
			delay := openRouterRetryDelay * time.Duration(attempt)
			log.Printf("[openrouter] retry %d/%d in %s (error: %v)", attempt, openRouterMaxRetries, delay, lastErr)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, openRouterTimeout)
		t0 := time.Now()

		req, err := http.NewRequestWithContext(callCtx, "POST", openRouterBaseURL, bytes.NewReader(body))
		if err != nil {
			cancel()
			return "", fmt.Errorf("creating request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)

		resp, err := c.http.Do(req)
		if err != nil {
			cancel()
			lastErr = err
			log.Printf("[openrouter] request failed in %s: %v", time.Since(t0), err)
			if callCtx.Err() == context.DeadlineExceeded || strings.Contains(err.Error(), "timeout") {
				continue
			}
			return "", fmt.Errorf("sending request: %w", err)
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel() // cancel AFTER body is read
		log.Printf("[openrouter] response %d in %s (%d bytes)", resp.StatusCode, time.Since(t0), len(respBody))

		if resp.StatusCode == 429 || resp.StatusCode == 503 || resp.StatusCode == 502 {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
			continue
		}

		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}

		var orResp orResponse
		if err := json.Unmarshal(respBody, &orResp); err != nil {
			return "", fmt.Errorf("decoding response: %w", err)
		}

		if orResp.Error != nil {
			errMsg := orResp.Error.Message
			if strings.Contains(errMsg, "rate") || strings.Contains(errMsg, "capacity") {
				lastErr = fmt.Errorf("API error: %s", errMsg)
				continue
			}
			return "", fmt.Errorf("API error: %s", errMsg)
		}

		if len(orResp.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		return orResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("failed after retries: %w", lastErr)
}

// orStreamDelta represents a single SSE chunk from OpenRouter's streaming response.
type orStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

const openRouterStreamTimeout = 120 * time.Second

func (c *OpenRouterClient) sendStream(ctx context.Context, messages []orMessage, callback StreamCallback) (string, error) {
	reqBody := orRequest{
		Model:    c.model,
		Messages: messages,
		ResponseFormat: &orResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.8,
		MaxTokens:   2048,
		Stream:      true,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	callCtx, cancel := context.WithTimeout(ctx, openRouterStreamTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(callCtx, "POST", openRouterBaseURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var accumulated strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk orStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 {
			token := chunk.Choices[0].Delta.Content
			if token != "" {
				accumulated.WriteString(token)
				if callback != nil {
					callback(token)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading stream: %w", err)
	}

	return accumulated.String(), nil
}

// BuildStartHistory reconstructs OpenRouter-format conversation history from a cached first-turn response.
func (c *OpenRouterClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, responseText string) any {
	return []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
		{Role: "assistant", Content: responseText},
	}
}

func (c *OpenRouterClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, callback StreamCallback) (*api.GameMessage, any, error) {
	t0 := time.Now()
	log.Printf("[openrouter] StartScenarioStream model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	messages := []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
	}

	text, err := c.sendStream(ctx, messages, callback)
	if err != nil {
		return nil, nil, fmt.Errorf("streaming seed message: %w", err)
	}

	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	history := append(messages, orMessage{Role: "assistant", Content: text})
	log.Printf("[openrouter] StartScenarioStream total: %s", time.Since(t0))
	return msg, history, nil
}

func (c *OpenRouterClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, historyAny any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, error) {
	history, _ := historyAny.([]orMessage)
	t0 := time.Now()

	messages := append(history, orMessage{Role: "user", Content: input})

	text, err := c.sendStream(ctx, messages, callback)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("streaming message: %w", err)
	}

	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := append(messages, orMessage{Role: "assistant", Content: text})
	log.Printf("[openrouter] SendInputStream total: %s", time.Since(t0))
	return msg, correction, newHistory, nil
}

func (c *OpenRouterClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language) (*api.GameMessage, any, error) {
	t0 := time.Now()
	log.Printf("[openrouter] StartScenario model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	messages := []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
	}

	text, err := c.send(ctx, messages)
	if err != nil {
		return nil, nil, fmt.Errorf("sending seed message: %w", err)
	}

	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	// Store full message history for subsequent turns
	history := append(messages, orMessage{Role: "assistant", Content: text})
	log.Printf("[openrouter] StartScenario total: %s", time.Since(t0))
	return msg, history, nil
}

func (c *OpenRouterClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, historyAny any, input string) (*api.GameMessage, *api.Correction, any, error) {
	history, _ := historyAny.([]orMessage)
	t0 := time.Now()

	messages := append(history, orMessage{Role: "user", Content: input})

	text, err := c.send(ctx, messages)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("sending message: %w", err)
	}

	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := append(messages, orMessage{Role: "assistant", Content: text})
	log.Printf("[openrouter] SendInput total: %s", time.Since(t0))
	return msg, correction, newHistory, nil
}
