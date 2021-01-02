package tree

import "strconv"

type Id string

type Node struct {
	id         Id
	text       string
	subNodes   []*Node
	subNodeMap map[Id]*Node
	onSelect   func(Id, *Node)
	open       bool
	hasNodes   bool
}

func (node *Node) Select() {
	if node.onSelect != nil {
		node.onSelect(node.id, node)
	}
}

func (node *Node) Open() {
	node.open = true
	if node.onSelect != nil {
		node.onSelect(node.id, node)
	}
}

func (node *Node) Close() {
	node.open = false
}

func (node *Node) Toggle() {
	node.open = !node.open
	if node.open && node.onSelect != nil {
		node.onSelect(node.id, node)
	}
}

func (node *Node) ClearNodes() {
	node.subNodes = nil
	node.subNodeMap = nil
}

func (node *Node) AddNode(text string, lazyLoadSubNodes bool, onSelect func(nodeId Id, node *Node)) *Node {
	id := Id(strconv.Itoa(nodeId))
	nodeId++
	var newNode *Node

	if node.subNodes != nil {
		newNode = &Node{
			id:       id,
			text:     text,
			subNodes: nil,
			onSelect: onSelect,
			open:     false,
			hasNodes: lazyLoadSubNodes,
		}
		node.subNodes = append(node.subNodes, newNode)
		node.subNodeMap[id] = newNode
	} else {
		newNode = &Node{
			id:       id,
			text:     text,
			subNodes: nil,
			onSelect: onSelect,
			open:     false,
			hasNodes: lazyLoadSubNodes,
		}
		node.subNodes = append([]*Node{}, newNode)
		node.subNodeMap = map[Id]*Node{id: newNode}
		node.hasNodes = true
	}

	return newNode
}

func (node *Node) HasNodes() bool {
	return node.hasNodes
}

func (node *Node) GetNode(id Id) *Node {
	if node.subNodes != nil {
		return node.subNodeMap[id]
	}
	return nil
}

func (node *Node) GetNodes() []*Node {
	if node.subNodes != nil {
		return node.subNodes
	}
	return []*Node{}
}
