package uioperation

import (
	"PLViewer/backend"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/interop"
	"PLViewer/ui/page"
	"fmt"
	"github.com/rivo/tview"
	"path"
	"strings"
)

func MakeCsv(app *tview.Application, bg *backend.Backend, interopData *interop.InteropData, deleteFunc func()) *Operation {
	dataStore := operations.MakeCsv()
	operation := MakeOperation(bg, page.MakeLayout([][]string{
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"file", "file", "file"},
		{"file", "file", "file"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore, deleteFunc)

	operation.SetPreDraw(func() {
		newLabel := fmt.Sprintf("Export as Csv - %s as %s", dataStore.GetDataNames(), dataStore.GetUri())

		if operation.GetLabel() != newLabel {
			operation.SetLabel(newLabel)
		}
	})

	nameInput := input.MakeInput(app)
	nameElement := element.MakeElement()
	nameElement.SetDirection(tview.FlexRow)
	nameElement.SetBorders(true)
	nameElement.SetTitle(" Data Names (separated by ,) ")
	nameElement.AddItem(nameInput, 0, 1, false)
	nameElement.SetOnSelect(func() {
		nameInput.Activate(true)
	})
	nameElement.SetOnDeselect(func() {
		nameInput.Activate(false)
	})
	nameElement.SetOnKeyEvent(nameInput.HandleEvents)
	nameInput.SetOnSubmit(dataStore.SetDataNames)

	fileInput := input.MakeInput(app)
	fileElement := element.MakeElement()
	fileElement.SetDirection(tview.FlexRow)
	fileElement.SetBorders(true)
	fileElement.SetTitle(" Output File ")
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
		if strings.HasPrefix(entry[1:], ":") {
			return entry
		}

		dir, _ := path.Split(fileInput.GetText())
		return path.Join(dir, entry)
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		currentText = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(currentText)), "\\", "/")

		viewerUri := interopData.GetString("viewerUri")
		creatorUri := interopData.GetString("creatorUri")

		testUri := func(uri string) {
			cleanUri := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(uri)), "\\", "/")
			dir, _ := path.Split(uri)
			cleanDir, _ := path.Split(cleanUri)

			if uri != "" && !strings.HasPrefix(currentText, cleanDir) {
				entries = append(entries, dir)
			}
		}

		testUri(viewerUri)
		testUri(creatorUri)

		uris := interopData.GetStringArray("addZipUris")
		for _, uri := range uris.GetAll() {
			if uri != viewerUri && uri != creatorUri {
				entries = append(entries, uri)
			}
		}

		return entries
	})

	operation.
		AddElement(nameElement, "name").
		AddElement(fileElement, "file")

	operation.SetOnSelect(func() {
		operation.SelectElement("name")
	})

	return operation
}
