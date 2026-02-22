package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/joshburnsxyz/lark/api"
)

// Settings holds user preferences. Fields use "Hide" so that the zero value
// (false) means everything is shown by default.
type Settings struct {
	HideTranslations bool   `json:"hideTranslations"`
	HideChoices      bool   `json:"hideChoices"`
	HideVocabulary   bool   `json:"hideVocabulary"`
	HideGrammar      bool   `json:"hideGrammar"`
	ExplanationLang  string `json:"explanationLang,omitempty"` // empty = English
}

// SaveData holds all locally persisted state.
type SaveData struct {
	Completed       map[string]bool          `json:"completed"`
	Sessions        map[string]*SavedSession `json:"sessions"`
	CustomScenarios []CustomScenario         `json:"customScenarios"`
	Settings        Settings                 `json:"settings"`
}

// SavedSession stores an in-progress game session for resume.
type SavedSession struct {
	SessionID      string           `json:"sessionId"`
	ScenarioName   string           `json:"scenarioName"`
	LastMessage    *api.GameMessage `json:"lastMessage"`
	LastCorrection *api.Correction  `json:"lastCorrection,omitempty"`
	CustomPrompt   string           `json:"customPrompt,omitempty"`
}

// CustomScenario is a user-created scenario stored locally.
type CustomScenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

const saveFileName = "savestate.json"

func saveFilePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, saveFileName), nil
}

func loadSaveData() *SaveData {
	path, err := saveFilePath()
	if err != nil {
		return newSaveData()
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return newSaveData()
	}
	var sd SaveData
	if err := json.Unmarshal(data, &sd); err != nil {
		return newSaveData()
	}
	if sd.Completed == nil {
		sd.Completed = make(map[string]bool)
	}
	if sd.Sessions == nil {
		sd.Sessions = make(map[string]*SavedSession)
	}
	return &sd
}

func saveSaveData(sd *SaveData) error {
	path, err := saveFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func newSaveData() *SaveData {
	return &SaveData{
		Completed: make(map[string]bool),
		Sessions:  make(map[string]*SavedSession),
	}
}

func saveKey(scenarioID, langCode string) string {
	return scenarioID + ":" + langCode
}

func (s *SaveData) IsCompleted(scenarioID, langCode string) bool {
	return s.Completed[saveKey(scenarioID, langCode)]
}

func (s *SaveData) MarkCompleted(scenarioID, langCode string) {
	s.Completed[saveKey(scenarioID, langCode)] = true
}

func (s *SaveData) SaveSession(scenarioID, langCode string, session *SavedSession) {
	s.Sessions[saveKey(scenarioID, langCode)] = session
}

func (s *SaveData) GetSession(scenarioID, langCode string) *SavedSession {
	return s.Sessions[saveKey(scenarioID, langCode)]
}

func (s *SaveData) ClearSession(scenarioID, langCode string) {
	delete(s.Sessions, saveKey(scenarioID, langCode))
}

func (s *SaveData) AddCustomScenario(cs CustomScenario) {
	s.CustomScenarios = append(s.CustomScenarios, cs)
}
