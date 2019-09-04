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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
	"zli/core"
)

var submitWallet *core.Wallet
var submitKeyStore string
var submitCSV string
var notoverride string

type Recipient struct {
	Name    string
	Address string
	Amount  uint64
}

func init() {
	SubmitCmd.Flags().StringVarP(&api, "api", "u", "https://dev-api.zilliqa.com/", "api url")
	SubmitCmd.Flags().IntVarP(&chainId, "chainId", "c", 333, "the message version of the network")
	SubmitCmd.Flags().StringVarP(&walletAddress, "address", "a", "zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2", "address of the fundWallet contract")
	SubmitCmd.Flags().StringVarP(&gasPrice, "price", "p", "1000000000", "gas price")
	SubmitCmd.Flags().StringVarP(&gasLimit, "limit", "l", "10000", "gas limit")
	SubmitCmd.Flags().StringVarP(&amount, "amount", "m", "0", "token amount will be transfer to the smart contract")
	SubmitCmd.Flags().StringVarP(&submitKeyStore, "submitkeystore", "s", "", "submit key store")
	SubmitCmd.Flags().StringVarP(&submitCSV, "recipient", "r", ".recipients", "the path of recipient file")
	SubmitCmd.Flags().StringVarP(&notoverride, "notoverride", "n", "notoverride.csv", "not override")
	SubmitCmd.Flags().BoolVarP(&priority, "priority", "f", true, "setup priority of transaction")
	SwapCmd.AddCommand(SubmitCmd)
}

var SubmitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submit transactions",
	Long:  "submit transactions",
	PreRun: func(cmd *cobra.Command, args []string) {
		logfile, _ := os.Create("submit.log")
		log.SetOutput(logfile)
		if submitKeyStore == "" {
			panic("invalid submit keystore or password")
		}
		fmt.Println("please type password to decrypt your keystore: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			panic(err.Error())
		}
		submitPrivateKey, err := core.LoadPirvateKeyFromKeyStore(submitKeyStore, string(pass))
		if err != nil {
			panic("load submit private key error: " + err.Error())
		}
		sw, err := core.NewWallet(LaksaGo.DecodeHex(submitPrivateKey), chainId, api)
		if err != nil {
			panic("construct submit wallet error: " + err.Error())
		}
		submitWallet = sw

	},
	Run: func(cmd *cobra.Command, args []string) {
		// the main process of internal token swap
		fmt.Println("start to read submit csv file...")
		f, err := os.Open(submitCSV)
		if err != nil {
			panic("cannot read submit csv file = " + err.Error())
		}

		var recipients []Recipient

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("read line: " + line)
			fields := strings.Split(line, " ")
			if len(fields) != 3 {
				fmt.Println("the fields of this line is pretty error,please check")
				os.Exit(1)
			}

			zil, err := strconv.ParseFloat(fields[2], 64)
			qa := zil * 1000000000000
			if err != nil {
				panic("parse amount error")
			}
			r := Recipient{
				Name:    fields[0],
				Address: fields[1],
				Amount:  uint64(qa),
			}
			recipients = append(recipients, r)
		}

		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}
		fmt.Printf("total load %d recipients\n", len(recipients))

		fmt.Printf("start to read not override csv...")

		notoverrides := make(map[string]interface{})
		f, err = os.Open(notoverride)
		if err != nil {
			panic("cannot read notoverride csv file = " + err.Error())
		}
		scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println("read line: " + line)
			notoverrides[line] = nil
		}
		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}
		fmt.Printf("tital load %d notoverrides\n", len(notoverrides))

		normalWalletAddress, _ := bech32.FromBech32Addr(walletAddress)
		p := provider.NewProvider(api)


		for _, value := range recipients {
			_, ok := notoverrides[value.Address]
			if ok {
				fmt.Printf("skip %s, address = %s\n", value.Name, value.Address)
				continue
			}
			checksumAddress, _ := bech32.FromBech32Addr(value.Address)
			toAddress := fmt.Sprintf("0x%s", strings.ToLower(checksumAddress))
			fmt.Printf("start to submit transaction for %s, bech32 address is %s, normal address is %s, amount is %d\n", value.Name, value.Address, toAddress, value.Amount)
			fmt.Println("please type Y to confirm: ")
			var confirmed string
			_, err := fmt.Scanln(&confirmed)
			if err != nil {
				fmt.Printf("confirm failed, skip send to %s\n", value.Name)
				continue
			}
			if confirmed != "Y" {
				fmt.Printf("skip send to %s\n", value.Name)
				continue
			}
			response := p.GetSmartContractState(normalWalletAddress)
			if response.Error != nil {
				panic(response.Error.Message)
			}
			states := response.Result.([]interface{})
			transactionsBeforeSubmit := parseUnsignedTransactionsFromState(states)
			log.Printf("start to submit transaction for %s, bech32 address is %s, normal address is %s, amount is %d", value.Name, value.Address, toAddress, value.Amount)
			a := []contract2.Value{
				{
					VName: "recipient",
					Type:  "ByStr20",
					Value: toAddress,
				},
				{
					VName: "amount",
					Type:  "Uint128",
					Value: strconv.FormatUint(value.Amount, 10),
				},
				{
					VName: "tag",
					Type:  "String",
					Value: "",
				},
			}
			p := provider.NewProvider(api)
			result := p.GetBalance(submitWallet.DefaultAccount.Address)
			if result.Error != nil {
				panic(result.Error.Message)
			}
			balance := result.Result.(map[string]interface{})
			nonce, _ := balance["nonce"].(json.Number).Int64()
			signer := account.NewWallet()
			signer.AddByPrivateKey(submitWallet.DefaultAccount.PrivateKey)
			contract := contract2.Contract{
				Address:  walletAddress,
				Singer:   signer,
				Provider: p,
			}
			params := contract2.CallParams{
				Version:      strconv.FormatInt(int64(LaksaGo.Pack(chainId, 1)), 10),
				Nonce:        strconv.FormatInt(nonce+1, 10),
				GasPrice:     gasPrice,
				GasLimit:     gasLimit,
				SenderPubKey: strings.ToUpper(submitWallet.DefaultAccount.PublicKey),
				Amount:       "0",
			}

			err, tx := contract.Call("SubmitTransaction", a, params, false, 1000, 3)

			if err != nil {
				panic(err.Error())
			}

			log.Printf("start to poll transaction: %s\n", tx.ID)
			tx.Confirm(tx.ID, 1000, 3, p)
			err, recipients := getReceiptForTransaction(p, tx.ID)
			if err != nil {
				panic(err.Error())
			}
			log.Printf("get recipients for %s", tx.ID)
			log.Printf("recipients:\n%s\n", recipients)
			response = p.GetSmartContractState(normalWalletAddress)
			if response.Error != nil {
				panic(response.Error.Message)
			}
			states = response.Result.([]interface{})
			transactionsAfterSubmit := parseUnsignedTransactionsFromState(states)
			//so we compare transactionsBeforeSubmit and transactionsAfterSubmit to get transactions should be process this time
			transactions := compareTransactions(transactionsBeforeSubmit, transactionsAfterSubmit)
			err = core.WriteLines(transactions, "transactions.csv")
			if err != nil {
				fmt.Println("write transactions to file failed: ", err.Error())
				fmt.Println(transactions)
			}
		}
	},
}

type Txn struct {
	TxId   string
	ToAddr string
	Amount string
}

func getReceiptForTransaction(provider2 *provider.Provider, transactionId string) (error, string) {
	response := provider2.GetTransaction(transactionId)
	if response.Error != nil {
		return errors.New(response.Error.Message), ""
	}
	result := response.Result.(map[string]interface{})
	receipt := result["receipt"]
	receiptMap := receipt.(map[string]interface{})
	success := receiptMap["success"].(bool)
	if success == false {
		return errors.New("receipt failure"), ""
	}
	b, err := json.Marshal(receipt)
	if err != nil {
		return errors.New(err.Error()), ""
	}

	return nil, string(b)

}

func parseUnsignedTransactionsFromState(states []interface{}) map[string]Txn {
	transactions := make(map[string]Txn)
	for _, state := range states {
		stateMap := state.(map[string]interface{})
		vname := stateMap["vname"].(string)
		if vname == "transactions" {
			maps := stateMap["value"].([]interface{})
			for _, m := range maps {
				transactionMap := m.(map[string]interface{})
				k := transactionMap["key"]
				v := transactionMap["val"].(map[string]interface{})
				arguments := v["arguments"]
				args := arguments.([]interface{})
				tx := Txn{
					TxId:   k.(string),
					ToAddr: args[0].(string),
					Amount: args[1].(string),
				}
				transactions[tx.TxId] = tx
			}
		}
	}
	return transactions
}

func compareTransactions(before, after map[string]Txn) []string {
	var remain []string
	for key, value := range after {
		_, ok := before[key]
		if !ok {
			sb := strings.Builder{}
			sb.WriteString(key)
			sb.WriteString(" ")
			sb.WriteString(value.ToAddr)
			sb.WriteString(" ")
			sb.WriteString(value.Amount)
			remain = append(remain, sb.String())
		}
	}
	return remain
}
