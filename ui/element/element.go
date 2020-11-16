package element

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Element struct {
	flex       *tview.Flex
	onSelect   func()
	onDeselect func()
	onEvent    func(*tcell.EventKey)
	selectable bool
	escapable  bool
	borders    bool
}

func MakeElement() *Element {
	return &Element{
		flex:       tview.NewFlex(),
		selectable: true,
		escapable:  true,
		borders:    false,
	}
}

func (element *Element) SetSelectable(selectable bool) *Element {
	element.selectable = selectable
	return element
}

func (element *Element) Selectable() bool {
	return element.selectable
}

func (element *Element) SetEscapable(escapable bool) *Element {
	element.escapable = escapable
	return element
}

func (element *Element) Escapable() bool {
	return element.escapable
}

func (element *Element) SetBorders(borders bool) *Element {
	element.borders = borders
	element.flex.SetBorder(borders)
	return element
}

func (element *Element) Borders() bool {
	return element.borders
}

func (element *Element) Flex() *tview.Flex {
	return element.flex
}

func (element *Element) SetOnSelect(handler func()) *Element {
	element.onSelect = handler
	return element
}

func (element *Element) Select() {
	if element.onSelect != nil {
		element.onSelect()
	}
}

func (element *Element) SetOnDeselect(handler func()) *Element {
	element.onDeselect = handler
	return element
}

func (element *Element) Deselect() {
	if element.onDeselect != nil {
		element.onDeselect()
	}
}

func (element *Element) SetBorderColor(color tcell.Color) *Element {
	element.flex.SetBorderColor(color)
	return element
}

func (element *Element) SetOnEvent(handler func(*tcell.EventKey)) *Element {
	element.onEvent = handler
	return element
}

func (element *Element) HandleEvents(key *tcell.EventKey) {
	if element.onEvent != nil {
		element.onEvent(key)
	}
}
