package reports

import (
	"fmt"

	"github.com/sorucoder/budgetbuddy/budget"
)

func ReportBudget(budget *budget.Budget) {
	reportIncomeList(budget.Income)
	fmt.Println()
	reportExpenseList(budget.Expenses)
	fmt.Println()
	reportSummary(budget)
}
