FROM golang:1.12.9

LABEL maintainer="Ren xiaohuo <lulu@zilliqa.com>"

WORKDIR /app

COPY ./ .

RUN go build  -o $GOPATH/bin/go-zli main/main.go