import unittest
import ../src/contains as mudskipper_contains

suite "contains":
  test "needle found returns 0":
    check mudskipper_contains.contains(@["foo", "bar", "baz"]) == 1
    check mudskipper_contains.contains(@["bar", "foo", "bar", "baz"]) == 0

  test "needle not found returns 1":
    check mudskipper_contains.contains(@["missing", "foo", "bar", "baz"]) == 1

  test "index flag returns 0 if found":
    check mudskipper_contains.contains(@["-i", "bar", "foo", "bar", "baz"]) == 0

  test "index0 flag returns 0 if found":
    check mudskipper_contains.contains(@["-0", "bar", "foo", "bar", "baz"]) == 0

  test "help returns 0":
    check mudskipper_contains.contains(@["-h"]) == 0
    check mudskipper_contains.contains(@["--help"]) == 0

  test "key not specified returns 2":
    check mudskipper_contains.contains(@["-i"]) == 2
    check mudskipper_contains.contains(@["-0"]) == 2
    check mudskipper_contains.contains(@[]) == 0 # help is shown for empty args
