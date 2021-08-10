package reports

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sorucoder/budgetbuddy/budget"
)

func reportSummary(budget *budget.Budget) {
	tableWriter := table.NewWriter()

	tableWriter.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:      1,
			Align:       text.AlignLeft,
			AlignHeader: text.AlignLeft,
		},
		{
			Number:      2,
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
		},
		{
			Number:      3,
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
		},
	})
	tableWriter.SetStyle(table.StyleColoredBright)

	tableWriter.SetTitle("Summary")
	tableWriter.AppendHeader(table.Row{"Income", "Expenses", "Remaining"})
	tableWriter.AppendRow(table.Row{budget.Income.Sum(), budget.Expenses.Sum(), budget.Sum()})

	fmt.Println(tableWriter.Render())
}
