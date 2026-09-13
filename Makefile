VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build release clean test

build:
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o ansi-explain .

# release cross-compiles a static binary for each entry in PLATFORMS.
# CGO_ENABLED=0 guarantees a static binary even on platforms (linux) where
# cgo would otherwise pull in a dynamic link against libc.
release: clean
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		out=dist/ansi-explain-$(VERSION)-$$os-$$arch; \
		if [ "$$os" = "windows" ]; then out=$$out.exe; fi; \
		echo "building $$out"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags '$(LDFLAGS)' -o $$out . || exit 1; \
	done

test:
	go test ./...

clean:
	rm -rf dist ansi-explain
