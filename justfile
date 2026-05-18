default:
    @just --list

build:
    mkdir -p dist
    go build -o dist/string ./cmd/string
    go build -o dist/count ./cmd/count
    go build -o dist/contains ./cmd/contains
    go build -o dist/path ./cmd/path

run cmd *ARGS:
    go run ./cmd/{{cmd}} {{ARGS}}

test *ARGS:
    go test ./... {{ARGS}}

clean:
    rm -rf dist
    rm -f coverage.out coverage.html

release version:
    #!/usr/dist/env bash
    set -euo pipefail
    if [[ ! "{{version}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        echo "Error: version must be X.Y.Z (got '{{version}}')"
        exit 1
    fi
    git tag v{{version}}
    echo "Run: git push origin --tags"

install:
    go install ./cmd/string ./cmd/count ./cmd/contains ./cmd/path

tidy:
    go mod tidy

format:
    gofmt -w .

vet:
    go vet ./...

lint: tidy format vet
