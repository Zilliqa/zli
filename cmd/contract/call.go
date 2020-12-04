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
package contract

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	core2 "github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var invokeTransition string
var invokeArgs string
var invokePrice int64
var invokeLimit int32
var invokeAddress string
var invokeAmount string
var invokePriority bool

func init() {
	callCmd.Flags().StringVarP(&invokeTransition, "transition", "t", "", "transition will be called")
	callCmd.Flags().StringVarP(&invokeArgs, "args", "r", "", "args will be passed to transition")
	callCmd.Flags().Int64VarP(&invokePrice, "price", "p", 2000000000, "set gas price")
	callCmd.Flags().Int32VarP(&invokeLimit, "limit", "l", 10000, "set gas limit")
	callCmd.Flags().StringVarP(&invokeAddress, "address", "a", "", "smart contract address")
	callCmd.Flags().IntVarP(&chainId, "chainId", "d", 0, "chain id")
	callCmd.Flags().StringVarP(&api, "api", "u", "", "api url")
	callCmd.Flags().StringVarP(&invokeAmount, "amount", "m", "0", "token amount to transfer to the contract")
	callCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key used to call to the contract")
	callCmd.Flags().BoolVarP(&invokePriority, "priority", "f", false, "setup priority of transaction")
	callCmd.Flags().StringVarP(&keystore, "keystore", "s", "", "keystore used to deploy the contract")
	ContractCmd.AddCommand(callCmd)
}

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Call a exist contract",
	Long:  "Call a exist contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			panic(err.Error())
		}
		wallet = w
		if chainId != 0 && api != "" {
			wallet.API = api
			wallet.ChainID = chainId
		}

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
		if len(invokeTransition) == 0 {
			panic("invalid transition")
		}

		var a []core2.ContractValue
		err := json.Unmarshal([]byte(invokeArgs), &a)

		if !validator.IsBech32(invokeAddress) {
			invokeAddress, _ = bech32.ToBech32Address(invokeAddress)
		}

		if err != nil {
			panic(err.Error())
		}

		p := provider.NewProvider(wallet.API)
		balAndNonce, err := p.GetBalance(wallet.DefaultAccount.Address)
		if err != nil {
			panic(err)
		}

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Address:  invokeAddress,
			Signer:   signer,
			Provider: p,
		}

		params := contract2.CallParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(balAndNonce.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(invokeLimit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			Amount:       invokeAmount,
		}

		tx, err := contract.Call(invokeTransition, a, params, invokePriority)

		if err != nil {
			panic(err.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)

	},
}
