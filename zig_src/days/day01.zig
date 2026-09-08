const std = @import("std");
const common = @import("common");

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);
    return switch (part) {
        1 => part1(lines),
        2 => part2(lines),
        else => return error.InvalidInput,
    };
}

fn part1(lines: [][]const u8) !u64 {
    var position: u8 = 50;
    var zero_hit: u64 = 0;
    for (lines) |line| {
        if (line.len == 0) continue;
        const rotation_number = try parseLine(line);
        if (rotate_p1(&position, rotation_number)) zero_hit += 1;
    }
    return zero_hit;
}

fn parseLine(line: []const u8) !i16 {
    const number = try std.fmt.parseInt(i16, line[1..],10);
    return switch (line[0]) {
        'R' => number,
        'L' => -number,
        else => return error.InvalidInput,
    };
}

// rotate_p1 rotates by value. Returns true if stopped a position 0
fn rotate_p1(position: *u8, rot_val: i16) bool {
    position.* = @mod(position.* + @as(u8, @intCast(@mod(rot_val, 100))), 100);
    return position.* == 0;
}

fn part2(lines: [][]const u8) !u64 {
    var position: u8 = 50;
    var zero_hit: u64 = 0;
    for (lines) |line| {
        if (line.len == 0) continue;
        const rotation_number = try parseLine(line);
        zero_hit += rotate_p2(&position, rotation_number);
    }
    return zero_hit;
}

// rotate_p2 rotates by value. Returns times going over position 0
fn rotate_p2(position: *u8, rot_val: i16) u16 {
    var zero_hit: u16 = @abs(@divTrunc(rot_val, 100));
    const turn_val: i16 = @rem(rot_val, 100);
    if (turn_val != 0) {
        const new_pos: i16 = @as(i16, @intCast(position.*)) + turn_val;
        if (position.* != 0 and (new_pos <= 0 or new_pos >= 100)) zero_hit += 1;
        position.* = @intCast(@mod(new_pos, 100));
    }
    return zero_hit;
}

const testing = std.testing;

fn runCase(start: u8, rot: i16, expect_final: u8, expect_hits: u16) !void {
    var p: u8 = start;
    const hits = rotate_p2(&p, rot);
    try testing.expectEqual(expect_final, p);
    try testing.expectEqual(expect_hits, hits);
}

test "rotate_p2: start=50, +30" {
    try runCase(50, 30, 80, 0);
}

test "rotate_p2: start=50, +50" {
    try runCase(50, 50, 0, 1);
}

test "rotate_p2: start=50, +130" {
    try runCase(50, 130, 80, 1);
}

test "rotate_p2: start=50, +150" {
    try runCase(50, 150, 0, 2);
}

test "rotate_p2: start=50, -30" {
    try runCase(50, -30, 20, 0);
}

test "rotate_p2: start=50, -50" {
    try runCase(50, -50, 0, 1);
}

test "rotate_p2: start=50, -130" {
    try runCase(50, -130, 20, 1);
}

test "rotate_p2: start=50, -150" {
    try runCase(50, -150, 0, 2);
}

test "rotate_p2: start=0, +30" {
    try runCase(0, 30, 30, 0);
}

test "rotate_p2: start=0, +100" {
    try runCase(0, 100, 0, 1);
}

test "rotate_p2: start=0, +130" {
    try runCase(0, 130, 30, 1);
}

test "rotate_p2: start=0, -30" {
    try runCase(0, -30, 70, 0);
}

test "rotate_p2: start=0, -100" {
    try runCase(0, -100, 0, 1);
}

test "rotate_p2: start=0, -130" {
    try runCase(0, -130, 70, 1);
}
