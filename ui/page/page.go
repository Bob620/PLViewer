package page

import (
	"PLViewer/ui/element"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page struct {
	Title      string
	Id         string
	grid       *tview.Grid
	onSelect   func()
	onDeselect func()
	elements   map[string]*element.Element
	layout     *Layout
	selected   string
	hasActive  bool
}

func MakePage(title string, navigation *Layout) *Page {
	grid := tview.NewGrid()
	navigation.Setup(grid)

	return &Page{
		Title:    title,
		Id:       title,
		layout:   navigation,
		grid:     grid,
		elements: map[string]*element.Element{},
	}
}

func (page *Page) SetOnSelect(handler func()) {
	page.onSelect = handler
}

func (page *Page) Select() {
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

	page.selectElement("")
}

func (page *Page) Grid() *tview.Grid {
	return page.grid
}

func (page *Page) AddTo(pages *tview.Pages) {
	pages.AddPage(page.Id, page.grid, true, false)
}

func (page *Page) AddElement(ele *element.Element, position string) *Page {
	page.elements[position] = ele

	pos := page.layout.GetPos(position)
	page.grid.AddItem(ele.Flex(), pos.Pos[0], pos.Pos[1], pos.Height, pos.Width, 0, 0, false)
	return page
}

func (page *Page) selectElement(element string) {
	if element == "" {
		if page.elements[page.selected] != nil {
			page.elements[page.selected].SetBorderColor(tcell.ColorWhite)
		}
		page.selected = ""
	} else if page.elements[element] != nil && page.selected != element {
		if page.elements[element].Selectable() {
			if page.elements[page.selected] != nil {
				page.elements[page.selected].SetBorderColor(tcell.ColorWhite)
			}
			page.selected = element
			page.elements[page.selected].SetBorderColor(tcell.ColorBlue)
		}
	}
}

func (page *Page) SelectElement(element string) {
	if element != "" {
		page.selectElement(element)
	}
}

func (page *Page) navigate(direction string) bool {
	ele := page.selected
	for {
		switch direction {
		case "left":
			ele = page.layout.Left(ele)
			break
		case "right":
			ele = page.layout.Right(ele)
			break
		case "up":
			ele = page.layout.Up(ele)
			break
		case "down":
			ele = page.layout.Down(ele)
			break
		}

		if ele == "" {
			return false
		}

		if page.elements[ele].Selectable() {
			page.SelectElement(ele)
			return true
		}
	}
}

func (page *Page) HandleEvents(key *tcell.EventKey, switchToPage func(*Page), focusHeader func(bool)) {
	if page.selected != "" {
		if page.hasActive {
			switch key.Key() {
			case tcell.KeyEscape:
				page.hasActive = false
				page.elements[page.selected].Deselect()
				break
			default:
				page.elements[page.selected].HandleEvents(key)
				break
			}
		} else {
			switch key.Key() {
			case tcell.KeyEnter:
				page.hasActive = true
				page.elements[page.selected].Select()
				break
			case tcell.KeyLeft:
				page.navigate("left")
				break
			case tcell.KeyRight:
				page.navigate("right")
				break
			case tcell.KeyUp:
				if !page.navigate("up") {
					focusHeader(true)
				}
				break
			case tcell.KeyDown:
				page.navigate("down")
				break
			}
		}
	}
}
