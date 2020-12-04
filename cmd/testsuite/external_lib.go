package testsuite

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	core2 "github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var price int64
var limit int32

func init() {
	externalLib.Flags().Int64VarP(&price, "price", "p", 2000000000, "set gas price")
	externalLib.Flags().Int32VarP(&limit, "limit", "l", 10000, "set gas limit")
	TestSuite.AddCommand(externalLib)
}

var externalLib = &cobra.Command{
	Use:   "external",
	Short: "test external library",
	Long:  "test external library",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			panic(err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 1. deploy lib1
		fmt.Println("start to deploy lib1")
		lib1, err := ioutil.ReadFile("testsuite/contracts/ExternalLib/lib1.scilla")
		if err != nil {
			panic(err.Error())
		}

		type Constructor struct {
			Constructor string   `json:"constructor"`
			ArgTypes    []string `json:"argtypes"`
			Arguments   []string `json:"arguments"`
		}

		argtypes := make([]string, 0)
		arguments := make([]string, 0)

		cons := Constructor{
			Constructor: "True",
			ArgTypes:    argtypes,
			Arguments:   arguments,
		}

		initStruct1 := []core2.ContractValue{
			{
				VName: "_scilla_version",
				Type:  "Uint32",
				Value: "0",
			},
			{
				VName: "_library",
				Type:  "Bool",
				Value: cons,
			},
		}

		p := provider.NewProvider(wallet.API)
		fmt.Println(wallet.DefaultAccount.Address)
		result, err := p.GetBalance(wallet.DefaultAccount.Address)

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		contract := contract2.Contract{
			Code:     string(lib1),
			Init:     initStruct1,
			Signer:   signer,
			Provider: p,
		}

		deployParams := contract2.DeployParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(result.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		}

		tx, err := contract.Deploy(deployParams)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("contract address = ", tx.ContractAddress)

		tx.Confirm(tx.ID, 1000, 3, p)

		if tx.Status == core2.Rejected {
			fmt.Println("deploy lib1 failed")
			return
		}

		// 2. deploy lib3
		fmt.Println("start to deploy lib3")
		lib3, err2 := ioutil.ReadFile("testsuite/contracts/ExternalLib/lib3.scilla")
		if err2 != nil {
			panic(err2.Error())
		}

		ats := []string{
			"String",
			"ByStr20",
		}

		ars := []string{
			"TestLib1",
			"0x" + tx.ContractAddress,
		}

		initStruct3 := []core2.ContractValue{
			{
				VName: "_scilla_version",
				Type:  "Uint32",
				Value: "0",
			},
			{
				VName: "_extlibs",
				Type:  "List(Pair String ByStr20)",
				Value: []Constructor{
					{
						Constructor: "Pair",
						ArgTypes:    ats,
						Arguments:   ars,
					},
				},
			},
			{
				VName: "_library",
				Type:  "Bool",
				Value: cons,
			},
		}

		result, err = p.GetBalance(wallet.DefaultAccount.Address)
		if err != nil {
			panic(err.Error())
		}

		contract = contract2.Contract{
			Code:     string(lib3),
			Init:     initStruct3,
			Signer:   signer,
			Provider: p,
		}

		deployParams = contract2.DeployParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(result.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		}

		tx, err = contract.Deploy(deployParams)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("contract address = ", tx.ContractAddress)

		tx.Confirm(tx.ID, 1000, 3, p)

		if tx.Status == core2.Rejected {
			fmt.Println("deploy lib2 failed")
			return
		}

		// 3. deploy test contract
		fmt.Println("start to deploy hello contract")
		hello, err3 := ioutil.ReadFile("testsuite/contracts/ExternalLib/hello.scilla")
		if err3 != nil {
			panic(err3.Error())
		}

		ats = []string{
			"String",
			"ByStr20",
		}

		ars = []string{
			"TestLib3",
			"0x" + tx.ContractAddress,
		}

		initHello := []core2.ContractValue{
			{
				VName: "_scilla_version",
				Type:  "Uint32",
				Value: "0",
			},
			{
				VName: "_extlibs",
				Type:  "List(Pair String ByStr20)",
				Value: []Constructor{
					{
						Constructor: "Pair",
						ArgTypes:    ats,
						Arguments:   ars,
					},
				},
			},
		}

		result, err = p.GetBalance(wallet.DefaultAccount.Address)
		if err != nil {
			panic(err.Error())
		}

		contract = contract2.Contract{
			Code:     string(hello),
			Init:     initHello,
			Signer:   signer,
			Provider: p,
		}

		deployParams = contract2.DeployParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(result.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
		}

		tx, err = contract.Deploy(deployParams)

		if err != nil {
			panic(err.Error())
		}

		fmt.Println("contract address = ", tx.ContractAddress)

		tx.Confirm(tx.ID, 1000, 3, p)

		if tx.Status == core2.Rejected {
			fmt.Println("deploy hello contract failed")
			return
		}

		result, err = p.GetBalance(wallet.DefaultAccount.Address)
		if err != nil {
			panic(err)
		}

		addr, _ := bech32.ToBech32Address(tx.ContractAddress)

		// 4. call hello contract
		fmt.Println("start to call hello contract")
		contract = contract2.Contract{
			Address:  addr,
			Signer:   signer,
			Provider: p,
		}
		callParams := contract2.CallParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(result.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			Amount:       "0",
		}

		a := make([]core2.ContractValue, 0)

		tx, err4 := contract.Call("Hi", a, callParams, false)
		if err4 != nil {
			panic(err4.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)
		if tx.Status == core2.Rejected {
			fmt.Println("call hello contract failed")
			return
		}

		fmt.Println(tx.Receipt)
	},
}
