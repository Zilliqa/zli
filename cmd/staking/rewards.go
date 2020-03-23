package staking

import (
	"encoding/json"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/spf13/cobra"
)

var ssn string

func init() {
	rewardsCmd.Flags().StringVarP(&api, "api", "a", "https://staking7-l2api.dev.z7a.xyz", "zilliqa api endpoint")
	rewardsCmd.Flags().StringVarP(&contractAddress, "contract address", "c", "343407558c9bb1f7ae737af80b90e1edf741a37a", "taking contract address")
	rewardsCmd.Flags().StringVarP(&ssn, "ssn", "s", "", "ssn operator address")

	StakingCmd.AddCommand(rewardsCmd)
}

var rewardsCmd = &cobra.Command{
	Use:   "rewards",
	Short: "Get rewards for specific ssn operator",
	Long:  "Get rewards for specific ssn operator",
	Run: func(cmd *cobra.Command, args []string) {
		p := provider.NewProvider(api)
		response, err := p.GetSmartContractSubState(contractAddress, "ssnlist", []string{ssn})
		if err != nil {
			panic(err)
		}
		var rep Rep
		err1 := json.Unmarshal([]byte(response), &rep)
		if err1 != nil {
			panic(err1)
		}
		s := rep.Result.SSNList[ssn]
		rewards := s.Arguments[2]
		fmt.Println(rewards)
	},
}

type Rep struct {
	Id      string `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}

type Result struct {
	SSNList map[string]SSN
}

type SSN struct {
	ArgTypes    []interface{} `json:"argtypes"`
	Arguments   []interface{} `json:"arguments"`
	Constructor string        `json:"constructor"`
}
