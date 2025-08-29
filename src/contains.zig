const std = @import("std");

pub fn main() !void {
    const stdout = std.fs.File.stdout();
    const stderr = std.fs.File.stderr();
    const args = try std.process.argsAlloc(std.heap.page_allocator);

    defer std.process.argsFree(std.heap.page_allocator, args);

    if (args.len <= 1 or std.mem.eql(u8, args[1], "-h") or std.mem.eql(u8, args[1], "--help")) {
        try stdout.writeAll("Usage: contains [-i|--index|-0|--index0] NEEDLE HAYSTACK...\n" ++ "  -i, --index      Print the 1-based index of NEEDLE in HAYSTACK\n" ++ "  -0, --index0     Print the 0-based index of NEEDLE in HAYSTACK\n" ++ "  -h, --help       Show this help message\n" ++ "Returns 0 if NEEDLE is found, 1 otherwise.\n\n");
        return;
    }

    var index_flag = false;
    var index_start: usize = 1;
    var start: usize = 1;

    if (std.mem.eql(u8, args[1], "-i") or std.mem.eql(u8, args[1], "--index")) {
        index_flag = true;
        index_start = 1;
        start = 2;
    } else if (std.mem.eql(u8, args[1], "-0") or std.mem.eql(u8, args[1], "--index0")) {
        index_flag = true;
        index_start = 0;
        start = 2;
    }

    if (args.len <= start) {
        try stderr.writeAll("contains: Key not specified\n");
        std.process.exit(2);
    }

    const needle = args[start];
    var found = false;
    var i: usize = 0;
    while (i < args.len - (start + 1)) : (i += 1) {
        const item = args[start + 1 + i];
        if (std.mem.eql(u8, item, needle)) {
            found = true;
            if (index_flag) {
                // format index into buffer then write
                var buf: [32]u8 = undefined;
                const out = try std.fmt.bufPrint(&buf, "{}\n", .{i + index_start});
                try stdout.writeAll(out);
            }
            break;
        }
    }

    if (found) {
        std.process.exit(0);
    } else {
        std.process.exit(1);
    }
}
