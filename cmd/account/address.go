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

package account

import (
	"fmt"
	bech322 "github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
)

var privateKey string

func init() {
	addressCmd.Flags().StringVarP(&privateKey, "key", "k", "", "private key")
	AccountCmd.AddCommand(addressCmd)
}

var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "Get address from a private key",
	Long:  "Get address from a private key",
	Run: func(cmd *cobra.Command, args []string) {
		addr := keytools.GetAddressFromPrivateKey(util.DecodeHex(privateKey))
		bech32, _ := bech322.ToBech32Address(addr)
		fmt.Println("0x" + addr)
		fmt.Println(bech32)
	},
}
