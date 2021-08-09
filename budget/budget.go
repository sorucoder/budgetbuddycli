package budget

import (
	"encoding/json"
	"fmt"
	"os"
)

// Budget describes a named budget comprised of income and expenses
type Budget struct {
	name     string
	Income   IncomeList  `json:"income"`
	Expenses ExpenseList `json:"expenses"`
}

// Make makes a named budget
func Make(name string) *Budget {
	return &Budget{
		name:     name,
		Income:   make(IncomeList),
		Expenses: make(ExpenseList),
	}
}

// Load loads a budget from disk
func Load(name string) (*Budget, error) {
	budget := Make(name)

	fileReader, err := os.Open(fmt.Sprintf("%s.budget", name))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(fileReader)
	if err := decoder.Decode(&budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// Save saves a budget to disk
func (budget *Budget) Save() error {
	fileWriter, err := os.Create(fmt.Sprintf("%s.budget", budget.name))
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	encoder := json.NewEncoder(fileWriter)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(budget); err != nil {
		return err
	}

	return nil
}
