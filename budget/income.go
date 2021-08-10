package budget

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

// Income describes a source of monthly income
type Income interface {
	MonthlyIncome() quantity.Money
}

// IncomeList is a list of named monthly income sources
type IncomeList map[string]Income

// SortedNames sorts the names of income sources lexographically
func (list IncomeList) SortedNames() []string {
	names := make([]string, 0, len(list))
	for name := range list {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Sum adds all income sources together
func (list IncomeList) Sum() quantity.Money {
	var total quantity.Money
	for _, income := range list {
		total += income.MonthlyIncome()
	}
	return total
}

// unmarshalIncomeJSON unmarshals JSON into an Income
func unmarshalIncomeJSON(incomeJSON json.RawMessage) (Income, error) {
	// Try Wages
	var wagesJSON struct {
		Rate  *float64 `json:"rate"`
		Hours *float64 `json:"hours"`
	}
	json.Unmarshal(incomeJSON, &wagesJSON)
	if wagesJSON.Rate != nil && wagesJSON.Hours != nil {
		return &Wages{Rate: quantity.Money(*wagesJSON.Rate), Hours: quantity.Number(*wagesJSON.Hours)}, nil
	}

	// Try Salary
	var salaryJSON struct {
		Salary *float64 `json:"salary"`
	}
	json.Unmarshal(incomeJSON, &salaryJSON)
	if salaryJSON.Salary != nil {
		return &Salary{Salary: quantity.Money(*salaryJSON.Salary)}, nil
	}

	// Try Sales
	var salesJSON struct {
		Rate  *float64 `json:"rate"`
		Items *float64 `json:"items"`
	}
	json.Unmarshal(incomeJSON, &salesJSON)
	if salesJSON.Rate != nil && salesJSON.Items != nil {
		return &Sales{Rate: quantity.Money(*salesJSON.Rate), Items: quantity.Integer(*salesJSON.Items)}, nil
	}

	// Try Commissions
	var commissionsJSON struct {
		Rate   *float64  `json:"rate"`
		Volume []float64 `json:"volume"`
	}
	json.Unmarshal(incomeJSON, &commissionsJSON)
	if commissionsJSON.Rate != nil && len(commissionsJSON.Volume) > 0 {
		volumeMoney := make([]quantity.Money, 0, len(commissionsJSON.Volume))
		for _, volume := range commissionsJSON.Volume {
			volumeMoney = append(volumeMoney, quantity.Money(volume))
		}
		return &Commissions{Rate: quantity.Percentage(*commissionsJSON.Rate), Volume: volumeMoney}, nil
	}

	// Try supplemental
	var supplementalJSON struct {
		Money *float64 `json:"money"`
	}
	json.Unmarshal(incomeJSON, &supplementalJSON)
	if supplementalJSON.Money != nil {
		return &Supplemental{Money: quantity.Money(*supplementalJSON.Money)}, nil
	}

	return nil, errors.New("unknown income format")
}

// UnmarshalJSON implements json.Unmarshaler for IncomeList
func (list *IncomeList) UnmarshalJSON(data []byte) error {
	var incomeListJSON map[string]json.RawMessage

	if err := json.Unmarshal(data, &incomeListJSON); err != nil {
		return err
	}

	for name, incomeJSON := range incomeListJSON {
		if income, err := unmarshalIncomeJSON(incomeJSON); err == nil {
			(*list)[name] = income
		} else {
			return err
		}
	}
	return nil
}

// Wages describes an income source paid a fixed rate every hour, including overtime pay.
// Example: You are paid $9 per hour and work about 50 hours a week. Assumming the legal
// amount of time to exceed normal pay is 40 hours, you would receive $360 with an
// additional overtime amount of $135, netting a total of $25,740 per year, or $2,145
// per month.
type Wages struct {
	Rate  quantity.Money  `survey:"rate" json:"rate"`   // Rate paid per hour
	Hours quantity.Number `survey:"hours" json:"hours"` // Hours worked in one week
}

// MonthlyIncome implements Income for Wages
func (income *Wages) MonthlyIncome() quantity.Money {
	var normalHours, overtimeHours float64
	if income.Hours.ValueOf() > MinimumOvertimeHours {
		normalHours = MinimumOvertimeHours
		overtimeHours = income.Hours.ValueOf() - MinimumOvertimeHours
	} else {
		normalHours = income.Hours.ValueOf()
		overtimeHours = 0
	}
	return quantity.Money(NetPayPercentage * 52 * (income.Rate.ValueOf()*normalHours + 1.5*income.Rate.ValueOf()*overtimeHours) / 12)
}

// Salary describes an income source that is paid as a fixed amount per year over regular intervals.
// Example: You earn $50,000 a year as a Mathematics Professor, and earn $4,166.67 per month.
type Salary struct {
	Salary quantity.Money `survey:"salary" json:"salary"`
}

// MonthyIncome implements Income for Salary
func (income Salary) MonthlyIncome() quantity.Money {
	return quantity.Money(NetPayPercentage * income.Salary.ValueOf() / 12)
}

// Sales describes an income source that is paid a fixed amount per item sold or task completed.
// For simplicity, the user is asked the estimated average of items sold or completed.
// Example: You sell 50 cups of lemonade on average each month at a lemonade stand for $1 per cup, so your monthly
// income would roughly be $50 per month, or $600 per year.
type Sales struct {
	Rate  quantity.Money   `survey:"rate" json:"rate"`   // Amount paid per item
	Items quantity.Integer `survey:"items" json:"items"` // Average count of items sold/tasks completed per month
}

// MonthlyIncome implements Income for Sales
func (income Sales) MonthlyIncome() quantity.Money {
	return quantity.Money(income.Rate.ValueOf() * income.Items.ValueOf())
}

// Commissions describes an income source that earns a portion of the value of each item sold or task completed.
// Example: You are a realtor and you make 6% on each home you sell. You sold 2 homes - one for $25,000 and one for $75,000.
// Your monthly income for this month would be $6,000
type Commissions struct {
	Rate   quantity.Percentage `survey:"rate" json:"rate"`     // Percentage for each item sold/task completed
	Volume []quantity.Money    `survey:"volume" json:"volume"` // Value of each item sold/task completed
}

// MonthlyIncome implements Income for Commissions
func (income Commissions) MonthlyIncome() quantity.Money {
	var totalValue float64
	for _, volume := range income.Volume {
		totalValue += income.Rate.ValueOf() * volume.ValueOf()
	}
	return quantity.Money(totalValue)
}

// Supplemental describes a generic monthly income source.
// Example: You receive $100 per month in allowance.
type Supplemental struct {
	Money quantity.Money `survey:"money" json:"money"`
}

// MonthlyIncome implements Income for Supplemental
func (income Supplemental) MonthlyIncome() quantity.Money {
	return income.Money
}
