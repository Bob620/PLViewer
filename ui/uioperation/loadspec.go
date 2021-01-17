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
)

func MakeLoadSpec(app *tview.Application, bg *backend.Backend, interopData *interop.InteropData, deleteFunc func()) *Operation {
	dataStore := operations.MakeLoadSpec()
	operation := MakeOperation(bg, page.MakeLayout([][]string{
		{"zip", "zip", "zip"},
		{"zip", "zip", "zip"},
		{"uuid", "uuid", "type"},
		{"uuid", "uuid", "type"},
		{"name", "name", "name"},
		{"name", "name", "name"},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore, deleteFunc)

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
	zipElement.SetTitle(" Zip Name ")
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
	uuidElement.SetTitle(" UUID ")
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
	typeElement.SetTitle(" Type ")
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
		AddElement(zipElement, "zip").
		AddElement(uuidElement, "uuid").
		AddElement(typeElement, "type").
		AddElement(nameElement, "name")

	operation.SetOnSelect(func() {
		operation.SelectElement("zip")
	})

	return operation
}
