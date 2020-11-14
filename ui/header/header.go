package header

import (
	"PLViewer/ui/page"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Header struct {
	*tview.Flex
	elements      []*Option
	pages         []*page.Page
	currentOption int
	IsFocused     bool
}

func MakeHeader(options []*page.Page) *Header {
	var elements = make([]*Option, len(options))
	var pages = make([]*page.Page, len(options))

	for index, page := range options {
		elements[index] = makeOption(page.Title)
		pages[index] = page
	}

	header := &Header{
		Flex:          tview.NewFlex(),
		elements:      elements,
		pages:         pages,
		currentOption: 0,
		IsFocused:     false,
	}

	for _, item := range elements {
		header.AddItem(tview.NewBox(), 0, 1, false)
		header.AddItem(item, 0, 2, false)
		header.AddItem(tview.NewBox(), 0, 1, false)
	}

	return header
}

func (header *Header) Focused(yes bool) {
	header.IsFocused = yes
}

func (header *Header) HandleEvents(key *tcell.EventKey, switchToPage func(*page.Page), focusHeader func(bool)) {
	switch key.Key() {
	case tcell.KeyLeft:
		header.currentOption--
		if header.currentOption < 0 {
			header.currentOption = len(header.elements) - 1
		}

		switchToPage(header.pages[header.currentOption])
		break
	case tcell.KeyRight:
		header.currentOption++
		if header.currentOption >= len(header.elements) {
			header.currentOption = 0
		}

		switchToPage(header.pages[header.currentOption])
		break
	case tcell.KeyUp:
		break
	case tcell.KeyDown:
		focusHeader(false)
		break
	}
}

func (header *Header) Draw(screen tcell.Screen) {
	for index, element := range header.elements {
		if index == header.currentOption {
			if header.IsFocused {
				element.SetBorderColor(tcell.ColorBlue)
			} else {
				element.SetBorderColor(tcell.ColorLightBlue)
			}
		} else {
			element.SetBorderColor(tcell.ColorWhite)
		}
	}
	header.Flex.Draw(screen)
}
