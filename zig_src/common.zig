const std = @import("std");

pub fn readFile(io: std.Io, allocator: std.mem.Allocator, path: []const u8) ![]u8 {
    return std.Io.Dir.cwd().readFileAlloc(io,path, allocator, .unlimited);
}

pub fn splitLines(allocator: std.mem.Allocator, input: []const u8) ![][]const u8 {
    var list: std.ArrayList([]const u8) = .empty;
    errdefer list.deinit(allocator);
    var it = std.mem.splitScalar(u8, std.mem.trimEnd(u8, input, "\n"), '\n');
    while (it.next()) |line| {
        try list.append(allocator, std.mem.trimEnd(u8, line, "\r"));
    }
    return list.toOwnedSlice(allocator);
}

pub fn parseInt(s: []const u8) !i64 {
    return std.fmt.parseInt(i64, std.mem.trim(u8, s, " \t\r\n"), 10);
}

pub fn parseInts(allocator: std.mem.Allocator, input: []const u8, sep: u8) ![]i64 {
    var list: std.ArrayList(i64) = .empty;
    errdefer list.deinit(allocator);
    var it = std.mem.splitScalar(u8, std.mem.trim(u8, input, " \t\r\n"), sep);
    while (it.next()) |part| {
        if (part.len == 0) continue;
        try list.append(allocator, try parseInt(part));
    }
    return list.toOwnedSlice(allocator);
}
