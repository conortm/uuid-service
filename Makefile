package = github.com/conortm/uuid-service

.PHONY: install release test travis

install:
	go get -t -v ./...

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/uuid-service-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/uuid-service-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/uuid-service-linux-arm $(package)

test:
	go test -v

travis:
	$(HOME)/gopath/bin/goveralls -service=travis-ci
