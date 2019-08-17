package contract

import (
	"github.com/spf13/cobra"
)

var ContractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Deploy or call zilliqa smart contract",
	Long:  "Use deploy sub command to deploy fresh smart contract, use call sub command to invoke exist smart contract",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
