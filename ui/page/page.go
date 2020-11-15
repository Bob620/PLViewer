package page

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page struct {
	Title      string
	Id         string
	grid       *tview.Grid
	active     bool
	onActivate func()
	selected   bool
	onSelect   func()
	onDeselect func()
}

func MakePage(title string) *Page {
	return &Page{
		Title: title,
		Id:    title,
		grid:  tview.NewGrid(),
	}
}

func (page *Page) SetOnActivate(handler func()) {
	page.onActivate = handler
}

func (page *Page) GetActive() bool {
	return page.active
}

func (page *Page) Activate(active bool) {
	page.active = active
	if page.onActivate != nil {
		page.onActivate()
	}
}

func (page *Page) SetOnSelect(handler func()) {
	page.onSelect = handler
}

func (page *Page) IsSelected() bool {
	return page.selected
}

func (page *Page) Select() {
	page.selected = true
	if page.onSelect != nil {
		page.onSelect()
	}
}

func (page *Page) SetOnDeselect(handler func()) {
	page.onDeselect = handler
}

func (page *Page) Deselect() {
	if page.onDeselect != nil {
		page.onDeselect()
	}
	page.selected = false
}

func (page *Page) Grid() *tview.Grid {
	return page.grid
}

func (page *Page) AddTo(pages *tview.Pages) {
	pages.AddPage(page.Id, page.grid, true, false)
}

func (page *Page) HandleEvents(key *tcell.EventKey, switchToPage func(*Page), focusHeader func(bool)) {
	switch key.Key() {
	case tcell.KeyLeft:
		break
	case tcell.KeyRight:
		break
	case tcell.KeyUp:
		focusHeader(true)
		break
	case tcell.KeyDown:
		break
	}
}
