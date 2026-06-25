BINARY := dep

.PHONY: build test clean build-all

build:
	go build -ldflags="-s -w" -o $(BINARY) .

test:
	go test ./... -count=1

build-all:
	mkdir -p dist && \
	for p in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do \
		GOOS=$${p%/*} GOARCH=$${p#*/} go build -ldflags="-s -w" -o dist/$(BINARY)-$${p%/*}-$${p#*/} .; \
	done

clean:
	rm -rf dist/ $(BINARY) $(BINARY).exe
