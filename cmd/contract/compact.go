package contract

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var input string
var output string

func init() {
	compactCmd.Flags().StringVarP(&input, "input", "i", "", "the path of original smart contract to be compacted")
	compactCmd.Flags().StringVarP(&output, "output", "o", "", "the path of generated smart contract ")
	ContractCmd.AddCommand(compactCmd)
}

var compactCmd = &cobra.Command{
	Use:   "compact",
	Short: "Compact smart contract",
	Long:  "Compact smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		if input == "" {
			fmt.Println("input path cannot be empty")
			return
		}

		if output == "" {
			fmt.Println("output path cannot be empty")
			return
		}

		code, err := ioutil.ReadFile(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w, err1 := os.Create(output)
		if err1 != nil {
			fmt.Println(err1.Error())
			return
		}

		defer w.Close()

		// trim all comments
		c := string(code)
		m := regexp.MustCompile("\\(\\*(.*?)\\*\\)")
		newCode := m.ReplaceAllString(c, "")

		sb := strings.Builder{}
		m = regexp.MustCompile("^(\n)*$")
		scanner := bufio.NewScanner(strings.NewReader(newCode))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if !m.MatchString(line) {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}

		_ = ioutil.WriteFile(output, []byte(sb.String()), 0640)
	},
}
