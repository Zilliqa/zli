package testsuit

import "github.com/spf13/cobra"

var TestsultCmd = &cobra.Command{
	Use:   "testsuit",
	Short: "Automatically running tests of specific smart",
	Long: "Automatically running tests of specific smart",
	Run: func(cmd *cobra.Command, args []string) {
	},

}
