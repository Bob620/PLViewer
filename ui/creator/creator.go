package creator

import (
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Creator struct {
	*page.Page
}

func MakeCreatorPage(application *tview.Application) *Creator {
	creator := Creator{
		Page: page.MakePage("Creator", page.MakeLayout([][]string{
			{"input", "side"},
			{"log", "side"},
		})),
	}
	creator.SetOnSelect(func() {
		creator.Page.SelectElement("input")
	})

	inputField := input.MakeInput(application)
	inputElement := element.MakeElement().SetBorders(true)
	rightElement := element.MakeElement().SetBorders(true)
	bottomElement := element.MakeElement().SetBorders(true)

	rightElement.Flex().
		AddItem(tview.NewTextView().SetText("Test"), 0, 1, false)
	bottomElement.Flex().
		AddItem(tview.NewTextView().SetText("Ahhh"), 0, 1, false)
	inputElement.Flex().
		AddItem(inputField, 0, 1, false)
	inputElement.SetOnEvent(func(key *tcell.EventKey) {
		inputField.HandleEvents(key)
	})
	inputElement.SetOnSelect(func() {
		inputField.Activate(true)
	})
	inputElement.SetOnDeselect(func() {
		inputField.Activate(false)
	})

	inputField.SetLabel("Select a Directory to Pack: ")
	inputField.SetAutocompleteSelectionFunc(func(entry string) string {
		dir, _ := path.Split(inputField.GetText())
		return path.Join(dir, entry) + "/"
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		dir, wanted := path.Split(prefix)
		files := getDir(dir)

		for _, file := range files {
			if file.IsDir() && (wanted == "" || strings.Contains(strings.ToLower(file.Name()), wanted)) {
				entries = append(entries, file.Name())
			}
		}
		return entries
	})

	creator.
		AddElement(inputElement, "input").
		AddElement(rightElement, "side").
		AddElement(bottomElement, "log")

	return &creator
}

func getDir(uri string) []os.FileInfo {
	files, err := ioutil.ReadDir(uri)
	if err != nil {
		return []os.FileInfo{}
	}

	return files
}
