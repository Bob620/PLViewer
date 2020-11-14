package page

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page struct {
	Title string
	Id    string
	grid  *tview.Grid
}

func MakePage(title string) *Page {
	return &Page{
		Title: title,
		Id:    title,
		grid:  tview.NewGrid().AddItem(tview.NewTextView().SetText(title), 0, 0, 1, 1, 0, 0, true),
	}
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
