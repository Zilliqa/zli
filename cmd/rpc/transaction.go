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
package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/spf13/cobra"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var transactionId string
var api string
var wallet *core.Wallet

func init() {
	transactionCmd.Flags().StringVarP(&transactionId, "transaction", "t", "", "transaction id")
	transactionCmd.Flags().StringVarP(&api, "api", "u", "", "api endpoint")
	RPCCmd.AddCommand(transactionCmd)
}

var transactionCmd = &cobra.Command{
	Use:   "transaction",
	Short: "Get transaction detail by transaction id",
	Long:  "Get transaction detail by transaction id",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			fmt.Println("cannot load wallet = ", err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		if transactionId == "" {
			panic("transaction id cannot be empty")
		}
		if api == "" && wallet == nil {
			panic("wallet ==  nil && transaction id empty")
		}

		var a string
		if api != "" {
			a = api
		} else {
			a = wallet.API
		}

		p := provider.NewProvider(a)
		txn, err := p.GetTransaction(transactionId)

		res, err := json.Marshal(txn)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(string(res))
	},
}
