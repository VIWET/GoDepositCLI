default: build

NETWORK ?= bahamut

GIT_VERSION=$(shell git describe --tags --always)

LDFLAGS = -X github.com/viwet/GoDepositCLI/version.GitVersion=$(GIT_VERSION)

.PHONY: test
test: generate
	go test -tags '$(NETWORK)' ./...

.PHONY: run
run: generate
	go run -tags '$(NETWORK)' .

.PHONY: build
build: generate
	go build -ldflags='$(LDFLAGS)' -tags '$(NETWORK)' -o ./bin/staking-cli .
	@echo "\nRun staking-cli using './bin/staking-cli' command"

.PHONY: generate
generate:
	go generate -tags '$(NETWORK)' ./...

.PHONY: clean
clean:
	rm -rf ./bin
