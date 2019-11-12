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
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
	"zli/core"
)

var accounts string
var batch int
var amount string
var toAddr string
var interval int32
var chainId int
var api string

func init() {
	transferCmd.Flags().StringVarP(&amount, "amount", "a", "0", "amount to transfer")
	transferCmd.Flags().StringVarP(&toAddr, "toAddr", "t", "", "to address")
	transferCmd.Flags().StringVarP(&accounts, "accounts", "c", "./testAccounts.txt", "path of testAccounts.txt")
	transferCmd.Flags().IntVarP(&batch, "batch", "b", 0, "the number of each spam")
	transferCmd.Flags().Int32VarP(&interval, "interval", "i", 1, "interval time (second) between each batch request")
	transferCmd.Flags().IntVarP(&chainId, "chainId", "d", 333, "chain id")
	transferCmd.Flags().StringVarP(&api, "api", "p", "https://dev-api.zilliqa.com/", "api url")
	SpamCmd.AddCommand(transferCmd)
}

func send(wallet *core.Wallet, group *sync.WaitGroup) {
	defer group.Done()
	fmt.Println("start to use private key ", wallet.DefaultAccount.PrivateKey, " to generate transaction")
	p := provider.NewProvider(wallet.API)
	tx := &transaction.Transaction{
		Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
		SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		ToAddr:       toAddr,
		Amount:       amount,
		GasPrice:     "1000000000",
		GasLimit:     "1",
		Code:         "",
		Data:         "",
	}

	signer := account.NewWallet()
	signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)
	err := signer.Sign(tx, *p)
	if err != nil {
		fmt.Println(err.Error())
	}

	rsp := p.CreateTransaction(tx.ToTransactionPayload())
	if rsp == nil {
		fmt.Println("create transaction error")
		return
	}
	if rsp.Error != nil {
		fmt.Println(rsp.Error)
	} else {
		result := rsp.Result.(map[string]interface{})
		hash := result["TranID"].(string)
		fmt.Printf("hash is %s", hash)
		//tx.Confirm(hash, 20, 10, p)
	}
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "to specific account",
	Long:  "send a large number of transactions to a specific account ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load file from ", accounts)
		accs, err := core.LoadFrom(accounts)
		fmt.Println("account number is ",len(accs))
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
					send(wallet, wg)
				}()

			}
			wg.Wait()
		}
		time.Sleep(time.Duration(rand.Int31n(interval)) * time.Second)
	},
}
