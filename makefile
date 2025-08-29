# Makefile for mudskipper

# Aggressive Nim optimization flags (tunable via environment)
# Removed obsolete -s (macOS ld warns); use dead_strip + post-build strip instead.
NIM_FLAGS ?= -d:release --mm:orc --opt:speed --passC="-march=native -flto -fomit-frame-pointer" --passL="-flto -dead_strip" -d:lto

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
	nim c $(NIM_FLAGS) -o:bin/contains src/contains.nim

# Post-build strip to reduce size (macOS: -x strips local symbols)
.PHONY: strip
strip: build
	strip -x bin/contains || true

.PHONY: build-safe
build-safe:
	mkdir -p bin
	nim c -d:release -o:bin/contains src/contains.nim
