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
