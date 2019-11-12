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
package rpc

import "github.com/spf13/cobra"

var RPCCmd = &cobra.Command{
	Use:   "rpc",
	Short: "readonly json rpc of zilliqa",
	Long:  "readonly json rpc of zilliqa",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
