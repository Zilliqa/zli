//  This file is part of zli
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

package wallet

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"os"
	"zli/core"
)

var DefaultConfigName = ".zilliqa"

func init() {
	WalletCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate new wallet for zli to use",
	Long:  "Generate new wallet for zli to use, with random generated private key as default account, ca be modified later",
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := core.DefaultWallet()
		if err != nil {
			panic(err)
		}

		walletJson, err1 := json.Marshal(wallet)
		if err1 != nil {
			panic(err1)
		}

		home := core.UserHomeDir()
		path := home + "/" + DefaultConfigName

		_, err2 := os.Stat(path)
		if err2 == nil {
			panic("file exist")
		}

		f, err3 := os.Create(path)
		if err3 != nil {
			panic(err3)
		}
		defer f.Close()

		_, err4 := f.Write(walletJson)

		if err4 != nil {
			panic(err4)
		}
	},
}
