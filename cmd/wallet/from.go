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
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/spf13/cobra"
	"os"
	"zli/core"
)

var private string

func init() {
	fromCmd.Flags().StringVarP(&private, "private", "p", "", "from specific private key")
	WalletCmd.AddCommand(fromCmd)
}

var fromCmd = &cobra.Command{
	Use:   "from [OPTIONS]",
	Short: "Generate new wallet from specific private key",
	Long:  "Generate new wallet from specific private key",
	Run: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		path := home + "/" + DefaultConfigName

		_, err := os.Stat(path)
		if err == nil {
			panic("file exist")
		}

		if !validator.IsPrivateKey(private) {
			panic("invalid private key")
		}

		wallet, err := core.FromPrivateKey(util.DecodeHex(private))
		if err != nil {
			panic(err)
		}

		walletJson, err := json.Marshal(wallet)

		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		_, err4 := f.Write(walletJson)

		if err4 != nil {
			panic(err4)
		}

	},
}
