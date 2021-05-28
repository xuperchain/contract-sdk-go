
export GO111MODULE=on

unit-test:
	go test ./...

example-test:
	make -C example build
	make -C example test

test:unit-test

lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
