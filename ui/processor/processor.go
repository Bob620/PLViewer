package processor

import (
	"PLViewer/backend"
	"PLViewer/processor"
	"PLViewer/processor/operations"
	"PLViewer/ui/element"
	"PLViewer/ui/interop"
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

func MakeProcessorPage(application *tview.Application, bg *backend.Backend, interopData *interop.InteropData) *Processor {
	proc := Processor{
		Page: page.MakePage("Processor", page.MakeLayout([][]string{
			{"create", "create", "create", "create", "create", "create"},
			{"create", "create", "create", "create", "create", "create"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "side", "side"},
			{"options", "options", "options", "options", "run", "clear"},
		})),
		operations: map[string]*uioperation.Operation{},
	}
	proc.SetOnSelect(func() {
		proc.SelectElement("create")
	})

	optionsElement := element.MakeElement()
	optionsElement.SetInlay(true).SetEscapable(false).SetHoverable(false)

	operationList := orderedlist.MakeOrderedList()
	operationList.SetBorders(true)
	operationList.SetTitle(" Operation Order ")
	operationList.SetTitleAlign(tview.AlignCenter)

	addOperation := element.MakeElement()
	addOperation.SetDirection(tview.FlexRow)
	addOperation.SetBorders(true)
	addOperation.SetTitle(" Add Operation ")
	addOperationsElement := page.MakePage("", page.MakeLayout([][]string{
		{"addPLZip", "loadSpec", "makeSlice", "cubicSpline"},
		{"", "", "", ""},
		{"csv", "", "", ""},
		{"", "", "", ""},
		{"", "", "", ""},
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
		optionsElement.SetHoverable(false)
	}

	i := 0
	addOpFunc := func(maker func(*tview.Application, *backend.Backend, *interop.InteropData, func()) *uioperation.Operation) {
		id := strconv.Itoa(i)
		i++

		exists := true
		newOperation := maker(application, bg, interopData, func() {
			exists = false
			operationList.DeleteId(id)
		})
		newOperation.ListElement.SetOnSelect(func() {
			proc.selectedOperation = id
			optionsElement.Clear()
			optionsElement.AddItem(newOperation, 0, 1, false)
			optionsElement.SetHoverable(true)
			optionsElement.SetOnSelect(func() {
				newOperation.Select()
			})
			newOperation.SetOnNavigateOutside(func(direction string) {
				switch direction {
				case "right":
					proc.NavigateRight()
					break
				case "up":
					proc.NavigateUp()
					break
				}
			})
			optionsElement.SetOnKeyEvent(func(key *tcell.EventKey, f func()) {
				newOperation.HandleKeyEvent(key, func(p *page.Page) {}, func(b bool) {})
			})
			optionsElement.SetOnDeselect(func() {
				newOperation.Deselect()
				if exists {
					ops := map[string]operations.Interface{}
					for id, op := range proc.operations {
						ops[id] = op.Data
					}
					processor.RunOperations(operationList.Serialize(), ops)
				} else {
					proc.SelectElement("create")
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

	addOp := func(maker func(*tview.Application, *backend.Backend, *interop.InteropData, func()) *uioperation.Operation, label string, position string) {
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
	addOp(uioperation.MakeCubicSpline, "Cubic Spline", "cubicSpline")
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

		go func() {
			lines := []string{}
			updated := false
			var lock sync.RWMutex
			for {
				line := <-output
				if strings.Contains(line, "\b") {
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

	clearButton := element.MakeElement()

	clearButton.SetBorders(true)
	clearButton.SetTitle(" Clear ")
	clearButton.SetTitleAlign(tview.AlignCenter)
	clearButton.SetOnSelect(func() {
		proc.DeselectActiveElement()

	})

	proc.
		AddElement(optionsElement, "options").
		AddElement(addOperation, "create").
		AddElement(operationList, "side").
		AddElement(runButton, "run").
		AddElement(clearButton, "clear")

	return &proc
}
