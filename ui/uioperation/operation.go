package uioperation

import (
	"PLViewer/backend"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"github.com/rivo/tview"
)

type Operation struct {
	*page.Page
	Data        operations.Interface
	ListElement *element.Element
	bg          *backend.Backend
	label       string
}

func MakeOperation(bg *backend.Backend, layout *page.Layout, data operations.Interface, deleteFunc func()) *Operation {
	operation := &Operation{
		Page:        page.MakePage("", layout),
		Data:        data,
		ListElement: element.MakeElement(),
		bg:          bg,
		label:       "",
	}

	deleteElement := element.MakeElement()
	deleteElement.SetBorders(true)
	deleteElement.SetTitle(" Delete ")
	deleteElement.SetOnSelect(func() {
		deleteFunc()
	})

	operation.AddElement(deleteElement, "delete")
	return operation
}

func (operation *Operation) SetLabel(label string) {
	operation.label = label
	operation.ListElement.Clear()
	operation.ListElement.AddItem(tview.NewTextView().SetText(label), 0, 1, false).SetRect(0, 0, 0, 1)
}

func (operation *Operation) GetLabel() string {
	return operation.label
}
