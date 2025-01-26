default: build

NETWORK ?= bahamut
GOLANG_CROSS_VERSION ?= v1.23

GIT_VERSION=$(shell git describe --tags --always)
ROOT_DIR=$(shell pwd)

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

.PHONY: release
release: check_github_token
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(HOME)/.docker/config.json:/root/.docker/config.json \
		-v $(ROOT_DIR):/go/src/staking-cli \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e NETWORK=$(NETWORK) \
		-w /go/src/staking-cli \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION}  release \
		--clean \

check_github_token:
	@[ "${GITHUB_TOKEN}" ] || ( echo "GITHUB_TOKEN wasn't provided"; exit 1 )
