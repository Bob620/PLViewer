package viewer

import (
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"PLViewer/ui/tree"
)

type Viewer struct {
	*page.Page
}

func MakeViewerPage() *Viewer {
	viewer := Viewer{
		Page: page.MakePage("Viewer", page.MakeLayout([][]string{
			{"input", "metadata"},
			{"tree", "metadata"},
			{"tree", "metadata"},
			{"tree", "graph"},
		})),
	}
	viewer.SetOnSelect(func() {
		viewer.Page.SelectElement("input")
	})

	inputElement := element.MakeElement().SetBorders(true)
	metadataElement := element.MakeElement().SetBorders(true)
	treeElement := tree.MakeTree()
	graphElement := element.MakeElement().SetBorders(true).SetSelectable(false)

	treeElement.Element.SetBorders(true)

	treeElement.
		AddNode("0", true, nil).
		AddNode("1", false, nil).
		AddNode("2", true, nil).
		AddNode("3", true, nil).
		AddNode("4", true, nil).
		AddNode("5", true, nil).
		AddNode("6", true, nil).
		AddNode("7", true, nil).
		AddNode("8", true, nil).
		AddNode("9", true, nil).
		AddNode("10", true, nil).
		AddNode("11", true, nil).
		AddNode("12", true, nil).
		AddNode("13", true, nil).
		AddNode("14", true, nil).
		AddNode("15", true, nil).
		AddNode("16", true, nil).
		AddNode("17", true, nil).
		AddNode("18", true, nil).
		AddNode("19", true, nil).
		AddNode("20", true, nil).
		AddNode("21", true, nil).
		AddNode("22", true, nil).
		AddNode("23", true, nil).
		AddNode("24", true, nil).
		AddNode("25", true, nil).
		AddNode("26", true, nil).
		AddNode("27", true, nil).
		AddNode("28", true, nil).
		AddNode("29", true, nil).
		AddNode("30", true, nil).
		AddNode("31", true, nil).
		AddNode("32", true, nil)
	test := treeElement.
		AddNode("33", true, nil).
		AddNode("36", true, nil)
	test.
		AddNode("34", true, nil)
	test.
		AddNode("35", true, nil)
	treeElement.AddNode("37", true, nil)

	//	treeElement.GetNode("3").AddNode("test").AddNode("test1")
	//	treeElement.GetNode("3").GetNode().AddNode("ahhh")

	//	rightElement.Flex().
	//		AddItem(tview.NewTextView().SetText("Test"), 0, 1, false)
	//	bottomElement.Flex().
	//		AddItem(tview.NewTextView().SetText("Ahhh"), 0, 1, false)

	viewer.
		AddElement(inputElement, "input").
		AddElement(metadataElement, "metadata").
		AddElement(treeElement, "tree").
		AddElement(graphElement, "graph")

	return &viewer
}
