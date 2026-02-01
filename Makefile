.PHONY: build install clean test lint

VERSION ?= dev
LDFLAGS := -ldflags="-s -w -X github.com/OverseedAI/overpork/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o opork .

install:
	go install $(LDFLAGS) .

clean:
	rm -f opork
	rm -rf dist/

test:
	go test ./...

lint:
	golangci-lint run

# Build for all platforms
dist: clean
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/opork-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/opork-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/opork-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/opork-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/opork-windows-amd64.exe .
