package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/spf13/cobra"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var addr string

func init() {
	balanceCmd.Flags().StringVarP(&addr, "addr", "a", "", "address (base16 or bech32)")
	balanceCmd.Flags().StringVarP(&api, "api", "u", "", "api endpoint")
	RPCCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get balance and nonce by address(base16 or bech32)",
	Long:  "Get balance and nonce by address(base16 or bech32)",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			fmt.Println("cannot load wallet = ", err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		if addr == "" {
			panic("address cannot be empty")
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

		var normalAddr string
		if validator.IsBech32(addr) {
			normalAddr, _ = bech32.FromBech32Addr(addr)
		} else if strings.HasPrefix(addr, "0x") {
			normalAddr = strings.Trim(addr, "0x")
		} else {
			normalAddr = addr
		}
		balance, err := p.GetBalance(normalAddr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		res, err1 := json.Marshal(balance)
		if err1 != nil {
			fmt.Println(err1.Error())
			return
		}
		fmt.Println(string(res))
	},
}
