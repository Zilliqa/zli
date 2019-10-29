#!/usr/bin/env bash
 zli contract deploy -u https://dev-api.zilliqa.com/  -d 333 -c multisig.scilla -i init.json -l 100000 -p 1000000000000 -s owner.keystore
# zli contract deploy -u https://lhc-9-api.dev.z7a.xyz  -d 2 -c multisig.scilla -i init.json -l 100000 -p 1000000000000 -s o.keystore