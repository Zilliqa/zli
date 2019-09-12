package account

import (
	"bufio"
	"fmt"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/spf13/cobra"
	"os"
	"zli/core"
)

var number int64

func init() {
	generateCmd.Flags().Int64VarP(&number, "number", "n", 2, "the number of generated keys")
	AccountCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Randomly generate some private keys",
	Long:  "Randomly generate some private keys",
	Run: func(cmd *cobra.Command, args []string) {
		if number%2 != 0 {
			panic("number should be even")
		}
		fmt.Println("start to generate ", number, " accounts")
		f, err := os.Create("./testAccounts.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		keys, err := core.GeneratePrivateKeys(number)
		if err != nil {
			panic(err)
		}

		i := 0
		w := bufio.NewWriter(f)

		for i+1 < len(keys) {
			k1 := keys[i]
			k2 := keys[i+1]
			line := fmt.Sprintf("%s %s", util.EncodeHex(k1[:]), util.EncodeHex(k2[:]))
			_, err := fmt.Fprintln(w, line)
			if err != nil {
				panic(err.Error())
			}
			i += 2
		}

		err = w.Flush()
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("end generate ", number, " accounts")

	},
}
