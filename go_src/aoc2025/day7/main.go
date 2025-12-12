package main

import (
	"go_src/util"
)

const (
	day = 7
	//example = true
	example = false
)

func readIn[L string | []rune](reader func(day int, filename string) ([]L, error)) []L {
	return util.ReadIn(reader, day, example)
}
func main() {
	lines := readIn(util.ReadChars)
	parsed := parse(lines)
	result := util.Run(p1, parsed, 1)
	if !example {
		if result != 1649 {
			panic("wrong result")
		}
	}
	result = util.Run(p2, parsed, 2)
	if !example {
		if result != 16937871060075 {
			panic("wrong result")
		}
	}

}

type T struct {
	input [][]rune
	beam  []bool
	s     int
	e     int
}

func p2(input T) int {
	beam := make([]int, len(input.beam))
	beam[input.s] = 1
	splits := 1
	for row := 2; row < len(input.input); row += 2 {
		for col := input.s; col < input.e; col++ {
			if beam[col] > 0 && input.input[row][col] == '^' {
				if col == input.s {
					input.s--
				}
				if col == input.e-1 {
					input.e++
				}
				splits += beam[col]
				beam[col-1] += beam[col]
				beam[col+1] += beam[col]
				beam[col] = 0
			}
		}
	}
	return splits
}

func p1(input T) int {
	splits := 0
	for row := 2; row < len(input.input); row += 2 {
		for col := input.s; col < input.e; col++ {
			if input.beam[col] == true && input.input[row][col] == '^' {
				if col == input.s {
					input.s--
				}
				if col == input.e-1 {
					input.e++
				}
				input.beam[col] = false
				input.beam[col-1] = true
				input.beam[col+1] = true
				splits++
			}
		}
	}
	return splits
}

func parse(lines [][]rune) T {
	beam := make([]bool, len(lines[0]))
	start := 0
	end := 0
	for i, char := range lines[0] {
		if char == 'S' {
			beam[i] = true
			start = i
			end = i + 1
			break
		}
	}
	return T{lines, beam, start, end}
}
