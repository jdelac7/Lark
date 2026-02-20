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

func startKey(scenarioID, langCode string) string {
	return scenarioID + ":" + langCode
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
		Choices:              msg.Choices,
		Vocabulary:           msg.Vocabulary,
		Finished:             msg.Finished,
	}
	data, _ := json.Marshal(tr)
	return string(data)
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

func (c *CachedClient) StartScenario(ctx context.Context, scenario *api.Scenario, lang *api.Language) (*api.GameMessage, any, error) {
	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code)
		if entry := c.cache.get(key); entry != nil {
			log.Printf("[cache] HIT StartScenario scenario=%s lang=%s", scenario.ID, lang.Code)
			history := c.inner.BuildStartHistory(scenario, lang, entry.ResponseText)
			msg := entry.Message // struct copy
			return &msg, history, nil
		}
		log.Printf("[cache] MISS StartScenario scenario=%s lang=%s", scenario.ID, lang.Code)
	}

	msg, history, err := c.inner.StartScenario(ctx, scenario, lang)
	if err != nil {
		return nil, nil, err
	}

	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code)
		c.cache.set(key, &cachedStart{
			ResponseText: messageToJSON(msg),
			Message:      *msg,
			CreatedAt:    time.Now(),
		})
	}

	return msg, history, nil
}

func (c *CachedClient) StartScenarioStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, callback StreamCallback) (*api.GameMessage, any, error) {
	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code)
		if entry := c.cache.get(key); entry != nil {
			log.Printf("[cache] HIT StartScenarioStream scenario=%s lang=%s", scenario.ID, lang.Code)
			if callback != nil {
				callback(entry.ResponseText)
			}
			history := c.inner.BuildStartHistory(scenario, lang, entry.ResponseText)
			msg := entry.Message
			return &msg, history, nil
		}
		log.Printf("[cache] MISS StartScenarioStream scenario=%s lang=%s", scenario.ID, lang.Code)
	}

	msg, history, err := c.inner.StartScenarioStream(ctx, scenario, lang, callback)
	if err != nil {
		return nil, nil, err
	}

	if isCacheable(scenario.ID) {
		key := startKey(scenario.ID, lang.Code)
		c.cache.set(key, &cachedStart{
			ResponseText: messageToJSON(msg),
			Message:      *msg,
			CreatedAt:    time.Now(),
		})
	}

	return msg, history, nil
}

func (c *CachedClient) SendInput(ctx context.Context, scenario *api.Scenario, lang *api.Language, history any, input string) (*api.GameMessage, *api.Correction, any, error) {
	return c.inner.SendInput(ctx, scenario, lang, history, input)
}

func (c *CachedClient) SendInputStream(ctx context.Context, scenario *api.Scenario, lang *api.Language, history any, input string, callback StreamCallback) (*api.GameMessage, *api.Correction, any, error) {
	return c.inner.SendInputStream(ctx, scenario, lang, history, input, callback)
}

func (c *CachedClient) BuildStartHistory(scenario *api.Scenario, lang *api.Language, responseText string) any {
	return c.inner.BuildStartHistory(scenario, lang, responseText)
}
