package surveys

import (
	"fmt"
	"math"

	"github.com/AlecAivazis/survey/v2"
	"github.com/muesli/termenv"
	"github.com/sorucoder/budgetbuddy/budget"
	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

func askExpenseListSurvey(list budget.ExpenseList) error {
	fmt.Println(termenv.String("Expenses").Underline())
	var done bool
	for !done {
		expenseTitle := fmt.Sprintf("%s Expense", quantity.MakeInteger(len(list)+1).Ordinal())
		fmt.Println(termenv.String(expenseTitle).Italic())

		if name, expense, err := askExpenseSurvey(); err == nil {
			list[name] = expense
		} else {
			return err
		}

		if err := survey.AskOne(
			&survey.Confirm{
				Message: "Are you finished entering all of your expenses?",
				Default: false,
			},
			&done,
		); err != nil {
			return err
		}

		fmt.Println()
	}

	return nil
}

func askExpenseSurvey() (string, quantity.Money, error) {
	var expense struct {
		Name   string         `survey:"name"`
		Amount quantity.Money `survey:"amount"`
	}
	if err := survey.Ask(
		[]*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Message: "Name of Expense:",
				},
				Validate: survey.Required,
			},
			{
				Name: "amount",
				Prompt: &survey.Input{
					Message: fmt.Sprintf("Cost of Expense %s:", termenv.String("($)").Faint()),
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					moneyValidator,
					boundedMoneyValidator(0.01, nil),
				),
			},
		},
		&expense,
	); err != nil {
		return "", quantity.Money(math.NaN()), err
	}

	return expense.Name, expense.Amount, nil
}
