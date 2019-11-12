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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"zli/core"
)

func init() {
	WalletCmd.AddCommand(echoCmd)
}

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Echo exist wallet info",
	Long:  "Try to load wallet file from file system, then print it",
	Run: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		wallet, err := core.LoadFromFile(home + "/" + DefaultConfigName)
		if err != nil {
			panic(err)
		}
		jsonString, _ := json.Marshal(wallet)
		fmt.Println(string(jsonString))
	},
}
