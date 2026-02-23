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
	openRouterTimeout   = 30 * time.Second
	openRouterMaxRetries = 2
	openRouterRetryDelay = 2 * time.Second
)

// OpenRouterClient uses the OpenRouter API (OpenAI-compatible).
type OpenRouterClient struct {
	apiKey    string
	model     string
	http      *http.Client
	reasoning *orReasoning
}

// NewOpenRouterClient creates an OpenRouter client.
// Reasoning is disabled by default for models that support it (e.g. Grok).
func NewOpenRouterClient(apiKey, model string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:    apiKey,
		model:     model,
		http:      &http.Client{},
		reasoning: &orReasoning{Enabled: false},
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
	Reasoning      *orReasoning      `json:"reasoning,omitempty"`
}

type orReasoning struct {
	Enabled bool `json:"enabled"`
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
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		Cost             float64 `json:"cost"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Code    any    `json:"code"`
	} `json:"error,omitempty"`
}

func (c *OpenRouterClient) send(ctx context.Context, messages []orMessage) (string, float64, error) {
	reqBody := orRequest{
		Model:    c.model,
		Messages: messages,
		ResponseFormat: &orResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.8,
		MaxTokens:   8192,
		Reasoning:   c.reasoning,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("marshaling request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= openRouterMaxRetries; attempt++ {
		if attempt > 0 {
			delay := openRouterRetryDelay * time.Duration(attempt)
			log.Printf("[openrouter] retry %d/%d in %s (error: %v)", attempt, openRouterMaxRetries, delay, lastErr)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", 0, ctx.Err()
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, openRouterTimeout)
		t0 := time.Now()

		req, err := http.NewRequestWithContext(callCtx, "POST", openRouterBaseURL, bytes.NewReader(body))
		if err != nil {
			cancel()
			return "", 0, fmt.Errorf("creating request: %w", err)
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
			return "", 0, fmt.Errorf("sending request: %w", err)
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
			return "", 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}

		var orResp orResponse
		if err := json.Unmarshal(respBody, &orResp); err != nil {
			return "", 0, fmt.Errorf("decoding response: %w", err)
		}

		if orResp.Error != nil {
			errMsg := orResp.Error.Message
			if strings.Contains(errMsg, "rate") || strings.Contains(errMsg, "capacity") {
				lastErr = fmt.Errorf("API error: %s", errMsg)
				continue
			}
			return "", 0, fmt.Errorf("API error: %s", errMsg)
		}

		if len(orResp.Choices) == 0 {
			return "", 0, fmt.Errorf("no choices in response")
		}

		content := orResp.Choices[0].Message.Content
		finishReason := orResp.Choices[0].FinishReason

		// Detect output truncation: finish_reason=length or invalid JSON content
		truncated := finishReason == "length"
		if !truncated && content != "" {
			// Check if the content is valid JSON — truncated responses produce invalid JSON
			if !json.Valid([]byte(content)) {
				truncated = true
			}
		}
		if truncated {
			lastErr = fmt.Errorf("response truncated (finish_reason=%s)", finishReason)
			log.Printf("[openrouter] response truncated (finish_reason=%s, %d bytes content), retrying", finishReason, len(content))
			continue
		}

		var cost float64
		if orResp.Usage != nil {
			cost = orResp.Usage.Cost
		}

		return content, cost, nil
	}

	return "", 0, fmt.Errorf("failed after retries: %w", lastErr)
}

// orStreamDelta represents a single SSE chunk from OpenRouter's streaming response.
type orStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int     `json:"prompt_tokens"`
		CompletionTokens int     `json:"completion_tokens"`
		Cost             float64 `json:"cost"`
	} `json:"usage,omitempty"`
}

const openRouterStreamTimeout = 120 * time.Second

func (c *OpenRouterClient) sendStream(ctx context.Context, messages []orMessage, callback StreamCallback) (string, float64, error) {
	reqBody := orRequest{
		Model:    c.model,
		Messages: messages,
		ResponseFormat: &orResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.8,
		MaxTokens:   8192,
		Stream:      true,
		Reasoning:   c.reasoning,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("marshaling request: %w", err)
	}

	callCtx, cancel := context.WithTimeout(ctx, openRouterStreamTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(callCtx, "POST", openRouterBaseURL, bytes.NewReader(body))
	if err != nil {
		return "", 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var accumulated strings.Builder
	var cost float64
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
		if chunk.Usage != nil {
			cost = chunk.Usage.Cost
		}
	}
	if err := scanner.Err(); err != nil {
		return "", 0, fmt.Errorf("reading stream: %w", err)
	}

	return accumulated.String(), cost, nil
}

// BuildStartHistory reconstructs OpenRouter-format conversation history from a cached first-turn response.
func (c *OpenRouterClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, explanationLang string, responseText string) any {
	return []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang, explanationLang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
		{Role: "assistant", Content: TrimHistoryJSON(responseText)},
	}
}

func (c *OpenRouterClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, callback StreamCallback) (*api.GameMessage, any, float64, error) {
	t0 := time.Now()
	log.Printf("[openrouter] StartScenarioStream model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	messages := []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang, explanationLang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
	}

	text, cost, err := c.sendStream(ctx, messages, callback)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("streaming seed message: %w", err)
	}

	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("parsing response: %w", err)
	}

	history := append(messages, orMessage{Role: "assistant", Content: TrimHistoryJSON(text)})
	log.Printf("[openrouter] StartScenarioStream total: %s", time.Since(t0))
	return msg, history, cost, nil
}

func (c *OpenRouterClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, historyAny any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, float64, error) {
	history, _ := historyAny.([]orMessage)
	t0 := time.Now()

	messages := append(history, orMessage{Role: "user", Content: input})

	text, cost, err := c.sendStream(ctx, messages, callback)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("streaming message: %w", err)
	}

	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := append(messages, orMessage{Role: "assistant", Content: TrimHistoryJSON(text)})
	log.Printf("[openrouter] SendInputStream total: %s", time.Since(t0))
	return msg, correction, newHistory, cost, nil
}

func (c *OpenRouterClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string) (*api.GameMessage, any, float64, error) {
	t0 := time.Now()
	log.Printf("[openrouter] StartScenario model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	messages := []orMessage{
		{Role: "system", Content: SystemPrompt(scenario, lang, explanationLang)},
		{Role: "user", Content: ScenarioSeed(scenario, lang)},
	}

	text, cost, err := c.send(ctx, messages)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("sending seed message: %w", err)
	}

	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("parsing response: %w", err)
	}

	// Store trimmed history for subsequent turns (saves tokens)
	history := append(messages, orMessage{Role: "assistant", Content: TrimHistoryJSON(text)})
	log.Printf("[openrouter] StartScenario total: %s", time.Since(t0))
	return msg, history, cost, nil
}

func (c *OpenRouterClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, historyAny any, input string) (*api.GameMessage, *api.Correction, any, float64, error) {
	history, _ := historyAny.([]orMessage)
	t0 := time.Now()

	messages := append(history, orMessage{Role: "user", Content: input})

	text, cost, err := c.send(ctx, messages)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("sending message: %w", err)
	}

	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, 0, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := append(messages, orMessage{Role: "assistant", Content: TrimHistoryJSON(text)})
	log.Printf("[openrouter] SendInput total: %s", time.Since(t0))
	return msg, correction, newHistory, cost, nil
}
