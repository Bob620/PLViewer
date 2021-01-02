package main

import (
	"PLViewer/backend"
	"PLViewer/ui"
)

func main() {
	bg := backend.StartBackend()
	ui.Initialize(bg)
}
