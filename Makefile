.PHONY: install test travis

install:
	go get -t -v ./...

test:
	go test

travis:
	go test -v ./...
