BINARY := clickhouse-gen

.PHONY: linux
linux:
	mkdir -p build/linux
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/$(BINARY) *.go

.PHONY: darwin
darwin:
	mkdir -p build/osx
	GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/osx/$(BINARY) *.go

.PHONY: build
build:  linux darwin