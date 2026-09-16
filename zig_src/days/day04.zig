const std = @import("std");
const common = @import("common");
const map_util = @import("map_util");
const Coordinate = map_util.Coordinate;
const ImmutableMap = map_util.Map(u8);

const SURROUNDING_COORDINATES: []const Coordinate = &[8]Coordinate{
    .{ .x = -1, .y = -1 }, .{ .x = 0, .y = -1 }, .{ .x = 1, .y = -1 },
    .{ .x = -1, .y = 0 },  .{ .x = 1, .y = 0 },  .{ .x = -1, .y = 1 },
    .{ .x = 0, .y = 1 },   .{ .x = 1, .y = 1 },
};

const Input: type = Map;
const Cell = u8;

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);
    return switch (part) {
        1 => part1(Map{.immutable = ImmutableMap.init(lines)}),
        2 => part2(allocator, Map{.mutable = try MutableMap.init(allocator, lines)}),
        else => return error.InvalidInput,
    };
}

const Map = union(enum) {
    immutable: ImmutableMap,
    mutable: MutableMap,

    fn deinit(self: Map, allocator: std.mem.Allocator) void {
        switch (self) {
            .immutable => {},
            .mutable => |m| m.deinit(allocator),
        }
    }

    fn getCell(self: Map, coord: Coordinate) ?Cell {
        return switch (self) {
            inline else => |m| m.getCell(coord),
        };
    }
};

const MutableMap = struct {
    mapping: [][]Cell,
    height: usize,
    width: usize,

    fn init(allocator: std.mem.Allocator, lines: [][]const Cell) !MutableMap {
        const copy = try allocator.alloc([]Cell, lines.len);
        var line_count: usize = 0;
        errdefer {
            for (0..line_count) |i| allocator.free(copy[i]);
            allocator.free(copy);
        }
        while (line_count < lines.len) : (line_count+=1) {
            copy[line_count] = try allocator.dupe(Cell, lines[line_count]);
        }
        return .{ .mapping = copy, .height = lines.len, .width = lines[0].len };
    }

    fn deinit(map: MutableMap, allocator: std.mem.Allocator) void {
        for (map.mapping) |line| allocator.free(line);
        allocator.free(map.mapping);
    }

    fn getCell(map: MutableMap, coord: Coordinate) ?Cell {
        if (coord.y < 0 or coord.y >= map.height) return null;
        if (coord.x < 0 or coord.x >= map.width) return null;
        return map.mapping[@intCast(coord.y)][@intCast(coord.x)];
    }

    fn setCell(map: *MutableMap, coord: Coordinate, value: Cell) !void {
        if (coord.y < 0 or coord.y >= map.height) return error.InvalidInput;
        if (coord.x < 0 or coord.x >= map.width) return error.InvalidInput;
        map.mapping[@intCast(coord.y)][@intCast(coord.x)] = value;
    }
};

fn part1(map: Input) !u64 {
    var sum: u64 = 0;
    for (0..map.immutable.mapping.len) |i| {
        sum += p1CalculateRow(map, i);
    }
    return sum;
}

fn p1CalculateRow(map: Map, row: usize) u64 {
    var sum: u64 = 0;
    const line = map.immutable.mapping[row];
    var curr_coord = Coordinate{ .x = 0, .y = @as(isize, @intCast(row)) };
    for (0..line.len) |i| {
        if (line[i] != '@') continue;
        curr_coord.x = @as(isize, @intCast(i));
        if (isUnaccessable(map, curr_coord)) continue;
        sum += 1;
    }
    return sum;
}

fn isUnaccessable(map: Map, curr_coord: Coordinate) bool {
    var occ: u3 = 0;
    for (SURROUNDING_COORDINATES) |offset| {
        const next = map_util.addCoord(curr_coord, offset);
        const coordinate = map.getCell(next);
        if (coordinate) |coord| {
            if (coord != '.') occ += 1;
        }
        if (occ > 3) return true;
    }
    return false;
}

fn part2(allocator: std.mem.Allocator, map: Input) !u64 {
    defer map.deinit(allocator);
    var sum: u64 = 0;
    var map_impl = map.mutable;
    const upper_limit = map_impl.height * map_impl.width;
    for (0..upper_limit) |_| {
        var change: u64 = 0;
        for (0..map_impl.height) |i| {
            change += p2CalculateRow(map, i);
        }
        sum += change;
        if (change == 0) break;

        var coord = Coordinate{.x=0, .y=0};
        for (0..map_impl.height) |i| {
            coord.y = @intCast(i);
            for (0..map_impl.width) |j| {
                coord.x = @intCast(j);
                if (map_impl.getCell(coord) == 'X') {
                    map_impl.setCell(coord, '.') catch {};
                }
            }
        }
    }
    return sum;
}

fn p2CalculateRow(map: Map, row: usize) u64 {
    var sum: u64 = 0;
    var map_impl = map.mutable;
    const line = map_impl.mapping[row];
    var curr_coord = Coordinate{ .x = 0, .y = @as(isize, @intCast(row)) };
    for (0..line.len) |i| {
        if (line[i] != '@') continue;
        curr_coord.x = @as(isize, @intCast(i));
        if (isUnaccessable(map, curr_coord)) {
            continue;
        }
        map_impl.setCell(curr_coord, 'X') catch {};
        sum += 1;
    }
    return sum;
}

const testing = std.testing;

test "p1 first line" {
    var input = [_][]const Cell{
        "@@.@.@.@@@",
        "@@@@@@@@@@",
        ".@@@@@@@@.",
    };
    const map = Map{.immutable = ImmutableMap.init(input[0..])};
    try std.testing.expectEqual(4, part1(map));
}

test "p1 last line" {
    var input = [_][]const Cell{
        ".@@@@@@@@.",
        "@@@@@@@@@@",
        "@@.@.@.@@@",
    };
    const map = Map{.immutable = ImmutableMap.init(input[0..])};
    try std.testing.expectEqual(4, part1(map));
}

test "p1 1" {
    var input = [_][]const Cell{
        ".....",
        ".@@..",
        ".@@@.",
        "..@@.",
        ".....",
    };
    const map = Map{.immutable = ImmutableMap.init(input[0..])};
    try std.testing.expectEqual(2, part1(map));
}
