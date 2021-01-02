package radio

import (
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
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
	Page             *page.Page
	label            string
	options          []*internalOption
	selectedOption   string
	selectedPosition int
	width            int
	height           int
}

func MakeRadioSelector(label string, options []*Option) *Radio {
	radio := Radio{
		Element:          element.MakeElement(),
		Page:             page.MakePage("", page.MakeLayout([][]string{{}})),
		label:            label,
		options:          []*internalOption{},
		selectedOption:   "",
		selectedPosition: -1,
		width:            len(label),
		height:           1,
	}

	radio.AddItem(radio.Page, 0, 1, false)

	radio.SetOnHover(func() {
		if radio.selectedPosition != -1 {
			radio.options[radio.selectedPosition].element.Hover()
		}
	})

	radio.SetOffHover(func() {
		if radio.selectedPosition != -1 {
			radio.options[radio.selectedPosition].element.OffHover()
		}
	})

	radio.SetSelectable(false)

	radio.SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		switch key.Key() {
		case tcell.KeyUp:
			radio.MoveUp()
			break
		case tcell.KeyDown:
			radio.MoveDown()
			break
		}
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
	for i := range radio.options {
		newLayout = append(newLayout, []string{strconv.Itoa(i)})
	}

	radio.Page.Clear()
	radio.Page.SetLayout(page.MakeLayout(newLayout))

	maxWidth := len(radio.label)
	for i, option := range radio.options {
		radio.Page.AddElement(option.element, strconv.Itoa(i))

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

	for i, option := range radio.options {
		if option.id == id {
			radio.selectedOption = id
			radio.selectedPosition = i
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
	ele.SetOnSelect(func() {
		radio.SelectOption(id)
		radio.Page.DeselectActiveElement()
	}).SetOnHover(func() {
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
				radio.selectedPosition = -1
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

func (radio *Radio) MoveDown() {
	if radio.selectedPosition != -1 {
		radio.options[radio.selectedPosition].element.OffHover()
	}

	radio.selectedPosition++
	maxPos := len(radio.options) - 1
	if radio.selectedPosition > maxPos {
		radio.selectedPosition = maxPos
	}

	radio.SelectOption(radio.options[radio.selectedPosition].id)
	radio.options[radio.selectedPosition].element.Hover()
}

func (radio *Radio) MoveUp() {
	if radio.selectedPosition != -1 {
		radio.options[radio.selectedPosition].element.OffHover()
	}

	radio.selectedPosition--
	if radio.selectedPosition == -1 {
		radio.selectedPosition = 0
	}

	radio.SelectOption(radio.options[radio.selectedPosition].id)
	radio.options[radio.selectedPosition].element.Hover()
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
