package main

import (
	"fmt"
	"go_src/util"
)

const (
	day = 4
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
		if result != 1626 {
			panic("wrong result")
		}
	}
	result = util.Run(p2, parsed, 2)
	if !example {
		if result != 9173 {
			panic("wrong result")
		}
	}
}
func printMap(input [][]rune) {
	for _, row := range input {
		for _, char := range row {
			fmt.Print(string(char))
		}
		fmt.Println("")
	}
}

func p2(input [][]rune) int {
	accessible := 0
	height := len(input) - 2
	width := len(input[0]) - 2
	for {
		removed := make([]util.Coordinate, 0, 1626)
		accessibleBefore := accessible
		for row := 1; row < height+1; row++ {
			for col := 1; col < width+1; col++ {
				if input[row][col] == '@' {
					surroundedAetts := surrAetts(row, col, input)
					if surroundedAetts < 4 {
						removed = append(removed, util.Coordinate{X: col, Y: row})
						accessible++
					}
				}
			}
		}
		if accessibleBefore == accessible {
			break
		}
		remove(removed, &input)
	}
	return accessible
}

func remove(toRemove []util.Coordinate, input *[][]rune) {
	for _, coord := range toRemove {
		(*input)[coord.Y][coord.X] = '.'
	}
}

func p1(input [][]rune) int {
	accessible := 0
	height := len(input) - 2
	width := len(input[0]) - 2
	for row := 1; row < height+1; row++ {
		for col := 1; col < width+1; col++ {
			if input[row][col] == '@' {
				surroundedAetts := surrAetts(row, col, input)
				if surroundedAetts < 4 {
					accessible++
				}
			}
		}
	}
	return accessible
}
func surrAetts(row, col int, input [][]rune) int {
	aetts := -1
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if input[r][c] == '@' {
				aetts++
			}
		}
	}
	return aetts
}

func parse(lines [][]rune) [][]rune {
	ret := make([][]rune, len(lines)+2)
	width := len(lines[0])
	outerLine := make([]rune, width+2)
	for i := 0; i < width+2; i++ {
		outerLine[i] = '.'
	}
	ret[0] = outerLine
	for i, line := range lines {
		innerLine := make([]rune, width+2)
		innerLine[0] = '.'
		for j := 1; j < len(innerLine)-1; j++ {
			innerLine[j] = line[j-1]
		}
		innerLine[len(innerLine)-1] = '.'
		ret[i+1] = innerLine
	}
	ret[len(ret)-1] = outerLine
	return ret
}
