package inputoptions

import (
	"PLViewer/backend"
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputOptions struct {
	*element.Element
	Page    *page.Page
	Options backend.Options
}

func MakeInputOptions(pageLayout *page.Layout, options backend.Options) *InputOptions {
	inputOptions := InputOptions{
		Element: element.MakeElement(),
		Page:    page.MakePage("", pageLayout),
		Options: options,
	}

	inputOptions.AddItem(inputOptions.Page, 0, 1, false)

	inputOptions.SetOnDeselect(func() {
		inputOptions.Page.Deselect()
	})

	inputOptions.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		inputOptions.Page.HandleKeyEvent(key, func(p *page.Page) {
		}, func(b bool) {
		})
	})

	return &inputOptions
}

func (inputOptions *InputOptions) AddCheckbox(position string, label string, option *bool) {
	checkbox := tview.NewCheckbox().SetLabel(label)
	checkbox.SetChecked(*option)
	ele := element.MakeElement()
	ele.SetOnSelect(func() {
		isChecked := !*option
		*option = isChecked
		checkbox.SetChecked(isChecked)
		inputOptions.Page.DeselectActiveElement()
	}).SetOnHover(func() {
		checkbox.SetFieldBackgroundColor(tcell.ColorWhite)
		checkbox.SetFieldTextColor(tcell.ColorBlue)
	}).SetOffHover(func() {
		checkbox.SetFieldBackgroundColor(tcell.ColorBlue)
		checkbox.SetFieldTextColor(tcell.ColorWhite)
	})

	ele.AddItem(checkbox, 0, 1, false)

	inputOptions.Page.AddElement(ele, position)
}

func (inputOptions *InputOptions) Serialize() []string {
	return inputOptions.Options.Serialize()
}
