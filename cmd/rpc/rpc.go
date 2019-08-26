package rpc

import "github.com/spf13/cobra"

var RPCCmd = &cobra.Command{
	Use:   "rpc",
	Short: "readonly json rpc of zilliqa",
	Long:  "readonly json rpc of zilliqa",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
