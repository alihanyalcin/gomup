package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	app       *tview.Application
	pages     *tview.Pages
	table     *tview.Table
	quitModal *tview.Modal
}

func Start() {
	v := &view{
		app:       tview.NewApplication(),
		pages:     tview.NewPages(),
		table:     tview.NewTable(),
		quitModal: tview.NewModal(),
	}

	v.setQuitModal()
	v.setTable()

	if err := v.app.SetRoot(v.pages, true).SetFocus(v.pages).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (v *view) setQuitModal() {
	v.quitModal.SetText("Do you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				v.app.Stop()
			} else {
				v.pages.SendToBack("quit")
			}
		})

	v.pages.AddPage("quit", v.quitModal, true, true)
	v.pages.SendToBack("quit")
}

func (v *view) setTable() {
	v.table.SetBorders(true)

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
				v.table.SetCell(r, c,
					tview.NewTableCell(colsHeader[c]).
						SetTextColor(color).
						SetAlign(tview.AlignCenter).
						SetSelectable(false))

				continue
			}

			v.table.SetCell(r, c,
				tview.NewTableCell(getString(r-1, c)).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
		}
	}

	v.table.SetSelectable(true, false)
	v.table.Select(1, 0).SetFixed(1, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				v.pages.SendToFront("quit")
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			// TODO: add update modal
			v.table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		})

	v.pages.AddPage("table", v.table, true, true)
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
