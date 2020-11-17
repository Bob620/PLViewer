package element

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Interface interface {
	GetFlex() *tview.Flex
	SetSelectable(selectable bool) Interface
	Selectable() bool
	SetEscapable(escapable bool) Interface
	Escapable() bool
	SetBorders(borders bool) Interface
	Borders() bool
	SetOnSelect(handler func()) Interface
	Select()
	SetOnDeselect(handler func()) Interface
	Deselect()
	SetBorderColor(color tcell.Color) Interface
	SetOnEvent(handler func(*tcell.EventKey)) Interface
	HandleEvents(key *tcell.EventKey)

	// Other interface shit for tview
	Draw(screen tcell.Screen)
	GetRect() (int, int, int, int)
	SetRect(x, y, width, height int)
	InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive))
	Focus(delegate func(p tview.Primitive))
	Blur()
	GetFocusable() tview.Focusable
	MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive)
}

type Element struct {
	*tview.Flex
	onSelect   func()
	onDeselect func()
	onEvent    func(*tcell.EventKey)
	selectable bool
	escapable  bool
	borders    bool
}

func MakeElement() *Element {
	return &Element{
		Flex:       tview.NewFlex(),
		selectable: true,
		escapable:  true,
		borders:    false,
	}
}

func (element *Element) GetFlex() *tview.Flex {
	return element.Flex
}

func (element *Element) SetSelectable(selectable bool) Interface {
	element.selectable = selectable
	return element
}

func (element *Element) Selectable() bool {
	return element.selectable
}

func (element *Element) SetEscapable(escapable bool) Interface {
	element.escapable = escapable
	return element
}

func (element *Element) Escapable() bool {
	return element.escapable
}

func (element *Element) SetBorders(borders bool) Interface {
	element.borders = borders
	element.Flex.SetBorder(borders)
	return element
}

func (element *Element) Borders() bool {
	return element.borders
}

func (element *Element) SetOnSelect(handler func()) Interface {
	element.onSelect = handler
	return element
}

func (element *Element) Select() {
	if element.onSelect != nil {
		element.onSelect()
	}
}

func (element *Element) SetOnDeselect(handler func()) Interface {
	element.onDeselect = handler
	return element
}

func (element *Element) Deselect() {
	if element.onDeselect != nil {
		element.onDeselect()
	}
}

func (element *Element) SetBorderColor(color tcell.Color) Interface {
	element.Flex.SetBorderColor(color)
	return element
}

func (element *Element) SetOnEvent(handler func(*tcell.EventKey)) Interface {
	element.onEvent = handler
	return element
}

func (element *Element) HandleEvents(key *tcell.EventKey) {
	if element.onEvent != nil {
		element.onEvent(key)
	}
}
