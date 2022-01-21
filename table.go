package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func drawTable() {
	var data [][]string
	for _, v := range dependencies {
		data = append(data, []string{
			v.path,
			v.name,
			v.version,
			v.updateVersion,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Path", "Name", "Version", "Update Version"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	var color = []tablewriter.Colors{
		{tablewriter.Bold, tablewriter.FgMagentaColor},
		{tablewriter.FgCyanColor},
		{tablewriter.Bold},
		{tablewriter.Bold},
	}

	for _, v := range data {
		table.Rich(v, color)
	}
	table.Render()
}
