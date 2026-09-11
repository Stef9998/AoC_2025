const std = @import("std");
const common = @import("common");

const Iterator = std.mem.TokenIterator(u8, .scalar);
const Input: type = []Iterator;

pub fn run(allocator: std.mem.Allocator, input: []const u8, part: u2) !u64 {
    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);
    var iterators: []Iterator = try allocator.alloc(Iterator, lines.len);
    defer allocator.free(iterators);
    for (lines, 0..) |line, i| {
        iterators[i] = std.mem.tokenizeScalar(u8, line, ' ');
    }
    return switch (part) {
        1 => part1(iterators),
        2 => part2(allocator, lines),
        else => return error.InvalidInput,
    };
}

const calculator = struct {
    const BinaryOp = *const fn (u64, u64) u64;

    fn add(a: u64, b: u64) u64 {
        return a + b;
    }

    fn mul(a: u64, b: u64) u64 {
        return a * b;
    }

    pub fn startNum(operator: u8) !u64 {
        return switch (operator) {
            '+' => 0,
            '*' => 1,
            else => return error.InvalidInput,
        };
    }

    pub fn calculation(operator: u8) !BinaryOp {
        return switch (operator) {
            '+' => add,
            '*' => mul,
            else => error.InvalidInput,
        };
    }
};

const Calculator = struct {
    value: u64,
    operation: calculator.BinaryOp,

    fn init(operator: u8) !Calculator {
        return Calculator{
            .value = try calculator.startNum(operator),
            .operation = try calculator.calculation(operator),
        };
    }

    fn calc(self: *Calculator, number: u64) void {
        self.value = self.operation(self.value, number);
    }

    fn getResult(self: Calculator) u64 {
        return self.value;
    }
};

fn part1(iterators: Input) !u64 {
    var sum: u64 = 0;
    var op_it = &iterators[iterators.len - 1];
    const cal_its = iterators[0 .. iterators.len - 1];
    while (op_it.next()) |operator_str| {
        std.debug.assert(operator_str.len == 1);
        const operator = operator_str[0];
        sum += try p1Calc(cal_its, operator);
    }
    return sum;
}

fn p1Calc(cal_its: []Iterator, operator: u8) !u64 {
    var calcer = try Calculator.init(operator);
    for (cal_its) |*it| {
        if (it.next()) |num| {
            const number = try std.fmt.parseInt(u64, num, 10);
            calcer.calc(number);
        } else return error.InvalidInput;
    }
    return calcer.getResult();
}

fn part2(allocator: std.mem.Allocator, lines: [][]const u8) !u64 {
    var sum: u64 = 0;
    var operator_pos: usize = 0;
    const operators_line = lines[lines.len - 1];
    const calculations_lines = lines[0 .. lines.len - 1];
    var calculation_lines = try allocator.alloc([]const u8, lines.len - 1);
    defer allocator.free(calculation_lines);
    for (1..operators_line.len) |i| {
        const possible_op = operators_line[i];
        if (possible_op == ' ') continue;
        for (calculations_lines, 0..) |line, j| {
            calculation_lines[j] = line[operator_pos..i];
        }
        sum += try p2Calc(calculation_lines, operators_line[operator_pos], i - operator_pos - 1);
        operator_pos = i;
    }
    std.debug.assert(operator_pos == operators_line.len - 1);
    std.debug.assert(operators_line[operator_pos] != ' ');
    var longest_line: usize = 0;
    for (calculations_lines, 0..) |line, j| {
        calculation_lines[j] = line[operator_pos..];
        if (calculation_lines[j].len > longest_line) longest_line = calculation_lines[j].len;
    }
    sum += try p2Calc(calculation_lines, operators_line[operator_pos], longest_line);
    return sum;
}

fn p2Calc(lines: [][]const u8, operator: u8, length: usize) !u64 {
    var calcer = try Calculator.init(operator);
    for (0..length) |i| {
        var num: u64 = 0;
        for (lines) |line| {
            if (i >= line.len) continue;
            const char = line[i];
            if (char == ' ') continue;
            num *= 10;
            if (!std.ascii.isDigit(char)) return error.InvalidInput;
            num += try std.fmt.charToDigit(char, 10);
        }
        calcer.calc(num);
    }
    return calcer.getResult();
}

test "day06 p2Calc" {
    var lines = [_][]const u8{
        "123 ",
        " 45 ",
        "  6 ",
    };

    const result = try p2Calc(lines[0..], '*', 3);
    try std.testing.expectEqual(@as(u64, 8544), result);
}

test "day06 parsing example operators" {
    const allocator = std.testing.allocator;
    const io = std.testing.io;

    const input = try common.readFile(io, allocator, "../input/day6/example.txt");
    defer allocator.free(input);

    const lines = try common.splitLines(allocator, input);
    defer allocator.free(lines);

    try std.testing.expectEqual(@as(usize, 4), lines.len);

    var iterators: []Iterator = try allocator.alloc(Iterator, lines.len);
    defer allocator.free(iterators);
    for (lines, 0..) |line, i| {
        iterators[i] = std.mem.tokenizeScalar(u8, line, ' ');
    }

    var ops_it = &iterators[iterators.len - 1];
    var ops_buf: [8]u8 = undefined; // more than enough for example
    var idx: usize = 0;
    while (ops_it.next()) |op| {
        if (op.len == 0) continue;
        try std.testing.expectEqual(@as(usize, 1), op.len);
        ops_buf[idx] = op[0];
        idx += 1;
    }

    try std.testing.expectEqual(@as(usize, 4), idx);
    try std.testing.expectEqualStrings("*+*+", ops_buf[0..idx]);
}
