package main

import (
	"fmt"
	"go_src/util"
	"os"
)

const day = 4

func main() {
	//const textfile = "example"
	const textfile = "input"
	input, err := util.ReadChars(day, textfile+".txt")
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	parsed := parse(input)
	//printMap(input)
	//fmt.Println("")
	//printMap(parsed)
	result := -1
	result = p1(parsed)
	fmt.Println("Part 1:", result)
	result = p2(parsed)
	fmt.Println("Part 2:", result)

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
