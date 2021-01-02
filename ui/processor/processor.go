package processor

import (
	"PLViewer/backend"
	"PLViewer/processor"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/orderedlist"
	"PLViewer/ui/page"
	"PLViewer/ui/uioperation"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
	"sync"
)

type Processor struct {
	*page.Page
	operations        map[string]*uioperation.Operation
	selectedOperation string
}

func MakeProcessorPage(application *tview.Application, bg *backend.Backend) *Processor {
	proc := Processor{
		Page: page.MakePage("Processor", page.MakeLayout([][]string{
			{"create", "create", "create", "create", "create", "create"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "run", "run"},
		})),
		operations: map[string]*uioperation.Operation{},
	}
	proc.SetOnSelect(func() {
		proc.SelectElement("create")
	})

	optionsElement := element.MakeElement()
	optionsElement.SetBorders(true)
	optionsElement.SetEscapable(false)
	optionsElement.SetSelectable(false)

	operationList := orderedlist.MakeOrderedList()
	operationList.SetBorders(true)
	operationList.SetTitle(" Operation Order ")
	operationList.SetTitleAlign(tview.AlignCenter)

	addOperation := element.MakeElement()
	addOperation.SetDirection(tview.FlexRow)
	addOperation.SetBorders(true)
	addOperation.SetTitle(" Add Operation ")
	addOperationsElement := page.MakePage("", page.MakeLayout([][]string{
		{"addPLZip", "loadSpec", "makeSlice"},
		{"", "", ""},
		{"csv", "", ""},
		{"", "", ""},
		{"", "", ""},
	}))

	addOperation.AddItem(addOperationsElement, 0, 1, false)
	addOperation.SetOnSelect(func() {
		addOperationsElement.SelectElement("loadSpec")
	})
	addOperation.SetOnDeselect(func() {
		addOperationsElement.Deselect()
	})
	addOperation.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		addOperationsElement.HandleKeyEvent(key, func(p *page.Page) {}, func(b bool) {})
	})

	clearOperation := func() {
		proc.selectedOperation = ""
		optionsElement.Clear()
		optionsElement.SetOnSelect(func() {})
		optionsElement.SetOnDeselect(func() {})
		optionsElement.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {})
		optionsElement.SetSelectable(false)
	}

	i := 0
	addOpFunc := func(maker func(*tview.Application, *backend.Backend, func()) *uioperation.Operation) {
		id := strconv.Itoa(i)
		i++

		newOperation := maker(application, bg, func() {
			operationList.DeleteId(id)
		})
		newOperation.ListElement.SetOnSelect(func() {
			proc.selectedOperation = id
			optionsElement.Clear()
			optionsElement.AddItem(newOperation, 0, 1, false)
			optionsElement.SetSelectable(true)
			optionsElement.SetOnSelect(func() {
				newOperation.Select()
			})
			optionsElement.SetOnDeselect(func() {
				newOperation.Deselect()
				ops := map[string]operations.Interface{}
				for id, op := range proc.operations {
					ops[id] = op.Data
				}
				processor.RunOperations(operationList.Serialize(), ops)
			})
			optionsElement.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
				switch key.Key() {
				case tcell.KeyEscape:
					if newOperation.SelectedElement() == "" {
						proc.DeselectActiveElement()
					} else {
						newOperation.DeselectActiveElement()
					}
					break
				default:
					newOperation.HandleKeyEvent(key, func(p *page.Page) {}, func(b bool) {})
					break
				}
			})
		})
		newOperation.ListElement.SetOnDeselect(clearOperation)

		pos := operationList.AddItem(newOperation.ListElement, -1, id)
		proc.operations[id] = newOperation

		newOperation.ListElement.SetEvent("delete", func() {
			delete(proc.operations, id)
			if proc.selectedOperation == id {
				proc.DeselectActiveElement()
				clearOperation()
			}
		})

		operationList.SelectPosition(pos)
		operationList.Deselect()
	}

	addOp := func(maker func(*tview.Application, *backend.Backend, func()) *uioperation.Operation, label string, position string) {
		elem := element.MakeElement()
		elemText := tview.NewTextView().SetText(label).SetTextAlign(tview.AlignCenter)
		elem.AddItem(elemText, 0, 1, false)
		addOperationsElement.AddElement(elem, position)
		elem.SetOnSelect(func() {
			addOpFunc(maker)
			addOperationsElement.DeselectActiveElement()
		})
		elem.SetOnHover(func() {
			elemText.SetBackgroundColor(tcell.ColorBlue)
		})
		elem.SetOffHover(func() {
			elemText.SetBackgroundColor(tcell.ColorBlack)
		})
	}

	addOp(uioperation.MakeAddPLZip, "Add PLZip", "addPLZip")
	addOp(uioperation.MakeLoadSpec, "Load Spectrum", "loadSpec")
	addOp(uioperation.MakeMakeSlice, "Make Slice", "makeSlice")
	addOp(uioperation.MakeCsv, "Export to Csv", "csv")

	runButton := element.MakeElement()
	loggingSpace := tview.NewTextView()

	runButton.GetFlex().
		AddItem(loggingSpace, 0, 1, false)
	runButton.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		loggingSpace.InputHandler()(key, func(p tview.Primitive) {
		})
	})

	runButton.SetBorders(true)
	runButton.SetTitle(" Run ")
	runButton.SetTitleAlign(tview.AlignCenter)
	runButton.SetOnSelect(func() {
		proc.DeselectActiveElement()

		ops := map[string]operations.Interface{}
		for id, op := range proc.operations {
			ops[id] = op.Data
		}
		command := processor.SerializeOperations(operationList.Serialize(), ops)

		output := make(chan string)
		lines := []string{}
		updated := false
		var lock sync.RWMutex

		go func() {
			for {
				line := <-output
				if line == "\b" {
					break
				} else {
					lock.Lock()
					lines = append(lines, line)
					updated = true
					lock.Unlock()
				}

				go func() {
					application.QueueUpdateDraw(func() {
						lock.RLock()
						if updated {
							updated = false
							loggingSpace.SetText(strings.Join(lines, "\n"))
						}
						lock.RUnlock()
					})
				}()
			}
		}()
		go backend.RunProcessor(output, command)
	})

	proc.
		AddElement(optionsElement, "options").
		AddElement(addOperation, "create").
		AddElement(operationList, "side").
		AddElement(runButton, "run")

	return &proc
}
