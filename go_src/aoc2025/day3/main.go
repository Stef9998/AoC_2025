package main

import (
	"fmt"
	"go_src/util"
	"os"
	"strconv"
)

const day = 3

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
	if result != 17109 {
		panic("wrong result")
	}
	result = p2(parsed)
	fmt.Println("Part 2:", result)
	if result != 169347417057382 {
		panic("wrong result")
	}
}

type T string

func p2(banks []T) int {
	sum := 0
	for _, bank := range banks {
		p2res := p2calc(bank)
		fmt.Println(p2res)
		sum += p2res
	}
	return sum
}
func p2calc(bank T) int {
	//highestStillFoundable := '9'
	str := ""
	lastDigitIndex := -1
	for i := 11; i >= 0; i-- {
		highestFound := '0'
		for j := lastDigitIndex + 1; j < len(bank)-i; j++ {
			char := rune(bank[j])
			if char > highestFound {
				highestFound = char
				lastDigitIndex = j
				//if highestFound == highestStillFoundable {
				//	break
				//}
			}
		}
		//highestStillFoundable = highestFound

		str = str + string(highestFound)
	}
	ret, err := strconv.Atoi(str)
	if err != nil {
		panic("not parsable")
	}
	return ret
}

func p1(banks []T) int {
	sum := 0
	for _, bank := range banks {
		p1res := p1calc(bank)
		//fmt.Println(p1res)
		sum += p1res
	}
	return sum
}
func p1calc(bank T) int {
	highestBefore := '0'
	highest := '0'
	second := '0'
	secondAfter := false
	for _, val := range bank {
		if val > highest {
			highestBefore = highest
			second = '0'
			secondAfter = false
			highest = val
		} else if val > second {
			second = val
			secondAfter = true
		}
	}
	ten, err := strconv.Atoi(string(highest))
	if err != nil {
		panic("error parsing")
	}
	one, err := strconv.Atoi(string(second))
	if err != nil {
		panic("error parsing")
	}
	if secondAfter == true {
		return ten*10 + one
	}
	prev, err := strconv.Atoi(string(highestBefore))
	if err != nil {
		panic("error parsing")
	}
	return prev*10 + ten

}

func parse(lines []string) []T {
	ret := make([]T, len(lines))
	for i, line := range lines {
		ret[i] = parseLine(line)
	}
	return ret
}
func parseLine(line string) T {
	return T(line)
}
