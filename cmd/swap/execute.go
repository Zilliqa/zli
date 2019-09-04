package swap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/FireStack-Lab/LaksaGo"
	"github.com/FireStack-Lab/LaksaGo/account"
	"github.com/FireStack-Lab/LaksaGo/bech32"
	contract2 "github.com/FireStack-Lab/LaksaGo/contract"
	"github.com/FireStack-Lab/LaksaGo/provider"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"zli/core"
)

var executeWallet *core.Wallet
var executeKeyStore string
var executeCSV string

func init() {
	ExecuteCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	ExecuteCmd.Flags().IntVarP(&chainId, "chainId", "c", 333, "the message version of the network")
	ExecuteCmd.Flags().StringVarP(&walletAddress, "address", "a", "zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2", "address of the fundWallet contract")
	ExecuteCmd.Flags().StringVarP(&gasPrice, "price", "p", "1000000000", "gas price")
	ExecuteCmd.Flags().StringVarP(&gasLimit, "limit", "l", "10000", "gas limit")
	ExecuteCmd.Flags().StringVarP(&amount, "amount", "m", "0", "token amount will be transfer to the smart contract")
	ExecuteCmd.Flags().StringVarP(&executeKeyStore, "executekeystore", "w", "", "execute key store")
	ExecuteCmd.Flags().StringVarP(&executeCSV, "signed", "r", "", "the path of signed file")
	ExecuteCmd.Flags().BoolVarP(&priority, "priority", "g", true, "setup priority of transaction")
	SwapCmd.AddCommand(ExecuteCmd)
}

type Pair struct {
	ID     string
	ToAddr string
}

var ExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute transactions",
	Long:  "execute transactions",
	PreRun: func(cmd *cobra.Command, args []string) {
		logfile, _ := os.Create("execute.log")
		log.SetOutput(logfile)
		if executeKeyStore == "" {
			panic("invalid execute keystore or password")
		}
		fmt.Println("please type password to decrypt your keystore: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			panic(err.Error())
		}
		executePrivateKey, err := core.LoadPirvateKeyFromKeyStore(executeKeyStore, string(pass))
		if err != nil {
			panic(err.Error())
		}
		siw, err := core.NewWallet(LaksaGo.DecodeHex(executePrivateKey), chainId, api)
		if err != nil {
			panic("construct exit wallet error: " + err.Error())
		}
		executeWallet = siw
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to read execute csv file...")
		f, err := os.Open(executeCSV)
		if err != nil {
			panic("cannot read execute csv file = " + err.Error())
		}

		scanner := bufio.NewScanner(f)
		var shouldBeExecuted []Pair

		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("read line: " + line)
			fields := strings.Split(line, " ")
			if len(fields) != 2 {
				fmt.Println("the fields of this line is pretty error,please check")
				os.Exit(1)
			}

			p := Pair{
				ID:     fields[0],
				ToAddr: fields[1],
			}
			shouldBeExecuted = append(shouldBeExecuted, p)
		}

		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}

		signer := account.NewWallet()
		signer.AddByPrivateKey(executeWallet.DefaultAccount.PrivateKey)
		p := provider.NewProvider(api)
		contract := contract2.Contract{
			Address:  walletAddress,
			Singer:   signer,
			Provider: p,
		}

		for _, value := range shouldBeExecuted {
			bech32, _ := bech32.ToBech32Address(value.ToAddr)
			fmt.Printf("start to execute id = %s, toAddr = %s, bech32 address = %s\n", value.ID, value.ToAddr, bech32)
			fmt.Println("please type Y to confirm: ")
			var confirmed string
			_, err = fmt.Scanln(&confirmed)
			if err != nil {
				fmt.Printf("confirm failed, skip execute tx %s\n", value.ID)
				continue
			}
			if confirmed != "Y" {
				fmt.Printf("confirm failed, skip execute tx %s\n", value.ID)
				continue
			}
			log.Printf("start to execute id = %s, toAddr = %s, bech32 address = %s\n", value.ID, value.ToAddr,bech32)
			result := p.GetBalance(executeWallet.DefaultAccount.Address)
			if result.Error != nil {
				panic(result.Error.Message)
			}
			balance := result.Result.(map[string]interface{})
			nonce, _ := balance["nonce"].(json.Number).Int64()
			params := contract2.CallParams{
				Version:      strconv.FormatInt(int64(LaksaGo.Pack(chainId, 1)), 10),
				Nonce:        strconv.FormatInt(nonce+1, 10),
				GasPrice:     gasPrice,
				GasLimit:     gasLimit,
				SenderPubKey: strings.ToUpper(executeWallet.DefaultAccount.PublicKey),
				Amount:       "0",
			}
			a := []contract2.Value{
				{
					VName: "transactionId",
					Type:  "Uint32",
					Value: value.ID,
				},
			}
			err, tx := contract.Call("ExecuteTransaction", a, params, priority, 1000, 3)
			if err != nil {
				log.Printf("execute transaction error %s, please check\n", err.Error())
				continue
			}

			log.Printf("start to poll execution transaction: %s\n", tx.ID)
			tx.Confirm(tx.ID, 1000, 3, p)
			err, recipients := getReceiptForTransaction(p, tx.ID)
			if err != nil {
				log.Printf("transaction failed")
				continue
			}
			log.Printf("get recipients for %s: %s\n", tx.ID, recipients)
			_ = core.AppendLine(bech32, "notoverride.csv")

		}

	},
}
