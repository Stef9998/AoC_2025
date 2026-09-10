const std = @import("std");
const common = @import("common");

const LineIterator = std.mem.SplitIterator(u8, .scalar);
const Input: type = struct { []Constrain, *ParseIntIterator(LineIterator) };
const Cell = u8;

const Constrain = struct {
    min: u64,
    max: u64,
};

fn ParseIntIterator(comptime It: type) type {
    return struct {
        split_it: It,

        pub fn next(self: *@This()) !?u64 {
            const line = self.split_it.next() orelse return null;
            return try std.fmt.parseInt(u64, std.mem.trimEnd(u8, line, "\r"), 10);
        }
    };
}

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    var constrain_list: std.ArrayList(Constrain) = try std.ArrayList(Constrain).initCapacity(allocator, 100);
    defer constrain_list.deinit(allocator);
    var it: LineIterator = std.mem.splitScalar(u8, std.mem.trimEnd(u8, input, "\n"), '\n');
    while (it.next()) |line| {
        const line_trimmed = std.mem.trimEnd(u8, line, "\r");
        if (std.mem.eql(u8, line_trimmed, "")) break;
        var constr_it = std.mem.splitScalar(u8, line_trimmed, '-');
        const min = try std.fmt.parseInt(u64, constr_it.next().?, 10);
        const max = try std.fmt.parseInt(u64, constr_it.next().?, 10);
        std.debug.assert(constr_it.next() == null);
        try constrain_list.append(allocator, Constrain{ .min = min, .max = max });
    }
    const constrains = try constrain_list.toOwnedSlice(allocator);
    defer allocator.free(constrains);
    std.mem.sort(Constrain, constrains, {}, struct {
        fn lessThan(_: void, a: Constrain, b: Constrain) bool {
            return a.min < b.min or (a.min == b.min and a.max < b.max);
        }
    }.lessThan);
    var ingredients: std.ArrayList(u64) = .empty;
    defer ingredients.deinit(allocator);
    var int_it = ParseIntIterator(LineIterator){ .split_it = it };
    return switch (part) {
        1 => part1(.{ constrains, &int_it }),
        2 => part2(.{ constrains, &int_it }),
        else => return error.InvalidInput,
    };
}

fn part1(map: Input) !u64 {
    const constrains, const ingredients = map;
    var sum: u64 = 0;
    while (try ingredients.next()) |ingredient| {
        if (p1Part(constrains, ingredient)) sum += 1;
    }
    return sum;
}

fn p1Part(constrains: []Constrain, ingredient: u64) bool {
    var i: usize = 0;
    while (i < constrains.len) : (i += 1) {
        if (constrains[i].min > ingredient) break;
    }
    if (i == 0) return false;
    for (0..i) |j_| {
        const j = i - 1 - j_;
        if (constrains[j].max >= ingredient) return true;
    }
    return false;
}

fn part2(map: Input) !u64 {
    const constrains, _ = map;
    var curr_idx: u64 = 0;
    var sum: u64 = 0;
    for (constrains) |constrain| {
        sum += p2Calc(constrain, &curr_idx);
    }
    return sum;
}

fn p2Calc(constrain: Constrain, index: *u64) u64 {
    if (index.* <= constrain.min) {
        index.* = constrain.max + 1;
        return constrain.max - constrain.min + 1;
    }
    if (index.* > constrain.max) {
        return 0;
    }
    const result = constrain.max - index.* + 1;
    index.* = constrain.max + 1;
    return result;
}
