const std = @import("std");
const common = @import("common");

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    _ = allocator;
    var it = pairs(input);
    return switch (part) {
        1 => part1(&it),
        2 => part2(&it),
        else => return error.InvalidInput,
    };
}

const Pair = struct { u64, u64 };

const PairIterator = struct {
    inner: std.mem.SplitIterator(u8, .scalar),

    fn next(self: *PairIterator) !?Pair {
        const raw = self.inner.next() orelse return null;
        return try parse(raw);
    }
};

fn pairs(input: []const u8) PairIterator {
    return .{ .inner = std.mem.splitScalar(u8, input, ',') };
}

fn parse(tuple: []const u8) !Pair {
    var it = std.mem.splitScalar(u8, tuple, '-');
    const first = it.next() orelse return error.InvalidInput;
    const second = it.next() orelse return error.InvalidInput;
    if (it.next()) |_| return error.InvalidInput;
    return .{ @as(u64, @intCast(try common.parseInt(first))), @as(u64, @intCast(try common.parseInt(second))) };
}

fn part1(it: *PairIterator) !u64 {
    var sum: u64 = 0;
    while (try it.next()) |pair| {
        const lower, const upper = pair;
        sum += try p1Calc(lower, upper);
    }
    return sum;
}

fn p1Calc(lower: u64, upper: u64) !u64 {
    if (upper < 11) return 0;
    var sum: u64 = 0;
    for (@max(11, lower)..upper + 1) |id| {
        const id_len: u64 = std.math.log10_int(id) + 1;
        if (id_len % 2 != 0) continue;
        const mask = try std.math.powi(u64, 10, id_len / 2);
        if (id / mask == id % mask) sum += id;
    }
    return sum;
}

fn part2(it: *PairIterator) !u64 {
    var sum: u64 = 0;
    while (try it.next()) |pair| {
        const lower, const upper = pair;
        sum += try p2Calc(lower, upper);
    }
    return sum;
}

fn p2Calc(lower: u64, upper: u64) !u64 {
    if (upper < 11) return 0;
    var sum: u64 = 0;
    for (@max(11, lower)..upper + 1) |id| {
        const id_len: u64 = std.math.log10_int(id) + 1;
        outer: for (2..id_len+1) |div| {
            if (id_len % div != 0) continue;
            const mask = try std.math.powi(u64, 10, id_len / div);
            var og = id;
            var shifted = id / mask;
            for (1..div) |_| {
                if (og % mask != shifted % mask) continue :outer;
                og = shifted;
                shifted /= mask;
            }
            sum += id;
            break;
        }
    }
    return sum;
}
