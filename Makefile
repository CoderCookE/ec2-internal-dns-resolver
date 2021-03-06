APP_NAME=ec2-internal-dns-resolver
OS=linux
ARCH=amd64
PKG_NAME=$(APP_NAME)_$(shell cat VERSION)_$(ARCH)
RELEASE=$$(git rev-parse HEAD)

default: bin

bin:
	mkdir -p bin
	cd ./ && go build -o ./bin/$(APP_NAME)
	shasum -a 1 ./bin/ec2-internal-dns-resolver > ./bin/shasum

test:
	go test -race -v ./...

.PHONY: bin default
