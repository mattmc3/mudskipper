#!/usr/bin/env fish
set -l haystack foo bar baz qux quux corge grault garply waldo fred plugh xyzzy thud
set -l needle grault

set -l scriptdir (status dirname)
set -l nim_contains $scriptdir/../bin/contains
set -l zig_contains $scriptdir/../bin/contains_zig
set -l go_contains $scriptdir/../bin/contains_go
# set -l lua_contains $scriptdir/../tools/contains
# set -l perl_contains $scriptdir/../tools/contains.pl

echo "Testing Fish builtin contains..."
time for i in (seq 1000)
    contains $needle $haystack >/dev/null
end

echo "Testing Nim contains..."
time for i in (seq 1000)
    $nim_contains $needle $haystack >/dev/null
end

echo "Testing Zig contains..."
time for i in (seq 1000)
    $zig_contains $needle $haystack >/dev/null
end

echo "Testing Go contains..."
time for i in (seq 1000)
    $go_contains $needle $haystack >/dev/null
end

# echo "Testing Perl contains..."
# time for i in (seq 1000)
#     $perl_contains $needle $haystack >/dev/null
# end

# echo "Testing Lua contains..."
# time for i in (seq 1000)
#     $lua_contains $needle $haystack >/dev/null
# end
