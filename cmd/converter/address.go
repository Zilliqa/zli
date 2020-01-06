package converter

import (
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/validator"
	"github.com/spf13/cobra"
	"strings"
)

var addr string

func init() {
	addressCmd.Flags().StringVarP(&addr, "addr", "a", "", "address to be converted")
	ConverterCmd.AddCommand(addressCmd)
}

var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "Convert the format of address",
	Long:  "Convert the format of address",
	Run: func(cmd *cobra.Command, args []string) {
		if addr == "" {
			return
		}

		if validator.IsBech32(addr) {
			checksum, err := bech32.FromBech32Addr(addr)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(checksum)
			}
			return
		}

		if !strings.Contains(addr, "0x") && !strings.Contains(addr, "0X") {
			addr = "0x" + addr
		}

		if validator.IsChecksumAddress(addr) {
			bech32, err := bech32.ToBech32Address(addr)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(bech32)
			}
			return
		}

		fmt.Println("not bech32 format or checksum format")
	},
}
