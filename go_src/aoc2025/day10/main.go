package main

import (
	"fmt"
	"go_src/util"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const day = 10

//const example = true

const example = false

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
	if !example {
		if result != 417 {
			fmt.Println("Wrong result of part 1")
		}
	}
	parsed = parse(lines)
	//result = run(p2, parsed, 2)
	//result = run(p2_par, parsed, 2)
	result = run(p2_ilp, parsed, 2)
	if !example {
		if result != 16765 {
			fmt.Println("wrong result")
		}
	}
	_ = result
}

type p1conMach struct {
	lights  []bool
	buttons [][]int
}

func (mach *p1conMach) pressButton(i int) {
	butChange := mach.buttons[i]
	for _, val := range butChange {
		mach.lights[val] = !mach.lights[val]
	}
}

type p2conMach struct {
	joltage []int
	buttons [][]int
}

func (mach *p2conMach) pressButton(i int) {
	butChange := mach.buttons[i]
	for _, val := range butChange {
		mach.joltage[val]++
	}
}
func (mach *p2conMach) backtrackButtonPress(i int) {
	butChange := mach.buttons[i]
	for _, val := range butChange {
		mach.joltage[val]--
	}
}

type machine struct {
	lights  []bool
	buttons [][]int
	joltage []int
}

func p2_ilp(input []machine) int {
	sum := 0
	for i, mach := range input {
		fewestPresses := p2_ilp_calc(i, mach)
		//fmt.Printf("%3d %d\n", i, fewestPresses)
		sum += fewestPresses
	}
	return sum
}
func p2_ilp_calc(i int, mach machine) int {
	path := "day10/ilp-files/"
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}
	inputFile := fmt.Sprintf("%smachine-%03d.lp", path, i)
	create_ilp(inputFile, mach)
	outputFile := fmt.Sprintf("%s%03d.out", path, i)
	cmd := exec.Command("glpsol", "--lp", inputFile, "-o", outputFile)
	err := cmd.Run()
	if err != nil {
		panic("couldn't call exec command")
	}
	data, err := os.ReadFile(outputFile)
	if err != nil {
		panic("couln't read in output file")
	}
	solution := get_ilp_val(string(data))
	_ = os.Remove(inputFile)
	_ = os.Remove(outputFile)
	return solution
}

func get_ilp_val(data string) int {
	for _, line := range strings.Split(data, "\n") {
		//line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Objective:") {
			var val int
			_, err := fmt.Sscanf(line, "Objective:  obj = %d (MINimum)", &val)
			if err == nil {
				return val
			}
		}
	}
	panic("Couldn't find solution")
}

func create_ilp(fileName string, mach machine) {
	targets := mach.joltage
	buttons := mach.buttons

	nCounters := len(targets)
	nButtons := len(buttons)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Objective
	fmt.Fprint(f, "Minimize\n  obj: ")
	for i := 0; i < nButtons; i++ {
		if i > 0 {
			fmt.Fprint(f, " + ")
		}
		fmt.Fprintf(f, "x%d", i)
	}
	fmt.Fprintln(f)

	// Constraints
	fmt.Fprintln(f, "Subject To")
	for c := 0; c < nCounters; c++ {
		fmt.Fprintf(f, "  c%d:", c)
		for b := 0; b < nButtons; b++ {
			val := 0
			for _, counter := range buttons[b] {
				if counter == c {
					val = 1
					break
				}
			}
			if b > 0 {
				fmt.Fprint(f, " + ")
			}
			//fmt.Fprintf(f, "%d*x%d", val, b)
			fmt.Fprintf(f, "%d x%d", val, b)
		}
		fmt.Fprintf(f, " = %d\n", targets[c])
	}

	// Bounds
	fmt.Fprintln(f, "Bounds")
	for i := 0; i < nButtons; i++ {
		fmt.Fprintf(f, "  x%d >= 0\n", i)
	}

	// Integer variables
	fmt.Fprintln(f, "Generals")
	for i := 0; i < nButtons; i++ {
		fmt.Fprintf(f, "  x%d ", i)
	}
	fmt.Fprintln(f, "\nEnd")
}

func p2_par(input []machine) int {
	var mu sync.Mutex
	var wg sync.WaitGroup
	sum := 0
	for i, mach := range input {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fewestPresses := p2calc(mach)
			fmt.Printf("%3d %d\n", i, fewestPresses)
			mu.Lock()
			sum += fewestPresses
			mu.Unlock()
		}()
	}
	wg.Wait()
	return sum
}

func p2(input []machine) int {
	sum := 0
	for i, mach := range input {
		fewestPresses := p2calc(mach)
		fmt.Printf("%3d %d\n", i, fewestPresses)
		sum += fewestPresses
	}
	return sum
}

func p2calc(goal machine) int {
	mach := p2conMach{
		joltage: make([]int, len(goal.joltage)),
		buttons: goal.buttons,
	}
	fewestPresses := recSolver(mach, goal.joltage)

	if fewestPresses == math.MaxInt {
		fmt.Println(goal, mach)
		panic("couldn't turn the light on")
	}
	return fewestPresses
}
func recSolver(mach p2conMach, goal []int) int {
	fewestPresses := math.MaxInt
	var rec func(int, int)
	rec = func(i, depth int) {
		if depth >= fewestPresses {
			return
		}
		mach.pressButton(i)
		for _, jolt := range mach.buttons[i] {
			if mach.joltage[jolt] > goal[jolt] {
				mach.backtrackButtonPress(i)
				return
			}
		}
		if slMatch(mach.joltage, goal) {
			fewestPresses = depth
			mach.backtrackButtonPress(i)
			return
		}
		for j := i; j < len(mach.buttons); j++ {
			rec(j, depth+1)
		}
		mach.backtrackButtonPress(i)
		return
	}
	for i, _ := range mach.buttons {
		rec(i, 1)
	}
	return fewestPresses
}

func p1(input []machine) int {
	sum := 0
	for _, mach := range input {
		fewestPresses := p1calc(mach)
		//fmt.Printf("%3d %d\n", i, fewestPresses)
		sum += fewestPresses
	}
	return sum
}

func p1calc(goal machine) int {
	mach := p1conMach{
		lights:  make([]bool, len(goal.lights)),
		buttons: goal.buttons,
	}
	fewestPresses := len(goal.buttons) + 1
	var rec func(int, int)
	rec = func(i, depth int) {
		if depth >= fewestPresses {
			return
		}
		mach.pressButton(i)
		if slMatch(mach.lights, goal.lights) {
			fewestPresses = depth
			mach.pressButton(i)
			return
		}
		for j := i + 1; j < len(mach.buttons); j++ {
			rec(j, depth+1)
		}
		mach.pressButton(i)
		return
	}
	for i, _ := range goal.buttons {
		rec(i, 1)
	}
	if fewestPresses == len(mach.buttons)+1 {
		fmt.Println(goal, mach)
		panic("couldn't turn the light on")
	}
	return fewestPresses
}

func slMatch[T comparable](sl1, sl2 []T) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i, val := range sl1 {
		if val != sl2[i] {
			return false
		}
	}
	return true
}

func parse(lines []string) []machine {
	ret := make([]machine, len(lines))
	for i, line := range lines {
		ret[i] = parseLine(line)
	}
	return ret
}
func parseLine(line string) machine {
	stuff := strings.Fields(line)
	lights := make([]bool, len(stuff[0])-2)
	if !(stuff[0][0] == '[' && stuff[0][len(stuff[0])-1] == ']') {
		panic("error parsing input 1")
	}
	for i, j := 1, 0; i < len(stuff[0])-1; i, j = i+1, j+1 {
		if stuff[0][i] == '.' {
			lights[j] = false
		} else if stuff[0][i] == '#' {
			lights[j] = true
		} else {
			panic("error parsing input 2")
		}
	}

	buttons := make([][]int, len(stuff)-2)
	for i, buttonStr := range stuff[1 : len(stuff)-1] {
		buttonArr := strings.Split(buttonStr[1:len(buttonStr)-1], ",")
		button := make([]int, len(buttonArr))
		for j, val := range buttonArr {
			num, err := strconv.Atoi(val)
			if err != nil {
				panic("error parsing input (button)")
			}
			button[j] = num
		}
		buttons[i] = button
	}

	joltageStr := strings.Split(stuff[len(stuff)-1][1:len(stuff[len(stuff)-1])-1], ",")
	joltage := make([]int, len(joltageStr))
	for i, val := range joltageStr {
		num, err := strconv.Atoi(val)
		if err != nil {
			panic("error parsing input (joltage)")
		}
		joltage[i] = num
	}

	return machine{
		lights:  lights,
		buttons: buttons,
		joltage: joltage,
	}
}
