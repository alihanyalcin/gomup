package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start() {

	app := tview.NewApplication()
	pages := tview.NewPages()

	// Quit Modal
	quitModal := tview.NewModal().
		SetText("Do you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				app.Stop()
			} else {
				pages.SendToBack("quit")
			}
		})

	pages.AddPage("quit", quitModal, true, true)
	pages.SendToBack("quit")

	table := tview.NewTable().SetBorders(true)
	pages.AddPage("table", table, true, true)

	colsHeader := []string{"Path", "Name", "Current Version", "Update Version"}
	cols, rows := 4, len(dependencies)+1
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}

			// Set Headers
			if r == 0 {
				table.SetCell(r, c,
					tview.NewTableCell(colsHeader[c]).
						SetTextColor(color).
						SetAlign(tview.AlignCenter).
						SetSelectable(false))

				continue
			}

			table.SetCell(r, c,
				tview.NewTableCell(getString(r-1, c)).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
		}
	}

	table.SetSelectable(true, false)
	table.Select(0, 0).SetFixed(1, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				pages.SendToFront("quit")
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			// TODO: add update modal
			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func getString(row, col int) string {
	value := dependencies[row]
	switch col {
	case 0:
		return value.path
	case 1:
		return value.name
	case 2:
		return value.version
	case 3:
		return value.updateVersion
	default:
		return ""
	}
}
