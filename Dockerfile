FROM golang:1.12.9 as build-stage
LABEL maintainer="Ren xiaohuo <lulu@zilliqa.com>"
WORKDIR /app
COPY ./ .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o zli main/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates bash
WORKDIR /app
COPY --from=build-stage /app/ .
RUN mv zli /bin/zli