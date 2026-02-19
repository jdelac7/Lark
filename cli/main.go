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

	RenderBanner()
	input := readInput()
	if input == ctrlC {
		return
	}

	for {
		if !playSession(client) {
			break
		}
		fmt.Print("\nPlay again? (y/n): ")
		again := readInput()
		if again == ctrlC || (again != "y" && again != "Y") {
			break
		}
	}
}

func playSession(client *Client) bool {
	scenarios, err := client.GetScenarios()
	if err != nil {
		PrintError("Failed to connect to server: " + err.Error())
		PrintError("Make sure the Lark server is running at " + client.baseURL)
		readInput()
		return false
	}

	languages, err := client.GetLanguages()
	if err != nil {
		PrintError("Failed to fetch languages: " + err.Error())
		readInput()
		return false
	}

	// Choose scenario
	var scenario api.Scenario
	var customPrompt string
	scenarioIdx, scenarioText := ReadScenarioChoice(scenarios)
	if scenarioIdx == -2 {
		return false
	}
	if scenarioIdx >= 0 {
		scenario = scenarios[scenarioIdx]
	} else {
		// Custom scenario from user text
		customPrompt = scenarioText
		scenario = api.Scenario{
			ID:          "custom",
			Name:        customPrompt,
			Description: customPrompt,
			Difficulty:  api.DifficultyBeginner,
		}
	}

	// Choose language
	RenderLanguageList(languages)
	langIdx := ReadMenuChoice(len(languages))
	if langIdx < 0 {
		return false
	}
	lang := languages[langIdx]

	// Start session (streaming)
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

	// Game loop
	for {
		if msg.Finished {
			RenderFinishedScreen(scenario.Name, lang.Name, msg)
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
				// Stay in input loop — warning is visible, prompt re-shown
				continue
			}
			break
		}

		ResetStreamState()
		rawJSON = ""
		RenderStreamingScreen(scenario.Name, lang.Name, lastCorrection, "")
		streamToken := func(token string) {
			rawJSON += token
			RenderStreamingScreen(scenario.Name, lang.Name, lastCorrection, rawJSON)
		}

		var resp *api.PlayerInputResponse
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
	}
}
