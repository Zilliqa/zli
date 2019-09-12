package core

import (
	"encoding/json"
	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"io/ioutil"
)

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

func NewAccount(privateKey string) (*Account, error) {
	publicKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	public := util.EncodeHex(publicKey)
	address := keytools.GetAddressFromPublic(publicKey)
	bech32, err := bech32.ToBech32Address(address)
	if err != nil {
		return nil, err
	}

	return &Account{
		PrivateKey:    privateKey,
		PublicKey:     public,
		Address:       address,
		Bech32Address: bech32,
	}, nil
}

func NewWallet(privateKey []byte, chainId int, api string) (*Wallet, error) {

	defaultAccount, err := NewAccount(util.EncodeHex(privateKey))
	if err != nil {
		return nil, err
	}

	accounts := []Account{*defaultAccount}

	return &Wallet{
		API:            api,
		ChainID:        chainId,
		DefaultAccount: *defaultAccount,
		Accounts:       accounts,
	}, nil
}

func FromPrivateKeyAndChain(privateKey []byte, chainId int, api string) (*Wallet, error) {
	return NewWallet(privateKey, chainId, api);
}
func FromPrivateKey(privateKey []byte) (*Wallet, error) {
	return FromPrivateKeyAndChain(privateKey, 333, "https://dev-api.zilliqa.com/")
}

func DefaultWallet() (*Wallet, error) {
	privateKey, err := keytools.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return FromPrivateKey(privateKey[:])
}

func LoadFromFile(file string) (*Wallet, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{}
	err1 := json.Unmarshal(b, wallet)

	if err1 != nil {
		return nil, err1
	}

	return wallet, nil

}
