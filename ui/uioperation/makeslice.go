package uioperation

import (
	"PLViewer/backend"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/page"
	"fmt"
	"github.com/rivo/tview"
	"strconv"
)

func MakeMakeSlice(app *tview.Application, bg *backend.Backend, deleteFunc func()) *Operation {
	dataStore := operations.MakeMakeSlice()
	operation := MakeOperation("Make Slice", bg, page.MakeLayout([][]string{
		{"inName", "inName", "inName"},
		{"inName", "inName", "inName"},
		{"outName", "outName", "outName"},
		{"outName", "outName", "outName"},
		{"", "", "low"},
		{"", "", "low"},
		{"", "", "high"},
		{"", "", "high"},
		{"", "", "delete"},
	}), dataStore)
	lowInput := input.MakeInput(app)
	highInput := input.MakeInput(app)

	operation.SetPreDraw(func() {
		newLabel := fmt.Sprintf("Make Slice - %s[%d:%d]", dataStore.GetInName(), dataStore.GetLow(), dataStore.GetHigh())

		if operation.GetLabel() != newLabel {
			operation.SetLabel(newLabel)
		}
	})

	inNameInput := input.MakeInput(app)
	inNameElement := element.MakeElement()
	inNameElement.SetDirection(tview.FlexRow)
	inNameElement.SetBorders(true)
	inNameElement.AddItem(tview.NewTextView().SetText("Input Name:").SetTextAlign(tview.AlignCenter), 0, 1, false)
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
	outNameElement.AddItem(tview.NewTextView().SetText("Output Name:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	outNameElement.AddItem(outNameInput, 0, 1, false)
	outNameElement.SetOnSelect(func() {
		outNameInput.Activate(true)
	})
	outNameElement.SetOnDeselect(func() {
		outNameInput.Activate(false)
	})
	outNameElement.SetOnKeyEvent(outNameInput.HandleEvents)
	outNameInput.SetOnSubmit(dataStore.SetName)

	lowElement := element.MakeElement()
	lowElement.SetDirection(tview.FlexRow)
	lowElement.SetBorders(true)
	lowElement.AddItem(tview.NewTextView().SetText("Low:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	lowElement.AddItem(lowInput, 0, 1, false)
	lowElement.SetOnSelect(func() {
		lowInput.Activate(true)
	})
	lowElement.SetOnDeselect(func() {
		lowInput.Activate(false)
	})
	lowElement.SetOnKeyEvent(lowInput.HandleEvents)
	lowInput.SetOnSubmit(func(lowString string) {
		dataStore.SetLow(lowString)
		low := dataStore.GetLow()
		if low != 0 {
			lowInput.SetText(strconv.Itoa(low))
		}
	})

	highElement := element.MakeElement()
	highElement.SetDirection(tview.FlexRow)
	highElement.SetBorders(true)
	highElement.AddItem(tview.NewTextView().SetText("High:").SetTextAlign(tview.AlignCenter), 0, 1, false)
	highElement.AddItem(highInput, 0, 1, false)
	highElement.SetOnSelect(func() {
		highInput.Activate(true)
	})
	highElement.SetOnDeselect(func() {
		highInput.Activate(false)
	})
	highElement.SetOnKeyEvent(highInput.HandleEvents)
	highInput.SetOnSubmit(func(highString string) {
		dataStore.SetHigh(highString)
		high := dataStore.GetHigh()
		if high != 0 {
			highInput.SetText(strconv.Itoa(high))
		}
	})

	deleteElement := element.MakeElement()
	deleteElement.SetBorders(true)
	deleteElement.AddItem(tview.NewTextView().SetText("Delete").SetTextAlign(tview.AlignCenter), 0, 1, false)
	deleteElement.SetOnSelect(func() {
		deleteFunc()
	})

	operation.
		AddElement(inNameElement, "inName").
		AddElement(outNameElement, "outName").
		AddElement(lowElement, "low").
		AddElement(highElement, "high").
		AddElement(deleteElement, "delete")

	operation.SetOnSelect(func() {
		operation.SelectElement("inName")
	})

	return operation
}
