# Makefile for mudskipper

.PHONY: pretty
pretty:
	stylua --glob '*' -- tools
	stylua --glob '*.lua' -- tests

.PHONY: test
test:
	./tests/test.lua .
