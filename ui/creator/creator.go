package creator

import (
	"PLViewer/ui/input"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type CreatorPage struct {
	*page.Page
	inputField    *input.Input
	rightElement  *tview.TextView
	bottomElement *tview.TextView
	activeElement int
}

func MakeCreatorPage(application *tview.Application) *CreatorPage {
	creator := CreatorPage{
		Page:          page.MakePage("Creator"),
		inputField:    input.MakeInput(application),
		rightElement:  tview.NewTextView(),
		bottomElement: tview.NewTextView(),
		activeElement: 1,
	}
	creator.SetOnSelect(creator.onSelect)
	creator.SetOnDeselect(creator.onDeselect)
	creator.rightElement.SetText("Test").SetBorder(true)
	creator.bottomElement.SetText("Ahhh").SetBorder(true)

	creator.inputField.SetLabel("Select a Directory to Pack: ").SetBorder(true)
	creator.inputField.SetAutocompleteSelectionFunc(func(entry string) string {
		dir, _ := path.Split(creator.inputField.GetText())
		return path.Join(dir, entry) + "/"
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		if prefix == "" {
			return nil
		}

		creator.inputField.Mutex.Lock()
		defer creator.inputField.Mutex.Unlock()

		if !creator.inputField.GetUpdated() {
			temp := creator.inputField.Entries
			creator.inputField.ClearEntries()
			return temp
		}

		go func() {
			dir, wanted := path.Split(prefix)
			files := getDir(dir)

			var auto []string
			for _, file := range files {
				if file.IsDir() && (wanted == "" || strings.Contains(strings.ToLower(file.Name()), wanted)) {
					auto = append(auto, file.Name())
				}
			}

			creator.inputField.Mutex.Lock()
			creator.inputField.SetEntries(auto)
			creator.inputField.Mutex.Unlock()

			creator.inputField.Autocomplete()
			application.Draw()
		}()

		return nil
	})

	creator.Grid().SetRows(0, 0).SetColumns(0, 0).
		AddItem(creator.inputField, 0, 0, 1, 1, 0, 0, false).
		AddItem(creator.rightElement, 0, 1, 2, 1, 0, 0, false).
		AddItem(creator.bottomElement, 1, 0, 1, 1, 0, 0, false)

	return &creator
}

func getDir(uri string) []os.FileInfo {
	files, err := ioutil.ReadDir(uri)
	if err != nil {
		return []os.FileInfo{}
	}

	return files
}

func (creator *CreatorPage) onSelect() {
	creator.setActive(1)
}

func (creator *CreatorPage) onDeselect() {
	creator.setActive(0)
}

func (creator *CreatorPage) setActive(element int) {
	creator.activeElement = element
	creator.inputField.SetBorderColor(tcell.ColorWhite)
	creator.rightElement.SetBorderColor(tcell.ColorWhite)
	creator.bottomElement.SetBorderColor(tcell.ColorWhite)

	switch element {
	case 1:
		creator.inputField.SetBorderColor(tcell.ColorBlue)
		break
	case 2:
		creator.rightElement.SetBorderColor(tcell.ColorBlue)
		break
	case 3:
		creator.bottomElement.SetBorderColor(tcell.ColorBlue)
		break
	}
}

func (creator *CreatorPage) HandleEvents(key *tcell.EventKey, switchToPage func(*page.Page), focusHeader func(bool)) {
	switch creator.activeElement {
	case 1:
		switch key.Key() {
		case tcell.KeyEscape:
			creator.inputField.Activate(false)
			break
		case tcell.KeyEnter:
			if creator.inputField.IsActive() {
				creator.inputField.Activate(false)
				creator.inputField.Autocomplete()
			} else {
				creator.inputField.Activate(true)
			}
			break
		case tcell.KeyLeft:
			if creator.inputField.IsActive() {
				creator.inputField.HandleEvents(key, switchToPage, focusHeader)
			}
			break
		case tcell.KeyRight:
			if creator.inputField.IsActive() {
				creator.inputField.HandleEvents(key, switchToPage, focusHeader)
			} else {
				creator.setActive(2)
			}
			break
		case tcell.KeyUp:
			if creator.inputField.IsActive() {
				creator.inputField.HandleEvents(key, switchToPage, focusHeader)
			} else {
				focusHeader(true)
				creator.Deselect()
			}
			break
		case tcell.KeyDown:
			if creator.inputField.IsActive() {
				creator.inputField.HandleEvents(key, switchToPage, focusHeader)
			} else {
				creator.setActive(3)
			}
			break
		default:
			creator.inputField.HandleEvents(key, switchToPage, focusHeader)
			break
		}
		break
	case 2:
		switch key.Key() {
		case tcell.KeyLeft:
			creator.setActive(1)
			break
		case tcell.KeyRight:
			break
		case tcell.KeyUp:
			focusHeader(true)
			creator.Deselect()
			break
		case tcell.KeyDown:
			break
		}
		break
	case 3:
		switch key.Key() {
		case tcell.KeyLeft:
			break
		case tcell.KeyRight:
			creator.setActive(2)
			break
		case tcell.KeyUp:
			creator.setActive(1)
			break
		case tcell.KeyDown:
			break
		}
		break
	}
}
