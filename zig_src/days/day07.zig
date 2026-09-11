const std = @import("std");
const common = @import("common");

const Input: type = [][]const u8;

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);
    return switch (part) {
        1 => part1(allocator, lines),
        2 => part2(allocator, lines),
        else => return error.InvalidInput,
    };
}

fn part1(allocator: std.mem.Allocator, lines: Input) !u64 {
    var sum: u64 = 0;
    const beams = try allocator.alloc(bool, lines[0].len);
    defer allocator.free(beams);
    for (lines[0], 0..) |char, i| {
        beams[i] = false;
        if (char == 'S') {
            beams[i] = true;
        }
    }

    var i: usize = 2;
    while (i < lines.len) : (i += 2) {
        const line = lines[i];
        sum += p1Calc(beams, line);
    }
    return sum;
}

fn p1Calc(beams: []bool, line: []const u8) u64 {
    var splits: u64 = 0;
    for (1..beams.len-1) |i| {
        if (beams[i] and line[i] == '^') {
            beams[i-1] = true;
            beams[i] = false;
            beams[i+1] = true;
            splits += 1;
        }
    }
    return splits;
}

fn part2(allocator: std.mem.Allocator, lines: Input) !u64 {
    const beams = try allocator.alloc(u64, lines[0].len);
    defer allocator.free(beams);
    for (lines[0], 0..) |char, i| {
        beams[i] = 0;
        if (char == 'S') {
            beams[i] = 1;
        }
    }

    var i: usize = 2;
    while (i < lines.len) : (i += 2) {
        const line = lines[i];
        p2Calc(beams, line);
    }

    var sum: u64 = 0;
    for (beams) |b| {
        sum += b;
    }
    return sum;
}

fn p2Calc(beams: []u64, line: []const u8) void {
    for (1..beams.len-1) |i| {
        if (beams[i] != 0 and line[i] == '^') {
            beams[i-1] += beams[i];
            beams[i+1] += beams[i];
            beams[i] = 0;
        }
    }
}