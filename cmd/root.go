/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"zli/cmd/account"
	"zli/cmd/contract"
	"zli/cmd/rpc"
	"zli/cmd/testsuit"
	"zli/cmd/transfer"
	"zli/cmd/wallet"
)

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(contract.ContractCmd)
	RootCmd.AddCommand(wallet.WalletCmd)
	RootCmd.AddCommand(account.AccountCmd)
	RootCmd.AddCommand(transfer.TransferCmd)
	RootCmd.AddCommand(rpc.RPCCmd)
	RootCmd.AddCommand(testsuit.TestsultCmd)
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
