package wallet

import (
	"github.com/spf13/cobra"
)


var WalletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Init a new wallet or get exist wallet info",
	Long:  `Init a new wallet or get exist wallet info`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
