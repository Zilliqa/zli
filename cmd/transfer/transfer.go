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
package transfer

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var amount string
var toAddr string
var wallet *core.Wallet
var privateKey string
var keystore string

func init() {
	TransferCmd.Flags().StringVarP(&amount, "amount", "a", "0", "amount to transfer")
	TransferCmd.Flags().StringVarP(&toAddr, "toAddr", "t", "", "to address")
	TransferCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key used to do this transfer")
	TransferCmd.Flags().StringVarP(&keystore, "keystore", "s", "", "keystore used to do this transfer")
}

var TransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer zilliqa token to a specific account",
	Long:  "Transfer zilliqa token to a specific account by using default account or a batch of accounts",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			panic(err.Error())
		}
		wallet = w
		if privateKey != "" {
			account, err := core.NewAccount(privateKey)
			if err != nil {
				panic(err.Error())
			}
			wallet.DefaultAccount = *account
		}
		if keystore != "" {
			fmt.Println("please type password to decrypt your keystore: ")
			pass, err := gopass.GetPasswd()
			if err != nil {
				panic(err.Error())
			}
			p, err := core.LoadPirvateKeyFromKeyStore(keystore, string(pass))
			if err != nil {
				panic(err.Error())
			}
			account, err := core.NewAccount(p)
			if err != nil {
				panic(err.Error())
			}
			wallet.DefaultAccount = *account
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)
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
			Priority:     false,
		}

		err := signer.Sign(tx, *p)
		if err != nil {
			panic(err.Error())
		}

		rsp := p.CreateTransaction(tx.ToTransactionPayload())
		if rsp.Error != nil {
			panic(rsp.Error)
		} else {
			result := rsp.Result.(map[string]interface{})
			hash := result["TranID"].(string)
			fmt.Printf("hash is %s", hash)
			tx.Confirm(hash, 1000, 3, p)
		}

	},
}
