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
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func MakeAddPLZip(app *tview.Application, bg *backend.Backend, interopData *interop.InteropData, deleteFunc func()) *Operation {
	dataStore := operations.MakeAddZip()
	operation := MakeOperation(bg, page.MakeLayout([][]string{
		{"file", "file", "file"},
		{"file", "file", "file"},
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore, deleteFunc)

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
	fileElement.SetTitle(" File ")
	fileElement.AddItem(fileInput, 0, 1, false)
	fileElement.SetOnSelect(func() {
		fileInput.Activate(true)
	})
	fileElement.SetOnDeselect(func() {
		fileInput.Activate(false)
	})
	fileElement.SetOnKeyEvent(fileInput.HandleEvents)
	fileInput.SetOnSubmit(func(uri string) {
		oldUri := dataStore.GetUri()
		dataStore.SetUri(uri)
		uris := interopData.GetStringArray("addZipUris")
		uris.Add(uri)
		if oldUri != "" {
			uris.Remove(oldUri)
		}
	})
	fileInput.SetAutocompleteSelectionFunc(func(entry string) string {
		if strings.HasPrefix(entry[1:], ":") {
			return entry
		}

		dir, _ := path.Split(fileInput.GetText())
		return path.Join(dir, entry)
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		currentText = strings.ReplaceAll(strings.TrimSpace(strings.ToLower(currentText)), "\\", "/")
		dir, wanted := path.Split(currentText)
		files := getDir(dir)

		viewerUri := interopData.GetString("viewerUri")
		creatorUri := interopData.GetString("creatorUri")

		testUri := func(uri string) {
			cleanUri := strings.ReplaceAll(strings.TrimSpace(strings.ToLower(uri)), "\\", "/")
			dir, _ := path.Split(cleanUri)

			if uri != "" && !strings.HasPrefix(currentText, dir) {
				entries = append(entries, uri)
			}
		}

		testUri(viewerUri)
		testUri(creatorUri)

		for _, file := range files {
			fileName := strings.ToLower(file.Name())
			if file.IsDir() || strings.HasSuffix(fileName, ".plzip") || strings.HasSuffix(fileName, ".pl7z") {
				if wanted == "" || (strings.Contains(fileName, wanted) && fileName != wanted) {
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
	nameElement.SetTitle(" Name ")
	nameElement.AddItem(nameInput, 0, 1, false)
	nameElement.SetOnSelect(func() {
		nameInput.Activate(true)
	})
	nameElement.SetOnDeselect(func() {
		nameInput.Activate(false)
	})
	nameElement.SetOnKeyEvent(nameInput.HandleEvents)
	nameInput.SetOnSubmit(dataStore.SetName)

	operation.
		AddElement(fileElement, "file").
		AddElement(nameElement, "name")

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
