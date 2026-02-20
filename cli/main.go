package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/joshburnsxyz/lark/api"
)

// lastCtrlC tracks when Ctrl-C was last pressed for double-press quit.
var lastCtrlC time.Time

// handleCtrlC implements two-press quit. Returns true if the app should exit.
func handleCtrlC() bool {
	now := time.Now()
	if now.Sub(lastCtrlC) < 2*time.Second {
		return true // second press within 2s — quit
	}
	lastCtrlC = now
	ShowCursor()
	PrintWarning("Press Ctrl-C again to quit")
	return false
}

func main() {
	// Subcommand routing
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "activate":
			handleActivateCommand(os.Args[1:])
			return
		case "deactivate":
			handleDeactivateCommand()
			return
		}
	}

	// License gate
	if err := checkLicense(); err != nil {
		fmt.Fprintf(os.Stderr, "\n  %s\n\n  Get your license at: %s\n\n", err, websiteURL)
		os.Exit(1)
	}

	EnterAltScreen()
	defer LeaveAltScreen()

	serverURL := os.Getenv("LARK_SERVER")
	if serverURL == "" {
		serverURL = "http://localhost:9292"
	}

	playerID := os.Getenv("LARK_PLAYER_ID")
	if playerID == "" {
		b := make([]byte, 16)
		rand.Read(b)
		playerID = hex.EncodeToString(b)
	}

	client := NewClient(serverURL, playerID)

	// Show banner while connecting
	RenderBanner()

	scenarios, err := client.GetScenarios()
	if err != nil {
		PrintError("Failed to connect to server: " + err.Error())
		PrintError("Make sure the Lark server is running at " + client.baseURL)
		readInput()
		return
	}

	languages, err := client.GetLanguages()
	if err != nil {
		PrintError("Failed to fetch languages: " + err.Error())
		readInput()
		return
	}

	saveData := loadSaveData()
	appSettings = &saveData.Settings

	for {
		if !playSession(client, scenarios, languages, saveData) {
			break
		}
	}
}

// playSession runs the selection flow (language → category → scenario → game)
// with go-back support at each step. Returns true if the game completed
// normally and the caller should offer "play again".
func playSession(client *Client, scenarios []api.Scenario, languages []api.Language, saveData *SaveData) bool {
	langNames := make([]string, len(languages))
	for i, l := range languages {
		langNames[i] = l.Name
	}

	var lang api.Language
	var selectedCategory api.Category
	var filtered []api.Scenario
	var scenario api.Scenario
	var customPrompt string

	step := 0
	for {
		switch step {
		// ── Step 0: Language selection (banner + two-column list) ──
		case 0:
			idx := ReadBannerLanguageChoice(len(langNames), func(cursor int) {
				RenderBannerLanguages(langNames, cursor)
			})
			if idx == -1 {
				return false
			}
			if idx == -2 {
				// Open settings, then return to this step
				ReadSettings(&saveData.Settings, func() { saveSaveData(saveData) })
				continue
			}
			lang = languages[idx]
			step = 1

		// ── Step 1: Category selection ──
		case 1:
			categories := []string{"Everyday Scenarios", "Adventure Scenarios", "Custom Scenario"}
			catIdx := ReadListChoice(len(categories), func(cursor int) {
				RenderListPage("Choose a Category", categories, cursor)
			})
			if catIdx == -1 {
				return false
			}
			if catIdx == -2 {
				step = 0
				continue
			}
			switch catIdx {
			case 0:
				selectedCategory = api.CategoryEveryday
				step = 2
			case 1:
				selectedCategory = api.CategoryAdventure
				step = 2
			case 2:
				step = 4 // custom scenario input
			}

		// ── Step 2: Scenario selection ──
		case 2:
			filtered = nil
			for _, s := range scenarios {
				if s.Category == selectedCategory {
					filtered = append(filtered, s)
				}
			}
			for _, cs := range saveData.CustomScenarios {
				if cs.Category == string(selectedCategory) {
					filtered = append(filtered, api.Scenario{
						ID:          cs.ID,
						Name:        cs.Name,
						Description: cs.Description,
						Difficulty:  api.DifficultyBeginner,
						Category:    selectedCategory,
					})
				}
			}

			scenarioIdx, scenarioText := ReadScenarioChoice(filtered, saveData.Completed, lang.Code)
			if scenarioIdx == -2 {
				return false // quit
			}
			if scenarioIdx == -3 {
				step = 1 // back to category
				continue
			}

			customPrompt = ""
			if scenarioIdx >= 0 {
				scenario = filtered[scenarioIdx]
			} else {
				customPrompt = scenarioText
				csID := "custom_" + randomHex(8)
				scenario = api.Scenario{
					ID:          csID,
					Name:        customPrompt,
					Description: customPrompt,
					Difficulty:  api.DifficultyBeginner,
					Category:    selectedCategory,
				}
				saveData.AddCustomScenario(CustomScenario{
					ID:          csID,
					Name:        customPrompt,
					Description: customPrompt,
					Category:    string(selectedCategory),
				})
				saveSaveData(saveData)
			}
			step = 3

		// ── Step 3: Check saved session / start game ──
		case 3:
			saved := saveData.GetSession(scenario.ID, lang.Code)
			if saved != nil {
				continueChoices := []string{"Continue where you left off", "Start again"}
				choice := ReadListChoice(len(continueChoices), func(cursor int) {
					RenderContinuePrompt(scenario.Name, continueChoices, cursor)
				})
				if choice == -1 {
					return false
				}
				if choice == -2 {
					step = 2 // back to scenario selection
					continue
				}
				if choice == 0 {
					if !resumeSession(client, scenario, lang, saveData, saved) {
						return false
					}
					step = 2
					continue
				}
				// Start again — clear saved session
				saveData.ClearSession(scenario.ID, lang.Code)
				saveSaveData(saveData)
			}
			if !startFreshSession(client, scenario, lang, customPrompt, saveData) {
				return false
			}
			step = 2

		// ── Step 4: Custom scenario input ──
		case 4:
			text := ReadCustomScenario()
			if text == ctrlC {
				return false
			}
			if text == "" {
				step = 1 // back to category
				continue
			}

			customPrompt = text
			csID := "custom_" + randomHex(8)
			scenario = api.Scenario{
				ID:          csID,
				Name:        customPrompt,
				Description: customPrompt,
				Difficulty:  api.DifficultyBeginner,
				Category:    api.CategoryEveryday,
			}
			saveData.AddCustomScenario(CustomScenario{
				ID:          csID,
				Name:        customPrompt,
				Description: customPrompt,
				Category:    string(api.CategoryEveryday),
			})
			saveSaveData(saveData)
			step = 3
		}
	}
}

func startFreshSession(client *Client, scenario api.Scenario, lang api.Language, customPrompt string, saveData *SaveData) bool {
	ResetStreamState()
	var rawJSON string
	RenderStreamingScreen(scenario.Name, lang.Name, nil, "")
	onToken := func(token string) {
		rawJSON += token
		RenderStreamingScreen(scenario.Name, lang.Name, nil, rawJSON)
	}

	startResp, err := client.StreamStartScenarioCustom(scenario.ID, lang.Code, customPrompt, onToken)
	if err != nil {
		// Fallback to non-streaming
		RenderThinkingScreen(scenario.Name, lang.Name)
		startResp, err = client.StartScenarioCustom(scenario.ID, lang.Code, customPrompt)
		if err != nil {
			ShowCursor()
			PrintError("Failed to start scenario: " + err.Error())
			readInput()
			return false
		}
	}

	sessionID := startResp.SessionID
	msg := &startResp.Message
	var lastCorrection *api.Correction

	// Save initial session state
	saveData.SaveSession(scenario.ID, lang.Code, &SavedSession{
		SessionID:    sessionID,
		ScenarioName: scenario.Name,
		LastMessage:  msg,
		CustomPrompt: customPrompt,
	})
	saveSaveData(saveData)

	return gameLoop(client, scenario, lang, sessionID, msg, lastCorrection, customPrompt, saveData)
}

func resumeSession(client *Client, scenario api.Scenario, lang api.Language, saveData *SaveData, saved *SavedSession) bool {
	sessionID := saved.SessionID
	msg := saved.LastMessage
	lastCorrection := saved.LastCorrection

	return gameLoop(client, scenario, lang, sessionID, msg, lastCorrection, saved.CustomPrompt, saveData)
}

func gameLoop(client *Client, scenario api.Scenario, lang api.Language, sessionID string, msg *api.GameMessage, lastCorrection *api.Correction, customPrompt string, saveData *SaveData) bool {
	for {
		if msg.Finished {
			RenderFinishedScreen(scenario.Name, lang.Name, msg)

			// Mark completed, clear session, save
			saveData.MarkCompleted(scenario.ID, lang.Code)
			saveData.ClearSession(scenario.ID, lang.Code)
			saveSaveData(saveData)

			readInput()
			return true
		}

		if len(msg.Choices) == 0 {
			PrintError("No choices available - scenario ended unexpectedly")
			readInput()
			return true
		}

		RenderGameScreen(&GameScreenData{
			ScenarioName: scenario.Name,
			Language:     lang.Name,
			Message:      msg,
			Correction:   lastCorrection,
		})

		var choiceIdx int
		var freeText string
		for {
			choiceIdx, freeText = ReadChoice(len(msg.Choices))
			if choiceIdx == -2 {
				if handleCtrlC() {
					return false
				}
				continue
			}
			break
		}

		ResetStreamState()
		var rawJSON string
		RenderStreamingScreen(scenario.Name, lang.Name, lastCorrection, "")
		streamToken := func(token string) {
			rawJSON += token
			RenderStreamingScreen(scenario.Name, lang.Name, lastCorrection, rawJSON)
		}

		var resp *api.PlayerInputResponse
		var err error
		if choiceIdx >= 0 {
			resp, err = client.StreamSendChoice(sessionID, choiceIdx, streamToken)
		} else {
			resp, err = client.StreamSendFreeText(sessionID, freeText, streamToken)
		}

		if err != nil {
			// Fallback to non-streaming
			RenderThinkingScreen(scenario.Name, lang.Name)
			if choiceIdx >= 0 {
				resp, err = client.SendChoice(sessionID, choiceIdx)
			} else {
				resp, err = client.SendFreeText(sessionID, freeText)
			}
		}

		if err != nil {
			ShowCursor()
			RenderGameScreen(&GameScreenData{
				ScenarioName: scenario.Name,
				Language:     lang.Name,
				Message:      msg,
				Correction:   lastCorrection,
			})
			PrintError("Server error, try again: " + err.Error())
			continue
		}

		lastCorrection = resp.Correction
		msg = &resp.Message

		// Save session state after each turn
		saveData.SaveSession(scenario.ID, lang.Code, &SavedSession{
			SessionID:      sessionID,
			ScenarioName:   scenario.Name,
			LastMessage:    msg,
			LastCorrection: lastCorrection,
			CustomPrompt:   customPrompt,
		})
		saveSaveData(saveData)
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
