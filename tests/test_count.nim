# filepath: /Users/matt/Projects/mattmc3/mudskipper/tests/test_count.nim
import unittest
import ../src/count

test "count with no args":
  check count(@[]) == 0

test "count with some args":
  check count(@["one"]) == 1
  check count(@["one", "two", "three"]) == 3
