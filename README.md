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

<h5> account </h5>

* go-zli account generate [flags]: randomly generate some private keys

<h5> transfer </h5>

* go-zli transfer [flags]: transfer zilliqa token to a specific account

<h5> spam </h5>

* go-zli spam transfer [flags]: send a large number of transactions to a specific account