package header

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Text struct {
	*tview.Box
	input string
}

func makeText(input string) *Text {
	text := &Text{
		Box:   tview.NewBox().SetBorder(true),
		input: input,
	}

	return text
}

func (text *Text) Draw(screen tcell.Screen) {
	text.Box.Draw(screen)
	x, y, width, _ := text.GetInnerRect()

	tview.Print(screen, text.input, x, y, width, tview.AlignCenter, tcell.ColorWhite)
}

func (text *Text) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return text.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

	})
}
