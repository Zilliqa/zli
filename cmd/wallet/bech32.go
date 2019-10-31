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
