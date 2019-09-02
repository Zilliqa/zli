package swap

import "github.com/spf13/cobra"

func init() {
	SwapCmd.AddCommand(submitCmd)
}

var submitCmd = &cobra.Command{
	Use:"submit",
	Short:"submit transaction",
	Long:"submit transaction",
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}
