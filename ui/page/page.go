package page

import (
	"PLViewer/ui/element"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Interface interface {
	SetOnSelect(handler func())
	Select()
	SetOnDeselect(handler func())
	Deselect()
	SetPreDraw(handler func())
	AddTo(pages *tview.Pages)
	Clear()
	AddElement(ele element.Interface, position string) *Page
	AddPageAsElement(newPage Interface, position string) *Page
	selectElement(element string)
	SelectElement(element string)
	HoveredElement() string
	SelectedElement() string
	navigate(direction string) bool
	DeselectActiveElement()
	HandleEvents(key *tcell.EventKey, switchToPage func(*Page), focusHeader func(bool))
}

type Page struct {
	*tview.Grid
	Title      string
	Id         string
	onSelect   func()
	onDeselect func()
	onPreDraw  func()
	elements   map[string]element.Interface
	layout     *Layout
	selected   string
	hasActive  bool
}

func MakePage(title string, layout *Layout) *Page {
	grid := tview.NewGrid()
	layout.Setup(grid)

	return &Page{
		Grid:     grid,
		Title:    title,
		Id:       title,
		layout:   layout,
		elements: map[string]element.Interface{},
	}
}

func (page *Page) SetLayout(layout *Layout) {
	page.layout = layout
	layout.Setup(page.Grid)
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

func (page *Page) SetPreDraw(handler func()) {
	page.onPreDraw = handler
}

func (page *Page) AddTo(pages *tview.Pages) {
	pages.AddPage(page.Id, page.Grid, true, false)
}

func (page *Page) Clear() {
	page.Deselect()
	page.hasActive = false
	page.elements = map[string]element.Interface{}
	page.Grid.Clear()
}

func (page *Page) AddElement(ele element.Interface, position string) *Page {
	page.elements[position] = ele

	pos := page.layout.GetPos(position)
	page.Grid.AddItem(ele, pos.Pos[0], pos.Pos[1], pos.Height, pos.Width, 0, 0, false)
	return page
}

func (page *Page) selectElement(element string) {
	if element == "" {
		if page.elements[page.selected] != nil {
			page.elements[page.selected].SetBorderColor(tcell.ColorWhite)
			page.elements[page.selected].OffHover()
		}
		page.selected = ""
	} else if page.elements[element] != nil && page.selected != element {
		if page.elements[element].Hoverable() {
			if page.elements[page.selected] != nil {
				page.elements[page.selected].SetBorderColor(tcell.ColorWhite)
				page.elements[page.selected].OffHover()
			}
			page.selected = element
			page.elements[page.selected].SetBorderColor(tcell.ColorBlue)
			page.elements[page.selected].Hover()
		}
	}
}

func (page *Page) SelectElement(element string) {
	if element != "" {
		page.selectElement(element)
	}
}

func (page *Page) HoveredElement() string {
	return page.selected
}

func (page *Page) SelectedElement() string {
	if page.hasActive {
		return page.selected
	}
	return ""
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

		if page.elements[ele] != nil && page.elements[ele].Hoverable() {
			page.SelectElement(ele)
			return true
		}
	}
}

func (page *Page) DeselectActiveElement() {
	page.hasActive = false
	page.elements[page.selected].Deselect()
}

func (page *Page) HandleKeyEvent(key *tcell.EventKey, switchToPage func(*Page), focusHeader func(bool)) {
	if page.selected != "" {
		if page.hasActive {
			page.elements[page.selected].HandleKeyEvent(key, page.DeselectActiveElement)
			switch key.Key() {
			case tcell.KeyEscape:
				if page.elements[page.selected].Escapable() {
					page.DeselectActiveElement()
				}
				break
			default:
				break
			}
		} else {
			switch key.Key() {
			case tcell.KeyEnter:
				if page.elements[page.selected].Selectable() {
					page.hasActive = true
					page.elements[page.selected].Select()
				}
				break
			case tcell.KeyLeft:
				if !page.navigate("left") && page.elements[page.selected].BleedThrough() {
					page.elements[page.selected].HandleKeyEvent(key, page.DeselectActiveElement)
				}
				break
			case tcell.KeyRight:
				if !page.navigate("right") && page.elements[page.selected].BleedThrough() {
					page.elements[page.selected].HandleKeyEvent(key, page.DeselectActiveElement)
				}
				break
			case tcell.KeyUp:
				if !page.navigate("up") {
					if page.elements[page.selected].BleedThrough() {
						page.elements[page.selected].HandleKeyEvent(key, page.DeselectActiveElement)
					}
					focusHeader(true)
				}
				break
			case tcell.KeyDown:
				if !page.navigate("down") && page.elements[page.selected].BleedThrough() {
					page.elements[page.selected].HandleKeyEvent(key, page.DeselectActiveElement)
				}
				break
			}
		}
	}
}

func (page *Page) Draw(screen tcell.Screen) {
	if page.onPreDraw != nil {
		page.onPreDraw()
	}
	page.Grid.Draw(screen)
}
