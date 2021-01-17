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
	SetOnNavigateOutside(handler func(direction string))
	SetPreDraw(handler func())
	AddTo(pages *tview.Pages)
	Clear()
	AddElement(ele element.Interface, position string) *Page
	AddPageAsElement(newPage Interface, position string) *Page
	deleteElement(position string) *Page
	selectElement(element string)
	SelectElement(element string)
	HoveredElement() string
	SelectedElement() string
	NavigateUp() bool
	NavigateDown() bool
	NavigateLeft() bool
	NavigateRight() bool
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
	onOverflow func(string)
	onPreDraw  func()
	elements   map[string]element.Interface
	layout     *Layout
	hovered    string
	isActive   bool
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

func (page *Page) SetOnNavigateOutside(handler func(direction string)) {
	page.onOverflow = handler
}

func (page *Page) SetPreDraw(handler func()) {
	page.onPreDraw = handler
}

func (page *Page) AddTo(pages *tview.Pages) {
	pages.AddPage(page.Id, page.Grid, true, false)
}

func (page *Page) Clear() {
	page.Deselect()
	page.isActive = false
	page.elements = map[string]element.Interface{}
	page.Grid.Clear()
}

func (page *Page) AddElement(ele element.Interface, position string) *Page {
	page.elements[position] = ele

	pos := page.layout.GetPos(position)
	page.Grid.AddItem(ele, pos.Pos[0], pos.Pos[1], pos.Height, pos.Width, 0, 0, false)
	return page
}

func (page *Page) DeleteElement(position string) *Page {
	page.Grid.RemoveItem(page.elements[position])
	page.elements[position] = nil
	return page
}

func (page *Page) selectElement(element string) {
	if element == "" {
		if page.elements[page.hovered] != nil {
			page.elements[page.hovered].SetBorderColor(tcell.ColorWhite)
			page.elements[page.hovered].OffHover()
		}
		page.hovered = ""
	} else if page.elements[element] != nil && page.hovered != element {
		if page.elements[element].Hoverable() {
			if page.elements[page.hovered] != nil {
				page.elements[page.hovered].SetBorderColor(tcell.ColorWhite)
				page.elements[page.hovered].OffHover()
			}
			page.hovered = element
			page.elements[page.hovered].SetBorderColor(tcell.ColorBlue)
			page.elements[page.hovered].Hover()
		}
	}
}

func (page *Page) SelectElement(element string) {
	if element != "" {
		page.selectElement(element)
	}
}

func (page *Page) HoveredElement() string {
	return page.hovered
}

func (page *Page) SelectedElement() string {
	if page.isActive {
		return page.hovered
	}
	return ""
}

func (page *Page) NavigateUp() bool {
	page.DeselectActiveElement()
	return page.navigate("up", "")
}

func (page *Page) NavigateDown() bool {
	page.DeselectActiveElement()
	return page.navigate("down", "")
}

func (page *Page) NavigateLeft() bool {
	page.DeselectActiveElement()
	return page.navigate("left", "")
}

func (page *Page) NavigateRight() bool {
	page.DeselectActiveElement()
	return page.navigate("right", "")
}

func (page *Page) navigate(direction string, ele string) bool {
	if ele == "" {
		ele = page.hovered
	}
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

		if page.elements[ele] != nil {
			if page.elements[ele].Hoverable() {
				page.SelectElement(ele)
				return true
			} else {
				if direction == "up" || direction == "down" {
					if !page.navigate("left", ele) {
						return page.navigate("right", ele)
					}
					return true
				} else if direction == "left" || direction == "right" {
					if !page.navigate("up", ele) {
						return page.navigate("down", ele)
					}
					return true
				}
			}
		}
	}
}

func (page *Page) DeselectActiveElement() {
	page.isActive = false
	page.elements[page.hovered].Deselect()
}

func (page *Page) HandleKeyEvent(key *tcell.EventKey, switchToPage func(*Page), focusHeader func(bool)) {
	if page.hovered != "" {
		if page.isActive {
			page.elements[page.hovered].HandleKeyEvent(key, page.DeselectActiveElement)
			switch key.Key() {
			case tcell.KeyEscape:
				if page.elements[page.hovered].Escapable() {
					page.DeselectActiveElement()
				}
				break
			default:
				break
			}
		} else {
			navigated := false
			switch key.Key() {
			case tcell.KeyEnter:
				if page.elements[page.hovered].Selectable() {
					page.isActive = true
					page.elements[page.hovered].Select()
				}
				break
			case tcell.KeyLeft:
				if !page.navigate("left", "") {
					if page.elements[page.hovered].BleedThrough() {
						page.elements[page.hovered].HandleKeyEvent(key, page.DeselectActiveElement)
					}
					if page.onOverflow != nil {
						page.onOverflow("left")
					}
				} else {
					navigated = true
				}
				break
			case tcell.KeyRight:
				if !page.navigate("right", "") {
					if page.elements[page.hovered].BleedThrough() {
						page.elements[page.hovered].HandleKeyEvent(key, page.DeselectActiveElement)
					}
					if page.onOverflow != nil {
						page.onOverflow("right")
					}
				} else {
					navigated = true
				}
				break
			case tcell.KeyUp:
				if !page.navigate("up", "") {
					if page.elements[page.hovered].BleedThrough() {
						page.elements[page.hovered].HandleKeyEvent(key, page.DeselectActiveElement)
					}
					if page.onOverflow != nil {
						page.onOverflow("up")
					}
					focusHeader(true)
				} else {
					navigated = true
				}
				break
			case tcell.KeyDown:
				if !page.navigate("down", "") {
					if page.elements[page.hovered].BleedThrough() {
						page.elements[page.hovered].HandleKeyEvent(key, page.DeselectActiveElement)
					}
					if page.onOverflow != nil {
						page.onOverflow("down")
					}
				} else {
					navigated = true
				}
				break
			}

			if navigated && page.elements[page.hovered] != nil && page.elements[page.hovered].Inlay() {
				page.isActive = true
				page.elements[page.hovered].Select()
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
