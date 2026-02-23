package api

// InputMode represents how the player provides input.
type InputMode string

const (
	InputModeChoice   InputMode = "choice"
	InputModeFreeText InputMode = "free_text"
)

// Difficulty level for scenarios.
type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

// Category groups scenarios into thematic collections.
type Category string

const (
	CategoryEveryday  Category = "everyday"
	CategoryAdventure Category = "adventure"
)

// Scenario represents a playable scenario in the catalog.
type Scenario struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Difficulty  Difficulty `json:"difficulty"`
	Category    Category   `json:"category"`
}

// Language represents a supported target language.
type Language struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// StartRequest is sent to POST /scenarios/start.
type StartRequest struct {
	ScenarioID      string `json:"scenarioId"`
	Language        string `json:"language"`
	CustomPrompt    string `json:"customPrompt,omitempty"`
	ExplanationLang string `json:"explanationLang,omitempty"`
}

// StartResponse is returned from POST /scenarios/start.
type StartResponse struct {
	SessionID string      `json:"sessionId"`
	Message   GameMessage `json:"message"`
}

// PlayerInputRequest is sent to POST /game/input.
type PlayerInputRequest struct {
	SessionID   string    `json:"sessionId"`
	Mode        InputMode `json:"mode"`
	Text        string    `json:"text,omitempty"`
	ChoiceIndex int       `json:"choiceIndex,omitempty"`
}

// PlayerInputResponse is returned from POST /game/input.
type PlayerInputResponse struct {
	Message    GameMessage `json:"message"`
	Correction *Correction `json:"correction,omitempty"`
}

// GameMessage is the core response from a game turn.
type GameMessage struct {
	Narrative   string       `json:"narrative"`
	Translation string       `json:"translation"`
	NPCDialog   string       `json:"npcDialog,omitempty"`
	NPCDialogTranslation string `json:"npcDialogTranslation,omitempty"`
	Vocabulary  []VocabItem  `json:"vocabulary,omitempty"`
	Choices     []Choice     `json:"choices,omitempty"`
	Finished    bool         `json:"finished"`
}

// Choice is an option the player can select.
type Choice struct {
	Text        string `json:"text"`
	Translation string `json:"translation"`
}

// VocabItem is a vocabulary word/phrase from the turn.
type VocabItem struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Usage       string `json:"usage,omitempty"`
}

// Correction is grammar/vocabulary feedback on free-text input.
type Correction struct {
	Original  string `json:"original"`
	Corrected string `json:"corrected"`
	Explanation string `json:"explanation"`
}

// GameStateResponse is returned from GET /game/state.
type GameStateResponse struct {
	SessionID  string      `json:"sessionId"`
	ScenarioID string      `json:"scenarioId"`
	Language   string      `json:"language"`
	TurnCount  int         `json:"turnCount"`
	Message    GameMessage `json:"message"`
}

// ProgressResponse is returned from GET /progress.
type ProgressResponse struct {
	PlayerID           string      `json:"playerId"`
	CompletedScenarios []CompletedScenario `json:"completedScenarios"`
	VocabBank          []VocabItem `json:"vocabBank"`
}

// CompletedScenario records a finished scenario.
type CompletedScenario struct {
	ScenarioID string `json:"scenarioId"`
	Language   string `json:"language"`
	TurnCount  int    `json:"turnCount"`
}

// ErrorResponse is a standard error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}
