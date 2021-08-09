/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/muesli/termenv"
	"github.com/sorucoder/budgetbuddy/budget"
	"github.com/sorucoder/budgetbuddy/surveys"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates budgets",
	Long:  `Interactively prompts the user to create a budget from scratch.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize configuration for budget
		budget.MinimumWage = viper.GetFloat64("minimum_wage")
		budget.MinimumOvertimeHours = viper.GetFloat64("minimum_overtime_hours")

		newBudget := budget.Make(args[0])
		if err := surveys.AskBudgetSurvey(newBudget); err != nil {
			switch err {
			case terminal.InterruptErr:
				fmt.Println(termenv.String("Aborted budget creation").Foreground(termenv.ANSIRed))
				os.Exit(0)
			default:
				panic(err)
			}
		}

		if err := newBudget.Save(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().Float64("minimum-wage", 7.25, "The legal minimum rate of pay for wages")
	viper.BindPFlag("minimum_wage", createCmd.Flags().Lookup("minimum-wage"))

	createCmd.Flags().Float64("minimum-overtime-hours", 40, "The legal minimum number of hours required for overtime pay")
	viper.BindPFlag("minimum_overtime_hours", createCmd.Flags().Lookup("minimum-overtime-hours"))
}
