package transfer

import (
	"fmt"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/account"
	"github.com/FireStack-Lab/LaksaGo/provider"
	"github.com/FireStack-Lab/LaksaGo/transaction"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var amount string
var toAddr string
var wallet *core.Wallet

func init() {
	TransferCmd.Flags().StringVarP(&amount, "amount", "a", "0", "amount to transfer")
	TransferCmd.Flags().StringVarP(&toAddr, "toAddr", "t", "", "to address")
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)
		p := provider.NewProvider(wallet.API)

		tx := &transaction.Transaction{
			Version:      strconv.FormatInt(int64(LaksaGo.Pack(wallet.ChainID, 1)), 10),
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
