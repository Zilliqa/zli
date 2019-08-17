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

#### Commands

<h5> wallet </h5>

* init: generate new wallet for zli to use, with random generated private key as default account, ca be modified later