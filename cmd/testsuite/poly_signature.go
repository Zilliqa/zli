package testsuite

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	contract2 "github.com/Zilliqa/gozilliqa-sdk/contract"
	core2 "github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"

	wallet2 "zli/cmd/wallet"
	"zli/core"
)

func init() {
	TestSuite.AddCommand(polysig)
}

var polysig = &cobra.Command{
	Use:   "polysig",
	Short: "test polynetwork signature",
	Long:  "test polynetwork signature",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			panic(err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 1. deploy simplified cross chain manager contract
		fmt.Println("start to deploy ccmc")
		ccm, err := ioutil.ReadFile("testsuite/contracts/Polynetwork/ccm.scilla")
		if err != nil {
			panic(err.Error())
		}

		p := provider.NewProvider(wallet.API)
		fmt.Println(wallet.DefaultAccount.Address)
		result, err := p.GetBalance(wallet.DefaultAccount.Address)

		signer := account.NewWallet()
		signer.AddByPrivateKey(wallet.DefaultAccount.PrivateKey)

		initStruct := []core2.ContractValue{
			{
				VName: "_scilla_version",
				Type:  "Uint32",
				Value: "0",
			},
		}

		contract := contract2.Contract{
			Code:     string(ccm),
			Init:     initStruct,
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

		// 2. call cmcc
		fmt.Println("start to call cmcc")
		addr, _ := bech32.ToBech32Address(tx.ContractAddress)
		contract = contract2.Contract{
			Address:  addr,
			Signer:   signer,
			Provider: p,
		}

		pairs := []core2.ParamConstructor{}
		proofArguments := make([]interface{}, 0)
		proofArguments = append(proofArguments, "0x20e0b21bb1950ee438fcf4aa858353352d0d47002235d6a40efaa9b1b3da8e6466050000000000000020190794747ba9de191888617877c325d4962d975f20bb2ee692b299df5abf6e450273a0141a785cfc5dbec2e1518e1b1d369154d0ce579640120000000000000014d73c6b871b4d0e130d64581993b745fc938a5be706756e6c6f636bc6117a62726b6c2e312e31382e6238633234661432339fa037f7ae1dfff25e13c6451a80289d61f414c38fc92d13f3d15feb364bdba2f89fd7ae491667f013fdd11c2023ce61110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001408d8f59e475830d9a1bb97d74285c4d34c6dac0814699362e4ceb7a6b6703741f889435b7cc9d74b4df572000000000000000000000000000000000000000000000000000000000000")
		proofArguments = append(proofArguments, pairs)
		proofConstructor := core2.ParamConstructor{
			Constructor: "Polynetwork.Proof",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   proofArguments,
		}

		headerPairs := []core2.ParamConstructor{}
		headProofArguments := make([]interface{}, 0)
		headProofArguments = append(headProofArguments, "0x")
		proofArguments = append(proofArguments, pairs)
		headProofArguments = append(headProofArguments, headerPairs)
		headProofConstructor := core2.ParamConstructor{
			Constructor: "Polynetwork.Proof",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   headProofArguments,
		}

		var sigs []core2.ParamConstructor

		sigs = append(sigs, core2.ParamConstructor{
			Constructor: "Polynetwork.Signature",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   []interface{}{"0x78f7166af83e5ec11666987bf14fc24e5567ae1cb36a28328d59b21de27585c8006fcfdf8831748bfed09be8d2a6d4704d6206040ebc5fd920a62fa6dc0b373101"},
		})

		sigs = append(sigs, core2.ParamConstructor{
			Constructor: "Polynetwork.Signature",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   []interface{}{"0x78f7166af83e5ec11666987bf14fc24e5567ae1cb36a28328d59b21de27585c8006fcfdf8831748bfed09be8d2a6d4704d6206040ebc5fd920a62fa6dc0b373101"},
		})

		sigs = append(sigs, core2.ParamConstructor{
			Constructor: "Polynetwork.Signature",
			ArgTypes:    make([]interface{}, 0),
			Arguments:   []interface{}{"0x78f7166af83e5ec11666987bf14fc24e5567ae1cb36a28328d59b21de27585c8006fcfdf8831748bfed09be8d2a6d4704d6206040ebc5fd920a62fa6dc0b373101"},
		})

		a := []core2.ContractValue{
			{
				VName: "curKeepers",
				Type:  "List ByStr20",
				Value: []string{"0x3dfccb7b8a6972cde3b695d3c0c032514b0f3825", "0x4c46e1f946362547546677bfa719598385ce56f2", "0xf81f676832f6dfec4a5d0671bd27156425fcef98", "0x51b7529137d34002c4ebd81a2244f0ee7e95b2c0"},
			},
			{
				"proof",
				"Polynetwork.Proof",
				proofConstructor,
			},
			{
				"rawHeader",
				"ByStr",
				"0x000000000000000000000000188a966e5475b590e4136ee50fbf380152d39887a5c09bac1b71327963359b8c0b5d84ffb4d549e93754646c27f22d61e3170c5ac45afc422643f36eba77a1c13b5d433f1856da7dd8bafa5ae21d77170958002a3f247616ac0da44fb6f41f32464dfcc84601a107cd977ca29c62741178b9b45f358eecac1914c121e082677ce4060f62d744fe00b44e499e65ce80a1fd13017b226c6561646572223a342c227672665f76616c7565223a22424e4a6a3371786b7a424730344869746d4950447857585275633474725a6574436f4d566d754d59592b763743765739364b71486a6e35394539684c53577a4a70304d634c4259516f555553504278552f6d675033746f3d222c227672665f70726f6f66223a227652732f4945486631574d327237415a55656468504261377950415a4d3230324461327149654f6d46796b5146446d6b4f5a65434d65753559424f357836775a66307431474e3644527658384735326a5459366d33773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a31363632303030302c226e65775f636861696e5f636f6e666967223a6e756c6c7d0000000000000000000000000000000000000000",
			},
			{
				"headerProof",
				"Polynetwork.Proof",
				headProofConstructor,
			},
			{
				"curRawHeader",
				"ByStr",
				"0x",
			},
			{
				"headerSig",
				"List Polynetwork.Signature",
				sigs,
			},
		}

		result, _ = p.GetBalance(wallet.DefaultAccount.Address)

		callParams := contract2.CallParams{
			Version:      strconv.FormatInt(int64(util.Pack(wallet.ChainID, 1)), 10),
			Nonce:        strconv.FormatInt(result.Nonce+1, 10),
			GasPrice:     strconv.FormatInt(price, 10),
			GasLimit:     strconv.FormatInt(int64(limit), 10),
			SenderPubKey: strings.ToUpper(wallet.DefaultAccount.PublicKey),
			Amount:       "0",
		}

		tx, err = contract.Call("VerifyHeaderAndExecuteTx", a, callParams, false)
		if err != nil {
			panic(err.Error())
		}

		tx.Confirm(tx.ID, 1000, 3, p)
		fmt.Println(tx.Receipt)

	},
}
