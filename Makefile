VERSION=$(shell git describe --tags)
GIT_REVISION=$(shell git rev-parse --short HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE=$$(date +%Y-%m-%d-%H:%M)
GO_FLAGS := -ldflags "-X main.BuildDate=$(BUILD_DATE) -X main.Branch=$(GIT_BRANCH) -X main.Revision=$(GIT_REVISION) -X main.Version=$(VERSION) -extldflags \"-static\" -s -w" -tags netgo

# Generate binaries for a powerproto release
.PHONY: build
build:
	rm -fr ./dist/
	mkdir -p ./dist
	GOOS="linux"  GOARCH="amd64" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-linux-amd64   ./cmd/powerproto
	GOOS="linux"  GOARCH="arm64" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-linux-arm64   ./cmd/powerproto
	GOOS="linux"  GOARCH="arm" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-linux-arm   ./cmd/powerproto
	GOOS="linux"  GOARCH="386" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-linux-x86   ./cmd/powerproto
	GOOS="windows"  GOARCH="386" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-windows-x86.exe   ./cmd/powerproto
	GOOS="windows"  GOARCH="amd64" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-windows-amd64.exe   ./cmd/powerproto
	GOOS="darwin"  GOARCH="amd64" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-darwin-amd64   ./cmd/powerproto
	GOOS="darwin"  GOARCH="arm64" CGO_ENABLED=0 go build $(GO_FLAGS) -o ./dist/powerproto-darwin-arm64   ./cmd/powerproto
