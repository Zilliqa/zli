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
