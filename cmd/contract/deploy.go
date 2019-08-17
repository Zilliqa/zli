package contract

import (
	"github.com/spf13/cobra"
)

func init() {
	ContractCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy new contract",
	Long:  "Deploy new contract",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
