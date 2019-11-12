//  This file is part of zli
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

package spam

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"zli/core"
)

var invokeTransition string
var invokeArgs string
var invokePrice int64
var invokeLimit int32
var invokeAddress string
var invokePriority bool

func init() {
	invokeCmd.Flags().StringVarP(&invokeTransition, "transition", "t", "", "transition will be called")
	invokeCmd.Flags().StringVarP(&invokeArgs, "args", "r", "", "args will be passed to transition")
	invokeCmd.Flags().Int64VarP(&invokePrice, "price", "p", 10000000000, "set gas price")
	invokeCmd.Flags().Int32VarP(&invokeLimit, "limit", "l", 10000, "set gas limit")
	invokeCmd.Flags().StringVarP(&invokeAddress, "address", "a", "", "smart contract address")
	invokeCmd.Flags().IntVarP(&batch, "batch", "b", 0, "the number of each spam")
	invokeCmd.Flags().IntVarP(&chainId, "chainId", "d", 333, "chain id")
	invokeCmd.Flags().StringVarP(&api, "api", "i", "https://dev-api.zilliqa.com/", "api url")
	invokeCmd.Flags().BoolVarP(&invokePriority, "priority", "f", false, "setup priority of transaction")
	invokeCmd.Flags().StringVarP(&accounts, "accounts", "c", "./testAccounts.txt", "path of testAccounts.txt")
	SpamCmd.AddCommand(invokeCmd)
}

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "on a specific contract",
	Long:  "send a large number of invocation transactions on a specific smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("contract address ", invokeAddress)
		fmt.Println("load file from ", accounts)
		accs, err := core.LoadFrom(accounts)
		if err != nil {
			panic(err.Error())
		}

		if len(accs) < batch {
			panic("the length of accounts should not less than batch number")
		}

		batchAccount := core.Split(accs, batch)
		for index, value := range batchAccount {
			fmt.Println("start to spam ", index)
			wg := &sync.WaitGroup{}
			wg.Add(len(value))
			for _, w := range value {
				wallet, err := core.FromPrivateKeyAndChain(util.DecodeHex(w.PrivateKey[:]), chainId, api)
				if err != nil {
					panic(err.Error())
				}
				go func() {
					invoke(wallet, wg)
				}()

			}
			wg.Wait()
		}
		time.Sleep(time.Duration(rand.Int31n(interval)) * time.Second)

	},
}

func invoke(wallet *core.Wallet, group *sync.WaitGroup) {
	defer group.Done()
	fmt.Println("start to use private key ", wallet.DefaultAccount.PrivateKey, " to generate ")

	var a []contract2.Value
	err := json.Unmarshal([]byte(invokeArgs), &a)
	//var sb strings.Builder
	//for i:=0;i<10;i++{
	//	u, _ := uuid.NewUUID()
	//	sb.WriteString(u.String())
	//}
	//
	//
	//a[0].Value = interface{}(sb.String())
	//a[1].Value = interface{}(sb.String())
	if err != nil {
		panic(err.Error())
	}

	p := provider.NewProvider(wallet.API)
	result := p.GetBalance(wallet.DefaultAccount.Address)
	if result.Error != nil {
		fmt.Println(result.Error.Message)
	}

	var nonce int64
	if result.Result == nil {
		nonce = 0
	} else {
		balance := result.Result.(map[string]interface{})
		nonce, _ = balance["nonce"].(json.Number).Int64()
	}

	signer := account.NewWallet()
	signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

	contract := contract2.Contract{
		Address:  invokeAddress,
		Singer:   signer,
		Provider: p,
	}
	params := contract2.CallParams{
		Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
		Nonce:        strconv.FormatInt(nonce+1, 10),
		GasPrice:     strconv.FormatInt(invokePrice, 10),
		GasLimit:     strconv.FormatInt(int64(invokeLimit), 10),
		SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		Amount:       "0",
	}

	err, tx := contract.Call(invokeTransition, a, params, invokePriority, 1000, 3)

	if err != nil {
		fmt.Println(err.Error())
	}

	if tx == nil {
		fmt.Println("create tx failed")
	} else {
		//tx.Confirm(tx.ID, 1000, 3, p)
	}
}
