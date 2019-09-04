package swap

import "github.com/spf13/cobra"

var api string
var chainId int
var walletAddress string
var gasPrice string
var gasLimit string
var amount string
var priority bool


var SwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Just for internal swap",
	Long:  "Just for internal swap",
	Run: func(cmd *cobra.Command, args []string) {

	},
}