package swap

import (
	"encoding/json"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/account"
	"github.com/FireStack-Lab/LaksaGo/bech32"
	contract2 "github.com/FireStack-Lab/LaksaGo/contract"
	"github.com/FireStack-Lab/LaksaGo/provider"
	"github.com/FireStack-Lab/LaksaGo/validator"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"zli/core"
)

var fundKeyStorePath string
var password string

func init() {
	fundCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	fundCmd.Flags().IntVarP(&chainId, "chainId", "c", 333, "the message version of the network")
	fundCmd.Flags().StringVarP(&walletAddress, "address", "a", "zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2", "address of the wallet contract")
	fundCmd.Flags().StringVarP(&gasPrice, "price", "p", "10000000000", "gas price")
	fundCmd.Flags().StringVarP(&gasLimit, "limit", "l", "1000", "gas limit")
	fundCmd.Flags().StringVarP(&amount, "amount", "m", "0", "token amount will be transfer to the smart contract")
	fundCmd.Flags().StringVarP(&password, "password", "s", "", "password to decrypt the keystore")
}

var fundCmd = &cobra.Command{
	Use:   "fund",
	Short: "Add funds to wallet contract",
	Long:  "Add funds to wallet contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		if fundKeyStorePath == "" {
			panic("the path of the key store should not be empty")
		}
		if password == "" {
			panic("password should not be empty")
		}
		p, err := core.LoadPirvateKeyFromKeyStore(fundKeyStorePath, password)
		if err != nil {
			panic("load private key from keystore error = " + err.Error())
		}
		w, err := core.NewWallet([]byte(p), chainId, api)
		if err != nil {
			panic("init wallet error = " + err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		var a []contract2.Value

		if !validator.IsBech32(walletAddress) {
			walletAddress, _ = bech32.ToBech32Address(walletAddress)
		}

		p := provider.NewProvider(api)
		result := p.GetBalance(wallet.API)
		if result.Error != nil {
			panic(result.Error.Message)
		}

		balance := result.Result.(map[string]interface{})
		nonce, _ := balance["nonce"].(json.Number).Int64()

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Address:  walletAddress,
			Singer:   signer,
			Provider: p,
		}

		params := contract2.CallParams{
			Version:      strconv.FormatInt(int64(LaksaGo.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(nonce+1, 10),
			GasPrice:     gasPrice,
			GasLimit:     gasLimit,
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			Amount:       amount,
		}

		err, tx := contract.Call("AddFunds", a, params, false, 1000, 3)
		if err != nil {
			panic(err.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)
	},
}
