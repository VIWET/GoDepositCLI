default: build

.PHONY: test
test: generate
	go test -tags '$(TAGS)' ./...

.PHONY: run
run: generate
	go run -tags '$(TAGS)' .

.PHONY: build
build: generate
	go build -tags '$(TAGS)' -o ./bin/staking-cli .

.PHONY: generate
generate:
	go generate -tags '$(TAGS)' ./...

.PHONY: clean
clean:
	rm -rf ./bin
