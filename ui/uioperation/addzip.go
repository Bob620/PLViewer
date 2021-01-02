package uioperation

import (
	"PLViewer/backend"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/page"
	"fmt"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func MakeAddPLZip(app *tview.Application, bg *backend.Backend, deleteFunc func()) *Operation {
	dataStore := operations.MakeAddZip()
	operation := MakeOperation("Add PLZip", bg, page.MakeLayout([][]string{
		{"file", "file", "file"},
		{"file", "file", "file"},
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore)

	operation.SetPreDraw(func() {
		newLabel := fmt.Sprintf("Add PLZip - %s", dataStore.GetName())

		if operation.GetLabel() != newLabel {
			operation.SetLabel(newLabel)
		}
	})

	fileInput := input.MakeInput(app)
	fileElement := element.MakeElement()
	fileElement.SetDirection(tview.FlexRow)
	fileElement.SetBorders(true)
	fileElement.AddItem(tview.NewTextView().SetText("File:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	fileElement.AddItem(fileInput, 0, 1, false)
	fileElement.SetOnSelect(func() {
		fileInput.Activate(true)
	})
	fileElement.SetOnDeselect(func() {
		fileInput.Activate(false)
	})
	fileElement.SetOnKeyEvent(fileInput.HandleEvents)
	fileInput.SetOnSubmit(dataStore.SetUri)
	fileInput.SetAutocompleteSelectionFunc(func(entry string) string {
		dir, _ := path.Split(fileInput.GetText())
		return path.Join(dir, entry)
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		dir, wanted := path.Split(prefix)
		files := getDir(dir)

		for _, file := range files {
			if file.IsDir() || strings.HasSuffix(file.Name(), ".plzip") || strings.HasSuffix(file.Name(), ".pl7z") {
				if wanted == "" || strings.Contains(strings.ToLower(file.Name()), wanted) {
					entries = append(entries, file.Name())
				}
			}
		}
		return entries
	})

	nameInput := input.MakeInput(app)
	nameElement := element.MakeElement()
	nameElement.SetDirection(tview.FlexRow)
	nameElement.SetBorders(true)
	nameElement.AddItem(tview.NewTextView().SetText("Name:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	nameElement.AddItem(nameInput, 0, 1, false)
	nameElement.SetOnSelect(func() {
		nameInput.Activate(true)
	})
	nameElement.SetOnDeselect(func() {
		nameInput.Activate(false)
	})
	nameElement.SetOnKeyEvent(nameInput.HandleEvents)
	nameInput.SetOnSubmit(dataStore.SetName)

	deleteElement := element.MakeElement()
	deleteElement.SetBorders(true)
	deleteElement.AddItem(tview.NewTextView().SetText("Delete").SetTextAlign(tview.AlignCenter), 0, 1, false)
	deleteElement.SetOnSelect(func() {
		deleteFunc()
	})

	operation.
		AddElement(fileElement, "file").
		AddElement(nameElement, "name").
		AddElement(deleteElement, "delete")

	operation.SetOnSelect(func() {
		operation.SelectElement("file")
	})

	return operation
}

func getDir(uri string) []os.FileInfo {
	files, err := ioutil.ReadDir(uri)
	if err != nil {
		return []os.FileInfo{}
	}

	return files
}
