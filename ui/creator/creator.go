package creator

import "PLViewer/ui/page"

type CreatorPage struct {
	*page.Page
}

func MakeCreatorPage() *CreatorPage {
	creator := CreatorPage{
		Page: page.MakePage("Creator"),
	}

	return &creator
}
