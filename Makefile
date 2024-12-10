default: build

NETWORK=bahamut

GIT_VERSION=$(shell git describe --tags | sed 's/[\.]/ /g' | tr -d 'v')

GIT_MAJOR=$(word 1, $(GIT_VERSION))
GIT_MINOR=$(word 2, $(GIT_VERSION))
GIT_PATCH=$(word 3, $(GIT_VERSION))

LDFLAGS = \
	-X 'github.com/viwet/GoDepositCLI/version.Major=$(GIT_MAJOR)' \
	-X 'github.com/viwet/GoDepositCLI/version.Minor=$(GIT_MINOR)' \
	-X 'github.com/viwet/GoDepositCLI/version.Patch=$(GIT_PATCH)'

.PHONY: test
test: generate
	go test -tags '$(NETWORK)' ./...

.PHONY: run
run: generate
	go run -tags '$(NETWORK)' .

.PHONY: build
build: generate
	go build -ldflags='$(LDFLAGS)' -tags '$(NETWORK)' -o ./bin/staking-cli .

.PHONY: generate
generate:
	go generate -tags '$(NETWORK)' ./...

.PHONY: clean
clean:
	rm -rf ./bin
