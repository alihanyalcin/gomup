package ui

import "github.com/rivo/tview"

func Start() {
	box := tview.NewBox().SetBorder(true).SetTitle("gomup")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
