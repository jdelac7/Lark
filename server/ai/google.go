package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/genai"

	"github.com/joshburnsxyz/lark/api"
)

const (
	googleCallTimeout  = 15 * time.Second
	googleMaxRetries   = 2
	googleRetryDelay   = 2 * time.Second
)

// GoogleClient uses the Google AI Studio (Gemini) SDK.
type GoogleClient struct {
	client *genai.Client
	model  string
}

// NewGoogleClient creates a Google AI Studio client.
func NewGoogleClient(ctx context.Context, apiKey, model string) (*GoogleClient, error) {
	c, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("creating genai client: %w", err)
	}
	return &GoogleClient{client: c, model: model}, nil
}

func (c *GoogleClient) generateConfig(scenario *api.Scenario, lang *api.Language) *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: SystemPrompt(scenario, lang)}},
		},
		ResponseMIMEType: "application/json",
		ResponseSchema:   googleTurnSchema(),
		Temperature:      genai.Ptr[float32](0.8),
		MaxOutputTokens:  2048,
	}
}

func googleIsRetryable(err error, callCtx context.Context) bool {
	if callCtx.Err() == context.DeadlineExceeded {
		return true
	}
	s := err.Error()
	return strings.Contains(s, "503") || strings.Contains(s, "UNAVAILABLE") ||
		strings.Contains(s, "429") || strings.Contains(s, "RESOURCE_EXHAUSTED")
}

func (c *GoogleClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, callback StreamCallback) (*api.GameMessage, any, error) {
	t0 := time.Now()
	config := c.generateConfig(scenario, lang)
	seed := ScenarioSeed(scenario, lang)
	log.Printf("[google] StartScenarioStream model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	chat, err := c.client.Chats.Create(ctx, c.model, config, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating chat: %w", err)
	}

	var accumulated strings.Builder
	for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: seed}) {
		if err != nil {
			return nil, nil, fmt.Errorf("streaming seed message: %w", err)
		}
		token := result.Text()
		if token != "" {
			accumulated.WriteString(token)
			if callback != nil {
				callback(token)
			}
		}
	}

	text := accumulated.String()
	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	history := chat.History(true)
	log.Printf("[google] StartScenarioStream total: %s", time.Since(t0))
	return msg, history, nil
}

func (c *GoogleClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, historyAny any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, error) {
	history, _ := historyAny.([]*genai.Content)
	t0 := time.Now()
	config := c.generateConfig(scenario, lang)

	chat, err := c.client.Chats.Create(ctx, c.model, config, history)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating chat: %w", err)
	}

	var accumulated strings.Builder
	for result, err := range chat.SendMessageStream(ctx, genai.Part{Text: input}) {
		if err != nil {
			return nil, nil, nil, fmt.Errorf("streaming message: %w", err)
		}
		token := result.Text()
		if token != "" {
			accumulated.WriteString(token)
			if callback != nil {
				callback(token)
			}
		}
	}

	text := accumulated.String()
	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := chat.History(true)
	log.Printf("[google] SendInputStream total: %s", time.Since(t0))
	return msg, correction, newHistory, nil
}

func (c *GoogleClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language) (*api.GameMessage, any, error) {
	t0 := time.Now()
	config := c.generateConfig(scenario, lang)
	seed := ScenarioSeed(scenario, lang)
	log.Printf("[google] StartScenario model=%s scenario=%s lang=%s", c.model, scenario.ID, lang.Code)

	var result *genai.GenerateContentResponse
	var chat *genai.Chat
	var lastErr error

	for attempt := 0; attempt <= googleMaxRetries; attempt++ {
		if attempt > 0 {
			delay := googleRetryDelay * time.Duration(attempt)
			log.Printf("[google] retry %d/%d in %s (error: %v)", attempt, googleMaxRetries, delay, lastErr)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			}
		}

		var err error
		chat, err = c.client.Chats.Create(ctx, c.model, config, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("creating chat: %w", err)
		}

		callCtx, cancel := context.WithTimeout(ctx, googleCallTimeout)
		t1 := time.Now()
		result, err = chat.SendMessage(callCtx, genai.Part{Text: seed})
		cancel()

		if err == nil {
			log.Printf("[google] SendMessage completed in %s", time.Since(t1))
			break
		}
		lastErr = err
		log.Printf("[google] SendMessage failed in %s: %v", time.Since(t1), err)
		if googleIsRetryable(err, callCtx) {
			result = nil
			continue
		}
		return nil, nil, fmt.Errorf("sending seed message: %w", err)
	}

	if result == nil {
		return nil, nil, fmt.Errorf("sending seed message after retries: %w", lastErr)
	}

	text := result.Text()
	msg, _, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	history := chat.History(true)
	log.Printf("[google] StartScenario total: %s", time.Since(t0))
	return msg, history, nil
}

func (c *GoogleClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, historyAny any, input string) (*api.GameMessage, *api.Correction, any, error) {
	history, _ := historyAny.([]*genai.Content)
	t0 := time.Now()
	config := c.generateConfig(scenario, lang)

	var result *genai.GenerateContentResponse
	var chat *genai.Chat
	var lastErr error

	for attempt := 0; attempt <= googleMaxRetries; attempt++ {
		if attempt > 0 {
			delay := googleRetryDelay * time.Duration(attempt)
			log.Printf("[google] SendInput retry %d/%d in %s (error: %v)", attempt, googleMaxRetries, delay, lastErr)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, nil, nil, ctx.Err()
			}
		}

		var err error
		chat, err = c.client.Chats.Create(ctx, c.model, config, history)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("creating chat: %w", err)
		}

		callCtx, cancel := context.WithTimeout(ctx, googleCallTimeout)
		t1 := time.Now()
		result, err = chat.SendMessage(callCtx, genai.Part{Text: input})
		cancel()

		if err == nil {
			log.Printf("[google] SendInput completed in %s", time.Since(t1))
			break
		}
		lastErr = err
		log.Printf("[google] SendInput failed in %s: %v", time.Since(t1), err)
		if googleIsRetryable(err, callCtx) {
			result = nil
			continue
		}
		return nil, nil, nil, fmt.Errorf("sending message: %w", err)
	}

	if result == nil {
		return nil, nil, nil, fmt.Errorf("sending message after retries: %w", lastErr)
	}

	text := result.Text()
	msg, correction, err := ParseTurnJSON(text)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing response: %w", err)
	}

	newHistory := chat.History(true)
	log.Printf("[google] SendInput total: %s", time.Since(t0))
	return msg, correction, newHistory, nil
}

// BuildStartHistory reconstructs Google-format conversation history from a cached first-turn response.
func (c *GoogleClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, responseText string) any {
	seed := ScenarioSeed(scenario, lang)
	return []*genai.Content{
		{Role: "user", Parts: []*genai.Part{{Text: seed}}},
		{Role: "model", Parts: []*genai.Part{{Text: responseText}}},
	}
}

func googleTurnSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"narrative":            {Type: genai.TypeString, Description: "Scene narration in the target language"},
			"translation":         {Type: genai.TypeString, Description: "English translation of the narrative"},
			"npcDialog":           {Type: genai.TypeString, Description: "NPC dialog in target language (empty if none)"},
			"npcDialogTranslation": {Type: genai.TypeString, Description: "English translation of NPC dialog (empty if none)"},
			"choices": {
				Type:        genai.TypeArray,
				Description: "Player choices (2-4 options)",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"text":        {Type: genai.TypeString, Description: "Choice text in target language"},
						"translation": {Type: genai.TypeString, Description: "English translation"},
					},
					Required: []string{"text", "translation"},
				},
			},
			"vocabulary": {
				Type:        genai.TypeArray,
				Description: "Key vocabulary from this turn",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"word":        {Type: genai.TypeString, Description: "Word or phrase in target language"},
						"translation": {Type: genai.TypeString, Description: "English translation"},
						"usage":       {Type: genai.TypeString, Description: "Brief usage note"},
					},
					Required: []string{"word", "translation"},
				},
			},
			"correction": {
				Type:        genai.TypeObject,
				Description: "Grammar correction if player used free text (null if not applicable)",
				Nullable:    genai.Ptr(true),
				Properties: map[string]*genai.Schema{
					"original":    {Type: genai.TypeString, Description: "Player's original text"},
					"corrected":   {Type: genai.TypeString, Description: "Corrected version"},
					"explanation": {Type: genai.TypeString, Description: "What was wrong and why"},
				},
				Required: []string{"original", "corrected", "explanation"},
			},
			"finished": {Type: genai.TypeBoolean, Description: "True if the scenario has reached its natural conclusion"},
		},
		Required: []string{"narrative", "translation", "choices", "vocabulary", "finished"},
	}
}
