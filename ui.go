package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	app    *tview.Application
	pages  *tview.Pages
	table  *tview.Table
	update *tview.Modal
	quit   *tview.Modal
}

func Start() {
	v := &view{
		app:    tview.NewApplication(),
		pages:  tview.NewPages(),
		table:  tview.NewTable(),
		update: tview.NewModal(),
		quit:   tview.NewModal(),
	}

	v.setUpdateModal()
	v.setQuitModal()
	v.setTable()

	if err := v.app.SetRoot(v.pages, true).SetFocus(v.table).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (v *view) setUpdateModal() {
	v.update.SetText("Do you want to update?").
		AddButtons([]string{"Yes", "No"})

	v.pages.AddPage("update", v.update, true, true)
	v.pages.SendToBack("update")
}

func (v *view) setQuitModal() {
	v.quit.SetText("Do you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				v.app.Stop()
			} else {
				v.pages.SendToBack("quit")
			}
		})

	v.pages.AddPage("quit", v.quit, true, true)
	v.pages.SendToBack("quit")
}

func (v *view) setTable() {

	colsHeader := []string{"PATH", "NAME", "CURRENT VERSION", "UPDATE VERSION"}
	cols, rows := 4, len(dependencies)+1
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c == 0 && r != 0 {
				color = tcell.ColorDarkCyan
			} else if c != 0 && r != 0 {
				color = tcell.ColorRed
			}

			align := tview.AlignLeft
			if r == 0 {
				align = tview.AlignCenter
			} else if c == 0 {
				align = tview.AlignRight
			}

			var tableCell *tview.TableCell
			// Set Headers
			if r == 0 {
				tableCell = tview.NewTableCell(colsHeader[c]).
					SetTextColor(color).
					SetAlign(align).
					SetSelectable(false)
			} else {
				tableCell = tview.NewTableCell(getString(r-1, c)).
					SetTextColor(color).
					SetAlign(align).
					SetSelectable(c != 0)
			}

			v.table.SetCell(r, c, tableCell)

			if c > 0 && c < 4 {
				tableCell.SetExpansion(1)
			}
		}
	}

	v.table.SetBorders(true)
	v.table.SetSelectable(true, false)
	v.table.Select(1, 0).SetFixed(1, 0).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				v.pages.SendToFront("quit")
			}
		}).
		SetSelectedFunc(func(row int, column int) {
			v.pages.SendToFront("update")
			v.update.SetText(fmt.Sprintf("Do you want to upgrade %s for %s?", dependencies[row-1].name, dependencies[row-1].path)).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "Yes" {
						// TODO:  update module

						for c := 1; c < cols; c++ {
							v.table.GetCell(row, c).SetTextColor(tcell.ColorWhite).SetSelectable(false)
						}
					}
					v.pages.SendToBack("update")

				})
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
