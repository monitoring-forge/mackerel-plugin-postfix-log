VERSION=0.0.10
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
GO111MODULE=on

all: mackerel-plugin-postfix-log

.PHONY: mackerel-plugin-postfix-log check lint linux

mackerel-plugin-postfix-log: main.go postfixlog/*.go
	go build $(LDFLAGS) -o mackerel-plugin-postfix-log

linux: main.go postfixlog/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-postfix-log

check:
	go test -v ./...

lint:
	golangci-lint run ./...