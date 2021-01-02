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

func MakeLoadSpec(app *tview.Application, bg *backend.Backend, deleteFunc func()) *Operation {
	dataStore := operations.MakeLoadSpec()
	operation := MakeOperation("Load Spectrum", bg, page.MakeLayout([][]string{
		{"zip", "zip", "zip"},
		{"zip", "zip", "zip"},
		{"uuid", "uuid", "type"},
		{"uuid", "uuid", "type"},
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore)

	operation.SetPreDraw(func() {
		newLabel := fmt.Sprintf("Load Spectrum - %s", dataStore.GetName())

		if operation.GetLabel() != newLabel {
			operation.SetLabel(newLabel)
		}
	})

	zipInput := input.MakeInput(app)
	zipElement := element.MakeElement()
	zipElement.SetDirection(tview.FlexRow)
	zipElement.SetBorders(true)
	zipElement.AddItem(tview.NewTextView().SetText("Zip Name:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	zipElement.AddItem(zipInput, 0, 1, false)
	zipElement.SetOnSelect(func() {
		zipInput.Activate(true)
	})
	zipElement.SetOnDeselect(func() {
		zipInput.Activate(false)
	})
	zipElement.SetOnKeyEvent(zipInput.HandleEvents)
	zipInput.SetOnSubmit(dataStore.SetZip)

	uuidInput := input.MakeInput(app)
	uuidElement := element.MakeElement()
	uuidElement.SetDirection(tview.FlexRow)
	uuidElement.SetBorders(true)
	uuidElement.AddItem(tview.NewTextView().SetText("UUID:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	uuidElement.AddItem(uuidInput, 0, 1, false)
	uuidElement.SetOnSelect(func() {
		uuidInput.Activate(true)
	})
	uuidElement.SetOnDeselect(func() {
		uuidInput.Activate(false)
	})
	uuidElement.SetOnKeyEvent(uuidInput.HandleEvents)
	uuidInput.SetOnSubmit(dataStore.SetUuid)

	typeInput := input.MakeInput(app)
	typeElement := element.MakeElement()
	typeElement.SetDirection(tview.FlexRow)
	typeElement.SetBorders(true)
	typeElement.AddItem(tview.NewTextView().SetText("Type:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	typeElement.AddItem(typeInput, 0, 1, false)
	typeElement.SetOnSelect(func() {
		typeInput.Activate(true)
	})
	typeElement.SetOnDeselect(func() {
		typeInput.Activate(false)
	})
	typeElement.SetOnKeyEvent(typeInput.HandleEvents)
	typeInput.SetOnSubmit(dataStore.SetType)

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
		AddElement(zipElement, "zip").
		AddElement(uuidElement, "uuid").
		AddElement(typeElement, "type").
		AddElement(nameElement, "name").
		AddElement(deleteElement, "delete")

	operation.SetOnSelect(func() {
		operation.SelectElement("uuid")
	})

	return operation
}
