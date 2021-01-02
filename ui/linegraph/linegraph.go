package linegraph

import (
	"PLViewer/ui/element"
	"github.com/exrook/drawille-go"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strings"
)

type LineGraph struct {
	*element.Element
	drawing drawille.Canvas
	height  int
	width   int
	line    []float64
	force   bool
}

func MakeLineGraph() *LineGraph {
	lineGraph := &LineGraph{
		Element: element.MakeElement(),
		drawing: drawille.NewCanvas(),
		line:    []float64{},
	}

	return lineGraph
}

func (lineGraph *LineGraph) SetLine(line []float64) {
	lineGraph.line = line
	if line == nil {
		lineGraph.drawing.Clear()
	} else {
		lineGraph.force = true
	}
}

func (lineGraph *LineGraph) Draw(screen tcell.Screen) {
	lineGraph.Flex.Draw(screen)
	x, y, width, height := lineGraph.GetInnerRect()
	height -= 1
	lineLength := len(lineGraph.line)

	if (lineGraph.force || (lineGraph.height != 4*height || lineGraph.width != 2*width)) && lineLength > 0 {
		lineGraph.force = false
		lineGraph.height = 4 * height
		lineGraph.width = 2 * width

		lineGraph.drawing.Clear()
		lineGraph.drawing.Frame(0, 0, lineGraph.width, lineGraph.height)

		poses := []float64{}
		min := -1.0
		max := 0.0

		for i := 0; i < lineGraph.width; i++ {
			lastPos := int(float64(lineLength/lineGraph.width) * (float64(i) - 0.5))
			if lastPos < 0 {
				lastPos = 0
			}
			pos := int(float64(lineLength/lineGraph.width) * (float64(i) + 0.5))
			if pos > len(lineGraph.line) {
				pos = len(lineGraph.line)
			}

			lineHeight := 0.0
			for _, height := range lineGraph.line[lastPos:pos] {
				lineHeight += height
			}
			lineHeight = lineHeight / (float64(pos) - float64(lastPos))

			poses = append(poses, lineHeight)

			if min > lineHeight || min < 0 {
				min = lineHeight
			}

			if max < lineHeight {
				max = lineHeight
			}
		}

		graphMin := (min / max) * float64(lineGraph.height)
		for i, lineHeight := range poses {
			lineGraph.drawing.Set(i, lineGraph.height-int(((lineHeight/max)*float64(lineGraph.height))-graphMin))
		}

	}

	output := strings.Split(lineGraph.drawing.String(), "\n")

	for i, line := range output {
		tview.Print(screen, line, x, y+i, width, tview.AlignLeft, tcell.ColorWhite)
	}
}
