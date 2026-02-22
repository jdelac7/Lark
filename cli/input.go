package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/joshburnsxyz/lark/api"
	"golang.org/x/term"
)

// ctrlC is a sentinel value returned by readInput when the user presses Ctrl-C.
const ctrlC = "\x03"

var stdinFd = int(os.Stdin.Fd())

// enableRawMode puts the terminal into raw mode and returns restore state.
func enableRawMode() (*term.State, error) {
	return term.MakeRaw(stdinFd)
}

// restoreTerminal restores the terminal to its previous state.
func restoreTerminal(state *term.State) {
	term.Restore(stdinFd, state)
}

// readByte reads a single byte from stdin.
func readByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := os.Stdin.Read(buf)
	return buf[0], err
}

// readInput reads a line of input in raw mode.
// Returns ctrlC sentinel if the user presses Ctrl-C.
func readInput() string {
	oldState, err := enableRawMode()
	if err != nil {
		return readLineFallback()
	}
	defer restoreTerminal(oldState)

	var buf []byte

	for {
		b, err := readByte()
		if err != nil {
			return string(buf)
		}

		switch {
		case b == '\r' || b == '\n':
			fmt.Print("\r\n")
			return string(buf)

		case b == 127 || b == 8: // backspace / delete
			if len(buf) > 0 {
				_, size := utf8.DecodeLastRune(buf)
				buf = buf[:len(buf)-size]
				fmt.Print("\b \b")
			}

		case b == 3: // Ctrl-C
			fmt.Print("\r\n")
			return ctrlC

		case b == '\033': // escape sequence — consume and discard
			b2, err := readByte()
			if err != nil {
				continue
			}
			if b2 == '[' || b2 == 'O' {
				for {
					b3, err := readByte()
					if err != nil {
						break
					}
					if b3 >= 0x40 && b3 <= 0x7E {
						break
					}
				}
			}

		case b >= 0x20 && b < 0x7F: // printable ASCII
			buf = append(buf, b)
			fmt.Print(string(rune(b)))

		case b >= 0xC0: // start of multi-byte UTF-8
			var seq []byte
			seq = append(seq, b)
			var n int
			switch {
			case b&0xE0 == 0xC0:
				n = 2
			case b&0xF0 == 0xE0:
				n = 3
			case b&0xF8 == 0xF0:
				n = 4
			default:
				continue
			}
			for i := 1; i < n; i++ {
				cb, err := readByte()
				if err != nil {
					break
				}
				seq = append(seq, cb)
			}
			if utf8.Valid(seq) {
				buf = append(buf, seq...)
				fmt.Print(string(seq))
			}
		}
	}
}

// readLineFallback is a simple line reader without raw mode (last resort).
func readLineFallback() string {
	var buf []byte
	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			return string(buf)
		}
		if b[0] == '\n' {
			return string(buf)
		}
		buf = append(buf, b[0])
	}
}

// ReadChoice reads a menu selection (1-based) or returns -1 for free text.
// Returns choiceIndex=-2 on Ctrl-C (quit signal).
// numChoices is the count of AI-provided choices; option numChoices+1 is "write your own".
// When dialog choices are hidden in settings, goes straight to free-text input.
func ReadChoice(numChoices int) (choiceIndex int, freeText string) {
	// Free-text only mode when choices are hidden
	if appSettings != nil && appSettings.HideChoices {
		for {
			input := readInput()
			if input == ctrlC {
				return -2, ""
			}
			if input == "" {
				fmt.Print("> ")
				continue
			}
			return -1, input
		}
	}

	freeOption := numChoices + 1
	for {
		input := readInput()
		if input == ctrlC {
			return -2, ""
		}
		if input == "" {
			fmt.Print("> ")
			continue
		}

		if n, err := strconv.Atoi(input); err == nil {
			if n >= 1 && n <= numChoices {
				return n - 1, ""
			}
			if n == freeOption {
				// Prompt for free text
				fmt.Print(dim + "  Write in the target language: " + reset)
				text := readInput()
				if text == ctrlC {
					return -2, ""
				}
				if text == "" {
					fmt.Print("> ")
					continue
				}
				return -1, text
			}
			PrintError(fmt.Sprintf("Please enter 1-%d", freeOption))
			continue
		}

		// Raw text typed directly (not a number) — treat as free text
		return -1, input
	}
}

// Arrow/key event constants for the scenario selector.
const (
	keyUp    = "\x1b[A"
	keyDown  = "\x1b[B"
	keyRight = "\x1b[C"
	keyLeft  = "\x1b[D"
	keyPgUp  = "\x1b[5~"
	keyPgDn  = "\x1b[6~"
	keyEnter = "\r"
	keyEsc   = "\x1b"
)

// byteAvailable is defined in input_select_*.go (platform-specific).

// readKeyEvent reads a single key event in raw mode.
// Returns arrow sentinels for arrow keys, keyEnter for bare Enter,
// ctrlC for Ctrl-C, or the first character typed as a string.
func readKeyEvent() string {
	oldState, err := enableRawMode()
	if err != nil {
		return readLineFallback()
	}
	defer restoreTerminal(oldState)

	b, err := readByte()
	if err != nil {
		return ""
	}

	switch {
	case b == '\r' || b == '\n':
		return keyEnter
	case b == 3:
		return ctrlC
	case b == '\033':
		// Bare Escape (no following bytes within 50ms)
		if !byteAvailable(50 * time.Millisecond) {
			return keyEsc
		}
		b2, err := readByte()
		if err != nil {
			return keyEsc
		}
		if b2 == '[' || b2 == 'O' {
			b3, err := readByte()
			if err != nil {
				return ""
			}
			if b2 == '[' {
				switch b3 {
				case 'A':
					return keyUp
				case 'B':
					return keyDown
				case 'C':
					return keyRight
				case 'D':
					return keyLeft
				case '5':
					// PgUp: \x1b[5~
					if b4, err := readByte(); err == nil && b4 == '~' {
						return keyPgUp
					}
					return ""
				case '6':
					// PgDn: \x1b[6~
					if b4, err := readByte(); err == nil && b4 == '~' {
						return keyPgDn
					}
					return ""
				}
			}
			// Consume rest of unknown escape sequence
			if b3 < 0x40 || b3 > 0x7E {
				for {
					b4, err := readByte()
					if err != nil {
						break
					}
					if b4 >= 0x40 && b4 <= 0x7E {
						break
					}
				}
			}
			return ""
		}
		return ""
	case b >= 0x20 && b < 0x7F:
		return string(rune(b))
	}
	return ""
}

// readInputSeeded enters raw mode with the first character already in
// the buffer (and echoed), then reads the rest of the line normally.
// This ensures backspace can delete every character including the first.
func readInputSeeded(first string) string {
	oldState, err := enableRawMode()
	if err != nil {
		return readLineFallback()
	}
	defer restoreTerminal(oldState)

	buf := []byte(first)
	fmt.Print(first)

	for {
		b, err := readByte()
		if err != nil {
			return string(buf)
		}

		switch {
		case b == '\r' || b == '\n':
			fmt.Print("\r\n")
			return string(buf)

		case b == 127 || b == 8: // backspace
			if len(buf) > 0 {
				_, size := utf8.DecodeLastRune(buf)
				buf = buf[:len(buf)-size]
				fmt.Print("\b \b")
			}
			// Buffer empty after backspace — return so caller can
			// resume arrow-key navigation.
			if len(buf) == 0 {
				return ""
			}

		case b == 3: // Ctrl-C
			fmt.Print("\r\n")
			return ctrlC

		case b == '\033': // escape sequence — consume and discard
			b2, err := readByte()
			if err != nil {
				continue
			}
			if b2 == '[' || b2 == 'O' {
				for {
					b3, err := readByte()
					if err != nil {
						break
					}
					if b3 >= 0x40 && b3 <= 0x7E {
						break
					}
				}
			}

		case b >= 0x20 && b < 0x7F:
			buf = append(buf, b)
			fmt.Print(string(rune(b)))

		case b >= 0xC0: // multi-byte UTF-8
			var seq []byte
			seq = append(seq, b)
			var n int
			switch {
			case b&0xE0 == 0xC0:
				n = 2
			case b&0xF0 == 0xE0:
				n = 3
			case b&0xF8 == 0xF0:
				n = 4
			default:
				continue
			}
			for i := 1; i < n; i++ {
				cb, err := readByte()
				if err != nil {
					break
				}
				seq = append(seq, cb)
			}
			if utf8.Valid(seq) {
				buf = append(buf, seq...)
				fmt.Print(string(seq))
			}
		}
	}
}

// ReadScenarioChoice reads a scenario selection with paginated display.
// Supports arrow keys (up/down to highlight, left/right to page) and typing.
// completedSet marks which scenarios have been completed for the given language.
// Returns choiceIndex=-2 on Ctrl-C, choiceIndex=-1 with freeText for custom,
// choiceIndex=-3 for go back (Escape, or Left on first page).
func ReadScenarioChoice(scenarios []api.Scenario, completedSet map[string]bool, langCode string) (choiceIndex int, freeText string) {
	if len(scenarios) == 0 {
		return -2, ""
	}

	_, h := getTermSize()
	perPage := scenariosPerPage(h)
	pages := totalPages(len(scenarios), perPage)
	pageIdx := 0
	cursorIdx := -1 // -1 = nothing highlighted
	total := len(scenarios)

	RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)

	for {
		key := readKeyEvent()

		// Compute items on current page
		start := pageIdx * perPage
		end := start + perPage
		if end > total {
			end = total
		}
		pageCount := end - start

		switch key {
		case ctrlC:
			return -2, ""

		case keyEsc:
			return -3, ""

		case keyDown:
			if cursorIdx < 0 {
				cursorIdx = 0
			} else if cursorIdx < pageCount-1 {
				cursorIdx++
			}
			RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			continue

		case keyUp:
			if cursorIdx < 0 {
				cursorIdx = pageCount - 1
			} else if cursorIdx > 0 {
				cursorIdx--
			}
			RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			continue

		case keyRight:
			if pageIdx < pages-1 {
				pageIdx++
				cursorIdx = -1
				RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			}
			continue

		case keyLeft:
			if pageIdx > 0 {
				pageIdx--
				cursorIdx = -1
				RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			}
			continue

		case keyEnter:
			if cursorIdx >= 0 {
				return start + cursorIdx, ""
			}
			fmt.Print("> ")
			continue

		case "":
			continue
		}

		// User typed a character — clear highlight and read the full line.
		cursorIdx = -1
		RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)

		input := readInputSeeded(key)
		if input == ctrlC {
			return -2, ""
		}
		if input == "" {
			RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			continue
		}

		// Numeric selection (global numbering)
		if n, err := strconv.Atoi(input); err == nil {
			if n >= 1 && n <= total {
				return n - 1, ""
			}
			// Re-render clean, then show error
			RenderScenarioPage(scenarios, pageIdx, cursorIdx, completedSet, langCode)
			PrintError(fmt.Sprintf("Please enter 1-%d, or type a scenario", total))
			continue
		}

		// Free text = custom scenario
		return -1, input
	}
}

// ReadMenuChoice reads a numeric menu selection (1-based).
// Returns -1 on Ctrl-C (quit signal).
func ReadMenuChoice(max int) int {
	for {
		input := readInput()
		if input == ctrlC {
			return -1
		}
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > max {
			PrintError(fmt.Sprintf("Please enter a number between 1 and %d", max))
			continue
		}
		return n - 1
	}
}

// ReadListChoice renders a single-page arrow-key selector using the given
// render function and returns the selected index.
// Returns -1 on Ctrl-C (quit), -2 on Escape (go back).
// renderFn is called with the current cursor index to redraw.
func ReadListChoice(count int, renderFn func(cursor int)) int {
	cursor := 0
	renderFn(cursor)

	for {
		key := readKeyEvent()
		switch key {
		case ctrlC:
			return -1
		case keyEsc:
			return -2
		case keyDown:
			if cursor < count-1 {
				cursor++
			}
			renderFn(cursor)
		case keyUp:
			if cursor > 0 {
				cursor--
			}
			renderFn(cursor)
		case keyEnter:
			return cursor
		case "":
			continue
		default:
			// Number key shortcut
			if n, err := strconv.Atoi(key); err == nil && n >= 1 && n <= count {
				return n - 1
			}
		}
	}
}

// ReadBannerLanguageChoice is a two-column arrow-key selector for the
// banner language screen. Up/Down move within a column, Left/Right switch
// columns. Returns the selected index, -1 on Ctrl-C/Escape (quit),
// or -2 when the user presses 's' (open settings).
func ReadBannerLanguageChoice(count int, renderFn func(cursor int)) int {
	// count is number of popular languages; total items = count + 1 ("Other Languages")
	total := count + 1
	cursor := 0
	renderFn(cursor)

	for {
		key := readKeyEvent()
		switch key {
		case ctrlC, keyEsc:
			return -1
		case keyDown:
			// In two-column layout, down moves by 2; last item is "Other" at index=count
			next := cursor + 2
			if next >= count {
				next = count // jump to "Other Languages"
			}
			if next < total {
				cursor = next
			}
			renderFn(cursor)
		case keyUp:
			if cursor == count {
				// From "Other", go to last row of the grid
				cursor = count - 1
				if cursor < 0 {
					cursor = 0
				}
			} else if cursor-2 >= 0 {
				cursor -= 2
			}
			renderFn(cursor)
		case keyRight:
			if cursor < count && cursor%2 == 0 && cursor+1 < count {
				cursor++
			}
			renderFn(cursor)
		case keyLeft:
			if cursor < count && cursor%2 == 1 {
				cursor--
			}
			renderFn(cursor)
		case keyEnter:
			return cursor
		case "":
			continue
		default:
			if key == "s" || key == "S" {
				return -2
			}
			if n, err := strconv.Atoi(key); err == nil && n >= 1 && n <= total {
				return n - 1
			}
		}
	}
}

// ReadAllLanguagesChoice runs the paginated all-languages selector.
// Returns selected index into the full languages slice, or -1 on Esc.
func ReadAllLanguagesChoice(count int, renderFn func(cursor, page int)) int {
	_, h := getTermSize()
	itemsPerPage := allLangsPerPage(h)
	pages := totalPages(count, itemsPerPage)
	cursor := 0
	page := 0
	renderFn(cursor, page)

	for {
		key := readKeyEvent()
		switch key {
		case ctrlC, keyEsc:
			return -1
		case keyDown:
			if cursor+2 < count {
				cursor += 2
			}
			page = cursor / itemsPerPage
			renderFn(cursor, page)
		case keyUp:
			if cursor-2 >= 0 {
				cursor -= 2
			}
			page = cursor / itemsPerPage
			renderFn(cursor, page)
		case keyRight:
			if cursor%2 == 0 && cursor+1 < count {
				cursor++
			}
			page = cursor / itemsPerPage
			renderFn(cursor, page)
		case keyLeft:
			if cursor%2 == 1 {
				cursor--
			}
			page = cursor / itemsPerPage
			renderFn(cursor, page)
		case keyPgDn:
			if page < pages-1 {
				page++
				cursor = page * itemsPerPage
			}
			renderFn(cursor, page)
		case keyPgUp:
			if page > 0 {
				page--
				cursor = page * itemsPerPage
			}
			renderFn(cursor, page)
		case keyEnter:
			return cursor
		case "":
			continue
		}
	}
}

// ReadSettings runs the settings toggle screen. Arrow keys navigate,
// Enter toggles the selected setting, Esc returns to the previous screen.
// saveFn is called after each toggle to persist changes.
func ReadSettings(settings *Settings, saveFn func()) {
	cursor := 0
	RenderSettingsPage(settings, cursor)

	for {
		key := readKeyEvent()
		switch key {
		case ctrlC, keyEsc:
			return
		case keyDown:
			if cursor < len(settingsLabels)-1 {
				cursor++
			}
		case keyUp:
			if cursor > 0 {
				cursor--
			}
		case keyEnter:
			toggleSetting(settings, cursor)
			saveFn()
		case "":
			continue
		}
		RenderSettingsPage(settings, cursor)
	}
}

// ReadCustomScenario renders the custom scenario page and reads the user's
// scenario description. Returns the text, or "" if the user pressed Escape
// (go back), or ctrlC on quit.
func ReadCustomScenario() string {
	RenderCustomScenarioPage()

	// Read first key event to detect Escape before committing to line input
	for {
		key := readKeyEvent()
		switch key {
		case ctrlC:
			return ctrlC
		case keyEsc:
			return ""
		case keyEnter, "":
			continue
		default:
			// User started typing — read the rest of the line
			input := readInputSeeded(key)
			if input == ctrlC {
				return ctrlC
			}
			if input == "" {
				// Backspaced everything — re-render and wait again
				RenderCustomScenarioPage()
				continue
			}
			return input
		}
	}
}

// ReadLine prints a prompt and reads a line.
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	return readInput()
}
