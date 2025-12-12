package main

import (
	"fmt"
	"go_src/util"
	"os"
	"strings"
	"time"
)

const (
	day = 12
	//example = true
	example = false
)

func readIn[L string | []rune](reader func(day int, filename string) ([]L, error)) []L {
	var textfile string
	if example {
		textfile = "example"
	} else {
		textfile = "input"
	}
	lines, err := reader(day, textfile+".txt")
	if err != nil {
		fmt.Printf("Error when reading in file %s from day %d\n%s", textfile, day, err)
		os.Exit(-1)
	}
	return lines
}
func run[PT any](code func(PT) int, parsed PT, partNo int) int {
	start := time.Now()
	result := code(parsed)
	duration := time.Since(start)
	fmt.Printf("Part %d: %d  time: %dms\n", partNo, result, duration.Milliseconds())
	return result
}
func main() {
	lines := readIn(util.ReadLines)
	parsed := parse(lines)
	var result int
	result = run(p1, parsed, 1)
	//if !example {
	//if result !=  {
	//	fmt.Println("Wrong result of part 1")
	//}
	//}
	//parsed = parse(lines)
	//result = run(p2, parsed, 2)
	//if !example {
	//	if result !=  {
	//		fmt.Println("wrong result")
	//	}
	//}
	_ = result
}

type presentShape [3][3]bool
type present struct {
	shape    presentShape
	tilesNum int
	//flipChange only true if flipping x- and y-axis does change the shape
	flipChange   bool
	rot90Change  bool
	rot180Change bool
}

// rotate 90 degrees clock-wise
func (ps presentShape) rotate90() presentShape {
	return [3][3]bool{
		{ps[2][0], ps[1][0], ps[0][0]},
		{ps[2][1], ps[1][1], ps[0][1]},
		{ps[2][2], ps[1][2], ps[0][2]},
	}
}

func (ps presentShape) flipX() presentShape {
	return [3][3]bool{
		{ps[0][2], ps[0][1], ps[0][0]},
		{ps[1][2], ps[1][1], ps[1][0]},
		{ps[2][2], ps[2][1], ps[2][0]},
	}
}

type Rectangle struct {
	width  int
	height int
	area   int
}

type christmasTree struct {
	space            Rectangle
	presentTypeCount []int
	numOfPresents    int
}

type T struct {
	presents []present
	trees    []christmasTree
}

func p2(input T) int {
	return -1
}

func p1(input T) int {
	sum := 0
	for _, tree := range input.trees {
		fits := p1calc(input.presents, tree)
		//fmt.Printf("input num: %3d result: %t\n", i, fits)
		if fits {
			sum++
		}
	}
	return sum
}

func p1calc(presents []present, tree christmasTree) bool {
	if belowLowerBound(tree) {
		//fmt.Println("below lower bound")
		return true
	}
	if easyAboveHigherBound(tree) {
		//fmt.Println("above easy higher bound")
		return false
	}
	if realAboveHigherBound(tree, presents) {
		//fmt.Println("above real higher bound")
		return false
	}
	// (Not) TODO implement real fitting check
	fmt.Println("real fitting check not implemented yet")
	return false
}

func belowLowerBound(tree christmasTree) bool {
	presentCount := (tree.space.width / 3) * (tree.space.height / 3)
	if tree.numOfPresents <= presentCount {
		return true
	}
	return false
}

func easyAboveHigherBound(tree christmasTree) bool {
	presentCount := tree.space.area / 9
	if tree.numOfPresents > presentCount {
		return true
	}
	return false
}

func realAboveHigherBound(tree christmasTree, presents []present) bool {
	tilesNum := 0
	for i, count := range tree.presentTypeCount {
		tilesNum += presents[i].tilesNum * count
	}
	if tilesNum > tree.space.area {
		return true
	}
	return false
}

func parse(lines []string) T {
	frst, err := util.SplitAtAllEmptyLines(lines)
	if err != nil {
		panic("error parsing (splitting at empty lines)")
	}
	presentsAsLines := frst[:len(frst)-1]
	presents := make([]present, len(presentsAsLines))
	for i, presentLines := range presentsAsLines {
		presents[i] = parsePresent(i, presentLines)
	}
	treeLines := frst[len(frst)-1]
	trees := make([]christmasTree, len(treeLines))
	for i, line := range treeLines {
		trees[i] = parseTree(line)
	}
	return T{
		presents: presents,
		trees:    trees,
	}
}
func parsePresent(i int, lines []string) present {
	if lines[0] != fmt.Sprintf("%d:", i) {
		panic(fmt.Errorf("error parsing tree numbers: %d <-> %s", i, string(lines[0][0])))
	}
	var shape presentShape
	tileNum := 0
	for j := 0; j < 3; j++ {
		for k := 0; k < 3; k++ {
			char := lines[j+1][k]
			if char == '#' {
				shape[j][k] = true
				tileNum++
			} else if char == '.' {
				shape[j][k] = false
			} else {
				panic("error parsing present")
			}
		}
	}
	rot90 := shape.rotate90()
	rot90change := true
	if shape == rot90 {
		rot90change = false
	}
	rot180 := rot90.rotate90()
	rot180change := true
	if rot90change == false {
		rot180change = false
	} else {
		if shape == rot180 {
			rot180change = false
		}
	}
	flipChange := false
	xFlip := shape.flipX()
	if !(shape == xFlip) {
		yFlip := rot90.flipX()
		if !(rot90 == yFlip) {
			flipChange = true
		}
	}
	//TODO see if the rot/flip stuff is correctly determined
	return present{
		shape:        shape,
		tilesNum:     tileNum,
		flipChange:   flipChange,
		rot90Change:  rot90change,
		rot180Change: rot180change,
	}
}
func parseTree(line string) christmasTree {
	frst := strings.Split(line, ": ")

	dim := strings.Split(frst[0], "x")
	width := util.PanicAtoi(dim[0])
	height := util.PanicAtoi(dim[1])
	rec := Rectangle{
		width:  width,
		height: height,
		area:   width * height,
	}

	numStrs := strings.Fields(frst[1])
	nums := make([]int, len(numStrs))
	numOfPres := 0
	for i, numStr := range numStrs {
		nums[i] = util.PanicAtoi(numStr)
		numOfPres += nums[i]
	}

	return christmasTree{
		space:            rec,
		presentTypeCount: nums,
		numOfPresents:    numOfPres,
	}
}
