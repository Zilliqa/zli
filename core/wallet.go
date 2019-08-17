package core

import (
	"github.com/FireStack-Lab/LaksaGo"
	bech322 "github.com/FireStack-Lab/LaksaGo/bech32"
	"github.com/FireStack-Lab/LaksaGo/keytools"
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