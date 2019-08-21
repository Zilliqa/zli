package contract

import (
	"encoding/json"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/account"
	contract2 "github.com/FireStack-Lab/LaksaGo/contract"
	"github.com/FireStack-Lab/LaksaGo/provider"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var code string
var initJson string
var price int64
var limit int32
var wallet *core.Wallet
var chainId int
var api string

func init() {
	deployCmd.Flags().StringVarP(&code, "code", "c", "", "file that contains contract code")
	deployCmd.Flags().StringVarP(&initJson, "init", "i", "", "file that contains init json")
	deployCmd.Flags().Int64VarP(&price, "price", "p", 10000000000, "set gas price")
	deployCmd.Flags().Int32VarP(&limit, "limit", "l", 10000, "set gas limit")
	deployCmd.Flags().IntVarP(&chainId, "chainId", "d", 333, "chain id")
	deployCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	ContractCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy new contract",
	Long:  "Deploy new contract",
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		c, err := ioutil.ReadFile(code)
		if err != nil {
			panic(err.Error())
		}

		i, err := ioutil.ReadFile(initJson)
		if err != nil {
			panic(err.Error())
		}

		var initArray []contract2.Value
		_ = json.Unmarshal(i, &initArray)

		p := provider.NewProvider(wallet.API)

		result := p.GetBalance(wallet.DefaultAccount.Address)
		if result.Error != nil {
			panic(result.Error.Message)
		}

		balance := result.Result.(map[string]interface{})
		nonce, _ := balance["nonce"].(json.Number).Int64()

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)
		contract := contract2.Contract{
			Code:     string(c),
			Init:     initArray,
			Singer:   signer,
			Provider: p,
		}

		deployParams := contract2.DeployParams{
			Version:      strconv.FormatInt(int64(LaksaGo.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		}

		tx, err := contract.Deploy(deployParams)

		if err != nil {
			panic(err.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)

	},
}
