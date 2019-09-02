package core

import (
	"github.com/FireStack-Lab/LaksaGo/crypto"
	"io/ioutil"
)

func LoadPirvateKeyFromKeyStore(path, password string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	ks := crypto.NewDefaultKeystore()
	private, err := ks.DecryptPrivateKey(string(b), password)
	if err != nil {
		return "", err
	}
	return private, nil

}
