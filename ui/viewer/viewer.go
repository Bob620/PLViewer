package viewer

import (
	"PLViewer/backend"
	"PLViewer/sxes"
	"PLViewer/ui/element"
	"PLViewer/ui/input"
	"PLViewer/ui/interop"
	"PLViewer/ui/linegraph"
	"PLViewer/ui/page"
	"PLViewer/ui/tree"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Viewer struct {
	*page.Page
	bg *backend.Backend
}

func getDir(uri string) []os.FileInfo {
	files, err := ioutil.ReadDir(uri)
	if err != nil {
		return []os.FileInfo{}
	}

	return files
}

func MakeViewerPage(application *tview.Application, bg *backend.Backend, interopData *interop.InteropData) *Viewer {
	viewer := Viewer{
		Page: page.MakePage("Viewer", page.MakeLayout([][]string{
			{"input", "metadata"},
			{"tree", "metadata"},
			{"tree", "metadata"},
			{"tree", "graph"},
		})),
		bg: bg,
	}
	viewer.SetOnSelect(func() {
		viewer.Page.SelectElement("input")
	})

	inputElement := element.MakeElement().SetBorders(true)
	metadataPages := createMetadata()
	metadataPages.SetBorders(true).SetSelectable(true)

	treeElement := tree.MakeTree()
	treeElement.SetBorders(true)
	treeElement.SetTitle(" PLZip Data Tree ")

	graphElement := linegraph.MakeLineGraph()
	graphElement.SetBorders(true).SetSelectable(false).SetHoverable(false)
	graphElement.SetTitle(" QLW Graph ")

	inputBox := input.MakeInput(application)
	inputElement.GetFlex().AddItem(inputBox, 0, 1, false)
	inputElement.SetOnSelect(func() {
		inputBox.Activate(true)
	})
	inputElement.SetOnKeyEvent(inputBox.HandleEvents)
	inputElement.SetOnDeselect(func() {
		inputBox.Activate(false)
	})

	inputBox.SetOnSubmit(func(uri string) {
		treeElement.ClearNodes()

		if !strings.HasSuffix(uri, ".plzip") {
			return
		}

		interopData.SetString("viewerUri", uri)
		bg.Load(uri)
		projects, _ := bg.GetProjects()
		for _, project := range projects {
			project := project
			treeElement.AddNode(project.Name, true, func(nodeId tree.Id, node *tree.Node) {
				node.ClearNodes()
				graphElement.SetLine(nil)

				// Initialize Metadata view
				metadataPages.SetAsProject(project)

				// Set up downstream nodes
				for _, analysisUuid := range project.Analyses {
					analysis, _ := bg.GetAnalysis(analysisUuid.Uuid)
					text := analysis.Name
					if text == "" {
						text = analysis.Comment
					}
					if text == "" {
						text = analysis.AcquisitionDate
					}

					node.AddNode(text, true, func(nodeId tree.Id, node *tree.Node) {
						node.ClearNodes()
						graphElement.SetLine(nil)

						// Initialize Metadata view
						metadataPages.SetAsAnalysis(project, analysis)

						// Set up downstream nodes asynchronously
						var wg sync.WaitGroup
						positions := make([]*sxes.Position, len(analysis.Positions))

						for i, positionUuid := range analysis.Positions {
							wg.Add(1)
							go func(i int, uuid string) {
								position, _ := bg.GetPosition(uuid)
								positions[i] = position
								wg.Done()
							}(i, positionUuid.Uuid)
						}

						wg.Wait()
						for _, position := range positions {
							func(position *sxes.Position) {
								text := position.Comment
								if text == "" {
									text = position.Uuid
								}
								node.AddNode(text, false, func(nodeId tree.Id, node *tree.Node) {
									line, _ := bg.GetLine(position.Uuid, position.Types[0])
									graphElement.SetLine(line)

									// Initialize Metadata view
									metadataPages.SetAsPosition(project, analysis, position)
								})
							}(position)
						}
					})
				}
			})
		}
	})

	inputBox.SetAutocompleteSelectionFunc(func(entry string) string {
		dir, _ := path.Split(inputBox.GetText())
		return path.Join(dir, entry)
	}).SetAutocompleteFunc(func(currentText string) (entries []string) {
		prefix := strings.TrimSpace(strings.ToLower(currentText))
		dir, wanted := path.Split(prefix)
		files := getDir(dir)

		for _, file := range files {
			name := strings.ToLower(file.Name())
			if file.IsDir() || strings.HasSuffix(name, ".plzip") || strings.HasSuffix(name, ".pl7z") {
				if wanted == "" || strings.Contains(strings.ToLower(name), wanted) {
					if strings.ToLower(name) != wanted {
						entries = append(entries, file.Name())
					}
				}
			}
		}
		return entries
	})

	viewer.
		AddElement(inputElement, "input").
		AddElement(metadataPages, "metadata").
		AddElement(treeElement, "tree").
		AddElement(graphElement, "graph")

	return &viewer
}
