package staking

import "github.com/spf13/cobra"

var api string
var contractAddress string

var StakingCmd = &cobra.Command{
	Use:   "staking",
	Short: "tools to interact with zilliqa staking contract",
	Long:  "tools to interact with zilliqa staking contract",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
