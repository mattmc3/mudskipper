#!/usr/bin/env fish

set --local i
for i in (seq 100)
    set --append argv $i
    count $argv
end
