package budget

import (
	"sort"

	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

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

// SortedNames sorts the names of expenses lexographically
func (list ExpenseList) SortedNames() []string {
	names := make([]string, 0, len(list))
	for name := range list {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
