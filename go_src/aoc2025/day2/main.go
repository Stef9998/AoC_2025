package main

import (
	"fmt"
	"go_src/util"
	"os"
	"strconv"
	"strings"
)

const day = 2

func main() {
	//const textfile = "example"
	const textfile = "input"
	lines, err := util.ReadLines(day, textfile+".txt")
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	parsed := parse(lines)
	result := -1
	result = p1(parsed)
	fmt.Println("Part 1:", result)
	result = p2(parsed)
	fmt.Println("Part 2:", result)

}

type Range struct {
	lower int
	upper int
}

func p2(ranges []Range) int {
	sum := 0
	for _, rang := range ranges {
		sum += p2calc(rang)
	}
	return sum
}
func p2calc(rang Range) int {
	foundNums := 0
	for i := rang.lower; i <= rang.upper; i++ {
		iStr := strconv.Itoa(i)
		iLen := len(iStr)
		foundNums += isNumberPattern(iLen, iStr, i)
	}
	return foundNums
}

func isNumberPattern(iLen int, iStr string, i int) int {
	for repLeng := 1; repLeng <= iLen/2; repLeng++ {
		if iLen%repLeng == 0 {
			numOfRep := iLen / repLeng
			found := true
			firstString := iStr[:repLeng]
			for j := 1; j < numOfRep; j++ {
				if firstString != iStr[j*repLeng:(j+1)*repLeng] {
					found = false
					break
				}
			}
			if found {
				return i
			}
		}
	}
	return 0
}

func p1(ranges []Range) int {
	sum := 0
	for _, rang := range ranges {
		rangeFound := p1calc(rang)
		//fmt.Printf("Range %d-%d found sum of numbers: %d\n", rang.lower, rang.upper, rangeFound)
		sum += rangeFound
	}
	return sum
}
func p1calc(rang Range) int {
	foundNums := 0
	for i := rang.lower; i <= rang.upper; i++ {
		iStr := strconv.Itoa(i)
		iLen := len(iStr)
		if (iLen % 2) == 1 {
			continue
		}
		if iStr[:iLen/2] == iStr[iLen/2:] {
			foundNums += i
		}
	}
	return foundNums
}

func parse(lines []string) []Range {
	line := lines[0]
	temp := strings.Split(line, ",")
	ret := make([]Range, len(temp))

	for i, rang := range temp {
		ret[i] = parseLine(rang)
	}
	return ret
}
func parseLine(rang string) Range {
	strNums := strings.Split(rang, "-")
	num1, err := strconv.Atoi(strNums[0])
	if err != nil {
		panic("while parsing input")
	}
	num2, err := strconv.Atoi(strNums[1])
	if err != nil {
		panic("while parsing input")
	}
	return Range{num1, num2}
}
