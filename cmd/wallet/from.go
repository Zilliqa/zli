package wallet

import (
	"encoding/json"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/validator"
	"github.com/spf13/cobra"
	"os"
	"zli/core"
)

func init() {
	fromCmd.Flags().StringP("from", "f", "", "from specific private key")
	WalletCmd.AddCommand(fromCmd)
}

var fromCmd = &cobra.Command{
	Use:   "from",
	Short: "Generate new wallet from specific private key",
	Long:  "Generate new wallet from specific private key",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		home := core.UserHomeDir()
		path := home + "/" + defaultConfigName

		_, err := os.Stat(path)
		if err == nil {
			panic("file exist")
		}

		private := args[0]
		if !validator.IsPrivateKey(private) {
			panic("invalid private key")
		}

		wallet, err := core.FromPrivateKey(LaksaGo.DecodeHex(private))
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
