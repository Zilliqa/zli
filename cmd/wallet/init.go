package wallet

import (
	"encoding/json"
	"github.com/FireStack-Lab/LaksaGo"
	bech322 "github.com/FireStack-Lab/LaksaGo/bech32"
	"github.com/FireStack-Lab/LaksaGo/keytools"
	"github.com/spf13/cobra"
	"os"
	"zli/core"
)

var defaultConfigName = ".zilliqa"

func init() {
	WalletCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate new wallet for zli to use",
	Long:  "Generate new wallet for zli to use, with random generated private key as default account, ca be modified later",
	Run: func(cmd *cobra.Command, args []string) {
		wallet, err := DefaultWallet()
		if err != nil {
			panic(err)
		}

		walletJson, err1 := json.Marshal(wallet)
		if err1 != nil {
			panic(err1)
		}

		home := core.UserHomeDir()
		f, err3 := os.Create(home + "/" + defaultConfigName)
		if err3 != nil {
			panic(err3)
		}
		defer f.Close()

		_, err4 := f.Write(walletJson)

		if err4 != nil {
			panic(err4)
		}
	},
}

type Wallet struct {
	API            string    `json:"api"`
	ChainID        int       `json:"chain_id"`
	DefaultAccount Account   `json:"default_account"`
	Accounts       []Account `json:"accounts"`
}

type Account struct {
	PrivateKey    string `json:"private_key"`
	PublicKey     string `json:"public_key"`
	Address       string `json:"address"`
	Bech32Address string `json:"bech_32_address"`
}

func DefaultWallet() (*Wallet, error) {
	privateKey, err := keytools.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	private := LaksaGo.EncodeHex(privateKey[:])
	publicKey := keytools.GetPublicKeyFromPrivateKey(privateKey[:], true)
	public := LaksaGo.EncodeHex(publicKey)
	address := keytools.GetAddressFromPublic(publicKey)
	bech32, err2 := bech322.ToBech32Address(address)
	if err2 != nil {
		return nil, err2
	}

	defaultAccount := Account{
		PrivateKey:    private,
		PublicKey:     public,
		Address:       address,
		Bech32Address: bech32,
	}

	accounts := []Account{defaultAccount}

	return &Wallet{
		API:            "https://dev-api.zilliqa.com/",
		ChainID:        333,
		DefaultAccount: defaultAccount,
		Accounts:       accounts,
	}, nil
}
