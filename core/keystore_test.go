//  This file is part of zli
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.
//
//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
