import unittest
import ../src/contains as mudskipper_contains

type CaptureBuffers = ref object
  outBuf: string
  errBuf: string

proc reset(self: CaptureBuffers) =
  self.outBuf.setLen(0)
  self.errBuf.setLen(0)

proc newContainsCmd(bufs: CaptureBuffers): mudskipper_contains.ContainsCommand =
  result.outWriter = proc (s: string) = bufs.outBuf.add(s)
  result.errWriter = proc (s: string) = bufs.errBuf.add(s)

suite "contains":
  test "needle found returns 0 and not found returns 1":
    check mudskipper_contains.contains(@["foo"]) == 1
    check mudskipper_contains.contains(@["foo", "foo"]) == 0
    check mudskipper_contains.contains(@["bar", "foo", "bar", "baz"]) == 0
    check mudskipper_contains.contains(@["missing", "foo", "bar", "baz"]) == 1

  test "index flag prints first occurrence (1-based)":
    let bufs = CaptureBuffers()
    var cmd = newContainsCmd(bufs)
    var code = cmd.contains(@["-i", "bar", "foo", "bar", "baz"])
    check code == 0
    check bufs.outBuf == "2\n"
    bufs.reset()
    code = cmd.contains(@["-i", "foo", "foo", "bar"])
    check code == 0
    check bufs.outBuf == "1\n"
    bufs.reset()
    code = cmd.contains(@["-i", "baz", "foo", "bar", "baz"])
    check code == 0
    check bufs.outBuf == "3\n"

  test "index flag no output when not found":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    let code = cmd.contains(@["-i", "nope", "foo", "bar"])
    check code == 1
    check bufs.outBuf == ""

  test "index0 flag prints first occurrence (0-based)":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    var code = cmd.contains(@["-0", "bar", "foo", "bar", "baz"])
    check code == 0
    check bufs.outBuf == "1\n"
    bufs.reset()
    code = cmd.contains(@["-0", "foo", "foo", "bar"])
    check code == 0
    check bufs.outBuf == "0\n"
    bufs.reset()
    code = cmd.contains(@["-0", "baz", "foo", "bar", "baz"])
    check code == 0
    check bufs.outBuf == "2\n"

  test "index0 flag no output when not found":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    let code = cmd.contains(@["-0", "nope", "foo", "bar"])
    check code == 1
    check bufs.outBuf == ""

  test "help returns 0":
    check mudskipper_contains.contains(@["-h"]) == 0
    check mudskipper_contains.contains(@["--help"]) == 0
    check mudskipper_contains.contains(@[]) == 0

  test "key not specified returns 2":
    check mudskipper_contains.contains(@["-i"]) == 2
    check mudskipper_contains.contains(@["-0"]) == 2
    check mudskipper_contains.contains(@["--"]) == 2

  test "double dash allows dash-leading needle":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    var code = cmd.contains(@["--", "-x", "foo", "-x", "bar"])
    check code == 0
    bufs.reset()
    code = cmd.contains(@["-i", "--", "-x", "foo", "-x", "bar"])
    check code == 0
    check bufs.outBuf == "2\n"
    bufs.reset()
    code = cmd.contains(@["-0", "--", "-x", "foo", "-x", "bar"])
    check code == 0
    check bufs.outBuf == "1\n"

  test "double dash dash-leading needle not found":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    let code = cmd.contains(@["--", "-missing", "foo", "bar"])
    check code == 1

  test "empty string needle cases":
    check mudskipper_contains.contains(@["", ""]) == 0
    check mudskipper_contains.contains(@["", "a", "b"]) == 1
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    var code = cmd.contains(@["-i", "", "", "a"])
    check code == 0
    check bufs.outBuf == "1\n"
    bufs.reset()
    code = cmd.contains(@["-0", "", "", "a"])
    check code == 0
    check bufs.outBuf == "0\n"

  test "needle equals double dash after sentinel":
    let bufs = CaptureBuffers(); var cmd = newContainsCmd(bufs)
    let code = cmd.contains(@["-i", "--", "--", "a", "b", "c", "--", "v"])
    check code == 0
    check bufs.outBuf == "4\n"
