package main

import (
	"go_src/util"
	"strconv"
)

const (
	day = 1
	//example = true
	example = false
)

func readIn[L string | []rune](reader func(day int, filename string) ([]L, error)) []L {
	return util.ReadIn(reader, day, example)
}
func main() {
	lines := readIn(util.ReadLines)
	parsed := parse(lines)
	result := util.Run(p1, parsed, 1)
	if !example {
		if result != 997 {
			panic("wrong result")
		}
	}
	result = util.Run(p2, parsed, 2)
	if !example {
		if result != 5978 {
			panic("wrong result")
		}
	}
}

func p2(lines []int) int {
	dial := 50
	timesZero := 0
	for _, val := range lines {
		newDial, zeroHit := p2calc(dial, val)
		//fmt.Printf("From %2d with %3d to %2d. Zero Hit %1d times.\n", dial, val, newDial, zeroHit)
		timesZero += zeroHit
		dial = newDial
	}
	return timesZero
}
func p2calc(dial int, rot int) (int, int) {
	zeroHit := 0
	hundreds := rot / 100
	if hundreds > 0 {
		zeroHit += hundreds
	} else {
		zeroHit -= hundreds
	}
	rot = rot - 100*hundreds
	// now rot should be -100 < rot < 100
	if rot == 0 {
		return dial, zeroHit
	}
	newDial := dial + rot
	if newDial > 0 && newDial < 100 {
		return newDial, zeroHit
	}
	if dial == 0 {
		zeroHit--
	}
	if newDial < 0 {
		return newDial + 100, zeroHit + 1
	}
	if newDial >= 100 {
		return newDial - 100, zeroHit + 1
	}
	return 0, zeroHit + 1
}

func p1(lines []int) int {
	dial := 50
	timesZero := 0
	for _, val := range lines {
		dial = (dial + val) % 100
		if dial == 0 {
			timesZero++
		}
	}
	return timesZero
}

func parse(lines []string) []int {
	ret := make([]int, len(lines))
	parseLine := func(l string) int {
		switch l[0] {
		case 'L':
			num, err := strconv.Atoi(l[1:])
			if err != nil {
				panic(err)
			}
			return -num
		case 'R':
			num, err := strconv.Atoi(l[1:])
			if err != nil {
				panic(err)
			}
			return num
		default:
			panic("First character of input line has to be L or R!")
		}
	}
	for i, line := range lines {
		ret[i] = parseLine(line)
	}
	return ret
}
