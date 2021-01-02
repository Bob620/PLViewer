package element

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Interface interface {
	GetFlex() *tview.Flex
	SetEvent(event string, handler func()) Interface
	EmitEvent(event string) Interface
	SetSelectable(selectable bool) Interface
	Selectable() bool
	SetHoverable(hoverable bool) Interface
	Hoverable() bool
	SetBleedThrough(bleedThrough bool) Interface
	BleedThrough() bool
	SetEscapable(escapable bool) Interface
	Escapable() bool
	SetBorders(borders bool) Interface
	Borders() bool
	SetOnSelect(handler func()) Interface
	Select()
	SetOnDeselect(handler func()) Interface
	Deselect()
	SetOnHover(handler func()) Interface
	Hover()
	SetOffHover(handler func()) Interface
	OffHover()
	SetBorderColor(color tcell.Color) Interface
	SetOnKeyEvent(handler func(*tcell.EventKey, func())) Interface
	HandleKeyEvent(key *tcell.EventKey, deSelector func())

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
	customEvents map[string]func()
	onEvent      func(*tcell.EventKey, func())
	selectable   bool
	hoverable    bool
	bleedThrough bool
	escapable    bool
	borders      bool
}

func MakeElement() *Element {
	return &Element{
		Flex:         tview.NewFlex(),
		customEvents: map[string]func(){},
		selectable:   true,
		hoverable:    true,
		bleedThrough: false,
		escapable:    true,
		borders:      false,
	}
}

func (element *Element) GetFlex() *tview.Flex {
	return element.Flex
}

func (element *Element) SetEvent(event string, handler func()) Interface {
	element.customEvents[event] = handler
	return element
}

func (element *Element) EmitEvent(event string) Interface {
	handler := element.customEvents[event]
	if handler != nil {
		handler()
	}
	return element
}

func (element *Element) SetSelectable(selectable bool) Interface {
	element.selectable = selectable
	return element
}

func (element *Element) Selectable() bool {
	return element.selectable
}

func (element *Element) SetHoverable(hoverable bool) Interface {
	element.hoverable = hoverable
	return element
}

func (element *Element) Hoverable() bool {
	return element.hoverable
}

func (element *Element) SetBleedThrough(bleedThrough bool) Interface {
	element.bleedThrough = bleedThrough
	return element
}

func (element *Element) BleedThrough() bool {
	return element.bleedThrough
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
	element.SetEvent("select", handler)
	return element
}

func (element *Element) Select() {
	element.EmitEvent("select")
}

func (element *Element) SetOnDeselect(handler func()) Interface {
	element.SetEvent("deselect", handler)
	return element
}

func (element *Element) Deselect() {
	element.EmitEvent("deselect")
}

func (element *Element) SetOnHover(handler func()) Interface {
	element.SetEvent("onHover", handler)
	return element
}

func (element *Element) Hover() {
	element.EmitEvent("onHover")
}

func (element *Element) SetOffHover(handler func()) Interface {
	element.SetEvent("offHover", handler)
	return element
}

func (element *Element) OffHover() {
	element.EmitEvent("offHover")
}

func (element *Element) SetBorderColor(color tcell.Color) Interface {
	element.Flex.SetBorderColor(color)
	return element
}

func (element *Element) SetOnKeyEvent(handler func(*tcell.EventKey, func())) Interface {
	element.onEvent = handler
	return element
}

func (element *Element) HandleKeyEvent(key *tcell.EventKey, deSelector func()) {
	if element.onEvent != nil {
		element.onEvent(key, deSelector)
	}
}
