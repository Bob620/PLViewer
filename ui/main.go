package ui

import (
	"PLViewer/backend"
	"PLViewer/ui/creator"
	Header "PLViewer/ui/header"
	"PLViewer/ui/page"
	"PLViewer/ui/processor"
	"PLViewer/ui/viewer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pages *tview.Pages

var header *Header.Header

var PageViewer *viewer.Viewer
var PageCreator *creator.Creator
var PageProcessor *processor.Processor

var activePage *page.Page

func Initialize(bg *backend.Backend) {
	PageViewer = viewer.MakeViewerPage(app, bg)
	PageCreator = creator.MakeCreatorPage(app)
	PageProcessor = processor.MakeProcessorPage(app, bg)
	activePage = PageViewer.Page

	pages = tview.NewPages()

	header = Header.MakeHeader([]*page.Page{PageViewer.Page, PageCreator.Page, PageProcessor.Page})

	grid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(pages, 1, 0, 1, 1, 0, 0, false)

	PageViewer.AddTo(pages)
	PageCreator.AddTo(pages)
	PageProcessor.AddTo(pages)

	switchToPage(PageViewer.Page)
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
		case PageViewer.Page:
			PageViewer.HandleKeyEvent(key, switchToPage, focusHeader)
			break
		case PageCreator.Page:
			PageCreator.HandleKeyEvent(key, switchToPage, focusHeader)
			break
		case PageProcessor.Page:
			PageProcessor.HandleKeyEvent(key, switchToPage, focusHeader)
			break
		}
	}

	return key
}
