
all: build

export GO111MODULE=on
install:
	go install github.com/xuperchain/xdev

unit-test:
	go test ./...

test:unit-test

lint:
	go vet ./...

coverage:g
	go test -coverprofile=coverage.txt -covermode=atomic ./...
