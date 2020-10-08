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
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var amount string
var toAddr string
var wallet *core.Wallet
var privateKey string
var keystore string
var unit string

func init() {
	TransferCmd.Flags().StringVarP(&amount, "amount", "a", "0", "amount to transfer")
	TransferCmd.Flags().StringVarP(&toAddr, "toAddr", "t", "", "to address")
	TransferCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key used to do this transfer")
	TransferCmd.Flags().StringVarP(&keystore, "keystore", "s", "", "keystore used to do this transfer")
	TransferCmd.Flags().StringVarP(&unit, "unit", "u", "qa", "provide either zil, li or qa")
}

var TransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer zilliqa token to a specific account",
	Long:  "Transfer zilliqa token to a specific account by using default account or a batch of accounts",
	PreRun: func(cmd *cobra.Command, args []string) {
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			panic(err.Error())
		}

		unit = strings.ToLower(unit)
		if unit != "zil" && unit != "li" && unit != "qa" {
			panic(fmt.Errorf("unit needs to be either zil, li or qa"))
		}

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
		strAmount, _ := strconv.ParseFloat(amount, 64)

		var qaFloat float64
		switch unit {
		case "zil":
			qaFloat = util.ToQa(strAmount, util.ZIL)
		case "li":
			qaFloat = util.ToQa(strAmount, util.LI)
		default:
			qaFloat = util.ToQa(strAmount, util.QA)
		}

		qaAmount := strconv.FormatFloat(qaFloat, 'f', 0, 64)

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)
		p := provider.NewProvider(wallet.API)

		tx := &transaction.Transaction{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			ToAddr:       toAddr,
			Amount:       qaAmount,
			GasPrice:     "2000000000",
			GasLimit:     "1",
			Code:         "",
			Data:         "",
			Priority:     false,
		}

		err := signer.Sign(tx, *p)
		if err != nil {
			panic(err.Error())
		}

		rsp, err := p.CreateTransaction(tx.ToTransactionPayload())
		if err != nil {
			panic(err.Error())
		}
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
