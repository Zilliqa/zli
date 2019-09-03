#!/usr/bin/env bash
 go-zli contract deploy -u https://dev-api.zilliqa.com/  -d 333 -c multisig.scilla -i init.json -l 100000 -p 10000000000 -s owner.keystore