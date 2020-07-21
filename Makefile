all: build

install:
	@true

depend:
	@go mod download

build: check
	@true

check:
	@go fmt *.go
	@go vet ./...
	@golint
	@go test -v -cover ./...

serve:
	@reflex -r '.*\.go' -s \
	-- sh -c \
	'make check'

.PHONY: all install depend build check serve
