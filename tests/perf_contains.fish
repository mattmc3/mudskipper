#!/usr/bin/env fish
set -l haystack foo bar baz qux quux corge grault garply waldo fred plugh xyzzy thud
set -l needle grault

set -l scriptdir (status dirname)
set -l nim_contains $scriptdir/../bin/contains

echo "Testing Fish builtin contains..."
time for i in (seq 1000)
    contains $needle $haystack >/dev/null
end

echo "Testing Nim contains..."
time for i in (seq 1000)
    $nim_contains $needle $haystack >/dev/null
end
