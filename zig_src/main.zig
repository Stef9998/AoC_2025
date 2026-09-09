const std = @import("std");
const clap = @import("clap");
const common = @import("common");

const std_input_file_name = "input";

const Modules = [_]type{
    struct {},
    @import("days/day01.zig"),
    @import("days/day02.zig"),
    // @import("days/day03.zig"),
    // @import("days/day04.zig"),
    // @import("days/day05.zig"),
    // @import("days/day06.zig"),
    // @import("days/day07.zig"),
    // @import("days/day08.zig"),
    // @import("days/day09.zig"),
    // @import("days/day10.zig"),
    // @import("days/day11.zig"),
    // @import("days/day12.zig"),
};

const Day: type = *const fn (std.mem.Allocator, []const u8, u2) anyerror!u64;

const DayTable = blk: {
    var tmp: [Modules.len]?Day = undefined;
    var i: usize = 0;
    while (i < Modules.len) : (i += 1) {
        tmp[i] = if (@hasDecl(Modules[i], "run")) @field(Modules[i], "run") else null;
    }
    break :blk tmp;
};

fn lookupDay(day: usize) ?Day {
    if (day >= DayTable.len) return null;
    return DayTable[day];
}

pub fn main(init: std.process.Init) !void {
    const allocator = init.gpa;

    const params = comptime clap.parseParamsComptime(
        \\-h, --help             Display this help and exit.
        \\-d, --day <usize>      Day number (defaults to newest).
        \\-p, --part <usize>     Part number (defaults to both).
        \\-i, --input <str>      Path to the input file (without file ending) (defaults to input[.txt]).
        \\-a, --all              Run all implemented days
        \\
    );

    var diag = clap.Diagnostic{};
    const res = clap.parse(clap.Help, &params, clap.parsers.default, init.minimal.args, .{
        .diagnostic = &diag,
        .allocator = allocator,
    }) catch |err| {
        try diag.reportToFile(init.io, .stderr(), err);
        return err;
    };
    defer res.deinit();

    if (res.args.help != 0) {
        try clap.helpToFile(init.io, .stdout(), clap.Help, &params, .{});
        std.process.exit(0);
    }

    const io = init.io;

    if (res.args.all != 0) {
        for (1..DayTable.len) |day| {
            std.debug.print("Day {d}:\n", .{day});
            try run(io, allocator, day, null, std_input_file_name);
        }
    } else {
        const day = res.args.day orelse DayTable.len - 1;
        if (res.args.input) |file_name| {
            try run(io, allocator, day, res.args.part, file_name);
        } else {
            try run(io, allocator, day, res.args.part, std_input_file_name);
        }
    }
}

fn run(io: std.Io, allocator: std.mem.Allocator, day: usize, part: ?usize, file_name: []const u8) !void {
    const runner = lookupDay(day) orelse {
        std.debug.print("day {d} not implemented\n", .{day});
        return;
    };

    var path_buf: [64]u8 = undefined;
    const path = try std.fmt.bufPrint(&path_buf, "../input/day{d}/{s}.txt", .{ day, file_name });
    const input = try common.readFile(io, allocator, path);
    defer allocator.free(input);

    if (part == null or part.? == 1) {
        _ = try runPart(io, allocator, runner, input, 1);
    }
    if (part == null or part.? == 2) {
        _ = try runPart(io, allocator, runner, input, 2);
    }
}

fn runPart(io: std.Io, allocator: std.mem.Allocator, runner: *const fn (std.mem.Allocator, []const u8, u2) anyerror!u64, input: []const u8, part: u2) !u64 {
    const start = std.Io.Timestamp.now(io, .awake);
    const result = try runner(allocator, input, part);
    std.debug.print("Part {d}: {d} ({f})\n", .{ part, result, start.untilNow(io, .awake) });
    return result;
}

test "implemented days match expected.csv" {
    const allocator = std.testing.allocator;
    const io = std.testing.io;

    const csv = try common.readFile(io, allocator, "../expected.csv");
    defer allocator.free(csv);

    var lines = std.mem.splitScalar(u8, csv, '\n');
    _ = lines.next(); // header

    while (lines.next()) |raw_line| {
        const line = std.mem.trimEnd(u8, raw_line, "\r");
        if (line.len == 0) continue;

        var fields = std.mem.splitScalar(u8, line, ',');
        const day = try std.fmt.parseInt(usize, fields.next().?, 10);
        const part = try std.fmt.parseInt(u2, fields.next().?, 10);
        const type_field = fields.next().?;
        const input_type: u2 = if (type_field.len == 0) 0 else try std.fmt.parseInt(u2, type_field, 10);
        const answer_field = fields.next().?;
        const expected: u64 = if (answer_field.len == 0) continue else try std.fmt.parseInt(u64, answer_field, 10);

        const runner = lookupDay(day) orelse continue;

        var path_buf: [64]u8 = undefined;
        const path = switch (input_type) {
            0 => try std.fmt.bufPrint(&path_buf, "../input/day{d}/{s}.txt", .{ day, std_input_file_name }),
            else => return error.NotImplemented,
        };
        const input = try common.readFile(io, allocator, path);
        defer allocator.free(input);

        const actual = try runner(allocator, input, part);
        std.testing.expectEqual(expected, actual) catch |err| {
            std.debug.print("day {d} part {d}: expected {d}, got {d}\n", .{ day, part, expected, actual });
            return err;
        };
    }
}
