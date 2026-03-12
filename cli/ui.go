package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/joshburnsxyz/lark/api"
)

// ANSI escape sequences
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	italic = "\033[3m"

	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	clearSeq    = "\033[2J\033[H"
	hideCursor  = "\033[?25l"
	showCursor  = "\033[?25h"
	altScreenOn = "\033[?1049h"
	altScreenOff = "\033[?1049l"
)

// winsize matches the C struct winsize for TIOCGWINSZ.
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// getTermSize returns terminal width and height, defaulting to 80x24.
func getTermSize() (int, int) {
	ws := &winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)
	if errno == 0 && ws.Col > 0 && ws.Row > 0 {
		return int(ws.Col), int(ws.Row)
	}
	return 80, 24
}

// EnterAltScreen switches to the alternate screen buffer.
func EnterAltScreen() {
	fmt.Print(altScreenOn + hideCursor + clearSeq)
}

// LeaveAltScreen switches back to the main screen buffer.
func LeaveAltScreen() {
	fmt.Print(showCursor + altScreenOff)
}

// ShowCursor makes the cursor visible.
func ShowCursor() {
	fmt.Print(showCursor)
}

// appSettings points to the active settings from save data. Checked by
// render functions to conditionally hide translations, choices, etc.
var appSettings *Settings

// --- box drawing helpers ---

func hline(w int) string {
	return strings.Repeat("─", w-2)
}

func boxTop(w int) string    { return "┌" + hline(w) + "┐" }
func boxBottom(w int) string { return "└" + hline(w) + "┘" }
func boxDiv(w int) string    { return "├" + hline(w) + "┤" }
func boxEmpty(w int) string  { return boxLine("", w) }

func boxLine(content string, w int) string {
	inner := w - 4
	content = truncateVisible(content, inner)
	return "│ " + padRight(content, inner) + " │"
}

func padRight(s string, w int) string {
	vis := visibleLen(s)
	if vis >= w {
		return s
	}
	return s + strings.Repeat(" ", w-vis)
}

// truncateVisible cuts a string to maxWidth visible characters,
// preserving ANSI escape sequences and appending a reset if truncated.
func truncateVisible(s string, maxWidth int) string {
	vis := visibleLen(s)
	if vis <= maxWidth {
		return s
	}
	var b strings.Builder
	count := 0
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			b.WriteRune(r)
			continue
		}
		if inEsc {
			b.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		if count >= maxWidth {
			break
		}
		b.WriteRune(r)
		count++
	}
	b.WriteString(reset)
	return b.String()
}

// visibleLen returns the display width of a string, ignoring ANSI escapes.
func visibleLen(s string) int {
	n := 0
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}

// wrapText word-wraps a plain string to maxWidth columns.
func wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) > maxWidth {
			lines = append(lines, cur)
			cur = w
		} else {
			cur += " " + w
		}
	}
	return append(lines, cur)
}

// writeLines appends a slice of pre-built lines to a builder, each followed by \n.
func writeLines(b *strings.Builder, lines []string) {
	for _, l := range lines {
		b.WriteString(l)
		b.WriteByte('\n')
	}
}

// headerLine builds the header content string.
func headerLine(scenario, lang string) string {
	return bold + cyan + "  Lark" + reset + dim + "  ·  " + reset +
		bold + scenario + reset + dim + "  ·  " + reset +
		bold + lang + reset
}

// --- section builders (each returns a []string of box lines) ---

func buildHeader(scenario, lang string, w int) []string {
	return []string{
		boxTop(w),
		boxLine(headerLine(scenario, lang), w),
		boxDiv(w),
	}
}

func buildNarrative(msg *api.GameMessage, w, inner int) []string {
	var lines []string
	lines = append(lines, boxEmpty(w))
	for _, l := range wrapText(msg.Narrative, inner) {
		lines = append(lines, boxLine(bold+blue+l+reset, w))
	}
	if appSettings == nil || !appSettings.HideTranslations {
		for _, l := range wrapText(msg.Translation, inner) {
			lines = append(lines, boxLine(dim+l+reset, w))
		}
	}
	return lines
}

func buildNPC(msg *api.GameMessage, w, inner int) []string {
	if msg.NPCDialog == "" {
		return nil
	}
	var lines []string
	lines = append(lines, boxEmpty(w))
	for i, l := range wrapText(msg.NPCDialog, inner-7) {
		pfx := bold + yellow + "  NPC: " + reset + yellow
		if i > 0 {
			pfx = "       " + yellow
		}
		lines = append(lines, boxLine(pfx+l+reset, w))
	}
	if appSettings == nil || !appSettings.HideTranslations {
		for _, l := range wrapText(msg.NPCDialogTranslation, inner-7) {
			lines = append(lines, boxLine("       "+dim+l+reset, w))
		}
	}
	return lines
}

func buildVocab(vocab []api.VocabItem, w int) []string {
	if len(vocab) == 0 || (appSettings != nil && appSettings.HideVocabulary) {
		return nil
	}
	var lines []string
	lines = append(lines, boxDiv(w))
	lines = append(lines, boxLine(bold+magenta+"Vocabulary:"+reset, w))
	for _, v := range vocab {
		entry := "  " + bold + v.Word + reset + " - " + v.Translation
		if v.Usage != "" {
			entry += "  " + dim + "(" + v.Usage + ")" + reset
		}
		lines = append(lines, boxLine(entry, w))
	}
	return lines
}

func buildGrammar(c *api.Correction, w, inner int) []string {
	if c == nil || c.Original == "" || (appSettings != nil && appSettings.HideGrammar) {
		return nil
	}
	var lines []string
	lines = append(lines, boxDiv(w))
	lines = append(lines, boxLine(bold+yellow+"Grammar Note:"+reset, w))
	lines = append(lines, boxLine("  You wrote:  "+red+c.Original+reset, w))
	lines = append(lines, boxLine("  Corrected:  "+green+c.Corrected+reset, w))
	for _, l := range wrapText(c.Explanation, inner-2) {
		lines = append(lines, boxLine("  "+dim+l+reset, w))
	}
	return lines
}

func buildChoices(msg *api.GameMessage, w int) []string {
	if len(msg.Choices) == 0 || msg.Finished {
		return nil
	}
	// When choices are hidden, show only a free-text prompt
	if appSettings != nil && appSettings.HideChoices {
		return []string{
			boxDiv(w),
			boxLine("  "+dim+italic+"Write your response in the target language..."+reset, w),
		}
	}
	var lines []string
	lines = append(lines, boxDiv(w))
	hideTranslations := appSettings != nil && appSettings.HideTranslations
	for i, c := range msg.Choices {
		var l string
		if hideTranslations {
			l = fmt.Sprintf("  %s%d)%s %s%s%s",
				bold+green, i+1, reset,
				cyan, c.Text, reset)
		} else {
			l = fmt.Sprintf("  %s%d)%s %s%s%s  %s(%s)%s",
				bold+green, i+1, reset,
				cyan, c.Text, reset,
				dim, c.Translation, reset)
		}
		lines = append(lines, boxLine(l, w))
	}
	// Add explicit free-text option
	freeIdx := len(msg.Choices) + 1
	l := fmt.Sprintf("  %s%d)%s %s%s%s",
		bold+green, freeIdx, reset,
		dim+italic, "Write your own response...", reset)
	lines = append(lines, boxLine(l, w))
	return lines
}

// --- screen render functions ---

// GameScreenData holds all data for a game turn render.
type GameScreenData struct {
	ScenarioName string
	Language     string
	Message      *api.GameMessage
	Correction   *api.Correction
}

// RenderGameScreen draws the full game turn, fitting within the terminal.
func RenderGameScreen(data *GameScreenData) {
	w, h := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}
	inner := w - 4

	header := buildHeader(data.ScenarioName, data.Language, w)
	narrative := buildNarrative(data.Message, w, inner)
	npc := buildNPC(data.Message, w, inner)
	vocab := buildVocab(data.Message.Vocabulary, w)
	grammar := buildGrammar(data.Correction, w, inner)
	choices := buildChoices(data.Message, w)

	// Footer: empty line + bottom border + prompt line = 3 lines
	footer := []string{boxEmpty(w), boxBottom(w)}
	const promptLines = 1 // the "> " line after the box

	// Budget: total available lines in terminal
	fixed := len(header) + len(choices) + len(footer) + promptLines
	budget := h - fixed
	if budget < 3 {
		budget = 3 // always show at least something
	}

	// Allocate budget in priority order: narrative > npc > vocab > grammar
	sections := []*[]string{&narrative, &npc, &vocab, &grammar}
	for _, sec := range sections {
		n := len(*sec)
		if n <= budget {
			budget -= n
		} else if budget > 0 {
			*sec = (*sec)[:budget]
			budget = 0
		} else {
			*sec = nil
		}
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	writeLines(&b, header)
	writeLines(&b, narrative)
	writeLines(&b, npc)
	writeLines(&b, vocab)
	writeLines(&b, grammar)
	writeLines(&b, choices)
	writeLines(&b, footer)

	b.WriteString(showCursor)
	b.WriteString("> ")

	fmt.Print(b.String())
}

// RenderThinkingScreen shows a loading state while waiting for the API.
func RenderThinkingScreen(scenarioName, language string) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	writeLines(&b, []string{
		boxTop(w),
		boxLine(headerLine(scenarioName, language), w),
		boxDiv(w),
		boxEmpty(w),
		boxLine(dim+italic+"  Thinking..."+reset, w),
		boxEmpty(w),
		boxBottom(w),
	})

	fmt.Print(b.String())
}

// lastParsedMsg holds the most recent successfully parsed GameMessage during
// streaming, so that when incremental parsing fails on a new token we keep
// showing the last good state instead of dumping raw JSON.
var lastParsedMsg *api.GameMessage

// extractJSONString extracts the value of a JSON string field from raw text.
// It looks for "key":"..." and returns the unescaped content, handling
// escape sequences. Returns the value and true if found, even if the
// closing quote hasn't arrived yet (partial string).
func extractJSONString(raw, key string) (string, bool) {
	needle := `"` + key + `":"`
	idx := strings.Index(raw, needle)
	if idx < 0 {
		return "", false
	}
	start := idx + len(needle)

	var b strings.Builder
	escaped := false
	for i := start; i < len(raw); i++ {
		c := raw[i]
		if escaped {
			switch c {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '"':
				b.WriteByte('"')
			case '\\':
				b.WriteByte('\\')
			case '/':
				b.WriteByte('/')
			default:
				b.WriteByte('\\')
				b.WriteByte(c)
			}
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '"' {
			return b.String(), true // complete string
		}
		b.WriteByte(c)
	}
	return b.String(), true // partial string (still streaming)
}

// extractJSONArray finds a JSON array value for "key":[...] and tries to parse
// the completed items inside it. Returns whatever items parsed successfully.
func extractJSONArray[T any](raw, key string) []T {
	needle := `"` + key + `":[`
	idx := strings.Index(raw, needle)
	if idx < 0 {
		return nil
	}
	arrStart := idx + len(needle) - 1 // points at '['

	// Find matching ] by tracking nesting
	depth := 0
	inStr := false
	esc := false
	arrEnd := -1
	for i := arrStart; i < len(raw); i++ {
		c := raw[i]
		if esc {
			esc = false
			continue
		}
		if inStr {
			if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				arrEnd = i + 1
			}
		}
		if arrEnd > 0 {
			break
		}
	}

	// If array is complete, parse it directly
	if arrEnd > 0 {
		var items []T
		if json.Unmarshal([]byte(raw[arrStart:arrEnd]), &items) == nil {
			return items
		}
	}

	// Array is incomplete — try to parse individual complete objects within it
	// by finding each {...} at depth 1
	var items []T
	depth = 0
	objStart := -1
	inStr = false
	esc = false
	for i := arrStart + 1; i < len(raw); i++ {
		c := raw[i]
		if esc {
			esc = false
			continue
		}
		if inStr {
			if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		switch c {
		case '"':
			inStr = true
		case '{':
			if depth == 0 {
				objStart = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 && objStart >= 0 {
				var item T
				if json.Unmarshal([]byte(raw[objStart:i+1]), &item) == nil {
					items = append(items, item)
				}
				objStart = -1
			}
		}
	}
	return items
}

// tryParsePartialJSON extracts fields directly from the raw JSON string
// as it streams in, without requiring valid JSON.
func tryParsePartialJSON(rawJSON string) *api.GameMessage {
	raw := strings.TrimSpace(rawJSON)
	if raw == "" {
		return nil
	}

	// Try parsing complete JSON first
	var msg api.GameMessage
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		return &msg
	}

	// Extract fields progressively from the raw text
	narrative, hasNarrative := extractJSONString(raw, "narrative")
	if !hasNarrative {
		return nil
	}

	msg.Narrative = narrative
	msg.Translation, _ = extractJSONString(raw, "translation")
	msg.NPCDialog, _ = extractJSONString(raw, "npcDialog")
	msg.NPCDialogTranslation, _ = extractJSONString(raw, "npcDialogTranslation")
	msg.Choices = extractJSONArray[api.Choice](raw, "choices")
	msg.Vocabulary = extractJSONArray[api.VocabItem](raw, "vocabulary")
	msg.Finished = strings.Contains(raw, `"finished":true`)

	return &msg
}

// ResetStreamState should be called before starting a new stream to clear
// the cached last-parsed message.
func ResetStreamState() {
	lastParsedMsg = nil
}

// RenderStreamingScreen renders the TUI progressively as JSON tokens arrive.
// It attempts incremental parsing and shows whatever fields are available so far.
// When parsing fails, it keeps showing the last successfully parsed state.
func RenderStreamingScreen(scenario, lang string, correction *api.Correction, rawJSON string) {
	w, h := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}
	inner := w - 4

	header := buildHeader(scenario, lang, w)

	// Try to parse; if it fails, reuse the last good parse
	parsed := tryParsePartialJSON(rawJSON)
	if parsed != nil {
		lastParsedMsg = parsed
	}
	msg := lastParsedMsg

	var narrative []string
	var npc []string
	var vocab []string
	var grammar []string

	if msg != nil && msg.Narrative != "" {
		narrative = append(narrative, boxEmpty(w))
		for _, l := range wrapText(msg.Narrative, inner) {
			narrative = append(narrative, boxLine(bold+blue+l+reset, w))
		}
		if msg.Translation != "" && (appSettings == nil || !appSettings.HideTranslations) {
			for _, l := range wrapText(msg.Translation, inner) {
				narrative = append(narrative, boxLine(dim+l+reset, w))
			}
		}

		npc = buildNPC(msg, w, inner)
		vocab = buildVocab(msg.Vocabulary, w)
		grammar = buildGrammar(correction, w, inner)
	} else {
		narrative = append(narrative, boxEmpty(w))
		narrative = append(narrative, boxLine(dim+italic+"  Waiting for response..."+reset, w))
	}

	// Typing indicator
	typingLine := boxLine(dim+"  "+reset+bold+"▌"+reset, w)

	footer := []string{boxEmpty(w), boxBottom(w)}

	// Budget: total available lines in terminal
	fixed := len(header) + 1 + len(footer) // +1 for typing indicator
	budget := h - fixed
	if budget < 3 {
		budget = 3
	}

	sections := []*[]string{&narrative, &npc, &vocab, &grammar}
	for _, sec := range sections {
		n := len(*sec)
		if n <= budget {
			budget -= n
		} else if budget > 0 {
			*sec = (*sec)[:budget]
			budget = 0
		} else {
			*sec = nil
		}
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	writeLines(&b, header)
	writeLines(&b, narrative)
	writeLines(&b, npc)
	writeLines(&b, vocab)
	writeLines(&b, grammar)
	b.WriteString(typingLine)
	b.WriteByte('\n')
	writeLines(&b, footer)

	fmt.Print(b.String())
}

// RenderFinishedScreen shows the scenario completion screen.
func RenderFinishedScreen(scenarioName, language string, msg *api.GameMessage) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}
	inner := w - 4

	var b strings.Builder
	b.WriteString(clearSeq)

	lines := []string{
		boxTop(w),
		boxLine(headerLine(scenarioName, language), w),
		boxDiv(w),
		boxEmpty(w),
	}
	for _, l := range wrapText(msg.Narrative, inner) {
		lines = append(lines, boxLine(bold+blue+l+reset, w))
	}
	if appSettings == nil || !appSettings.HideTranslations {
		for _, l := range wrapText(msg.Translation, inner) {
			lines = append(lines, boxLine(dim+l+reset, w))
		}
	}
	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxEmpty(w),
		boxLine(bold+green+"  ★  Scenario Complete!  ★"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)

	fmt.Print(b.String())
}

// bannerArt returns the ASCII art lines for the Lark logo.
func bannerArt() []string {
	return []string{
		bold + cyan + "  _          _    " + reset,
		bold + cyan + " | |   __ _ | |__ " + reset,
		bold + cyan + " | |  / _` || '__|| |/ /" + reset,
		bold + cyan + " | |_| (_| || |   |   < " + reset,
		bold + cyan + " |____\\__,_||_|   |_|\\_\\" + reset,
	}
}

// RenderBanner renders a static title screen (used as a loading state).
func RenderBanner() {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	lines := []string{boxTop(w), boxEmpty(w)}
	for _, a := range bannerArt() {
		lines = append(lines, boxLine(a, w))
	}
	lines = append(lines,
		boxEmpty(w),
		boxLine(dim+"  A text-adventure language learning game"+reset, w),
		boxEmpty(w),
		boxLine(dim+italic+"  Connecting..."+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)

	fmt.Print(b.String())
}

// formatLangItem formats a single language item for two-column display,
// padded to colWidth visible characters.
func formatLangItem(languages []string, idx, cursorIdx, colWidth int) string {
	if idx >= len(languages) {
		return padRight("", colWidth)
	}
	name := languages[idx]
	num := idx + 1
	if idx == cursorIdx {
		content := fmt.Sprintf("  %s%s▸ %s%s", highlight+bold, white, name, reset)
		return padRight(content, colWidth)
	}
	content := fmt.Sprintf("  %s%d)%s %s", bold+green, num, reset, name)
	return padRight(content, colWidth)
}

// RenderBannerLanguages renders the banner with an integrated two-column
// language selector (popular languages + "Other Languages" option).
func RenderBannerLanguages(languages []string, cursorIdx int) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}
	inner := w - 4

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	lines := []string{boxTop(w), boxEmpty(w)}
	for _, a := range bannerArt() {
		lines = append(lines, boxLine(a, w))
	}
	lines = append(lines,
		boxEmpty(w),
		boxLine(dim+"  A text-adventure language learning game"+reset, w),
		boxEmpty(w),
		boxDiv(w),
		boxLine(bold+cyan+"  Choose a Language"+reset, w),
		boxEmpty(w),
	)

	// Two-column layout
	colWidth := inner / 2
	numRows := (len(languages) + 1) / 2
	for row := 0; row < numRows; row++ {
		leftIdx := row * 2
		rightIdx := row * 2 + 1
		left := formatLangItem(languages, leftIdx, cursorIdx, colWidth)
		right := ""
		if rightIdx < len(languages) {
			right = formatLangItem(languages, rightIdx, cursorIdx, colWidth)
		}
		lines = append(lines, boxLine(left+right, w))
	}

	// "Other Languages" option
	otherLabel := fmt.Sprintf("  %s%d)%s %s", bold+green, len(languages)+1, reset, "Other Languages...")
	if cursorIdx == len(languages) {
		otherLabel = fmt.Sprintf("  %s%s▸ %s%s", highlight+bold, white, "Other Languages...", reset)
	}
	lines = append(lines, boxLine(otherLabel, w))

	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+bold+"↑↓←→"+reset+dim+" select  ·  "+bold+"Enter"+reset+dim+" confirm  ·  "+bold+"s"+reset+dim+" settings"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	fmt.Print(b.String())
}

// RenderAllLanguagesPage renders a paginated two-column list of all languages.
func RenderAllLanguagesPage(languages []string, cursorIdx, pageIdx int) {
	w, h := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}
	inner := w - 4

	itemsPerPage := allLangsPerPage(h)
	pages := totalPages(len(languages), itemsPerPage)
	if pageIdx >= pages {
		pageIdx = pages - 1
	}

	start := pageIdx * itemsPerPage
	end := start + itemsPerPage
	if end > len(languages) {
		end = len(languages)
	}
	pageLangs := languages[start:end]

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	pageInfo := fmt.Sprintf("  %d/%d", pageIdx+1, pages)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  All Languages"+pageInfo+reset, w),
		boxDiv(w),
		boxEmpty(w),
	}

	colWidth := inner / 2
	numRows := (len(pageLangs) + 1) / 2
	for row := 0; row < numRows; row++ {
		leftIdx := row * 2
		rightIdx := row * 2 + 1
		globalLeft := start + leftIdx
		globalRight := start + rightIdx
		left := formatLangItem(languages, globalLeft, cursorIdx, colWidth)
		right := ""
		if rightIdx < len(pageLangs) {
			right = formatLangItem(languages, globalRight, cursorIdx, colWidth)
		}
		lines = append(lines, boxLine(left+right, w))
	}

	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+bold+"↑↓←→"+reset+dim+" select  ·  "+bold+"PgUp/PgDn"+reset+dim+" page  ·  "+bold+"Enter"+reset+dim+" confirm  ·  "+bold+"Esc"+reset+dim+" back"+reset, w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	fmt.Print(b.String())
}

// difficultyColor returns the ANSI color for a difficulty level.
func difficultyColor(d api.Difficulty) string {
	switch d {
	case api.DifficultyBeginner:
		return green
	case api.DifficultyIntermediate:
		return yellow
	case api.DifficultyAdvanced:
		return red
	default:
		return dim
	}
}

// allLangsPerPage calculates how many language items (two-column) fit on one page.
func allLangsPerPage(termHeight int) int {
	// Overhead: boxTop + header + boxDiv + boxEmpty(top) + boxEmpty(bot) + boxDiv + nav + boxBottom + 2 extra
	const overhead = 10
	avail := termHeight - overhead
	if avail < 3 {
		avail = 3
	}
	return avail * 2 // two columns per row
}

// scenariosPerPage calculates how many scenarios fit on one page given
// the terminal height. Each scenario takes 2 lines (name + description).
// Overhead: top + header + divider + empty(top) + empty(bot) + divider + nav + empty + bottom + prompt = 10 lines.
func scenariosPerPage(termHeight int) int {
	const overhead = 10
	avail := termHeight - overhead
	perItem := 2 // name line + description line
	count := avail / perItem
	if count < 3 {
		count = 3
	}
	return count
}

// totalPages returns how many pages are needed for n scenarios.
func totalPages(n, perPage int) int {
	if n <= 0 {
		return 1
	}
	return (n + perPage - 1) / perPage
}

// highlight is the ANSI sequence for reverse video (selected item).
const highlight = "\033[7m"

// RenderListPage renders a generic single-page arrow-key selector.
// title is shown in the header, items are the choices, cursorIdx is highlighted.
func RenderListPage(title string, items []string, cursorIdx int) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  "+title+reset, w),
		boxDiv(w),
		boxEmpty(w),
	}
	for i, item := range items {
		if i == cursorIdx {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s▸ %s%s", highlight+bold, white, item, reset), w),
			)
		} else {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%d)%s %s", bold+green, i+1, reset, item), w),
			)
		}
	}
	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+bold+"↑↓"+reset+dim+" select  ·  "+bold+"Enter"+reset+dim+" to confirm  ·  "+bold+"Esc"+reset+dim+" back"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	fmt.Print(b.String())
}

// RenderScenarioPage renders a page of the scenario selection menu.
// pageIdx is 0-based, scenarios is the full list.
// cursorIdx is the index within the current page of the highlighted scenario,
// or -1 if nothing is highlighted.
// completedSet and langCode are used to show checkmarks on finished scenarios.
func RenderScenarioPage(scenarios []api.Scenario, pageIdx, cursorIdx int, completedSet map[string]bool, langCode string) {
	w, h := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	perPage := scenariosPerPage(h)
	pages := totalPages(len(scenarios), perPage)
	if pageIdx >= pages {
		pageIdx = pages - 1
	}

	start := pageIdx * perPage
	end := start + perPage
	if end > len(scenarios) {
		end = len(scenarios)
	}
	pageScenarios := scenarios[start:end]

	var b strings.Builder
	b.WriteString(clearSeq)

	pageInfo := fmt.Sprintf("  %d/%d", pageIdx+1, pages)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  Choose a Scenario"+reset+dim+pageInfo+reset, w),
		boxDiv(w),
		boxEmpty(w),
	}
	for i, s := range pageScenarios {
		num := start + i + 1
		dc := difficultyColor(s.Difficulty)
		selected := i == cursorIdx

		// Check completion
		check := "  "
		key := saveKey(s.ID, langCode)
		if completedSet[key] {
			check = green + "✓ " + reset
		}

		if selected {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s%s▸ %2d) %-24s [%s]%s",
					check, highlight+bold, white, num, s.Name, s.Difficulty, reset), w),
				boxLine(fmt.Sprintf("  %s       %s%s",
					highlight, s.Description, reset), w),
			)
		} else {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s%2d)%s %s%-24s%s %s[%s]%s",
					check, bold+green, num, reset, cyan, s.Name, reset, dc, s.Difficulty, reset), w),
				boxLine("       "+dim+s.Description+reset, w),
			)
		}
	}
	lines = append(lines, boxEmpty(w), boxDiv(w))

	// Navigation hints
	var navParts []string
	if pageIdx > 0 {
		navParts = append(navParts, bold+"←"+reset+dim+" prev")
	}
	if pageIdx < pages-1 {
		navParts = append(navParts, bold+"→"+reset+dim+" next")
	}
	navParts = append(navParts, bold+"↑↓"+reset+dim+" select")
	navParts = append(navParts, bold+"Esc"+reset+dim+" back")
	navHint := "  " + dim + strings.Join(navParts, "  ·  ") + reset
	navHint += dim + "  ·  " + reset
	navHint += dim + italic + "or type a custom scenario" + reset

	lines = append(lines,
		boxLine(navHint, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	b.WriteString("\n> ")
	fmt.Print(b.String())
}

// RenderContinuePrompt renders a prompt asking if the user wants to continue
// a previous session or start fresh.
func RenderContinuePrompt(scenarioName string, items []string, cursorIdx int) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  "+scenarioName+reset, w),
		boxDiv(w),
		boxEmpty(w),
		boxLine(bold+yellow+"  You have a saved session for this scenario."+reset, w),
		boxEmpty(w),
	}
	for i, item := range items {
		if i == cursorIdx {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s▸ %s%s", highlight+bold, white, item, reset), w),
			)
		} else {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%d)%s %s", bold+green, i+1, reset, item), w),
			)
		}
	}
	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+bold+"↑↓"+reset+dim+" select  ·  "+bold+"Enter"+reset+dim+" to confirm  ·  "+bold+"Esc"+reset+dim+" back"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	fmt.Print(b.String())
}

// RenderCustomScenarioPage renders the custom scenario creation screen
// with writing tips and a text input prompt.
func RenderCustomScenarioPage() {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  Custom Scenario"+reset, w),
		boxDiv(w),
		boxEmpty(w),
		boxLine(bold+yellow+"  Tips for a great scenario:"+reset, w),
		boxEmpty(w),
		boxLine("  "+cyan+"1."+reset+" Set a scene with a location, goal, and an NPC to talk to", w),
		boxLine("     "+dim+"e.g. \"Haggle with a grumpy merchant at a Moroccan souk\""+reset, w),
		boxEmpty(w),
		boxLine("  "+cyan+"2."+reset+" Pick a situation that forces back-and-forth conversation", w),
		boxLine("     "+dim+"e.g. \"Your train is cancelled and you need alternatives\""+reset, w),
		boxEmpty(w),
		boxLine("  "+cyan+"3."+reset+" Mention the language level you want", w),
		boxLine("     "+dim+"e.g. \"Use only simple beginner-friendly vocabulary\""+reset, w),
		boxEmpty(w),
		boxLine("  "+cyan+"4."+reset+" Add a fun twist to keep things interesting", w),
		boxLine("     "+dim+"e.g. \"Order food at a medieval banquet as a time traveler\""+reset, w),
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+"Describe your scenario below, or press "+bold+"Esc"+reset+dim+" to go back"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	}
	writeLines(&b, lines)

	b.WriteString(showCursor)
	b.WriteString("> ")
	fmt.Print(b.String())
}

// settingsLabels are the display names for each setting toggle.
var settingsLabels = []string{
	"Show translations",
	"Show dialog choices",
	"Show vocabulary hints",
	"Show grammar corrections",
	"Explanation language",
}

// explanationLangs is the cycle list for the explanation language setting.
var explanationLangs = []string{
	"English", "Español", "Français", "Deutsch",
	"日本語", "中文", "Português", "한국어",
}

// explanationLangDisplay returns the current explanation language display name.
func explanationLangDisplay(s *Settings) string {
	if s.ExplanationLang == "" {
		return "English"
	}
	return s.ExplanationLang
}

// cycleExplanationLang advances to the next language in the cycle.
func cycleExplanationLang(s *Settings) {
	cur := explanationLangDisplay(s)
	for i, l := range explanationLangs {
		if l == cur {
			s.ExplanationLang = explanationLangs[(i+1)%len(explanationLangs)]
			return
		}
	}
	s.ExplanationLang = explanationLangs[0]
}

// settingValue returns whether the i-th setting is enabled (shown).
// For the explanation language row (index 4), this always returns true (not a toggle).
func settingValue(s *Settings, i int) bool {
	switch i {
	case 0:
		return !s.HideTranslations
	case 1:
		return !s.HideChoices
	case 2:
		return !s.HideVocabulary
	case 3:
		return !s.HideGrammar
	case 4:
		return true // not a boolean toggle
	}
	return true
}

// toggleSetting flips the i-th setting. For index 4 (explanation language),
// it cycles through the language list instead.
func toggleSetting(s *Settings, i int) {
	switch i {
	case 0:
		s.HideTranslations = !s.HideTranslations
	case 1:
		s.HideChoices = !s.HideChoices
	case 2:
		s.HideVocabulary = !s.HideVocabulary
	case 3:
		s.HideGrammar = !s.HideGrammar
	case 4:
		cycleExplanationLang(s)
	}
}

// RenderSettingsPage renders the settings toggle screen.
func RenderSettingsPage(settings *Settings, cursorIdx int) {
	w, _ := getTermSize()
	if w > 120 {
		w = 120
	}
	if w < 40 {
		w = 40
	}

	var b strings.Builder
	b.WriteString(clearSeq)
	b.WriteString(hideCursor)

	lines := []string{
		boxTop(w),
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  Settings"+reset, w),
		boxDiv(w),
		boxEmpty(w),
	}
	for i, label := range settingsLabels {
		if i == 4 {
			// Explanation language — cycle picker, not a toggle
			langVal := explanationLangDisplay(settings)
			tag := bold + cyan + "[" + langVal + "]" + reset
			if i == cursorIdx {
				lines = append(lines,
					boxLine(fmt.Sprintf("  %s%s▸ %-28s %s%s", highlight+bold, white, label, reset+highlight+bold+cyan, "["+langVal+"]"+reset), w),
				)
			} else {
				lines = append(lines,
					boxLine(fmt.Sprintf("    %-28s %s", label, tag), w),
				)
			}
			continue
		}
		on := settingValue(settings, i)
		var tag string
		if on {
			tag = bold + green + "[ON]" + reset
		} else {
			tag = bold + red + "[OFF]" + reset
		}
		if i == cursorIdx {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s▸ %-28s %s%s", highlight+bold, white, label, reset+highlight+bold, func() string {
					if on {
						return green + "[ON]"
					}
					return red + "[OFF]"
				}()+reset), w),
			)
		} else {
			lines = append(lines,
				boxLine(fmt.Sprintf("    %-28s %s", label, tag), w),
			)
		}
	}
	lines = append(lines,
		boxEmpty(w),
		boxDiv(w),
		boxLine("  "+dim+bold+"↑↓"+reset+dim+" select  ·  "+bold+"Enter"+reset+dim+" toggle  ·  "+bold+"Esc"+reset+dim+" back"+reset, w),
		boxEmpty(w),
		boxBottom(w),
	)
	writeLines(&b, lines)
	fmt.Print(b.String())
}

// PrintError displays an error message inline.
func PrintError(msg string) {
	fmt.Printf("\n%sError: %s%s\n> ", red, msg, reset)
}

// PrintWarning displays a warning message inline (yellow, no "Error" prefix).
func PrintWarning(msg string) {
	fmt.Printf("\n%s%s%s\n> ", yellow, msg, reset)
}
