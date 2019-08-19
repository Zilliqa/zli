package core

import (
	"fmt"
	"testing"
)

func TestLoadFrom(t *testing.T) {
	accounts, err := LoadFrom("../testsuit/accounts/testAccounts.txt")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(accounts)

	if len(accounts) != 20 {
		panic("load from file failed")
	}
}
