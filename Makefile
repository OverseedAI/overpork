.PHONY: build install clean test lint

VERSION ?= dev
LDFLAGS := -ldflags="-s -w -X github.com/OverseedAI/overpork/cmd.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o overpork .

install:
	go install $(LDFLAGS) .

clean:
	rm -f overpork
	rm -rf dist/

test:
	go test ./...

lint:
	golangci-lint run

# Build for all platforms
dist: clean
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/overpork-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/overpork-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/overpork-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/overpork-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/overpork-windows-amd64.exe .
