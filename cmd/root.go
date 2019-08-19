package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"zli/cmd/account"
	"zli/cmd/contract"
	"zli/cmd/transfer"
	"zli/cmd/wallet"
)

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(contract.ContractCmd)
	RootCmd.AddCommand(wallet.WalletCmd)
	RootCmd.AddCommand(account.AccountCmd)
	RootCmd.AddCommand(transfer.TransferCmd)
}

var RootCmd = &cobra.Command{
	Use:   "zli",
	Short: "Zli is a command line tool based on zilliqa golang sdk",
	Long:  `A convenient command line tool to generate accounts, run integration testings or run http server .etc`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
