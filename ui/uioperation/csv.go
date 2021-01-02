package uioperation

import (
	"PLViewer/backend"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/page"
	"fmt"
	"github.com/rivo/tview"
)

func MakeCsv(app *tview.Application, bg *backend.Backend, deleteFunc func()) *Operation {
	dataStore := operations.MakeCsv()
	operation := MakeOperation("Export as Csv", bg, page.MakeLayout([][]string{
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"file", "file", "file"},
		{"file", "file", "file"},
		{"", "", ""},
		{"", "", ""},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore)

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
	nameElement.AddItem(tview.NewTextView().SetText("Data Names (separated by ,):").SetTextAlign(tview.AlignCenter), 0, 1, false)
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
	fileElement.AddItem(tview.NewTextView().SetText("Output File:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	fileElement.AddItem(fileInput, 0, 1, false)
	fileElement.SetOnSelect(func() {
		fileInput.Activate(true)
	})
	fileElement.SetOnDeselect(func() {
		fileInput.Activate(false)
	})
	fileElement.SetOnKeyEvent(fileInput.HandleEvents)
	fileInput.SetOnSubmit(dataStore.SetUri)

	deleteElement := element.MakeElement()
	deleteElement.SetBorders(true)
	deleteElement.AddItem(tview.NewTextView().SetText("Delete").SetTextAlign(tview.AlignCenter), 0, 1, false)
	deleteElement.SetOnSelect(func() {
		deleteFunc()
	})

	operation.
		AddElement(nameElement, "name").
		AddElement(fileElement, "file").
		AddElement(deleteElement, "delete")

	operation.SetOnSelect(func() {
		operation.SelectElement("name")
	})

	return operation
}
