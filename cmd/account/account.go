package account

import "github.com/spf13/cobra"

var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Generate a large number of accounts",
	Long:  "Generate a large number of accounts for stress test and .etc",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
