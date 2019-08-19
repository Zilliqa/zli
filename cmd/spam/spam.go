package spam

import "github.com/spf13/cobra"

var SpamCmd = &cobra.Command{
	Use:   "spam",
	Short: "Send a large number of transactions",
	Long:  "Send a large number of transactions to a specific account or contract",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
