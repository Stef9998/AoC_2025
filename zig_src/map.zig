const std = @import("std");

pub const Coordinate = struct {
    x: isize,
    y: isize,

    pub fn offset(self: *Coordinate, c2: Coordinate) void {
        self.x += c2.x;
        self.y += c2.y;
    }
};

pub fn addCoord(c1: Coordinate, c2: Coordinate) Coordinate {
    return Coordinate{
        .x = c1.x + c2.x,
        .y = c1.y + c2.y,
    };
}

pub fn Map(Cell: type) type {
    return struct {
        mapping: [][]const Cell,
        height: usize,
        width: usize,

        const Self = @This();

        pub fn init(lines: [][]const Cell) Self {
            return Self{
                .mapping = lines,
                .height = lines.len,
                .width = lines[0].len,
            };
        }

        pub fn getCell(map: Self, coord: Coordinate) ?Cell {
            if (coord.y < 0 or coord.y >= map.height) return null;
            if (coord.x < 0 or coord.x >= map.width) return null;
            return map.mapping[@intCast(coord.y)][@intCast(coord.x)];
        }
    };
}

pub fn printMap(mapping: []const []const u8) void {
    for (mapping) |row| {
        std.debug.print("{s}\n", .{row});
    }
}
