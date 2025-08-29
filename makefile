# Makefile for mudskipper

# Aggressive Nim optimization flags (tunable via environment)
# Removed obsolete -s (macOS ld warns); use dead_strip + post-build strip instead.
NIM_FLAGS ?= -d:release --mm:orc --opt:speed --passC="-march=native -flto -fomit-frame-pointer" --passL="-flto -dead_strip" -d:lto
SAFE_NIM_FLAGS ?= -d:release
NIM ?= nim
BIN_DIR := bin
PROGS := contains count

.PHONY: test build build-safe strip clean

test:
	nimble test tests/test_contains.nim

# Pattern rule for building each program with optimized flags
$(BIN_DIR)/%: src/%.nim
	@echo "=== Building $* ==="
	@mkdir -p $(BIN_DIR)
	@$(NIM) c $(NIM_FLAGS) -o:$@ $<

# Aggregate build (optimized)
build: $(PROGS:%=$(BIN_DIR)/%)

# Safe (baseline) build target using reduced flags (rebuilds everything)
build-safe: NIM_FLAGS := $(SAFE_NIM_FLAGS)
build-safe: clean build

# Strip just 'contains' (extend by adding others to the loop if desired)
strip: $(BIN_DIR)/contains
	strip -x $(BIN_DIR)/contains || true

clean:
	rm -rf $(BIN_DIR)
