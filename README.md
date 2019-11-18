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

we use the following command to generate a binary called zli:

```go
go build -o zli main.go
```

<h5> install </h5>

Or you can install the binary:

```go
go build -o $GOPATH/bin/zli main.go
```

or just run the script:

```bash
sh install.sh
```

#### Commands

Currently, we provide four kinds of command, You can use `zli -h` to see all help messages:

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
  transfer    Transfer zilliqa token to a specific account
  version     Print the version number of zli
  wallet      Init a new wallet or get exist wallet info

Flags:
  -h, --help   help for zli

Use "zli [command] --help" for more information about a command.

```

<h5> wallet </h5>

* zli wallet init: generate new wallet for zli to use, with random generated private key as default account, ca be modified later
* zli wallet echo: try to load wallet from file system, then print it
* **zli wallet from [flags]** : generate new wallet from specific private key


<h5> contract </h5>

* zli contract deploy [flags]: deploy new contract
* zli contract call [flags]: call a exist contract
* zli contract state [flags]: get state data for a specific smart contract

First, `contract` command will load config file from default configuration directory (~/.zilliqa)

<h5> account </h5>

* zli account generate [flags]: randomly generate some private keys

<h5> transfer </h5>

* zli transfer [flags]: transfer zilliqa token to a specific account

<h5> rpc </h5>

* zli rpc transaction [flags]: get transaction detail by transaction id

#### Example

<h5> test tiny contract for corner cases </h5>

1. Prepare ~/.zilliqa by using `zli wallet init` or `zli wallet from -p [private_key]`:

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

3. Run `zli testsuit tiny -a [contract_address]` like `zli testsuit tiny -a zil1yvnhvcage9w0yncuqj3wjp3vkg5qw5yuw4j6p5` or `sh scripts/test-tiny-contract.sh` to do the tests.

if the receipt of any transaction returns false, the whole tests will be stopped.

<h5> invoke contract </h5>

1. First, you have to generate a wallet configuration (~/.zilliqa) which contains a private key, zli will use this private key to sign
transactions, there are two ways to generate wallet file:

    * zli wallet init: randomly generate a private key with no balance
    * zli wallet from -p <private key>: using a exist private key (may have balance) to generate wallet file

like `zli wallet init` and `zli wallet from  -p  3B6674116AF2B954675E6373AC27E6A5CE03BCC8675ECDB7915AC8EE68B7ADCF`

2. Then, you can use following command to invoke a smart contract:

```bash
zli contract call -a <smart contract address> -t <transition name> -r <parameter>
```

for instance:

```bash
zli contract call -a 305d5b3acaa2f4a56b5e149400466c58194e695f -t SubmitTransaction -r "[{\"vname\":\"recipient\",\"type\":\"ByStr20\",\"value\":\"0x381f4008505e940ad7681ec3468a719060caf796\"},{\"vname\":\"amount\",\"type\":\"Uint128\",\"value\":\"10\"},{\"vname\":\"tag\",\"type\":\"String\",\"value\":\"a\"}]"
```

**warning**

Currently, `zli` now support pass private key as a parameter to `zli contract deploy or call` command, so, every time
you want to switch a different private key to send transactions, just use `-k private_key` option.

#### Run `zli` inside a docker container

If you do not want to run `zli` using the native binary, just build a docker image(or download from our repository later) 
then run it! But make sure you have installed docker environment correctly, if not, just refer `https://docs.docker.com/install/`

1. Build image:

```bash
sh build_docker_image.sh
```

2. Run inner a container environment

```bash
docker run --rm  -it -v ~/contract:/contract docker.pkg.github.com/zilliqa/zli/zli bash
```

