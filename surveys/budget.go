package surveys

import "github.com/sorucoder/budgetbuddy/budget"

func AskBudgetSurvey(budget *budget.Budget) error {
	// Ask for income
	if err := askIncomeListSurvey(budget.Income); err != nil {
		return err
	}

	// Ask for expenses
	if err := askExpenseListSurvey(budget.Expenses); err != nil {
		return err
	}

	return nil
}
