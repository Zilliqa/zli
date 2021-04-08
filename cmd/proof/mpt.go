package proof

import (
	"fmt"
	core2 "github.com/Zilliqa/gozilliqa-sdk/core"
	"github.com/Zilliqa/gozilliqa-sdk/mpt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"strings"
	wallet2 "zli/cmd/wallet"
	"zli/core"
)

var wallet *core.Wallet
var api string
var contractAddress string
var key string
var block string

func init() {
	mptCmd.Flags().StringVarP(&contractAddress, "contract", "c", "", "smart contract address")
	mptCmd.Flags().StringVarP(&api, "api", "u", "", "api endpoint")
	mptCmd.Flags().StringVarP(&key, "key", "k", "", "for now only support simple key without nesting")
	mptCmd.Flags().StringVarP(&block, "block", "b", "", "block num")
	ProofCmd.AddCommand(mptCmd)
}

var mptCmd = &cobra.Command{
	Use:   "mpt",
	Short: "test mpt verification",
	Long:  "test mpt verification",
	PreRun: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		w, err := core.LoadFromFile(home + "/" + wallet2.DefaultConfigName)
		if err != nil {
			fmt.Println("cannot load wallet = ", err.Error())
		}
		wallet = w
	},
	Run: func(cmd *cobra.Command, args []string) {
		var a string
		if api != "" {
			a = api
		} else {
			a = wallet.API
		}
		p := provider.NewProvider(a)
		contractAddress = strings.TrimPrefix(contractAddress, "0x")
		storageKey := core2.GenerateStorageKey(contractAddress, key, []string{})
		hashedStorageKey := util.Sha256(storageKey)

		proofResult, err := p.GetStateProof(contractAddress, util.EncodeHex(hashedStorageKey), &block)
		if err != nil {
			fmt.Println("get state proof error: " + err.Error())
			return
		}

		if proofResult == nil {
			fmt.Println("proof is nil")
			return
		}

		txBlock, err := p.GetTxBlock(block)
		if err != nil {
			fmt.Println("get tx block error: " + err.Error())
			return
		}

		if proofResult.AccountProof == nil {
			fmt.Println("account proof is nil")
			return
		}

		if proofResult.StateProof == nil {
			fmt.Println("state proof is nil")
			return
		}

		var proof [][]byte
		for _, p := range proofResult.AccountProof {
			bytes := util.DecodeHex(p)
			proof = append(proof, bytes)
		}

		db := mpt.NewFromProof(proof)
		stateRoot := util.DecodeHex(txBlock.Header.StateRootHash)
		value, err := mpt.Verify([]byte(contractAddress), db, stateRoot)
		if err != nil {
			fmt.Println("account proof failed, err: ", err.Error())
			return
		}

		accountBase, _ := core2.AccountBaseFromBytes(value)

		var proof2 [][]byte
		for _, p := range proofResult.StateProof {
			bytes := util.DecodeHex(p)
			proof2 = append(proof2, bytes)
		}

		db3 := mpt.NewFromProof(proof2)
		value2, err2 := mpt.Verify([]byte((util.EncodeHex(hashedStorageKey))), db3, accountBase.StorageRoot)
		if err2 != nil {
			fmt.Println("state proof failed, err: " + err2.Error())
			return
		}

		fmt.Println(string(value2))

	},
}
