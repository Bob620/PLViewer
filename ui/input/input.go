package input

import (
	"PLViewer/ui/page"
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
	lastEntries            []string
	autocompleteSelectFunc func(string) string
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
	}

	go func() {
		for {
			_ = <-input.ticker.C

			input.blink = !input.blink
			application.Draw()
		}
	}()

	return input
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
}

func (input *Input) SetAutocompleteSelectionFunc(handler func(string) string) *Input {
	input.autocompleteSelectFunc = handler
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
	input.Entries = entries
	input.lastEntries = entries

}

func (input *Input) Draw(screen tcell.Screen) {
	input.InputField.Draw(screen)
	x, y, width, _ := input.GetInnerRect()
	label := input.label

	tview.Print(screen, label, x, y, len(label), tview.AlignLeft, tcell.ColorWhite)
	x += len(label)

	activeColor := tcell.ColorGreen
	if input.active {
		activeColor = tcell.ColorLightBlue
	}

	text := input.GetText()
	cursorPos := input.cursorPos

	if cursorPos > len(text) || len(text) == 0 {
		input.cursorPos = -1
		cursorPos = -1
	}

	if cursorPos == -1 {
		cursorPos = len(text)
	}

	if cursorPos > 0 {
		tview.Print(screen, text[0:cursorPos], x, y, width, tview.AlignLeft, tcell.ColorWhite)
	}

	if !input.active {
		if cursorPos < len(text) {
			tview.Print(screen, string(text[cursorPos]), x+cursorPos, y, 1, tview.AlignLeft, activeColor)
		} else {
			tview.Print(screen, "˾", x+cursorPos, y, 1, tview.AlignLeft, activeColor)
		}
	} else {
		if cursorPos < len(text) {
			if input.blink {
				tview.Print(screen, "˾", x+cursorPos, y, 1, tview.AlignLeft, activeColor)
			} else {
				tview.Print(screen, string(text[cursorPos]), x+cursorPos, y, 1, tview.AlignLeft, activeColor)
			}
		} else {
			if input.blink {
				tview.Print(screen, "˾", x+cursorPos, y, 1, tview.AlignLeft, activeColor)
			} else {
				tview.Print(screen, " ", x+cursorPos, y, 1, tview.AlignLeft, activeColor)
			}
		}
	}

	if cursorPos < len(text) {
		tview.Print(screen, text[cursorPos+1:], x+cursorPos+1, y, width, tview.AlignLeft, tcell.ColorWhite)
	}
}

func (input *Input) HandleEvents(key *tcell.EventKey, switchToPage func(*page.Page), focusHeader func(bool)) {
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
			input.SetText(text[0 : length-1])
		}
		break
	case tcell.KeyRune:
		input.entriesUpdated = true
		input.SetText(input.GetText() + string(key.Rune()))
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
