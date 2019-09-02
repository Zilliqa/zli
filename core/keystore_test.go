package core

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoadPirvateKeyFromKeyStore(t *testing.T) {
	privateKey, err := LoadPirvateKeyFromKeyStore(".keystore", "MyKey")
	if err != nil {
		t.Error(err.Error())
	}

	if strings.Compare(strings.ToUpper(privateKey), "3B6674116AF2B954675E6373AC27E6A5CE03BCC8675ECDB7915AC8EE68B7ADCF") != 0 {
		t.Error("decrypt private key failed")
	}

	fmt.Println("private key = " + privateKey)

}
