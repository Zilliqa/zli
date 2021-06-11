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

package keystore

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/spf13/cobra"
	"strings"
)

var privateKey string
var password string

func init() {
	toCmd.Flags().StringVarP(&privateKey, "privateKey", "k", "", "private key to be converted")
	toCmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt the json keystore")
	KeystoreCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:   "to",
	Short: "convert private key to keystore",
	Long:  "convert private key to keystore",
	Run: func(cmd *cobra.Command, args []string) {
		privateKey = strings.TrimPrefix(privateKey,"0x")
		keystore,err := account.ToFile(privateKey,password,0)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(keystore)
	},
}
