package header

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Option struct {
	*tview.Flex
	text *Text
}

func makeOption(input string) *Option {
	text := makeText(input)

	option := &Option{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow).AddItem(text, 0, 1, false),
		text: text,
	}

	return option
}

func (option *Option) SetBorderColor(color tcell.Color) {
	option.text.Box.SetBorderColor(color)
}
