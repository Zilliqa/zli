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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"zli/core"

	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
)

var number int64

func init() {
	generateCmd.Flags().Int64VarP(&number, "number", "n", 2, "the number of generated keys")
	AccountCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Randomly generate some private keys",
	Long:  "Randomly generate some private keys",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to generate ", number, " accounts")
		file, err := os.Create("./testAccounts.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		keys, err := core.GeneratePrivateKeys(number)
		if err != nil {
			panic(err)
		}

		accounts := []core.Account{}
		for _, key := range keys {
			account, err := core.NewAccount(util.EncodeHex(key[:]))
			if err != nil {
				panic(err.Error())
			}
			accounts = append(accounts, *account)
		}

		w := bufio.NewWriter(file)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(&accounts); err != nil {
			panic(err.Error())
		}

		err = w.Flush()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("end generate ", number, " accounts")

	},
}
