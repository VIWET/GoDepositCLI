default: build

NETWORK=bahamut

.PHONY: test
test: generate
	go test -tags '$(NETWORK)' ./...

.PHONY: run
run: generate
	go run -tags '$(NETWORK)' .

.PHONY: build
build: generate
	go build -tags '$(NETWORK)' -o ./bin/staking-cli .

.PHONY: generate
generate:
	go generate -tags '$(NETWORK)' ./...

.PHONY: clean
clean:
	rm -rf ./bin
