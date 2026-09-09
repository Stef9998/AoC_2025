const std = @import("std");
const common = @import("common");

const Input: type = [][]const u8;
const Parsed: type = []const u8;

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);
    return switch (part) {
        1 => part1(lines),
        2 => part2(lines),
        else => return error.InvalidInput,
    };
}

fn part1(lines: Input) !u64 {
    var sum: u64 = 0;
    for (lines) |line| {
        const temp = p1_calc(line);
        sum += temp;
    }
    return sum;
}

fn p1_calc(line: Parsed) u64 {
    var highest = line[0];
    var second: u8 = line[line.len - 1];
    for (1..line.len - 1) |i| {
        const joltage = line[i];
        if (joltage > highest) {
            highest = joltage;
            second = line[line.len - 1];
        } else if (joltage > second) {
            second = joltage;
        }
    }
    return (highest - '0') * 10 + (second - '0');
}

fn part2(lines: Input) !u64 {
    var sum: u64 = 0;
    var buffer: [12]u8 = undefined;
    for (lines) |line| {
        buffer = @splat('1');
        const temp = p2_calc(line, buffer[0..]);
        sum += temp;
    }
    return sum;
}

fn p2_calc(line: Parsed, buf: []u8) u64 {
    for (0..line.len - 12) |i| {
        const val = line[i];
        for (0..12) |j| {
            if (val > buf[j]) {
                buf[j] = val;
                @memset(buf[j + 1 .. 12], '1');
                break;
            }
        }
    }
    for (0..12) |i_| {
        const i = line.len-12+i_;
        const val = line[i];
        for (i_..12) |j| {
            if (val > buf[j]) {
                buf[j] = val;
                @memset(buf[j + 1 .. 12], '1');
                break;
            }
        }
    }
    var ret: u64 = 0;
    for (buf) |char| {
        std.debug.assert(std.ascii.isDigit(char));
        ret *= 10;
        ret += char - '0';
    }
    return ret;
}

const testing = std.testing;

test "input line 1" {
    const line =
        "3233434223352253322244323562413222322422522622312422332123223123235422212196323222332232332242222211";

    var buffer: [12]u8 = @splat('1');

    try testing.expectEqual(@as(u64, 963342222211), p2_calc(line, buffer[0..]));
}
