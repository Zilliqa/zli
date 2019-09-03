package swap

import (
	"encoding/json"
	"fmt"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/account"
	"github.com/FireStack-Lab/LaksaGo/bech32"
	contract2 "github.com/FireStack-Lab/LaksaGo/contract"
	"github.com/FireStack-Lab/LaksaGo/provider"
	"github.com/FireStack-Lab/LaksaGo/validator"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"zli/core"
)

var fundWallet *core.Wallet
var fundKeyStorePath string

func init() {
	fundCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	fundCmd.Flags().IntVarP(&chainId, "chainId", "c", 333, "the message version of the network")
	fundCmd.Flags().StringVarP(&walletAddress, "address", "a", "zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2", "address of the fundWallet contract")
	fundCmd.Flags().StringVarP(&gasPrice, "price", "p", "10000000000", "gas price")
	fundCmd.Flags().StringVarP(&gasLimit, "limit", "l", "10000", "gas limit")
	fundCmd.Flags().StringVarP(&amount, "amount", "m", "0", "token amount will be transfer to the smart contract")
	fundCmd.Flags().StringVarP(&fundKeyStorePath, "fundkeystore", "f", "", "fund keystore")
	SwapCmd.AddCommand(fundCmd)
}

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Add funds to fundWallet contract",
	Long:  "Add funds to fundWallet contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logfile, _ := os.Create("fund.log")
		log.SetOutput(logfile)
		if fundKeyStorePath == "" {
			panic("the path of the key store should not be empty")
		}
		fmt.Println("please type password to decrypt your keystore: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			panic(err.Error())
		}
		p, err := core.LoadPirvateKeyFromKeyStore(fundKeyStorePath, string(pass))
		if err != nil {
			panic("load private key from keystore error = " + err.Error())
		}
		w, err := core.NewWallet(LaksaGo.DecodeHex(p), chainId, api)
		if err != nil {
			panic("init fundWallet error = " + err.Error())
		}
		fundWallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		a := []contract2.Value{}

		if !validator.IsBech32(walletAddress) {
			walletAddress, _ = bech32.ToBech32Address(walletAddress)
		}

		fmt.Printf("start to send %s qa to address %s\n", amount, walletAddress)
		fmt.Println("please type Y to confirm: ")
		var confirmed string
		_, err := fmt.Scanln(&confirmed)
		if err != nil || confirmed != "Y" {
			fmt.Printf("transfer cancled")
			return
		}
		log.Printf("start to send %s qa to address %s\n", amount, walletAddress)
		p := provider.NewProvider(api)
		result := p.GetBalance(fundWallet.DefaultAccount.Address)
		if result.Error != nil {
			panic(result.Error.Message)
		}

		balance := result.Result.(map[string]interface{})
		nonce, _ := balance["nonce"].(json.Number).Int64()

		signer := account.NewWallet()
		signer.AddByPrivateKey(fundWallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Address:  walletAddress,
			Singer:   signer,
			Provider: p,
		}

		params := contract2.CallParams{
			Version:      strconv.FormatInt(int64(LaksaGo.Pack(fundWallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(nonce+1, 10),
			GasPrice:     gasPrice,
			GasLimit:     gasLimit,
			SenderPubKey: strings.ToUpper(fundWallet.DefaultAccount.PublicKey),
			Amount:       amount,
		}

		err, tx := contract.Call("AddFunds", a, params, false, 1000, 3)
		if err != nil {
			panic(err.Error())
		}

		log.Printf("start to poll transaction: %s", tx.ID)
		tx.Confirm(tx.ID, 1000, 3, p)
		err, recipients := getReceiptForTransaction(p, tx.ID)
		log.Printf("get recipients for %s: %s\n", tx.ID, recipients)
	},
}
