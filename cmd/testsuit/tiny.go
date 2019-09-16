package testsuit

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var address string
var wallet *core.Wallet
var transitions = []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "t11", "t12", "t13", "t14", "t15", "t16", "t17", "t18"}

func init() {
	tinyCmd.Flags().StringVarP(&address, "address", "a", "", "the address of tiny contract")
	TestsultCmd.AddCommand(tinyCmd)
}

var tinyCmd = &cobra.Command{
	Use:   "tiny",
	Short: "test tiny contract for corner cases",
	Long:  "test tiny contract for corner cases",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			panic(err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		if address == "" {
			panic("invalid contract address")
		}

		p := provider.NewProvider(wallet.API)

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Address:  address,
			Singer:   signer,
			Provider: p,
		}
		a := []contract2.Value{}

		for index, value := range transitions {
			fmt.Println("start to invoke transition ", index+1)
			result := p.GetBalance(wallet.DefaultAccount.Address)
			if result.Error != nil {
				panic(result.Error.Message)
			}
			balance := result.Result.(map[string]interface{})
			nonce, _ := balance["nonce"].(json.Number).Int64()
			params := contract2.CallParams{
				Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
				Nonce:        strconv.FormatInt(nonce+1, 10),
				GasPrice:     "1000000000",
				GasLimit:     "100000",
				SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
				Amount:       "0",
			}
			err, tx := contract.Call(value, a, params, false, 1000, 3)
			if err != nil {
				panic(err.Error())
			}
			tx.Confirm(tx.ID, 1000, 3, p)
			r, ok := p.GetTransaction(tx.ID).Result.(map[string]interface{})
			if !ok {
				panic("get transaction result failed")
			}
			receipt := r["receipt"].(map[string]interface{})
			success := receipt["success"].(bool)
			if !success {
				fmt.Println("test failed at transition ", index+1)
				os.Exit(1)
			}

		}

	},
}
