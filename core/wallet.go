/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
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
	return NewWallet(privateKey, chainId, api)
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
