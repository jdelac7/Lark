package ai

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/joshburnsxyz/lark/api"
)

const defaultCacheTTL = 24 * time.Hour

// cachedStart holds a cached first-turn AI response.
type cachedStart struct {
	ResponseText string          // re-serialized JSON for history reconstruction
	Message      api.GameMessage // parsed response (stored as value for safe copying)
	CreatedAt    time.Time
}

// responseCache is a concurrent-safe in-memory cache with TTL expiration.
type responseCache struct {
	mu      sync.RWMutex
	entries map[string]*cachedStart
	ttl     time.Duration
}

func newResponseCache(ttl time.Duration) *responseCache {
	c := &responseCache{
		entries: make(map[string]*cachedStart),
		ttl:     ttl,
	}
	go c.reaper()
	return c
}

func (c *responseCache) get(key string) *cachedStart {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[key]
	if !ok || time.Since(e.CreatedAt) > c.ttl {
		return nil
	}
	return e
}

func (c *responseCache) set(key string, e *cachedStart) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = e
}

func (c *responseCache) reaper() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, e := range c.entries {
			if now.Sub(e.CreatedAt) > c.ttl {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}

func startKey(scenarioID, langCode, explanationLang string) string {
	return scenarioID + ":" + langCode + ":" + explanationLang
}

// isCacheable returns true if the scenario is premade (not a custom user scenario).
func isCacheable(scenarioID string) bool {
	return !strings.HasPrefix(scenarioID, "custom_")
}

// messageToJSON re-serializes a GameMessage into the TurnResponse JSON format
// that AI backends produce, so it can be used for history reconstruction.
func messageToJSON(msg *api.GameMessage) string {
	tr := TurnResponse{
		Narrative:            msg.Narrative,
		Translation:          msg.Translation,
		NPCDialog:            msg.NPCDialog,
		NPCDialogTranslation: msg.NPCDialogTranslation,
		Vocabulary:           msg.Vocabulary,
		Choices:              msg.Choices,
		Finished:             msg.Finished,
	}
	data, _ := json.Marshal(tr)
	return string(data)
}

// simulateStream sends cached text to a StreamCallback in one shot.
// No artificial delay — cached responses should appear instantly.
func simulateStream(text string, callback StreamCallback) {
	callback(text)
}

// CachedClient wraps an ai.Client and caches first-turn responses for premade scenarios.
// Subsequent turns (SendInput/SendInputStream) pass through to the inner client uncached.
type CachedClient struct {
	inner Client
	cache *responseCache
}

// NewCachedClient creates a caching wrapper around the given AI client.
func NewCachedClient(inner Client) *CachedClient {
	return &CachedClient{
		inner: inner,
		cache: newResponseCache(defaultCacheTTL),
	}
}

func (c *CachedClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string) (*api.GameMessage, any, float64, error) {
	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code, explanationLang)
		if entry := c.cache.get(key); entry != nil {
			log.Printf("[cache] HIT StartScenario scenario=%s lang=%s", scenario.ID, lang.Code)
			history := c.inner.BuildStartHistory(scenario, lang, explanationLang, entry.ResponseText)
			msg := entry.Message // struct copy
			return &msg, history, 0, nil
		}
		log.Printf("[cache] MISS StartScenario scenario=%s lang=%s", scenario.ID, lang.Code)
	}

	msg, history, cost, err := c.inner.StartScenario(ctx, scenario, lang, explanationLang)
	if err != nil {
		return nil, nil, 0, err
	}

	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code, explanationLang)
		c.cache.set(key, &cachedStart{
			ResponseText: messageToJSON(msg),
			Message:      *msg,
			CreatedAt:    time.Now(),
		})
	}

	return msg, history, cost, nil
}

func (c *CachedClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, callback StreamCallback) (*api.GameMessage, any, float64, error) {
	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code, explanationLang)
		if entry := c.cache.get(key); entry != nil {
			log.Printf("[cache] HIT StartScenarioStream scenario=%s lang=%s", scenario.ID, lang.Code)
			if callback != nil {
				simulateStream(entry.ResponseText, callback)
			}
			history := c.inner.BuildStartHistory(scenario, lang, explanationLang, entry.ResponseText)
			msg := entry.Message
			return &msg, history, 0, nil
		}
		log.Printf("[cache] MISS StartScenarioStream scenario=%s lang=%s", scenario.ID, lang.Code)
	}

	msg, history, cost, err := c.inner.StartScenarioStream(ctx, scenario, lang, explanationLang, callback)
	if err != nil {
		return nil, nil, 0, err
	}

	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code, explanationLang)
		c.cache.set(key, &cachedStart{
			ResponseText: messageToJSON(msg),
			Message:      *msg,
			CreatedAt:    time.Now(),
		})
	}

	return msg, history, cost, nil
}

func (c *CachedClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string) (*api.GameMessage, *api.Correction, any, float64, error) {
	return c.inner.SendInput(ctx, scenario, lang, explanationLang, history, input)
}

func (c *CachedClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, explanationLang string, history any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, float64, error) {
	return c.inner.SendInputStream(ctx, scenario, lang, explanationLang, history, input, callback)
}

func (c *CachedClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, explanationLang string, responseText string) any {
	return c.inner.BuildStartHistory(scenario, lang, explanationLang, responseText)
}
