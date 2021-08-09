package reports

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sorucoder/budgetbuddy/budget"
)

func reportExpenseList(list budget.ExpenseList) {
	tableWriter := table.NewWriter()

	tableWriter.SetColumnConfigs([]table.ColumnConfig{
		{
			Number: 1,
			Hidden: true,
		},
		{
			Number:      2,
			Align:       text.AlignLeft,
			AlignHeader: text.AlignLeft,
			WidthMin:    75,
			WidthMax:    75,
		},
		{
			Number:      3,
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
			WidthMin:    25,
			WidthMax:    25,
		},
	})
	tableWriter.SetStyle(table.StyleColoredBright)

	tableWriter.SetTitle("Expenses")
	tableWriter.AppendHeader(table.Row{"Index", "Name", "Amount"})
	index := 1
	for _, name := range list.SortedNames() {
		expense := list[name]
		tableWriter.AppendRow(table.Row{index, name, expense})
		index++
	}
	tableWriter.AppendFooter(table.Row{"Index", "Total", list.Sum()})

	fmt.Println(tableWriter.Render())
}
