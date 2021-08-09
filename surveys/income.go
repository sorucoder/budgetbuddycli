package surveys

import (
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/muesli/termenv"
	"github.com/sorucoder/budgetbuddy/budget"
	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

type incomeSurvey func() (budget.Income, error)

var incomeSurveys = map[string]incomeSurvey{
	"Wages":       askWagesSurvey,
	"Salary":      askSalarySurvey,
	"Sales":       askSalesSurvey,
	"Commissions": askCommissionsSurvey,
}

func askIncomeListSurvey(list budget.IncomeList) error {
	fmt.Println(termenv.String("Income").Underline())
	var done bool
	for !done {
		incomeTitle := fmt.Sprintf("%s Source Of Income", quantity.MakeInteger(len(list)+1).Ordinal())
		fmt.Println(termenv.String(incomeTitle).Italic())

		if name, income, err := askIncomeSurvey(); err == nil {
			list[name] = income
		} else {
			return err
		}

		if err := survey.AskOne(
			&survey.Confirm{
				Message: "Are you finished entering all of your sources of income?",
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

func askIncomeSurvey() (string, budget.Income, error) {
	incomeTypes := make([]string, 0, len(incomeSurveys))
	for incomeType := range incomeSurveys {
		incomeTypes = append(incomeTypes, incomeType)
	}
	sort.Strings(incomeTypes)

	var incomeNameAnswer string
	if err := survey.AskOne(
		&survey.Input{
			Message: "Name of Income:",
		},
		&incomeNameAnswer,
		survey.WithValidator(survey.Required),
	); err != nil {
		return "", nil, err
	}

	var incomeTypeAnswer string
	if err := survey.AskOne(
		&survey.Select{
			Message: "Type Of Income:",
			Options: incomeTypes,
		},
		&incomeTypeAnswer,
	); err != nil {
		return "", nil, err
	}

	var incomeAnswer budget.Income
	if income, err := incomeSurveys[incomeTypeAnswer](); err == nil {
		incomeAnswer = income
	} else {
		return "", nil, err
	}

	return incomeNameAnswer, incomeAnswer, nil
}

func askWagesSurvey() (budget.Income, error) {
	var wages budget.Wages
	if err := survey.Ask(
		[]*survey.Question{
			{
				Name: "rate",
				Prompt: &survey.Input{
					Message: fmt.Sprintf("Hourly Rate %s:", termenv.String("($)").Faint()),
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					moneyValidator,
					boundedMoneyValidator(budget.MinimumWage, nil),
				),
			},
			{
				Name: "hours",
				Prompt: &survey.Input{
					Message: fmt.Sprintf("Average Hours Per Week %s:", termenv.String("(#)").Faint()),
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					numberValidator,
					boundedNumberValidator(1, nil),
				),
			},
		},
		&wages,
	); err != nil {
		return nil, err
	}
	return &wages, nil
}

func askSalarySurvey() (budget.Income, error) {
	var salary budget.Salary
	if err := survey.AskOne(
		&survey.Input{
			Message: fmt.Sprintf("Salary %s:", termenv.String("($)").Faint()),
		},
		&salary,
		survey.WithValidator(
			survey.ComposeValidators(
				survey.Required,
				moneyValidator,
				boundedMoneyValidator(0.01, nil),
			),
		),
	); err != nil {
		return nil, err
	}
	return &salary, nil
}

func askSalesSurvey() (budget.Income, error) {
	var sales budget.Sales
	if err := survey.Ask(
		[]*survey.Question{
			{
				Name: "rate",
				Prompt: &survey.Input{
					Message: fmt.Sprintf("Selling Price %s:", termenv.String("($)").Faint()),
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					moneyValidator,
					boundedMoneyValidator(0.01, nil),
				),
			},
			{
				Name: "items",
				Prompt: &survey.Input{
					Message: fmt.Sprintf("Average Number of Items Sold %s:", termenv.String("(@)").Faint()),
				},
				Validate: survey.ComposeValidators(
					survey.Required,
					integerValidator,
					boundedIntegerValidator(1, nil),
				),
			},
		},
		&sales,
	); err != nil {
		return nil, err
	}
	return &sales, nil
}

func askCommissionsSurvey() (budget.Income, error) {
	var commissions budget.Commissions

	if err := survey.AskOne(
		&survey.Input{
			Message: fmt.Sprintf("Percentage %s:", termenv.String("(%)").Faint()),
		},
		&commissions.Rate,
		survey.WithValidator(survey.ComposeValidators(survey.Required, percentageValidator, boundedPercentageValidator(0, nil))),
	); err != nil {
		return nil, err
	}

	fmt.Printf("%s Please enter each item that made commissions:\n", termenv.String("?").Foreground(termenv.ANSIGreen))
	var done bool
	for !done {
		var commissionVolume quantity.Money
		if err := survey.AskOne(
			&survey.Input{
				Message: fmt.Sprintf("    Item #%d %s:", len(commissions.Volume)+1, termenv.String("($)").Faint()),
			},
			&commissionVolume,
			survey.WithValidator(survey.ComposeValidators(survey.Required, moneyValidator, boundedMoneyValidator(0.01, nil))),
		); err != nil {
			return nil, err
		}
		commissions.Volume = append(commissions.Volume, commissionVolume)

		if err := survey.AskOne(
			&survey.Confirm{
				Message: fmt.Sprintf("    Are you finished entering all items?"),
				Default: false,
			},
			&done,
		); err != nil {
			return nil, err
		}
	}

	return &commissions, nil
}
