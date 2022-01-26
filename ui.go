package main

import (
	"fmt"
	"os/exec"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	app          *tview.Application
	pages        *tview.Pages
	table        *tview.Table
	update       *tview.Modal
	quit         *tview.Modal
	info         *tview.Modal
	dependencies []dependency
}

func startUI(d []dependency) {
	sort.Slice(d, func(i, j int) bool {
		return d[i].path < d[j].path
	})

	v := &view{
		app:          tview.NewApplication(),
		pages:        tview.NewPages(),
		table:        tview.NewTable(),
		update:       tview.NewModal(),
		quit:         tview.NewModal(),
		info:         tview.NewModal(),
		dependencies: d,
	}

	v.setInfoModal()
	v.setUpdateModal()
	v.setQuitModal()
	v.setTable()

	if err := v.app.SetRoot(v.pages, true).SetFocus(v.table).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func (v *view) setInfoModal() {
	v.info.AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.pages.SendToBack("info")
		})

	v.pages.AddPage("info", v.info, true, true)
	v.pages.SendToBack("info")
}

func (v *view) setUpdateModal() {
	v.update.AddButtons([]string{"Yes", "No"})

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
	v.table.SetBorder(true).SetTitle("gomUP - ESC: quit - ENTER: update").SetTitleAlign(tview.AlignCenter)

	colsHeader := []string{"PATH", "NAME", "CURRENT VERSION", "UPDATE VERSION"}
	cols, rows := len(colsHeader), len(v.dependencies)+1
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c == 0 && r != 0 {
				color = tcell.ColorAquaMarine
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
			if r == 0 {
				// Set Headers
				tableCell = tview.NewTableCell(colsHeader[c]).
					SetTextColor(color).
					SetAlign(align).
					SetSelectable(false).
					SetAttributes(tcell.AttrBold)
			} else {
				tableCell = tview.NewTableCell(v.getString(r-1, c)).
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
			dependency := v.dependencies[row-1]
			v.pages.SendToFront("update")
			v.update.SetText(fmt.Sprintf("Do you want to upgrade %s for %s?", dependency.name, dependency.path)).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					v.pages.SendToBack("update")
					if buttonLabel == "Yes" {
						v.upgrade(row, cols)
					}

				})
		})

	v.pages.AddPage("table", v.table, true, true)
}

func (v *view) upgrade(row, cols int) {
	v.pages.SendToFront("info")

	dependency := v.dependencies[row-1]

	cmd := exec.Command("go", "get", dependency.name)
	cmd.Dir = dependency.path
	_, err := cmd.Output()
	if err != nil {
		v.info.SetText("Error occured: " + err.Error())
	} else {
		v.info.SetText("Success!")

		for c := 1; c < cols; c++ {
			v.table.GetCell(row, c).SetTextColor(tcell.ColorAquaMarine).SetSelectable(false)
		}
	}
}

func (v *view) getString(row, col int) string {
	value := v.dependencies[row]
	return []string{value.path, value.name, value.version, value.updateVersion}[col]
}
