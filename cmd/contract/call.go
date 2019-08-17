package contract

import (
	"github.com/spf13/cobra"
)

func init() {
	ContractCmd.AddCommand(callCmd)
}

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Call a exist contract",
	Long:  "Call a exist contract",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
