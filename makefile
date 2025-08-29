# Makefile for mudskipper

.PHONY: pretty
pretty:
	stylua --glob '*' -- tools
	stylua --glob '*.lua' -- tests

.PHONY: test
test:
	nimble test tests/test_contains.nim

.PHONY: build
build:
	mkdir -p bin
	nim c -d:release -o:bin/contains src/contains.nim
