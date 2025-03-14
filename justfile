set shell := ["bash", "-cu"]

default: lint test

install-tools:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
    golangci-lint run --timeout 1m

test *ARGS="./...":
    go test -v -count=1 -race -covermode=atomic -coverprofile=coverage.out {{ARGS}}

# run `just test` first, then run this to view test coverage
cover:
    go tool cover -html coverage.out
