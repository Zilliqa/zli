package swap

import (
	"github.com/spf13/cobra"
	"zli/core"
)

var api string
var chainId int
var walletAddress string
var gasPrice string
var gasLimit string
var amount string
var wallet *core.Wallet

var SwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Just for internal swap",
	Long:  "Just for internal swap",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
