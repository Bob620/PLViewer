package creator

import (
	"PLViewer/backend"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/inputoptions"
	"PLViewer/ui/page"
	"PLViewer/ui/radio"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Creator struct {
	*page.Page
}

func MakeCreatorPage(application *tview.Application) *Creator {
	creator := Creator{
		Page: page.MakePage("Creator", page.MakeLayout([][]string{
			{"input", "input", "side"},
			{"options", "options", "side"},
			{"options", "options", "side"},
			{"log", "log", "log"},
			{"log", "log", "log"},
			{"log", "log", "log"},
		})),
	}
	creator.SetOnSelect(func() {
		creator.Page.SelectElement("input")
	})

	optionsElement := inputoptions.MakeInputOptions(page.MakeLayout([][]string{
		{"selector", "xes", "qmap"},
		{"selector", "qlw", "loose"},
		{"selector", "map", "recover"},
		{"selector", "line", "debug"},
	}), backend.MakeCreatorOptions())
	optionsElement.SetBorder(true)
	optionsElement.SetOnSelect(func() {
		optionsElement.Page.SelectElement("xes")
	})

	radioSelector := radio.MakeRadioSelector("Output Format:", []*radio.Option{{"Csv", "csv"}, {"Jeol", "jeol"}, {"Json", "json"}, {"PlZip ", "plzip"}})
	radioSelector.SetBleedThrough(true)
	optionsElement.Page.AddElement(radioSelector, "selector")
	radioSelector.SelectOption("plzip")

	options := optionsElement.Options.(*backend.CreatorOptions)

	optionsElement.AddCheckbox("xes", "Convert XES:      ", &options.ConvertXes)
	optionsElement.AddCheckbox("qlw", "Convert QLW:      ", &options.ConvertQlw)
	optionsElement.AddCheckbox("map", "Convert Maps:     ", &options.ConvertMap)
	optionsElement.AddCheckbox("line", "Convert Lines:    ", &options.ConvertLine)

	optionsElement.AddCheckbox("qmap", "Convert qMaps:  ", &options.ConvertQMap)
	optionsElement.AddCheckbox("loose", "Search Loosely: ", &options.Loose)
	optionsElement.AddCheckbox("recover", "Recover Data:   ", &options.Recover)
	optionsElement.AddCheckbox("debug", "Debug:          ", &options.Debug)

	inputField := input.MakeInput(application)
	inputElement := element.MakeElement()
	sideElement := element.MakeElement().SetBorders(true).SetHoverable(false).SetSelectable(false)
	logElement := element.MakeElement().SetBorders(true)

	inputElement.SetBorders(true)
	loggingSpace := tview.NewTextView()

	logElement.GetFlex().
		AddItem(loggingSpace, 0, 1, false)
	logElement.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		loggingSpace.InputHandler()(key, func(p tview.Primitive) {
		})
	})

	inputElement.GetFlex().
		AddItem(inputField, 0, 1, false)
	inputElement.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		inputField.HandleEvents(key, deSelector)
	})
	inputElement.SetOnSelect(func() {
		inputField.Activate(true)
	})
	inputElement.SetOnDeselect(func() {
		inputField.Activate(false)
	})

	inputField.SetOnSubmit(func(uri string) {
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
		options.ExportTo = radioSelector.Serialize()
		go backend.RunCreator(output, uri, options)
	})

	inputField.SetAutocompleteSelectionFunc(func(entry string) string {
		dir, _ := path.Split(inputField.GetText())
		return path.Join(dir, entry) + "/"
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		dir, wanted := path.Split(prefix)
		files := getDir(dir)

		for _, file := range files {
			if file.IsDir() && (wanted == "" || strings.Contains(strings.ToLower(file.Name()), wanted)) {
				entries = append(entries, file.Name())
			}
		}
		return entries
	})

	creator.
		AddElement(inputElement, "input").
		AddElement(optionsElement, "options").
		AddElement(sideElement, "side").
		AddElement(logElement, "log")

	return &creator
}

func getDir(uri string) []os.FileInfo {
	files, err := ioutil.ReadDir(uri)
	if err != nil {
		return []os.FileInfo{}
	}

	return files
}
