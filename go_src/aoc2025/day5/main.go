package main

import (
	"fmt"
	"go_src/util"
	"os"
	"slices"
	"strconv"
	"strings"
)

const day = 5

func main() {
	//const textfile = "example"
	const textfile = "input"
	lines, err := util.ReadLines(day, textfile+".txt")
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	freshRange, ingredientNums, err := util.SplitAtEmptyLine(lines)
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	parsed := parse(freshRange, ingredientNums)
	result := -1
	result = p1(parsed)
	fmt.Println("Part 1:", result)
	if textfile == "input" && result != 739 {
		panic("wrong result for part 1")
	}
	result = p2(parsed)
	fmt.Println("Part 2:", result)
	if textfile == "input" && result != 344486348901788 {
		panic("wrong result for part 2")
	}
}

type ran struct {
	b int
	e int
}
type T struct {
	fresh      []ran
	ingredient []int
}

func p2(input T) int {
	llBegin := p2createLL(input.fresh)
	p2llMergeRanges(llBegin)
	return p2sumRanges(llBegin)
}

func p2sumRanges(llBegin *D5ll) int {
	sum := 0
	cur := llBegin
	for cur != nil {
		sum += cur.e - cur.b + 1
		cur = cur.next
	}
	return sum
}

func p2llMergeRanges(llBegin *D5ll) {
	cur1 := llBegin
	for cur1.next != nil {
		merged := cur1.mergeNext()
		if !merged {
			cur1 = cur1.next
		}
	}
}

func p2createLL(freshRanges []ran) *D5ll {
	slices.SortFunc(freshRanges, func(a, b ran) int { return a.b - b.b })

	llInit := NewLL(freshRanges[len(freshRanges)-1])
	currLeft := &llInit
	for i := len(freshRanges) - 2; i >= 0; i-- {
		prevLL := NewLL(freshRanges[i])
		prevLL.next = currLeft
		currLeft = &prevLL
	}

	return currLeft
}

func p1(input T) int {
	freshCount := 0
	for _, val := range input.ingredient {
		for _, rang := range input.fresh {
			if val >= rang.b && val <= rang.e {
				freshCount++
				break
			}
		}
	}
	return freshCount
}

func parse(fre, inc []string) T {
	retFresh := make([]ran, len(fre))
	retInc := make([]int, len(inc))
	for i, freshRange := range fre {
		retFresh[i] = parseFresh(freshRange)
	}
	for i, ingredientNum := range inc {
		retInc[i] = parseIngredient(ingredientNum)
	}
	return T{retFresh, retInc}
}
func parseFresh(line string) ran {
	ret := strings.Split(line, "-")
	num1, err := strconv.Atoi(ret[0])
	if err != nil {
		panic("error parsing input")
	}
	num2, err := strconv.Atoi(ret[1])
	if err != nil {
		panic("error parsing input")
	}
	return ran{num1, num2}
}
func parseIngredient(line string) int {
	num, err := strconv.Atoi(line)
	if err != nil {
		panic("error parsing input")
	}
	return num
}
