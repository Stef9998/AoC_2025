package main

import (
	"fmt"
	"go_src/util"
	"slices"
	"strconv"
	"strings"
)

type Coordinate = util.Coordinate

const (
	day = 9
	//const example = true
	example = false
)

func readIn[L string | []rune](reader func(day int, filename string) ([]L, error)) []L {
	return util.ReadIn(reader, day, example)
}
func main() {
	lines := readIn(util.ReadLines)
	var result int
	result = util.Run(p1, parse(lines), 1)
	if !example {
		if result != 4763509452 {
			fmt.Println("wrong result")
		}
	}
	result = util.Run(p2, parse(lines), 2)
	if !example {
		if result != 1516897893 {
			fmt.Println("wrong result")
		}
	}
}

type Rectangle struct {
	no   Coordinate
	so   Coordinate
	sw   Coordinate
	nw   Coordinate
	area int
}

func NewRectangle(a, b Coordinate) Rectangle {
	var no, so, sw, nw Coordinate
	if a.X >= b.X {
		if a.Y >= b.Y {
			no = Coordinate{a.X, b.Y}
			so = a
			sw = Coordinate{b.X, a.Y}
			nw = b
		} else {
			no = a
			so = Coordinate{a.X, b.Y}
			sw = b
			nw = Coordinate{b.X, a.Y}
		}
	} else {
		if a.Y >= b.Y {
			no = b
			so = Coordinate{b.X, a.Y}
			sw = a
			nw = Coordinate{a.X, b.Y}
		} else {
			no = Coordinate{b.X, a.Y}
			so = b
			sw = Coordinate{a.X, b.Y}
			nw = a
		}
	}
	return Rectangle{no, so, sw, nw, recArea(so, nw)}
}

func (rec Rectangle) inside(point Coordinate) (int, int) {
	x := 0
	if point.X <= rec.nw.X {
		x = -1
	} else if point.X >= rec.so.X {
		x = 1
	}
	y := 0
	if point.Y <= rec.nw.Y {
		y = -1
	} else if point.Y >= rec.so.Y {
		y = 1
	}
	return x, y
}

func (rec Rectangle) cuts(p1, p2 Coordinate) bool {
	p1x, p1y := rec.inside(p1)
	p2x, p2y := rec.inside(p2)

	if p1.X == p2.X {
		//vertical line
		if p1x != 0 {
			return false
		}
		//x inside
		if p1y != p2y {
			return true
		}
		if p1y == 0 {
			return true
		}
		return false
	} else if p1.Y == p2.Y {
		//horizontal line
		if p1y != 0 {
			return false
		}
		//y inside
		if p1x != p2x {
			return true
		}
		if p1x == 0 {
			return true
		}
		return false
	} else {
		// diagonal line
		panic("should not be possible with input data")
	}
}

type line struct {
	p1 Coordinate
	p2 Coordinate
}

func p2(input []Coordinate) int {
	lines := genLines(input)
	maxArea := 0
	for i := 0; i < len(input)-1; i++ {
	newrec:
		for j := i + 1; j < len(input); j++ {
			rec := NewRectangle(input[i], input[j])
			if rec.area <= maxArea {
				continue
			}
			for _, l := range lines {
				if rec.cuts(l.p1, l.p2) {
					continue newrec
				}
			}
			if rec.area > maxArea {
				maxArea = rec.area
			}
		}
	}
	return maxArea
}

func genLines(input []Coordinate) []line {
	lines := make([]line, len(input))
	for i := 0; i < len(input)-1; i++ {
		lines[i] = line{input[i], input[i+1]}
	}
	lines[len(input)-1] = line{input[len(input)-1], input[0]}
	return lines
}

func pos(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}
func p1(input []Coordinate) int {
	slices.SortFunc(input, func(a, b Coordinate) int {
		return a.X - b.X
	})
	half := len(input) / 2
	slices.SortFunc(input[:half], func(a, b Coordinate) int {
		return a.Y - b.Y
	})
	slices.SortFunc(input[half:], func(a, b Coordinate) int {
		return a.Y - b.Y
	})
	quart := half / 2
	nw := input[:quart]
	sw := input[quart:half]
	no := input[half : half+quart]
	so := input[half+quart:]

	maxArea := 0
	for _, c1 := range nw {
		for _, c2 := range so {
			area := recArea(c2, c1)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	for _, c1 := range no {
		for _, c2 := range sw {
			area := recArea(c1, c2)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

func recArea(so Coordinate, nw Coordinate) int {
	return (so.X - nw.X + 1) * (so.Y - nw.Y + 1)
}

func parse(lines []string) []Coordinate {
	ret := make([]Coordinate, len(lines))
	for i, line := range lines {
		ret[i] = parseLine(line)
	}
	return ret
}
func parseLine(line string) Coordinate {
	numStr := strings.Split(line, ",")
	nums := sliceAtoi(numStr)
	return Coordinate{
		X: nums[0],
		Y: nums[1],
	}
}
func sliceAtoi(in []string) []int {
	out := make([]int, len(in))
	for i, str := range in {
		num, err := strconv.Atoi(str)
		if err != nil {
			panic("error while parsing")
		}
		out[i] = num
	}
	return out
}
