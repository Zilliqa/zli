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

var signWallet *core.Wallet
var signKeyStore string
var signCSV string

func init() {
	SignCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	SignCmd.Flags().IntVarP(&chainId, "chainId", "c", 333, "the message version of the network")
	SignCmd.Flags().StringVarP(&walletAddress, "address", "a", "zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2", "address of the fundWallet contract")
	SignCmd.Flags().StringVarP(&gasPrice, "price", "p", "1000000000", "gas price")
	SignCmd.Flags().StringVarP(&gasLimit, "limit", "l", "10000", "gas limit")
	SignCmd.Flags().StringVarP(&amount, "amount", "m", "0", "token amount will be transfer to the smart contract")
	SignCmd.Flags().StringVarP(&signKeyStore, "signkeystore", "w", "", "sign and execute key store")
	SignCmd.Flags().StringVarP(&signCSV, "recipient", "r", "", "the path of transaction file")
	SwapCmd.Flags().BoolVarP(&priority, "priority", "f", true, "setup priority of transaction")
	SwapCmd.AddCommand(SignCmd)
}

var SignCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign transactions",
	Long:  "sign transactions",
	PreRun: func(cmd *cobra.Command, args []string) {
		logfile, _ := os.Create("sign.log")
		log.SetOutput(logfile)
		if signKeyStore == "" {
			panic("invalid sign keystore or password")
		}

		fmt.Println("please type password to decrypt your keystore: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			panic(err.Error())
		}
		signPrivateKey, err := core.LoadPirvateKeyFromKeyStore(signKeyStore, string(pass))
		if err != nil {
			panic(err.Error())
		}
		siw, err := core.NewWallet(LaksaGo.DecodeHex(signPrivateKey), chainId, api)
		if err != nil {
			panic("construct sign wallet error: " + err.Error())
		}
		signWallet = siw
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start to read sign csv file...")
		fmt.Println(signCSV)
		f, err := os.Open(signCSV)
		if err != nil {
			panic("cannot read sign csv file = " + err.Error())
		}

		scanner := bufio.NewScanner(f)
		var shouldBeProcess []Txn
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("read line: " + line)
			fields := strings.Split(line, " ")
			if len(fields) != 3 {
				fmt.Println("the fields of this line is pretty error,please check")
				os.Exit(1)
			}
			tx := Txn{
				TxId:   fields[0],
				ToAddr: fields[1],
				Amount: fields[2],
			}
			shouldBeProcess = append(shouldBeProcess, tx)
		}

		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}

		signer := account.NewWallet()
		signer.AddByPrivateKey(signWallet.DefaultAccount.PrivateKey)
		p := provider.NewProvider(api)

		contract := contract2.Contract{
			Address:  walletAddress,
			Singer:   signer,
			Provider: p,
		}

		for _, value := range shouldBeProcess {
			fmt.Printf("transaction id = %s should be process on, toAddr = %s, amount = %s\n", value.TxId, value.ToAddr, value.Amount)
			fmt.Printf("start to sign id = %s, toAddr = %s, value = %s\n", value.TxId, value.ToAddr, value.Amount)
			fmt.Println("please type Y to confirm: ")
			var confirmed string
			_, err := fmt.Scanln(&confirmed)
			if err != nil {
				fmt.Printf("confirm failed, skip sign tx %s\n", value.TxId)
				continue
			}
			if confirmed != "Y" {
				fmt.Printf("confirm failed, skip sign tx %s\n", value.TxId)
				continue
			}
			log.Printf("start to sign id = %s, toAddr = %s, value = %s\n", value.TxId, value.ToAddr, value.Amount)
			result := p.GetBalance(signWallet.DefaultAccount.Address)
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
				SenderPubKey: strings.ToUpper(signWallet.DefaultAccount.PublicKey),
				Amount:       "0",
			}
			a := []contract2.Value{
				{
					VName: "transactionId",
					Type:  "Uint32",
					Value: value.TxId,
				},
			}
			err, tx := contract.Call("SignTransaction", a, params, false, 1000, 3)
			if err != nil {
				log.Printf("sign transaction error %s, please check\n", err.Error())
				continue
			}
			log.Printf("start to poll sign transaction: %s\n", tx.ID)
			tx.Confirm(tx.ID, 1000, 3, p)
			err, recipients := getReceiptForTransaction(p, tx.ID)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("get recipients for %s: %s\n", tx.ID, recipients)

			fmt.Printf("start to execute id = %s, toAddr = %s, value = %s\n", value.TxId, value.ToAddr, value.Amount)
			fmt.Println("please type Y to confirm: ")
			_, err = fmt.Scanln(&confirmed)
			if err != nil {
				fmt.Printf("confirm failed, skip execute tx %s\n", value.TxId)
				continue
			}
			if confirmed != "Y" {
				fmt.Printf("confirm failed, skip execute tx %s\n", value.TxId)
				continue
			}
			log.Printf("start to execute id = %s, toAddr = %s, value = %s\n", value.TxId, value.ToAddr, value.Amount)
			result = p.GetBalance(signWallet.DefaultAccount.Address)
			if result.Error != nil {
				panic(result.Error.Message)
			}
			balance = result.Result.(map[string]interface{})
			nonce, _ = balance["nonce"].(json.Number).Int64()
			params = contract2.CallParams{
				Version:      strconv.FormatInt(int64(LaksaGo.Pack(chainId, 1)), 10),
				Nonce:        strconv.FormatInt(nonce+1, 10),
				GasPrice:     gasPrice,
				GasLimit:     gasLimit,
				SenderPubKey: strings.ToUpper(signWallet.DefaultAccount.PublicKey),
				Amount:       "0",
			}
			a = []contract2.Value{
				{
					VName: "transactionId",
					Type:  "Uint32",
					Value: value.TxId,
				},
			}
			err, tx = contract.Call("ExecuteTransaction", a, params, false, 1000, 3)
			if err != nil {
				log.Printf("execute transaction error %s, please check\n", err.Error())
				continue
			}

			log.Printf("start to poll execution transaction: %s\n", tx.ID)
			tx.Confirm(tx.ID, 1000, 3, p)
			err, recipients = getReceiptForTransaction(p, tx.ID)
			if err != nil {
				log.Printf("transaction failed")
				continue
			}
			log.Printf("get recipients for %s: %s\n", tx.ID, recipients)

			bech32Address, _ := bech32.ToBech32Address(value.ToAddr)
			_ = core.AppendLine(bech32Address, "notoverride.csv")
		}
	},
}
