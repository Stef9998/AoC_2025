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

pub fn GrowableIndexList(T: type) type {
    return struct {
        list: []T,

        const Self = @This();

        pub fn init(allocator: std.mem.Allocator, initial_capacity: usize) std.mem.Allocator.Error!Self {
            return Self{
                .list = try allocator.alloc(T, initial_capacity),
            };
        }

        pub fn fill(self: *Self, value: T) void {
            for (self.list) |*item| {
                item.* = value;
            }
        }

        pub fn get(self: Self, index: usize) ?T {
            if (index >= self.list.len) return null;
            return self.list[index];
        }

        pub fn set(self: *Self, allocator: std.mem.Allocator, pos: usize, value: T) std.mem.Allocator.Error!void {
            while (pos >= self.list.len) {
                const new_capacity = if (self.list.len == 0) 16 else self.list.len * 2;
                const new_list = try allocator.realloc(self.list, new_capacity);
                self.list = new_list;
            }
            self.list[pos] = value;
        }

        pub fn deinit(self: *Self, allocator: std.mem.Allocator) void {
            allocator.free(self.list);
        }

        pub fn toOwnedSlice(self: *Self) []T {
            const temp = self.list;
            self.list = undefined;
            return temp;
        }
    };
}