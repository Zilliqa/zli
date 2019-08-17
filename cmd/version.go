package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"zli"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of zli",
	Long:  `All software has versions. This is zli's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(zli.CurrentVersionNumber)
	},
}
