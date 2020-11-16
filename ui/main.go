package ui

import (
	"PLViewer/ui/creator"
	Header "PLViewer/ui/header"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pages *tview.Pages

var header *Header.Header

var PageViewer = page.MakePage("Viewer", page.MakeLayout([][]string{{""}}))
var PageCreator = creator.MakeCreatorPage(app)

var activePage = PageViewer

func Initialize() {
	pages = tview.NewPages()

	header = Header.MakeHeader([]*page.Page{PageViewer, PageCreator.Page})

	grid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(pages, 1, 0, 1, 1, 0, 0, false)

	PageViewer.AddTo(pages)
	PageCreator.AddTo(pages)

	switchToPage(PageViewer)
	header.Focused(true)

	app.SetInputCapture(eventHandler)
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

func switchToPage(page *page.Page) {
	activePage = page
	pages.SwitchToPage(page.Id)
}

func focusHeader(yes bool) {
	header.Focused(yes)
	if yes {
		activePage.Deselect()
	} else {
		activePage.Select()
	}
}

func eventHandler(key *tcell.EventKey) *tcell.EventKey {
	if header.IsFocused {
		if key.Rune() == 'q' {
			app.Stop()
			return nil
		}

		header.HandleEvents(key, switchToPage, focusHeader)
	} else {
		switch activePage {
		case PageViewer:
			PageViewer.HandleEvents(key, switchToPage, focusHeader)
			break
		case PageCreator.Page:
			PageCreator.HandleEvents(key, switchToPage, focusHeader)
			break
		}
	}

	return key
}
