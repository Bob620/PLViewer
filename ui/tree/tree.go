package tree

import (
	"PLViewer/ui/element"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

type nodePos struct {
	node         *Node
	prefix       string
	depth        int
	hasNext      bool
	previousNode *nodePos
}

type Tree struct {
	*element.Element
	//	application *tview.Application
	nodes          map[Id]*Node
	nodeOrder      []Id
	currentNodes   []*nodePos
	cursorPos      int
	availableNodes int
}

var nodeId int

func MakeTree( /*application *tview.Application*/ ) *Tree {
	tree := &Tree{
		Element:   element.MakeElement(),
		nodes:     map[Id]*Node{},
		nodeOrder: []Id{},
		cursorPos: 0,
		//		application: application,
	}

	tree.Element.SetOnSelect(func() {

	})
	tree.Element.SetOnDeselect(func() {

	})
	tree.Element.SetOnEvent(func(key *tcell.EventKey) {
		switch key.Key() {
		case tcell.KeyUp:
			if tree.cursorPos > 0 {
				tree.cursorPos--
			}
			break
		case tcell.KeyDown:
			if tree.cursorPos < tree.availableNodes {
				tree.cursorPos++
			}
			break
		case tcell.KeyEnter:
			cursorPos := tree.cursorPos
			_, _, _, height := tree.GetInnerRect()
			maxNodes := height - 2
			if cursorPos > (maxNodes / 2) {
				skipNodes := cursorPos - (maxNodes / 2)
				cursorPos = cursorPos - skipNodes
			}

			node := tree.currentNodes[cursorPos]
			if node.node.HasNodes() {
				node.node.Toggle()
			}
			break
		}
	})

	return tree
}

func (tree *Tree) GetNode(nodeId string) *Node {
	return tree.nodes[Id(nodeId)]
}

func (tree *Tree) DeleteNode(nodeId Id) {
	if tree.nodes[nodeId] != nil {
		delete(tree.nodes, nodeId)
		index := 0
		for i, id := range tree.nodeOrder {
			if id == nodeId {
				index = i
				break
			}
		}

		tree.nodeOrder = append(tree.nodeOrder[0:index], tree.nodeOrder[index+1:]...)
	}
}

func (tree *Tree) AddNode(text string, lazyLoadSubNodes bool, onSelect func(nodeId string, node Node)) *Node {
	id := Id(strconv.Itoa(nodeId))
	nodeId++

	if tree.nodes[id] == nil {
		tree.nodes[id] = &Node{
			id:       id,
			text:     text,
			onSelect: onSelect,
			hasNodes: lazyLoadSubNodes,
		}
		tree.nodeOrder = append(tree.nodeOrder, id)
	}

	return tree.nodes[id]
}

// Error with 2/maxHeight cursor depth
func (tree *Tree) Draw(screen tcell.Screen) {
	tree.Flex.Draw(screen)
	x, y, width, height := tree.GetInnerRect()

	cursorPos := tree.cursorPos
	totalPoses := len(tree.nodeOrder)
	totalNodes := len(tree.nodeOrder)
	maxNodes := height - 2
	skipNodes := 0
	availableNodes := 0

	if cursorPos > (maxNodes / 2) {
		skipNodes = cursorPos - (maxNodes / 2)
		cursorPos = cursorPos - skipNodes
	}
	extraAbove := skipNodes

	var nodes []*nodePos
	nodesToVisit := []*nodePos{}

	pos := 0
	for len(nodes) < maxNodes && (pos < totalPoses || len(nodesToVisit) > 0) {
		var node *nodePos
		if len(nodesToVisit) == 0 && pos < totalPoses {
			node = &nodePos{
				node:    tree.nodes[tree.nodeOrder[pos]],
				prefix:  "├─",
				depth:   0,
				hasNext: pos+1 < totalPoses,
			}
			pos++
			if pos == totalPoses {
				node.prefix = "└─"
			}
		} else {
			node = nodesToVisit[0]
			nodesToVisit = nodesToVisit[1:]

			prefix := ""
			thisNode := node
			for i := 0; i < node.depth; i++ {
				thisNode = thisNode.previousNode
				if thisNode.hasNext {
					prefix = "│ " + prefix
				} else {
					prefix = "  " + prefix
				}
			}
			node.prefix = prefix + node.prefix
		}

		if node.node.open && node.node.HasNodes() {
			newNodes := node.node.GetNodes()
			var temp []*nodePos
			for _, newNode := range newNodes {
				totalNodes++
				tree.nodes[newNode.id] = newNode
				temp = append(temp, &nodePos{
					node:         newNode,
					prefix:       "├─",
					depth:        node.depth + 1,
					hasNext:      true,
					previousNode: node,
				})
			}
			if len(temp) > 0 {
				temp[len(temp)-1].prefix = "└─"
				temp[len(temp)-1].hasNext = false
			}

			nodesToVisit = append(temp, nodesToVisit...)
		}

		availableNodes++
		if skipNodes > 0 {
			skipNodes--
		} else {
			nodes = append(nodes, node)
		}
	}

	tree.availableNodes = availableNodes - 1
	tree.currentNodes = nodes
	for i, node := range nodes {
		style := ""
		if node.node.HasNodes() {
			style = "[green]"
		}

		if i == cursorPos {
			style = "[black:green]"
		}

		tview.Print(screen, fmt.Sprintf("%s%s%s", node.prefix, style, node.node.text), x, y+i+1, width, tview.AlignLeft, tcell.ColorWhite)
	}

	extraBelow := totalPoses - pos
	if extraBelow > 0 {
		tview.Print(screen, fmt.Sprintf("[black:green]%d More", extraBelow), x, y+maxNodes+1, width, tview.AlignLeft, tcell.ColorWhite)
	}

	if extraAbove > 0 {
		tview.Print(screen, fmt.Sprintf("[black:green]%d More", extraAbove), x, y, width, tview.AlignLeft, tcell.ColorWhite)
	}
}
