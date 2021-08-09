package budget

import "github.com/sorucoder/budgetbuddy/budget/quantity"

// ExpenseList is a named list of monthly expenses
type ExpenseList map[string]quantity.Money

// Sum adds all expenses together
func (list ExpenseList) Sum() quantity.Money {
	var total quantity.Money
	for _, expense := range list {
		total += expense
	}
	return total
}
