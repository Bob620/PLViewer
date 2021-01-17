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
	"strconv"
)

func MakeCubicSpline(app *tview.Application, bg *backend.Backend, interopData *interop.InteropData, deleteFunc func()) *Operation {
	dataStore := operations.MakeCubicSpline()
	operation := MakeOperation(bg, page.MakeLayout([][]string{
		{"inName", "inName", "inName"},
		{"inName", "inName", "inName"},
		{"outName", "outName", "outName"},
		{"outName", "outName", "outName"},
		{"", "", "interpret"},
		{"", "", "interpret"},
		{"", "", ""},
		{"", "", ""},
		{"", "", "delete"},
	}), dataStore, deleteFunc)
	interpretInput := input.MakeInput(app)

	operation.SetPreDraw(func() {
		newLabel := fmt.Sprintf("Cubic Spline - %s[%d] as %s", dataStore.GetInName(), dataStore.GetInterpret(), dataStore.GetName())

		if operation.GetLabel() != newLabel {
			operation.SetLabel(newLabel)
		}
	})

	inNameInput := input.MakeInput(app)
	inNameElement := element.MakeElement()
	inNameElement.SetDirection(tview.FlexRow)
	inNameElement.SetBorders(true)
	inNameElement.SetTitle(" Input Name ")
	inNameElement.AddItem(inNameInput, 0, 1, false)
	inNameElement.SetOnSelect(func() {
		inNameInput.Activate(true)
	})
	inNameElement.SetOnDeselect(func() {
		inNameInput.Activate(false)
	})
	inNameElement.SetOnKeyEvent(inNameInput.HandleEvents)
	inNameInput.SetOnSubmit(dataStore.SetInName)

	outNameInput := input.MakeInput(app)
	outNameElement := element.MakeElement()
	outNameElement.SetDirection(tview.FlexRow)
	outNameElement.SetBorders(true)
	outNameElement.SetTitle(" Output Name ")
	outNameElement.AddItem(outNameInput, 0, 1, false)
	outNameElement.SetOnSelect(func() {
		outNameInput.Activate(true)
	})
	outNameElement.SetOnDeselect(func() {
		outNameInput.Activate(false)
	})
	outNameElement.SetOnKeyEvent(outNameInput.HandleEvents)
	outNameInput.SetOnSubmit(dataStore.SetName)

	interpretElement := element.MakeElement()
	interpretElement.SetDirection(tview.FlexRow)
	interpretElement.SetBorders(true)
	interpretElement.SetTitle(" Interpolation ")
	interpretElement.AddItem(interpretInput, 0, 1, false)
	interpretElement.SetOnSelect(func() {
		interpretInput.Activate(true)
	})
	interpretElement.SetOnDeselect(func() {
		interpretInput.Activate(false)
	})
	interpretElement.SetOnKeyEvent(interpretInput.HandleEvents)
	interpretInput.SetOnSubmit(func(interp string) {
		value, err := strconv.Atoi(interp)
		if err == nil {
			dataStore.SetInterpret(value)
		}
	})

	operation.
		AddElement(inNameElement, "inName").
		AddElement(outNameElement, "outName").
		AddElement(interpretElement, "interpret")

	operation.SetOnSelect(func() {
		operation.SelectElement("inName")
	})

	return operation
}
