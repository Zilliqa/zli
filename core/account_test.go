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

func TestSplit(t *testing.T) {
	a := Account{
		PrivateKey: "a",
	}

	b := Account{
		PrivateKey: "b",
	}

	c := Account{
		PrivateKey: "c",
	}

	d := Account{
		PrivateKey: "d",
	}

	e := Account{
		PrivateKey: "e",
	}

	accounts := Accounts{
		a, b, c, d, e,
	}

	accs := Split(accounts, 2)
	fmt.Println(accs)
	if len(accs) != 3 {
		panic("split error")
	}
}
