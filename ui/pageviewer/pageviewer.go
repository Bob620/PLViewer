package pageviewer

import (
	"PLViewer/ui/element"
	"github.com/rivo/tview"
)

type PageViewer struct {
	*element.Element
	Pages *tview.Pages
}

func MakePageViewer() *PageViewer {
	pageViewer := PageViewer{
		Element: element.MakeElement(),
		Pages:   tview.NewPages(),
	}

	pageViewer.AddItem(pageViewer.Pages, 0, 1, false)

	return &pageViewer
}
