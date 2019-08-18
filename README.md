### ZLI

Zli is a command line tool based on zilliqa golang sdk

#### Requirements

golang environment:

* [download golang](https://golang.org/dl/)
* [installation instructions](https://golang.org/doc/install)

#### Install

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
* go-zli wallet echo: try to load wallet file from file system, then print it
* go-zli wallet from [flags]: generate new wallet from specific private key

<h5> contract </h5>

* go-zli contract deploy [flags]: deploy new contract
* go-zli contract call [flags]: call a exist contract