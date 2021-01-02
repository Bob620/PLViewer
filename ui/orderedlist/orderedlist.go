package orderedlist

import (
	"PLViewer/ui/element"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type item struct {
	*element.Element
	check  *element.Element
	user   *element.Element
	id     string
	height int
}

type OrderedList struct {
	*element.Element
	hoveredPosition int
	selected        bool
	items           []*item
}

func MakeOrderedList() *OrderedList {
	orderedList := &OrderedList{
		Element: element.MakeElement(),
	}

	orderedList.SetOnSelect(func() {
		if orderedList.hoveredPosition == -1 {
			orderedList.hoverPosition(0)
		} else {
			orderedList.hoverPosition(orderedList.hoveredPosition)
		}
	}).SetOnDeselect(func() {
		if orderedList.hoveredPosition != -1 {
			orderedList.items[orderedList.hoveredPosition].OffHover()
			orderedList.items[orderedList.hoveredPosition].user.OffHover()
		}
	}).SetOnKeyEvent(func(key *tcell.EventKey, deSelector func()) {
		switch key.Key() {
		case tcell.KeyEscape:
			if orderedList.hoveredPosition != -1 {
				orderedList.items[orderedList.hoveredPosition].OffHover()
				orderedList.items[orderedList.hoveredPosition].user.OffHover()
			}
			break
		case tcell.KeyEnter:
			orderedList.SelectPosition(orderedList.hoveredPosition)
			break
		case tcell.KeyUp:
			if orderedList.selected {
				orderedList.MoveElementUp(orderedList.hoveredPosition)
			}
			orderedList.MoveUp()
			break
		case tcell.KeyDown:
			if orderedList.selected {
				orderedList.MoveElementDown(orderedList.hoveredPosition)
			}
			orderedList.MoveDown()
			break
		}
	})

	orderedList.hoveredPosition = -1
	return orderedList
}

func (orderedList *OrderedList) sanitizePosition(position int) int {
	length := len(orderedList.items)
	if length == -1 {
		length = 0
	}

	if position < 0 {
		return 0
	} else if position > length {
		return length
	} else {
		return position
	}
}

func (orderedList *OrderedList) SelectPosition(position int) {
	if orderedList.hoveredPosition != -1 {
		orderedList.items[orderedList.hoveredPosition].Deselect()
		orderedList.items[orderedList.hoveredPosition].user.Deselect()
	}

	if position == -1 || (orderedList.hoveredPosition == position && orderedList.selected) {
		orderedList.selected = false
	} else {
		position = orderedList.sanitizePosition(position)

		orderedList.selected = true
		orderedList.hoverPosition(position)
		orderedList.items[position].Select()
		orderedList.items[position].user.Select()
	}
}

func (orderedList *OrderedList) hoverPosition(position int) {
	if orderedList.hoveredPosition != -1 && len(orderedList.items) > orderedList.hoveredPosition {
		orderedList.items[orderedList.hoveredPosition].OffHover()
		orderedList.items[orderedList.hoveredPosition].user.OffHover()
	}

	if position == -1 || len(orderedList.items) == 0 {
		orderedList.hoveredPosition = -1
	} else {
		orderedList.hoveredPosition = orderedList.sanitizePosition(position)
		orderedList.items[orderedList.hoveredPosition].Hover()
		orderedList.items[orderedList.hoveredPosition].user.Hover()
	}
}

func (orderedList *OrderedList) AddItem(elem *element.Element, position int, id string) int {
	if position == -1 {
		position = len(orderedList.items)
	}
	position = orderedList.sanitizePosition(position)

	checkbox := tview.NewCheckbox().SetLabel(" ")
	checkEle := element.MakeElement()
	ele := element.MakeElement()
	ele.SetOnSelect(func() {
		checkbox.SetChecked(true)
	}).SetOnDeselect(func() {
		checkbox.SetChecked(false)
	}).SetOnHover(func() {
		checkbox.SetFieldBackgroundColor(tcell.ColorWhite)
		checkbox.SetFieldTextColor(tcell.ColorBlue)
	}).SetOffHover(func() {
		checkbox.SetFieldBackgroundColor(tcell.ColorBlue)
		checkbox.SetFieldTextColor(tcell.ColorWhite)
	})

	checkEle.AddItem(checkbox, 0, 1, false).SetDirection(tview.FlexRow)
	ele.AddItem(checkEle, 3, 1, false)
	ele.AddItem(elem, 0, 1, false)

	var newItems []*item
	if position != 0 {
		newItems = append(newItems, orderedList.items[:position]...)
	}

	_, _, _, userHeight := elem.GetRect()
	newItems = append(newItems, &item{
		Element: ele,
		check:   checkEle,
		user:    elem,
		id:      id,
		height:  userHeight,
	})

	if len(orderedList.items) > position {
		newItems = append(newItems, orderedList.items[position:]...)
	}

	orderedList.items = newItems
	if orderedList.hoveredPosition >= position {
		orderedList.MoveDown()
	}

	return position
}

func (orderedList *OrderedList) DeleteId(id string) {
	for i, item := range orderedList.items {
		if item.id == id {
			orderedList.DeleteItem(i)
		}
	}
}

func (orderedList *OrderedList) DeleteItem(position int) {
	if len(orderedList.items) == 0 {
		return
	}

	position = orderedList.sanitizePosition(position)

	orderedList.items[position].user.EmitEvent("delete")

	var newItems []*item
	newItems = append(newItems, orderedList.items[:position]...)
	if len(orderedList.items) > position+1 {
		newItems = append(newItems, orderedList.items[position+1:]...)
	}
	orderedList.items = newItems

	if orderedList.hoveredPosition == position {
		orderedList.selected = false
	}

	if orderedList.hoveredPosition >= position {
		orderedList.MoveUp()
	}
}

func (orderedList *OrderedList) MoveDown() {
	newPos := orderedList.hoveredPosition + 1
	maxPos := len(orderedList.items) - 1
	if newPos > maxPos {
		newPos = maxPos
	}

	orderedList.hoverPosition(newPos)
}

func (orderedList *OrderedList) MoveUp() {
	newPos := orderedList.hoveredPosition - 1
	if newPos == -1 {
		newPos = 0
	}

	orderedList.hoverPosition(newPos)
}

func (orderedList *OrderedList) MoveElementUp(position int) {
	position = orderedList.sanitizePosition(position)
	if position != 0 {
		temp := orderedList.items[position-1]
		orderedList.items[position-1] = orderedList.items[position]
		orderedList.items[position] = temp
	}
}

func (orderedList *OrderedList) MoveElementDown(position int) {
	position = orderedList.sanitizePosition(position)
	if position != len(orderedList.items)-1 {
		temp := orderedList.items[position]
		orderedList.items[position] = orderedList.items[position+1]
		orderedList.items[position+1] = temp
	}
}

func (orderedList *OrderedList) MoveElementTop(position int) {
	position = orderedList.sanitizePosition(position)
	if position != 0 {
		temp := orderedList.items[position]
		orderedList.DeleteItem(position)
		items := orderedList.items
		var newItems []*item
		orderedList.items = append(newItems, temp)
		orderedList.items = append(newItems, items...)
	}
}

func (orderedList *OrderedList) MoveElementBottom(position int) {
	position = orderedList.sanitizePosition(position)
	if position != len(orderedList.items)-1 {
		temp := orderedList.items[position]
		orderedList.DeleteItem(position)
		orderedList.items = append(orderedList.items, temp)
	}
}

func (orderedList *OrderedList) Serialize() []string {
	var items []string

	for _, item := range orderedList.items {
		items = append(items, item.id)
	}

	return items
}

func (orderedList *OrderedList) Draw(screen tcell.Screen) {
	orderedList.Element.Draw(screen)
	x, y, width, height := orderedList.GetInnerRect()

	for _, item := range orderedList.items {
		if height > 0 {
			itemHeight := item.height
			if item.Borders() {
				itemHeight = itemHeight + 2
			}
			height = height - itemHeight

			if height < 0 {
				itemHeight = height + itemHeight
			}

			item.check.SetRect(x, y, 3, itemHeight)
			item.SetRect(x, y, width, itemHeight)
			//			item.SetBorders(true)
			item.Draw(screen)

			y += itemHeight
		}
	}
}
