package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Start() {
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(true)
	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 10, 40
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}

			cell := tview.NewTableCell(lorem[word]).
				SetTextColor(color).
				SetAlign(tview.AlignCenter)

			if r == 0 {
				cell.SetSelectable(false)
			}
			table.SetCell(r, c, cell)
			word = (word + 1) % len(lorem)
		}
	}
	table.SetSelectable(true, false)
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		//if key == tcell.KeyEnter {
		//table.SetSelectable(true, true)
		//}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	})

	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
