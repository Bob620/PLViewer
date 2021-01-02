package viewer

import (
	"PLViewer/sxes"
	"PLViewer/ui/element"
	"PLViewer/ui/page"
	"PLViewer/ui/pageviewer"
	"github.com/rivo/tview"
)

type Metadata struct {
	*pageviewer.PageViewer
	project  *page.Page
	analysis *page.Page
	position *page.Page
}

func createMetadata() *Metadata {
	metadata := Metadata{
		pageviewer.MakePageViewer(),
		page.MakePage("project", page.MakeLayout([][]string{
			{"uuid", "uuid", "uuid", "uuid"},
			{"project", "project", "", "date"},
			{"name", "name", "", "operator"},
			{"comment", "comment", "", "instrument"},
		})),
		page.MakePage("analysis", page.MakeLayout([][]string{
			{"uuid", "uuid", "uuid", "uuid"},
			{"project", "project", "", "date"},
			{"name", "name", "", "operator"},
			{"comment", "comment", "", "instrument"},
		})),
		page.MakePage("position", page.MakeLayout([][]string{
			{"uuid", "uuid", "uuid", "uuid", "uuid", "uuid"},
			{"project", "project", "project", "", "date", "date"},
			{"analysis", "analysis", "analysis", "", "operator", "operator"},
			{"comment", "comment", "comment", "", "instrument", "instrument"},
			{"types", "types", "types", "types", "types", "types"},
		})),
	}

	metadata.project.AddTo(metadata.Pages)
	metadata.analysis.AddTo(metadata.Pages)
	metadata.position.AddTo(metadata.Pages)

	return &metadata
}

func (metadata *Metadata) SetAsProject(project *sxes.Project) {
	metadata.project.Clear()
	metadata.PageViewer.SetTitle(project.Name)

	elem := element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(project.Uuid), 0, 1, false)
	metadata.project.AddElement(elem, "uuid")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(project.Name), 0, 1, false)
	metadata.project.AddElement(elem, "name")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(project.Comment), 0, 1, false)
	metadata.project.AddElement(elem, "comment")

	metadata.Pages.SwitchToPage("project")
}

func (metadata *Metadata) SetAsAnalysis(project *sxes.Project, analysis *sxes.Analysis) {
	text := analysis.Name
	if text == "" {
		text = analysis.Comment
	}
	if text == "" {
		text = analysis.AcquisitionDate
	}

	metadata.analysis.Clear()
	metadata.PageViewer.SetTitle(text)

	elem := element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(project.Name), 0, 1, false)
	metadata.analysis.AddElement(elem, "project")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.Uuid), 0, 1, false)
	metadata.analysis.AddElement(elem, "uuid")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.Name), 0, 1, false)
	metadata.analysis.AddElement(elem, "name")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.Comment), 0, 1, false)
	metadata.analysis.AddElement(elem, "comment")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.AcquisitionDate), 0, 1, false)
	metadata.analysis.AddElement(elem, "date")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.Instrument), 0, 1, false)
	metadata.analysis.AddElement(elem, "instrument")

	elem = element.MakeElement()
	elem.AddItem(tview.NewTextView().SetText(analysis.Operator), 0, 1, false)
	metadata.analysis.AddElement(elem, "operator")

	metadata.Pages.SwitchToPage("analysis")
}

func (metadata *Metadata) SetAsPosition(project *sxes.Project, analysis *sxes.Analysis, position *sxes.Position) {
	text := position.Comment
	if text == "" {
		text = position.Uuid
	}

	analysisText := analysis.Name
	if analysisText == "" {
		analysisText = analysis.Comment
	}
	if analysisText == "" {
		analysisText = analysis.AcquisitionDate
	}

	metadata.position.Clear()
	metadata.PageViewer.SetTitle(text)

	elem := element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Project:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(project.Name), 0, 1, false)
	metadata.position.AddElement(elem, "project")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Analysis:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(analysisText), 0, 1, false)
	metadata.position.AddElement(elem, "analysis")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("UUID:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(position.Uuid), 0, 1, false)
	metadata.position.AddElement(elem, "uuid")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Comment:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(position.Comment), 0, 1, false)
	metadata.position.AddElement(elem, "comment")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Date:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(analysis.AcquisitionDate), 0, 1, false)
	metadata.position.AddElement(elem, "date")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Instrument:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(analysis.Instrument), 0, 1, false)
	metadata.position.AddElement(elem, "instrument")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Operator:"), 1, 1, false)
	elem.AddItem(tview.NewTextView().SetText(position.Operator), 0, 1, false)
	metadata.position.AddElement(elem, "operator")

	elem = element.MakeElement()
	elem.SetDirection(tview.FlexRow)
	elem.AddItem(tview.NewTextView().SetText("Types:"), 1, 1, false)
	typeElem := element.MakeElement()
	for _, lineType := range position.Types {
		typeElem.AddItem(tview.NewTextView().SetText(lineType), 0, 1, false)
	}
	elem.AddItem(typeElem, 0, 1, false)
	metadata.position.AddElement(elem, "types")

	metadata.Pages.SwitchToPage("position")
}
