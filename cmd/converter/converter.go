package converter

import "github.com/spf13/cobra"

var ConverterCmd = &cobra.Command{
	Use:   "converter",
	Short: "conversion tool",
	Long:  "conversion tool",
	Run: func(cmd *cobra.Command, args []string) {
	},
}