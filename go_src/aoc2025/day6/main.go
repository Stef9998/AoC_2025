package main

import (
	"fmt"
	"go_src/util"
	"os"
	"strconv"
	"strings"
)

const day = 6

func main() {
	//const textfile = "example"
	const textfile = "input"
	lines, err := util.ReadLines(day, textfile+".txt")
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	parsed := parse(lines)
	var result int
	result = p1(parsed)
	fmt.Println("Part 1:", result)
	if result != 6891729672676 {
		fmt.Println(fmt.Errorf("wrong result"))
	}
	chars, err := util.LinesToChars(lines)
	if err != nil {
		panic("error parsing input")
	}
	result = p2(chars)
	fmt.Println("Part 2:", result)

}

type T struct {
	numbers   [][]int
	operators []string
}

func p2(chars [][]rune) int {
	longestLine := 0
	for _, line := range chars {
		if len(line) > longestLine {
			longestLine = len(line)
		}
	}
	sum := 0
	opRow := len(chars) - 1
	var function func(int) int
	var ret int
	for col := 0; col < longestLine; col++ {
		if col < len(chars[opRow]) {
			switch chars[opRow][col] {
			case '*':
				function = getMulter()
			case '+':
				function = getSummer()
			case ' ':
			default:
				panic("wrong input string parsed")
			}
		}
		digitExisting := false
		num := ""
		for row := 0; row < opRow; row++ {
			if col < len(chars[row]) {
				char := chars[row][col]
				if char >= '0' && char <= '9' {
					digitExisting = true
					num += string(char)
				}
			}
		}
		if digitExisting {
			val, err := strconv.Atoi(num)
			if err != nil {
				panic("error parsing")
			}
			ret = function(val)
		} else {
			sum += ret
		}
	}
	return sum + ret
}

func getMulter() func(int) int {
	mult := 1
	return func(a int) int {
		mult *= a
		return mult
	}
}
func getSummer() func(int) int {
	sum := 0
	return func(a int) int {
		sum += a
		return sum
	}
}

func p1(input T) int {
	sum := 0
	for i, op := range input.operators {
		var function func(int) int
		switch op {
		case "*":
			function = getMulter()
		case "+":
			function = getSummer()
		default:
			panic("wrong input string parsed")
		}
		var col int
		for _, row := range input.numbers {
			col = function(row[i])
		}
		sum += col
	}
	return sum
}

func parse(lines []string) T {
	numbers := make([][]int, len(lines)-1)
	for i := 0; i < len(lines)-1; i++ {
		numbers[i] = parseLine(lines[i])
	}
	op := strings.Fields(lines[len(lines)-1])
	return T{numbers, op}
}

func parseLine(line string) []int {
	parts := strings.Fields(line) // splits on any whitespace
	nums := make([]int, len(parts))

	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			panic("input not parsable")
		}
		nums[i] = n
	}
	return nums
}
