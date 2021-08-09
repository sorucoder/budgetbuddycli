package budget

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sorucoder/budgetbuddy/budget/quantity"
)

// Income describes a source of monthly income
type Income interface {
	MonthlyIncome() quantity.Money
}

// IncomeList is a list of named monthly income sources
type IncomeList map[string]Income

// Sum adds all income sources together
func (list IncomeList) Sum() quantity.Money {
	var total quantity.Money
	for _, income := range list {
		total += income.MonthlyIncome()
	}
	return total
}

// marshalIncomeJSON marshals an Income into JSON
func marshalIncomeJSON(income Income) ([]byte, error) {
	switch incomeValue := income.(type) {
	case *Wages:
		if incomeJSON, err := json.Marshal(incomeValue); err == nil {
			return incomeJSON, nil
		} else {
			return nil, err
		}
	case *Salary:
		if incomeJSON, err := json.Marshal(incomeValue); err == nil {
			return incomeJSON, nil
		} else {
			return nil, err
		}
	case *Sales:
		if incomeJSON, err := json.Marshal(incomeValue); err == nil {
			return incomeJSON, nil
		} else {
			return nil, err
		}
	case *Commissions:
		if incomeJSON, err := json.Marshal(incomeValue); err == nil {
			return incomeJSON, nil
		} else {
			return nil, err
		}
	default:
		return nil, errors.New("unknown income type")
	}
}

// MarshalJSON implements json.Marshaler for IncomeList
func (list IncomeList) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("{")

	index := 0
	for name, income := range list {
		buffer.WriteString(fmt.Sprintf("\"%s\":", name))

		if incomeJSON, err := marshalIncomeJSON(income); err == nil {
			buffer.Write(incomeJSON)
		} else {
			return nil, err
		}

		index++
		if index < len(list)-1 {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

// unmarshalIncomeJSON unmarshals JSON into an Income
func unmarshalIncomeJSON(incomeJSON json.RawMessage) (Income, error) {
	var wages Wages
	if err := json.Unmarshal(incomeJSON, &wages); err == nil {
		return &wages, nil
	}

	var salary Salary
	if err := json.Unmarshal(incomeJSON, &salary); err == nil {
		return &salary, nil
	}

	var sales Sales
	if err := json.Unmarshal(incomeJSON, &sales); err == nil {
		return &sales, nil
	}

	var commissions Commissions
	if err := json.Unmarshal(incomeJSON, &commissions); err == nil {
		return &commissions, nil
	}

	return nil, errors.New("unknown income type")
}

// UnmarshalJSON implements json.Unmarshaler for IncomeList
func (list IncomeList) UnmarshalJSON(data []byte) error {
	var incomeMapJSON map[string]json.RawMessage
	if err := json.Unmarshal(data, &incomeMapJSON); err != nil {
		return err
	}
	for name, incomeJSON := range incomeMapJSON {
		if income, err := unmarshalIncomeJSON(incomeJSON); err == nil {
			list[name] = income
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
	return quantity.Money(52 * (income.Rate.ValueOf()*normalHours + 1.5*income.Rate.ValueOf()*overtimeHours) / 12)
}

// Salary describes an income source that is paid as a fixed amount per year over regular intervals.
// Example: You earn $50,000 a year as a Mathematics Professor, and earn $4,166.67 per month.
type Salary struct {
	quantity.Money `survey:"salary"`
}

// MonthyIncome implements Income for Salary
func (income Salary) MonthlyIncome() quantity.Money {
	return income.Money / 12
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
	Rate   quantity.Percentage // Percentage for each item sold/task completed
	Volume []quantity.Money    // Value of each item sold/task completed
}

// MonthlyIncome implements Income for Commissions
func (income Commissions) MonthlyIncome() quantity.Money {
	var totalValue float64
	for _, volume := range income.Volume {
		totalValue += income.Rate.ValueOf() * volume.ValueOf()
	}
	return quantity.Money(totalValue)
}
