package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/joshburnsxyz/lark/api"
	"github.com/joshburnsxyz/lark/server/ai"
)

const defaultBYOKModel = "x-ai/grok-4.1-fast"

// GameClient is the interface used by the game loop.
// Both the server-backed Client and the local-mode LocalClient implement it.
type GameClient interface {
	GetScenarios() ([]api.Scenario, error)
	GetLanguages() ([]api.Language, error)
	StartScenarioCustom(scenarioID, language, customPrompt, explanationLang string) (*api.StartResponse, error)
	StreamStartScenarioCustom(scenarioID, language, customPrompt, explanationLang string, onToken func(string)) (*api.StartResponse, error)
	SendChoice(sessionID string, choiceIndex int) (*api.PlayerInputResponse, error)
	SendFreeText(sessionID, text string) (*api.PlayerInputResponse, error)
	StreamSendChoice(sessionID string, choiceIndex int, onToken func(string)) (*api.PlayerInputResponse, error)
	StreamSendFreeText(sessionID, text string, onToken func(string)) (*api.PlayerInputResponse, error)
}

// LocalClient calls OpenRouter directly without going through the server.
// Used in BYOK mode when the user provides their own API key.
type LocalClient struct {
	aiClient        ai.Client
	scenario        *api.Scenario
	lang            *api.Language
	explanationLang string
	history         any
	lastMessage     *api.GameMessage
}

// NewLocalClient creates a LocalClient that talks to OpenRouter directly.
func NewLocalClient(apiKey, model string) *LocalClient {
	return &LocalClient{
		aiClient: ai.NewOpenRouterClient(apiKey, model),
	}
}

func (l *LocalClient) GetScenarios() ([]api.Scenario, error) {
	return api.Scenarios, nil
}

func (l *LocalClient) GetLanguages() ([]api.Language, error) {
	return api.Languages, nil
}

func (l *LocalClient) resolveScenario(scenarioID, language, customPrompt string) (*api.Scenario, *api.Language, error) {
	lang := api.LanguageByCode(language)
	if lang == nil {
		return nil, nil, fmt.Errorf("unsupported language: %s", language)
	}

	scenario := api.ScenarioByID(scenarioID)
	if scenario == nil {
		if customPrompt != "" {
			scenario = &api.Scenario{
				ID:          scenarioID,
				Name:        customPrompt,
				Description: customPrompt,
				Difficulty:  api.DifficultyBeginner,
			}
		} else {
			return nil, nil, fmt.Errorf("unknown scenario: %s", scenarioID)
		}
	}

	return scenario, lang, nil
}

func localSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "local-" + hex.EncodeToString(b)
}

func (l *LocalClient) StreamStartScenarioCustom(scenarioID, language, customPrompt, explanationLang string, onToken func(string)) (*api.StartResponse, error) {
	scenario, lang, err := l.resolveScenario(scenarioID, language, customPrompt)
	if err != nil {
		return nil, err
	}

	msg, history, _, err := l.aiClient.StartScenarioStream(context.Background(), scenario, lang, explanationLang, onToken)
	if err != nil {
		return nil, fmt.Errorf("starting scenario: %w", err)
	}

	l.scenario = scenario
	l.lang = lang
	l.explanationLang = explanationLang
	l.history = history
	l.lastMessage = msg

	return &api.StartResponse{
		SessionID: localSessionID(),
		Message:   *msg,
	}, nil
}

func (l *LocalClient) StartScenarioCustom(scenarioID, language, customPrompt, explanationLang string) (*api.StartResponse, error) {
	scenario, lang, err := l.resolveScenario(scenarioID, language, customPrompt)
	if err != nil {
		return nil, err
	}

	msg, history, _, err := l.aiClient.StartScenario(context.Background(), scenario, lang, explanationLang)
	if err != nil {
		return nil, fmt.Errorf("starting scenario: %w", err)
	}

	l.scenario = scenario
	l.lang = lang
	l.explanationLang = explanationLang
	l.history = history
	l.lastMessage = msg

	return &api.StartResponse{
		SessionID: localSessionID(),
		Message:   *msg,
	}, nil
}

func (l *LocalClient) StreamSendChoice(sessionID string, choiceIndex int, onToken func(string)) (*api.PlayerInputResponse, error) {
	if l.lastMessage == nil || choiceIndex < 0 || choiceIndex >= len(l.lastMessage.Choices) {
		return nil, fmt.Errorf("invalid choice index")
	}
	choiceText := l.lastMessage.Choices[choiceIndex].Text
	input := ai.FormatChoiceInput(choiceIndex, choiceText)
	return l.sendInput(input, onToken)
}

func (l *LocalClient) StreamSendFreeText(sessionID, text string, onToken func(string)) (*api.PlayerInputResponse, error) {
	input := ai.FormatFreeTextInput(text)
	return l.sendInput(input, onToken)
}

func (l *LocalClient) SendChoice(sessionID string, choiceIndex int) (*api.PlayerInputResponse, error) {
	if l.lastMessage == nil || choiceIndex < 0 || choiceIndex >= len(l.lastMessage.Choices) {
		return nil, fmt.Errorf("invalid choice index")
	}
	choiceText := l.lastMessage.Choices[choiceIndex].Text
	input := ai.FormatChoiceInput(choiceIndex, choiceText)
	return l.sendInputNoStream(input)
}

func (l *LocalClient) SendFreeText(sessionID, text string) (*api.PlayerInputResponse, error) {
	input := ai.FormatFreeTextInput(text)
	return l.sendInputNoStream(input)
}

func (l *LocalClient) sendInput(input string, onToken func(string)) (*api.PlayerInputResponse, error) {
	msg, correction, newHistory, _, err := l.aiClient.SendInputStream(
		context.Background(), l.scenario, l.lang, l.explanationLang, l.history, input, onToken,
	)
	if err != nil {
		return nil, fmt.Errorf("sending input: %w", err)
	}

	l.history = newHistory
	l.lastMessage = msg

	return &api.PlayerInputResponse{
		Message:    *msg,
		Correction: correction,
	}, nil
}

func (l *LocalClient) sendInputNoStream(input string) (*api.PlayerInputResponse, error) {
	msg, correction, newHistory, _, err := l.aiClient.SendInput(
		context.Background(), l.scenario, l.lang, l.explanationLang, l.history, input,
	)
	if err != nil {
		return nil, fmt.Errorf("sending input: %w", err)
	}

	l.history = newHistory
	l.lastMessage = msg

	return &api.PlayerInputResponse{
		Message:    *msg,
		Correction: correction,
	}, nil
}

// getBYOKModel returns the model to use for BYOK mode from LARK_MODEL env var.
func getBYOKModel() string {
	if m := os.Getenv("LARK_MODEL"); m != "" {
		return m
	}
	return defaultBYOKModel
}
