### ZLI

Zli is a command line tool based on zilliqa golang sdk

#### Requirements

golang environment:

* [download golang](https://golang.org/dl/)
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

<h5> wallet </h5>

* go-zli wallet init: generate new wallet for zli to use, with random generated private key as default account, ca be modified later
* go-zli wallet echo: try to load wallet from file system, then print it
* go-zli wallet from [flags]: generate new wallet from specific private key

<h5> contract </h5>

* go-zli contract deploy [flags]: deploy new contract
* go-zli contract call [flags]: call a exist contract
* go-zli contract state [flags]: get state data for a specific smart contract

<h5> account </h5>

* go-zli account generate [flags]: randomly generate some private keys

<h5> transfer </h5>

* go-zli transfer [flags]: transfer zilliqa token to a specific account

<h5> spam </h5>

* go-zli spam transfer [flags]: send a large number of transactions to a specific account
* go-zli spam invoke [flags]: invoke a large number of transactions on a exist smart contract

#### Example

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

Currently, `go-zli` doesn't support pass private key as a parameter to `go-zli contract call` command (but will complete this feature soon),so, every time
you want to switch a different private key to send transactions, just delete `~/.zilliqa` then generate a new one!