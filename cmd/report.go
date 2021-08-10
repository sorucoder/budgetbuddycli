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

	"github.com/muesli/termenv"
	"github.com/sorucoder/budgetbuddy/budget"
	"github.com/sorucoder/budgetbuddy/reports"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "generates reports on created budgets",
	Long:  `Generates reports on budgets`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		budget.NetPayPercentage = viper.GetFloat64("net_pay_percentage")

		reportBudget, err := budget.Load(args[0])
		if err != nil {
			fmt.Println(termenv.String(fmt.Sprintf(`Could not load budget "%s.budget"`, args[0])).Foreground(termenv.ANSIRed))
			os.Exit(1)
		}
		reports.ReportBudget(reportBudget)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.Flags().Float64("net-pay-percentage", 0.75, "The estimated percentage used to calculate net pay from gross pay")
	viper.BindPFlag("net_pay_percentage", createCmd.Flags().Lookup("net-pay-percentage"))
}
