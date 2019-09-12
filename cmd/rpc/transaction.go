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
		response := p.GetTransaction(transactionId)
		if response == nil {
			panic("cannot get response")
		}
		if response.Error != nil {
			panic(response.Error)
		}

		res, err := json.Marshal(response.Result)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(string(res))
	},
}
