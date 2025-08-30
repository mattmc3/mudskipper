import std/[osproc, strutils, unittest]

# Helper to run the path command (compiled from src/path.nim with its when isMainModule block).
# Adjust PATH_CMD if your binary name/location differs.
const PATH_CMD = when defined(windows): "path.exe" else: "path"

# Simplify runPath: execCmdEx does not provide stderr; we only need stdout & exit code.
proc runPath(args: seq[string]): tuple[code:int, stdout:string] =
  let full = @["./" & PATH_CMD] & args
  let res = execCmdEx(full.join(" "))
  (res.exitCode, res.output)

suite "path basename":
  test "single path":
    let r = runPath(@["basename", "/foo/bar"])
    check r.code == 0
    check r.stdout.strip == "bar"
  test "preserve extension by default":
    let r = runPath(@["basename", "/tmp/example.txt"])
    check r.stdout.strip == "example.txt"
  test "strip extension with -E":
    let r = runPath(@["basename", "-E", "/tmp/example.txt"])
    check r.stdout.strip == "example"
  test "multiple paths":
    let r = runPath(@["basename", "/a/b.txt", "/c/d"])
    check r.code == 0
    let lines = r.stdout.strip.splitLines
    check lines.len == 2
    check lines[0] == "b.txt"
    check lines[1] == "d"
  test "hidden file without extra extension":
    let r = runPath(@["basename", ".bashrc"])
    check r.stdout.strip == ".bashrc" # extension not removed because leading dot only
  test "hidden file strip with -E (no change)":
    let r = runPath(@["basename", "-E", ".bashrc"])
    check r.stdout.strip == ".bashrc"

suite "path extension":
  test "simple extension":
    let r = runPath(@["extension", "file.txt"])
    check r.stdout.strip == ".txt"
  test "multi-dot takes last":
    let r = runPath(@["extension", "archive.tar.gz"])
    check r.stdout.strip == ".gz"
  test "no extension gives empty line":
    let r = runPath(@["extension", "Makefile"])
    check r.stdout.strip == ""
  test "hidden file without ext":
    let r = runPath(@["extension", ".bashrc"])
    check r.stdout.strip == ""

suite "path dirname":
  test "basic dirname":
    let r = runPath(@["dirname", "/foo/bar.txt"])
    check r.stdout.strip == "/foo"
  test "trailing slash collapses":
    let r = runPath(@["dirname", "/foo/bar/"])
    check r.stdout.strip == "/foo/bar"
  test "root path":
    let r = runPath(@["dirname", "/"])
    check r.stdout.strip == "/"

# NOTE: Additional fish path subcommand behaviors (filter, is, etc.) will be
# added once those subcommands are implemented in src/path.nim.
