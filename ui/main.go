package ui

import (
	Header "PLViewer/ui/header"
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

var header *Header.Header

var PageViewer = page.MakePage("Viewer")
var PageCreator = page.MakePage("Creator")

var activePage = PageViewer

func Initialize() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	header = Header.MakeHeader([]*page.Page{PageViewer, PageCreator})

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

func eventHandler(key *tcell.EventKey) *tcell.EventKey {
	if key.Rune() == 'q' {
		app.Stop()
		return nil
	}

	if key.Key() == tcell.KeyEscape {
		header.Focused(true)
	}

	if header.IsFocused {
		header.HandleEvents(key, switchToPage, header.Focused)
	} else {
		switch activePage {
		case PageViewer:
			PageViewer.HandleEvents(key, switchToPage, header.Focused)
			break
		case PageCreator:
			PageCreator.HandleEvents(key, switchToPage, header.Focused)
			break
		}
	}

	return key
}
