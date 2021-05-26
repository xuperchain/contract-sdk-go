
all: build

export GO111MODULE=on

unit-test:
	go test ./...

test:unit-test

lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
