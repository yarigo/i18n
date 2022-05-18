all: check test

.PHONY: deend
depend:
	@go mod download

.PHONY: check
check:
	@gofmt -w ./
	@go vet ./...
	@golint
	@staticcheck ./...

.PHONY: test
test:
	@go test -failfast -v -race -cover ./...

.PHONY: serve
serve:
	@reflex -r '.*\.go' -s \
	-- sh -c \
	'make check'

.DEFAULT_GOAL := all
