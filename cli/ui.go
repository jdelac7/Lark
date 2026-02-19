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
	for _, l := range wrapText(msg.Translation, inner) {
		lines = append(lines, boxLine(dim+l+reset, w))
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
	for _, l := range wrapText(msg.NPCDialogTranslation, inner-7) {
		lines = append(lines, boxLine("       "+dim+l+reset, w))
	}
	return lines
}

func buildVocab(vocab []api.VocabItem, w int) []string {
	if len(vocab) == 0 {
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
	if c == nil {
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
	var lines []string
	lines = append(lines, boxDiv(w))
	for i, c := range msg.Choices {
		l := fmt.Sprintf("  %s%d)%s %s%s%s  %s(%s)%s",
			bold+green, i+1, reset,
			cyan, c.Text, reset,
			dim, c.Translation, reset)
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

	// Allocate budget in priority order: narrative > npc > grammar > vocab
	sections := []*[]string{&narrative, &npc, &grammar, &vocab}
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

// tryParsePartialJSON attempts to parse the accumulated raw JSON tokens into
// a partial GameMessage. It tries the raw text first, then appends closing
// brackets as a best-effort heuristic for incomplete JSON.
func tryParsePartialJSON(rawJSON string) *api.GameMessage {
	raw := strings.TrimSpace(rawJSON)
	if raw == "" {
		return nil
	}

	// Try parsing as-is first (maybe the JSON is already complete)
	var msg api.GameMessage
	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
		return &msg
	}

	// Best-effort: try closing open brackets/braces
	suffixes := []string{
		`"}`,
		`}`,
		`"]}`,
		`]}`,
		`"]}]}`,
		`]}]}`,
		`null}`,
		`false}`,
	}
	for _, suffix := range suffixes {
		attempt := raw + suffix
		if err := json.Unmarshal([]byte(attempt), &msg); err == nil {
			return &msg
		}
	}

	return nil
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
		if msg.Translation != "" {
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

	sections := []*[]string{&narrative, &npc, &grammar, &vocab}
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
	for _, l := range wrapText(msg.Translation, inner) {
		lines = append(lines, boxLine(dim+l+reset, w))
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

// RenderBanner renders the title screen.
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

	art := []string{
		bold + cyan + "  _          _    " + reset,
		bold + cyan + " | |   __ _ | |__ " + reset,
		bold + cyan + " | |  / _` || '__|| |/ /" + reset,
		bold + cyan + " | |_| (_| || |   |   < " + reset,
		bold + cyan + " |____\\__,_||_|   |_|\\_\\" + reset,
	}

	lines := []string{boxTop(w), boxEmpty(w)}
	for _, a := range art {
		lines = append(lines, boxLine(a, w))
	}
	lines = append(lines,
		boxEmpty(w),
		boxLine(dim+"  A text-adventure language learning game"+reset, w),
		boxEmpty(w),
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

// RenderScenarioPage renders a page of the scenario selection menu.
// pageIdx is 0-based, scenarios is the full list.
// cursorIdx is the index within the current page of the highlighted scenario,
// or -1 if nothing is highlighted.
func RenderScenarioPage(scenarios []api.Scenario, pageIdx, cursorIdx int) {
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
		if selected {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%s▸ %2d) %-24s [%s]%s",
					highlight+bold, white, num, s.Name, s.Difficulty, reset), w),
				boxLine(fmt.Sprintf("  %s     %s%s",
					highlight, s.Description, reset), w),
			)
		} else {
			lines = append(lines,
				boxLine(fmt.Sprintf("  %s%2d)%s %s%-24s%s %s[%s]%s",
					bold+green, num, reset, cyan, s.Name, reset, dc, s.Difficulty, reset), w),
				boxLine("     "+dim+s.Description+reset, w),
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

// RenderLanguageList renders the language selection menu.
func RenderLanguageList(languages []api.Language) {
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
		boxLine(bold+cyan+"  Lark"+reset+dim+"  ·  Choose a Language"+reset, w),
		boxDiv(w),
		boxEmpty(w),
	}
	for i, l := range languages {
		lines = append(lines,
			boxLine(fmt.Sprintf("  %s%d)%s %s", bold+green, i+1, reset, l.Name), w),
		)
	}
	lines = append(lines, boxEmpty(w), boxBottom(w))
	writeLines(&b, lines)
	b.WriteString("\n> ")
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
