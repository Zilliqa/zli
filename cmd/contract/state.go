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
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/spf13/cobra"
)

func init() {
	stateCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	stateCmd.Flags().StringVarP(&invokeAddress, "address", "a", "", "smart contract address")
	ContractCmd.AddCommand(stateCmd)
}

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "get state data for specific smart contract",
	Long:  "get state data for specific smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		p := provider.NewProvider(api)
		contractAddr := invokeAddress
		if validator.IsBech32(invokeAddress) {
			a, err := bech32.FromBech32Addr(invokeAddress)
			if err != nil {
				panic("invalid address format: " + err.Error())
			}
			contractAddr = a
		}
		response := p.GetSmartContractState(contractAddr)
		if response == nil {
			fmt.Println("get response error")
			return
		}
		if response.Error != nil {
			fmt.Println("get response error = ", response.Error)
			return
		}

		result := response.Result
		data, err := json.Marshal(result)
		if err != nil {
			fmt.Println("get state data error = ", err.Error())
		} else {
			fmt.Println(string(data))

		}
	},
}
