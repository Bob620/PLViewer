package radio

import (
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Option struct {
	Name string
	Id   string
}

type internalOption struct {
	element  *element.Element
	checkbox *tview.Checkbox
	name     string
	id       string
}

type Radio struct {
	*element.Element
	Page           *page.Page
	label          string
	options        []*internalOption
	selectedOption string
	width          int
	height         int
}

func MakeRadioSelector(label string, options []*Option) *Radio {
	radio := Radio{
		Element:        element.MakeElement(),
		Page:           page.MakePage("", page.MakeLayout([][]string{{}})),
		label:          label,
		options:        []*internalOption{},
		selectedOption: "",
		width:          len(label),
		height:         1,
	}

	radio.AddItem(radio.Page, 0, 1, false)

	radio.SetOnHover(func() {
		radio.Page.SelectElement(radio.selectedOption)
	})

	radio.SetOffHover(func() {
		radio.Page.Deselect()
	})

	radio.SetSelectable(false).SetInlay(true).SetEscapable(false)

	radio.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		radio.Page.HandleKeyEvent(key, func(p *page.Page) {}, func(b bool) {})
	})

	for _, option := range options {
		radio.AddOption(option.Name, option.Id)
	}

	return &radio
}

func (radio *Radio) SetLabel(label string) {
	radio.label = label
}

func (radio *Radio) reLayout() {
	newLayout := [][]string{{"label"}}
	for _, option := range radio.options {
		newLayout = append(newLayout, []string{option.id})
	}

	radio.Page.Clear()
	radio.Page.SetLayout(page.MakeLayout(newLayout))

	maxWidth := len(radio.label)
	for _, option := range radio.options {
		radio.Page.AddElement(option.element, option.id)

		nameLength := len(option.name)
		if maxWidth < nameLength {
			maxWidth = nameLength
		}
	}

	radio.width = maxWidth
	radio.height = len(radio.options) + 1
}

func (radio *Radio) SelectOption(id string) {
	if id == radio.selectedOption {
		return
	}
	oldId := radio.selectedOption

	for _, option := range radio.options {
		if option.id == id {
			radio.selectedOption = id
			option.checkbox.SetChecked(true)
		}

		if option.id == oldId {
			option.checkbox.SetChecked(false)
		}
	}
}

func (radio *Radio) AddOption(name string, id string) {
	label := tview.NewTextView().SetText(name)
	checkbox := tview.NewCheckbox().SetLabel(" ")
	ele := element.MakeElement()
	ele.SetOnHover(func() {
		radio.SelectOption(id)
		checkbox.SetFieldBackgroundColor(tcell.ColorWhite)
		checkbox.SetFieldTextColor(tcell.ColorBlue)
	}).SetOffHover(func() {
		checkbox.SetFieldBackgroundColor(tcell.ColorBlue)
		checkbox.SetFieldTextColor(tcell.ColorWhite)
	})

	ele.AddItem(checkbox, 0, 1, false)
	ele.AddItem(label, 0, 1, false)

	radio.options = append(radio.options, &internalOption{
		element:  ele,
		checkbox: checkbox,
		name:     name,
		id:       id,
	})

	radio.reLayout()
}

func (radio *Radio) DeleteOption(id string) {
	for i, option := range radio.options {
		if option.id == id {
			radio.options = append(radio.options[:i], radio.options[i+1:]...)

			if radio.selectedOption == option.id {
				radio.selectedOption = ""
			}
		}
	}

	radio.reLayout()
}

func (radio *Radio) Clear() {
	radio.options = []*internalOption{}
	radio.selectedOption = ""
	radio.reLayout()
}

func (radio *Radio) Serialize() string {
	return radio.selectedOption
}

func (radio *Radio) Draw(screen tcell.Screen) {
	//	radio.Element.Draw(screen)
	x, y, width, height := radio.GetRect()

	ele := tview.NewTextView().SetText(radio.label)
	ele.SetRect(x, y, width, 1)
	ele.Draw(screen)

	height--
	y++

	for _, item := range radio.options {
		if height > 0 {
			item.element.SetRect(x, y, width, 1)
			item.element.Draw(screen)

			y += 1
			height--
		}
	}
}
