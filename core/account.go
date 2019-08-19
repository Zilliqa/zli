package core

import (
	"bufio"
	"github.com/FireStack-Lab/LaksaGo"
	bech322 "github.com/FireStack-Lab/LaksaGo/bech32"
	"github.com/FireStack-Lab/LaksaGo/keytools"
	"os"
	"strings"
)

type Accounts []Account


func Split(accounts Accounts, lim int) [][]Account {
	var chunk []Account
	chunks := make([][]Account,0,len(accounts)/lim+1)
	for len(accounts) >= lim {
		chunk,accounts = accounts[:lim],accounts[lim:]
		chunks = append(chunks,chunk)
	}

	if len(accounts) > 0 {
		chunks = append(chunks,accounts[:len(accounts)])
	}

	return chunks
}

func LoadFrom(path string) (Accounts, error) {
	var accounts Accounts
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		privates := strings.Split(line, " ")
		accs, err := fromPrivateKeys(privates)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, accs...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func GeneratePrivateKeys(number int64) ([]keytools.PrivateKey, error) {
	var i int64
	var keys []keytools.PrivateKey
	for i = 0; i < number; i++ {
		key, err := keytools.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func fromPrivateKeys(privates []string) ([]Account, error) {
	var accounts Accounts
	for _, v := range privates {
		private := LaksaGo.DecodeHex(v)
		publicKey := keytools.GetPublicKeyFromPrivateKey(private, true)
		public := LaksaGo.EncodeHex(publicKey)
		address := keytools.GetAddressFromPublic(publicKey)
		bech32, err := bech322.ToBech32Address(address)
		if err != nil {
			return nil, err
		}
		account := Account{
			PrivateKey:    v,
			PublicKey:     public,
			Address:       address,
			Bech32Address: bech32,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
