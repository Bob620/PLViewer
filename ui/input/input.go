package input

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
	"time"
)

type Input struct {
	*tview.InputField
	Mutex                  sync.Mutex
	Entries                []string
	ticker                 *time.Ticker
	interval               time.Duration
	lastEntries            []string
	autocompleteSelectFunc func(string) string
	autocompleteFunc       func(string) []string
	cursorPos              int
	active                 bool
	entriesUpdated         bool
	blink                  bool
	label                  string
}

func MakeInput(application *tview.Application) *Input {
	input := &Input{
		InputField:  tview.NewInputField().Autocomplete(),
		Mutex:       sync.Mutex{},
		Entries:     []string{},
		lastEntries: []string{},
		cursorPos:   -1,
		ticker:      time.NewTicker(time.Second / 2),
		interval:    time.Second / 2,
	}

	input.InputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if currentText == "" {
			return nil
		}

		input.Mutex.Lock()
		defer input.Mutex.Unlock()

		if !input.GetUpdated() {
			var temp []string
			for _, item := range input.Entries {
				temp = append(temp, item)
			}
			input.ClearEntries()
			_, _, _, height := input.GetInnerRect()

			if len(temp) > height-2 {
				return append(temp[:height-2], "More...")
			}
			return temp
		}

		if input.autocompleteFunc != nil {
			go func() {
				input.SetEntries(input.autocompleteFunc(currentText))
				input.Autocomplete()
				application.Draw()
			}()
		}

		return nil
	})

	input.ticker.Stop()
	go func() {
		for {
			_ = <-input.ticker.C

			input.blink = !input.blink
			application.Draw()
		}
	}()

	return input
}

func (input *Input) GetCursorPos() int {
	text := input.GetText()
	cursorPos := input.cursorPos

	// Sanitize the cursor position
	if cursorPos > len(text) || len(text) == 0 {
		input.cursorPos = -1
		cursorPos = -1
	}

	if cursorPos == -1 {
		cursorPos = len(text)
	}

	return cursorPos
}

func (input *Input) GetUpdated() bool {
	if input.entriesUpdated {
		input.entriesUpdated = false
		return true
	}
	return false
}

func (input *Input) IsActive() bool {
	return input.active
}

func (input *Input) Activate(yes bool) {
	input.active = yes
	if yes {
		input.blink = true
		input.ticker.Reset(input.interval)
	} else {
		input.ticker.Stop()
	}
}

func (input *Input) SetAutocompleteSelectionFunc(handler func(string) string) *Input {
	input.autocompleteSelectFunc = handler
	return input
}

func (input *Input) SetAutocompleteFunc(handler func(string) []string) *Input {
	input.autocompleteFunc = handler
	return input
}

func (input *Input) GetEntries() []string {
	return input.Entries
}

func (input *Input) GetLastEntries() []string {
	return input.lastEntries
}

func (input *Input) ClearEntries() {
	input.Entries = nil
}

func (input *Input) SetEntries(entries []string) {
	input.Mutex.Lock()
	input.Entries = entries
	input.lastEntries = entries
	input.Mutex.Unlock()

}

func (input *Input) Draw(screen tcell.Screen) {
	input.InputField.Draw(screen)
	x, y, width, _ := input.GetInnerRect()
	label := input.label

	// Write Label into the input bg box
	tview.Print(screen, label, x, y, len(label), tview.AlignLeft, tcell.ColorWhite)
	x += len(label)

	text := input.GetText()
	cursorPos := input.GetCursorPos()

	// Determine cursor effects
	activeColor := tcell.ColorWhite

	cursorChar := " "
	if cursorPos < len(text) {
		cursorChar = string(text[cursorPos])
	}

	cursorEffect := ""
	if input.blink && input.active {
		cursorEffect = "[::r]"
	}

	// Write out the text before the cursor
	if cursorPos > 0 {
		tview.Print(screen, text[0:cursorPos], x, y, width, tview.AlignLeft, tcell.ColorWhite)
	}

	// Write out the cursor
	tview.Print(screen, cursorEffect+cursorChar, x+cursorPos, y, 1, tview.AlignLeft, activeColor)

	// Write out the text after the cursor
	if cursorPos < len(text) {
		tview.Print(screen, text[cursorPos+1:], x+cursorPos+1, y, width, tview.AlignLeft, tcell.ColorWhite)
	}
}

func (input *Input) HandleEvents(key *tcell.EventKey) {
	switch key.Key() {
	case tcell.KeyDEL:
		fallthrough
	case tcell.KeyDelete:
		fallthrough
	case tcell.KeyBackspace:
		text := input.GetText()
		length := len(text)

		if length == 0 {
			input.SetText("")
		} else {
			input.entriesUpdated = true
			cursorPos := input.GetCursorPos()

			newText := ""
			if cursorPos > 0 {
				newText += text[0 : cursorPos-1]
			}

			if cursorPos < len(text) {
				newText += text[cursorPos:]
				if cursorPos > 0 {
					input.cursorPos--
				}
			}

			input.SetText(newText)
		}
		break
	case tcell.KeyRune:
		input.entriesUpdated = true
		text := input.GetText()
		cursorPos := input.GetCursorPos()

		newText := ""
		if cursorPos >= 0 {
			newText += text[0:cursorPos] + string(key.Rune())
		}

		if cursorPos < len(text) {
			newText += text[cursorPos:]
			input.cursorPos++
		}

		input.SetText(newText)
		break
	case tcell.KeyUp:
		entries := input.GetLastEntries()
		if len(entries) > 1 {
			newEntries := []string{entries[len(entries)-1]}
			for _, entry := range entries[0 : len(entries)-1] {
				newEntries = append(newEntries, entry)
			}

			input.SetEntries(newEntries)
		}
		break
	case tcell.KeyDown:
		entries := input.GetLastEntries()
		if len(entries) > 1 {
			entries = append(entries, entries[0])[1:]

			input.SetEntries(entries)
		}
		break
	case tcell.KeyLeft:
		if input.cursorPos > 0 {
			input.cursorPos--
		} else if input.cursorPos == -1 {
			input.cursorPos = len(input.GetText()) - 1
		}
		break
	case tcell.KeyRight:
		if input.cursorPos != -1 {
			if input.cursorPos < len(input.GetText()) {
				input.cursorPos++
			} else {
				input.cursorPos = -1
			}
		} else {
			entries := input.GetLastEntries()
			if len(entries) > 0 {
				input.entriesUpdated = true
				if input.autocompleteSelectFunc != nil {
					input.SetText(input.autocompleteSelectFunc(entries[0]))
				} else {
					input.SetText(input.GetText() + entries[0])
				}
			}
		}
		break
	case tcell.KeyTab:
		entries := input.GetLastEntries()
		if len(entries) > 0 {
			input.entriesUpdated = true
			if input.autocompleteSelectFunc != nil {
				input.SetText(input.autocompleteSelectFunc(entries[0]))
			} else {
				input.SetText(input.GetText() + entries[0])
			}
		}
		break
	}

	input.Autocomplete()
}

func (input *Input) SetLabel(label string) *Input {
	input.label = label
	return input
}
