default:
    @just --list

build:
    go build -o string ./cmd/string
    go build -o count ./cmd/count
    go build -o contains ./cmd/contains

run cmd *ARGS:
    go run ./cmd/{{cmd}} {{ARGS}}

test *ARGS:
    go test ./... {{ARGS}}

clean:
    rm -f string count contains
    rm -f coverage.out coverage.html

release version:
    #!/usr/bin/env bash
    set -euo pipefail
    if [[ ! "{{version}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: version must be X.Y.Z (got '{{version}}')"
        exit 1
    fi
    git tag v{{version}}
    echo "Run: git push origin --tags"

install:
    go install ./cmd/string ./cmd/count ./cmd/contains

tidy:
    go mod tidy

format:
    gofmt -w .

vet:
    go vet ./...

lint: tidy format vet
