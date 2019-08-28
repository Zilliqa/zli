package contract

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
	callCmd.Flags().Int64VarP(&invokePrice, "price", "p", 10000000000, "set gas price")
	callCmd.Flags().Int32VarP(&invokeLimit, "limit", "l", 10000, "set gas limit")
	callCmd.Flags().StringVarP(&invokeAddress, "address", "a", "", "smart contract address")
	callCmd.Flags().IntVarP(&chainId, "chainId", "d", 0, "chain id")
	callCmd.Flags().StringVarP(&api, "api", "u", "", "api url")
	callCmd.Flags().StringVarP(&invokeAmount, "amount", "m", "0", "token amount to transfer to the contract")
	callCmd.Flags().StringVarP(&privateKey, "private_key", "k", "", "private key used to call to the contract")
	callCmd.Flags().BoolVarP(&invokePriority, "priority", "f", false, "setup priority of transaction")
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

		//fmt.Println(privateKey)
		//
		//if privateKey != "" {
		//	account, err := core.NewAccount(privateKey)
		//	if err != nil {
		//		panic(err.Error())
		//	}
		//	wallet.DefaultAccount = *account
		//}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(invokeTransition) == 0 {
			panic("invalid transition")
		}

		var a []contract2.Value
		err := json.Unmarshal([]byte(invokeArgs), &a)

		if !validator.IsBech32(invokeAddress) {
			invokeAddress, _ = bech32.ToBech32Address(invokeAddress)
		}

		if err != nil {
			panic(err.Error())
		}

		p := provider.NewProvider(wallet.API)
		//fmt.Println(wallet.DefaultAccount.Address)
		result := p.GetBalance(wallet.DefaultAccount.Address)
		if result.Error != nil {
			panic(result.Error.Message)
		}

		balance := result.Result.(map[string]interface{})
		nonce, _ := balance["nonce"].(json.Number).Int64()

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Address:  invokeAddress,
			Singer:   signer,
			Provider: p,
		}

		params := contract2.CallParams{
			Version:      strconv.FormatInt(int64(LaksaGo.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			Amount:       invokeAmount,
		}

		err, tx := contract.Call(invokeTransition, a, params, invokePriority, 1000, 3)

		if err != nil {
			panic(err.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)

	},
}
