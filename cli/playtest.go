package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	mathrand "math/rand"
	"os"
	"time"

	"github.com/joshburnsxyz/lark/api"
)


type playtestTurnOutput struct {
	Turn       int             `json:"turn"`
	InputMode  string          `json:"inputMode"`
	ChosenIdx  int             `json:"chosenIndex,omitempty"`
	ChosenText string          `json:"chosenText,omitempty"`
	FreeText   string          `json:"freeText,omitempty"`
	Message    api.GameMessage `json:"message"`
	Correction *api.Correction `json:"correction,omitempty"`
}

type playtestSummary struct {
	Type            string          `json:"type"`
	ScenarioID      string          `json:"scenarioId"`
	ScenarioName    string          `json:"scenarioName"`
	Language        string          `json:"language"`
	LanguageCode    string          `json:"languageCode"`
	TotalTurns      int             `json:"totalTurns"`
	Finished        bool            `json:"finished"`
	AllVocabulary   []api.VocabItem `json:"allVocabulary"`
	UniqueWords     int             `json:"uniqueWords"`
	DuplicateWords  []string        `json:"duplicateWords,omitempty"`
	FreeTextTurns   int             `json:"freeTextTurns"`
	CorrectionsGot  int             `json:"correctionsReceived"`
	AllCorrections  []api.Correction `json:"allCorrections,omitempty"`
}

func handlePlaytestCommand(args []string) {
	fs := flag.NewFlagSet("playtest", flag.ExitOnError)
	maxTurns := fs.Int("max-turns", 30, "maximum number of turns (0 = no limit)")
	seed := fs.Int64("seed", 0, "random seed (0 = current time)")
	freeTextRatio := fs.Float64("free-text-ratio", 0, "fraction of turns to use free text input (0.0-1.0)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: lark playtest <scenario-id> <language-code> [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nScenarios:\n")
		for _, s := range api.Scenarios {
			fmt.Fprintf(os.Stderr, "  %-25s %s\n", s.ID, s.Name)
		}
		fmt.Fprintf(os.Stderr, "\nPopular languages:\n")
		for _, l := range api.PopularLanguages {
			fmt.Fprintf(os.Stderr, "  %-6s %s\n", l.Code, l.Name)
		}
	}

	// Need at least scenario + language before flags
	if len(args) < 3 {
		fs.Usage()
		os.Exit(1)
	}

	scenarioID := args[1]
	langCode := args[2]

	if err := fs.Parse(args[3:]); err != nil {
		os.Exit(1)
	}

	scenario := api.ScenarioByID(scenarioID)
	if scenario == nil {
		emitPlaytestError(fmt.Sprintf("unknown scenario: %s", scenarioID))
		os.Exit(1)
	}

	lang := api.LanguageByCode(langCode)
	if lang == nil {
		emitPlaytestError(fmt.Sprintf("unknown language code: %s", langCode))
		os.Exit(1)
	}

	// Always use server mode so cost events are tracked
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
	var gameClient GameClient = NewClient(serverURL, playerID)

	seedVal := *seed
	if seedVal == 0 {
		seedVal = time.Now().UnixNano()
	}
	rng := mathrand.New(mathrand.NewSource(seedVal))

	runPlaytest(gameClient, scenario, lang, rng, *maxTurns, *freeTextRatio)
}

func runPlaytest(client GameClient, scenario *api.Scenario, lang *api.Language, rng *mathrand.Rand, maxTurns int, freeTextRatio float64) {
	enc := json.NewEncoder(os.Stdout)

	startResp, err := client.StartScenarioCustom(scenario.ID, lang.Code, "", "")
	if err != nil {
		emitPlaytestError("failed to start scenario: " + err.Error())
		os.Exit(1)
	}

	sessionID := startResp.SessionID
	msg := &startResp.Message

	var allVocab []api.VocabItem
	var allCorrections []api.Correction
	freeTextTurns := 0
	allVocab = append(allVocab, msg.Vocabulary...)

	enc.Encode(playtestTurnOutput{
		Turn:      0,
		InputMode: "start",
		Message:   *msg,
	})

	turn := 0
	finished := msg.Finished

	for !finished && (maxTurns <= 0 || turn < maxTurns) {
		if len(msg.Choices) == 0 {
			break
		}

		choiceIdx := rng.Intn(len(msg.Choices))
		choiceText := msg.Choices[choiceIdx].Text
		turn++

		// Decide whether to use free text on this turn
		useFreeText := freeTextRatio > 0 && rng.Float64() < freeTextRatio

		var resp *api.PlayerInputResponse
		var turnOutput playtestTurnOutput

		if useFreeText {
			// Send the choice text as free text (simulates player typing)
			resp, err = client.SendFreeText(sessionID, choiceText)
			if err != nil {
				emitPlaytestError(fmt.Sprintf("turn %d (free text): %s", turn, err.Error()))
				os.Exit(1)
			}
			freeTextTurns++
			turnOutput = playtestTurnOutput{
				Turn:       turn,
				InputMode:  "free_text",
				FreeText:   choiceText,
				Message:    resp.Message,
				Correction: resp.Correction,
			}
		} else {
			resp, err = client.SendChoice(sessionID, choiceIdx)
			if err != nil {
				emitPlaytestError(fmt.Sprintf("turn %d: %s", turn, err.Error()))
				os.Exit(1)
			}
			turnOutput = playtestTurnOutput{
				Turn:       turn,
				InputMode:  "choice",
				ChosenIdx:  choiceIdx,
				ChosenText: choiceText,
				Message:    resp.Message,
				Correction: resp.Correction,
			}
		}

		msg = &resp.Message
		allVocab = append(allVocab, msg.Vocabulary...)
		if resp.Correction != nil {
			allCorrections = append(allCorrections, *resp.Correction)
		}
		finished = msg.Finished

		enc.Encode(turnOutput)
	}

	// Compute vocabulary stats
	seen := make(map[string]int)
	for _, v := range allVocab {
		seen[v.Word]++
	}
	var dupes []string
	for word, count := range seen {
		if count > 1 {
			dupes = append(dupes, word)
		}
	}

	enc.Encode(playtestSummary{
		Type:           "summary",
		ScenarioID:     scenario.ID,
		ScenarioName:   scenario.Name,
		Language:       lang.Name,
		LanguageCode:   lang.Code,
		TotalTurns:     turn,
		Finished:       finished,
		AllVocabulary:  allVocab,
		UniqueWords:    len(seen),
		DuplicateWords: dupes,
		FreeTextTurns:  freeTextTurns,
		CorrectionsGot: len(allCorrections),
		AllCorrections: allCorrections,
	})
}

func emitPlaytestError(msg string) {
	j, _ := json.Marshal(map[string]string{"type": "error", "error": msg})
	fmt.Fprintln(os.Stdout, string(j))
}
