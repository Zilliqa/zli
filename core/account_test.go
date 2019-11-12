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
