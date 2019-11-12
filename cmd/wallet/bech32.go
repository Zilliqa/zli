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
package wallet

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/spf13/cobra"
)

var addr string

func init() {
	bech32Cmd.Flags().StringVarP(&addr,"address","a","","normal format or bech32 format address")
	WalletCmd.AddCommand(bech32Cmd)
}

var bech32Cmd = &cobra.Command{
	Use:   "bech32",
	Short: "translate normal address to bech32 format or vice versa",
	Long:  "translate normal address to bech32 format or vice versa",
	Run: func(cmd *cobra.Command, args []string) {
		if addr == "" {
			panic("address cannot been empty")
		}


		if validator.IsBech32(addr) {
			normal,err  := bech32.FromBech32Addr(addr)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(normal)
			}
		} else {
			b32,err := bech32.ToBech32Address(addr)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(b32)
			}
		}

	},
}
