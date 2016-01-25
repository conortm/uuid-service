.PHONY: install test travis

install:
	go get -t -v ./...

test:
	go test

travis:
	$(HOME)/gopath/bin/goveralls -service=travis-ci
