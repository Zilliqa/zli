### ZLI

Zli is a command line tool based on zilliqa golang sdk

#### Requirements

golang environment: 


* [download golang](https://golang.org/dl/) (minimum version: go.1.12)
* [installation instructions](https://golang.org/doc/install)

#### Installation

<h5> dependencies </h5>

We use `go module` to manage dependencies, following command will download all dependencies according to the module file `go.mod`

```
go get ./...
```

<h5> build </h5>

In order to distinguish it from the already existing zilliqa-cli, we use the following command to generate a binary called zli-golang:

```go
go build -o go-zli main.go
```

<h5> install </h5>

Or you can install the binary:

```go
go build -o $GOPATH/bin/go-zli main.go
```

or just run the script:

```bash
sh install.sh
```

#### Commands

Currently, we provide four kinds of command, You can use `go-zli -h` to see all help messages:

```bash
A convenient command line tool to generate accounts, run integration testings or run http server .etc

Usage:
  zli [flags]
  zli [command]

Available Commands:
  account     Generate or load a large number of accounts
  contract    Deploy or call zilliqa smart contract
  help        Help about any command
  rpc         readonly json rpc of zilliqa
  spam        Send a large number of transactions
  transfer    Transfer zilliqa token to a specific account
  version     Print the version number of zli
  wallet      Init a new wallet or get exist wallet info

Flags:
  -h, --help   help for zli

Use "zli [command] --help" for more information about a command.

```

<h5> wallet </h5>

* go-zli wallet init: generate new wallet for zli to use, with random generated private key as default account, ca be modified later
* go-zli wallet echo: try to load wallet from file system, then print it
* **go-zli wallet from [flags]** : generate new wallet from specific private key


<h5> contract </h5>

* go-zli contract deploy [flags]: deploy new contract
* go-zli contract call [flags]: call a exist contract
* go-zli contract state [flags]: get state data for a specific smart contract

First, `contract` command will load config file from default configuration directory (~/.zilliqa)

<h5> account </h5>

* go-zli account generate [flags]: randomly generate some private keys

<h5> transfer </h5>

* go-zli transfer [flags]: transfer zilliqa token to a specific account

<h5> spam </h5>

* go-zli spam transfer [flags]: send a large number of transactions to a specific account
* go-zli spam invoke [flags]: invoke a large number of transactions on a exist smart contract

<h5> rpc </h5>

* go-zli rpc transaction [flags]: get transaction detail by transaction id

#### Example

<h5> internal tokenswap </h5>

1. Prepare two kinds of keystore, one for submitting transactions, another for signing and executing transactions:

submitters.keystore

```json
{"address": "4d7f134042b1423bd4940792ab6494730fbf95d4", "crypto": {"cipher": "aes-128-ctr", "cipherparams": {"iv": "0c1f5a0f83377e684fe64b732386f324"}, "ciphertext": "2b6ff8ec89e73827f7d1657db7ed3493b1cf4e2eda2b4254823161ec7e835655", "kdf": "pbkdf2", "kdfparams": {"salt": "df97994dc19b8a4af8ed37912c67d83ec04207fc05a59d13f5c90c09a4da3094", "n": 8192, "c": 262144, "r": 8, "p": 1, "dklen": 32}, "mac": "3c68fe95c047e8cf8fb7c63c9c54d7261f0bcc80fdc7e71f6770607cfd0cbe56"}, "id": "8a422ca1-96ca-4f04-97de-93b796698009", "version": 3}
```

signer.keystore

```json
{"address": "98b1a91648c1097e4ca54a82811109a2fc42cb92", "crypto": {"cipher": "aes-128-ctr", "cipherparams": {"iv": "1a336982dc294beb538ba8727002ffa8"}, "ciphertext": "73ef4879dd2da1c573ed83e0914329a91c8c417969656b9c47d8b6728f39235a", "kdf": "pbkdf2", "kdfparams": {"salt": "b7cbd1bc3c33afebba03875cd20d68e92449acec066c988874fde14c933564b7", "n": 8192, "c": 262144, "r": 8, "p": 1, "dklen": 32}, "mac": "79b528f7c231dfca8cd66293370b1cf3fa0697de479e387dc6d6ee746604e6ea"}, "id": "532b702e-64c9-421e-8f5a-cdbb1f25461c", "version": 3}
```


2. Prepare one recipient file:

```text
amrit zil1vuvhslc7qmt2cgyn25ssqlz6d2a2ee9d0ku2re 1
juzar zil1sk02c9s846qslj7nt32fw56e3nx78p2lxusrt7 1
```

3. Run submit command:

```bash
go-zli swap submit -a zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2 -u https://dev-api.zilliqa.com/  -c 333  -r recipients.csv -s submitter.keystore
```

this command will generate a `transaction file` which like (the third field means `zil` not `qa`):

```text
21 0x6719787f1e06d6ac20935521007c5a6abaace4ad 2
```



4. Feed the above `transaction file` to sign command:

```bash
go-zli swap sign -a zil1xpw4kwk25t622667zj2qq3nvtqv5u62l3xv6f2 -u https://dev-api.zilliqa.com/  -c 333  -r ./transactions.csv -w signer.keystore
```

<h5> test tiny contract for corner cases </h5>

1. Prepare ~/.zilliqa by using `go-zli wallet init` or `go-zli wallet from -p [private_key]`:

```json
{
	"api": "https://ipc-ud-api.dev.z7a.xyz",
	"chain_id": 2,
	"default_account": {
		"private_key": "227159779c78c9a920cba73086cf73fb3ee15cdd95380aa3b93757669e345300",
		"public_key": "0324cdd72db3de0e9f570d550631438d581056fb0d9c4daddbad2928eaf49f54ee",
		"address": "31f33d13ad6aa724cde1f3d12d75fb344a1df9de",
		"bech_32_address": "zil1x8en6yadd2njfn0p70gj6a0mx39pm7w7lz3kpm"
	},
	"accounts": [{
		"private_key": "227159779c78c9a920cba73086cf73fb3ee15cdd95380aa3b93757669e345300",
		"public_key": "0324cdd72db3de0e9f570d550631438d581056fb0d9c4daddbad2928eaf49f54ee",
		"address": "31f33d13ad6aa724cde1f3d12d75fb344a1df9de",
		"bech_32_address": "zil1x8en6yadd2njfn0p70gj6a0mx39pm7w7lz3kpm"
	}]
}
```

2. Deploy `tiny contract` using `sh scripts/deploy-tiny-contract.sh`

3. Run `go-zli testsuit tiny -a [contract_address]` like `go-zli testsuit tiny -a zil1yvnhvcage9w0yncuqj3wjp3vkg5qw5yuw4j6p5` or `sh scripts/test-tiny-contract.sh` to do the tests.

if the receipt of any transaction returns false, the whole tests will be stopped.

<h5> invoke contract </h5>

1. First, you have to generate a wallet configuration (~/.zilliqa) which contains a private key, go-zli will use this private key to sign
transactions, there are two ways to generate wallet file:

    * go-zli wallet init: randomly generate a private key with no balance
    * go-zli wallet from -p <private key>: using a exist private key (may have balance) to generate wallet file

like `go-zli wallet init` and `go-zli wallet from  -p  3B6674116AF2B954675E6373AC27E6A5CE03BCC8675ECDB7915AC8EE68B7ADCF`

2. Then, you can use following command to invoke a smart contract:

```bash
go-zli contract call -a <smart contract address> -t <transition name> -r <parameter>
```

for instance:

```bash
go-zli contract call -a 305d5b3acaa2f4a56b5e149400466c58194e695f -t SubmitTransaction -r "[{\"vname\":\"recipient\",\"type\":\"ByStr20\",\"value\":\"0x381f4008505e940ad7681ec3468a719060caf796\"},{\"vname\":\"amount\",\"type\":\"Uint128\",\"value\":\"10\"},{\"vname\":\"tag\",\"type\":\"String\",\"value\":\"a\"}]"
```

**warning**

Currently, `go-zli` now support pass private key as a parameter to `go-zli contract deploy or call` command, so, every time
you want to switch a different private key to send transactions, just use `-k private_key` option.

#### Run `go-zli` as docker container

If you do not want to run `go-zli` using the native binary, just build a docker image(or download from our repository later) 
then run it! But make sure you have installed docker environment correctly, if not, just refer `https://docs.docker.com/install/`

1. Build image:

```bash
sh build_docker_image.sh
```

2. Run inner a container environment

```bash
docker run --rm  -it go-zli bash
```
