package main

import (
	"go_src/util"
	"slices"
	"strconv"
	"strings"
)

const (
	day = 8
	// const example = true
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
		if result != 115885 {
			panic("Wrong result of part 1")
		}
	}
	if !example {
		result = util.Run(p2, parsed, 2)
		if result != 274150525 {
			panic("Wrong result of part 2")
		}
	}
}

type edge struct {
	a      Coordinate
	b      Coordinate
	length int
}

type T = Coordinate

func p2(input []T) int {
	edges := computeEdges(input)
	set := util.NewSet[Coordinate]()
	set.Add(edges[0].a)
	set.Add(edges[0].b)
	highestI := 0
	for {
		added := false
		for i := 1; i < len(edges); i++ {
			e := edges[i]
			if set.Contains(e.a) {
				if !set.Contains(e.b) {
					set.Add(e.b)
					added = true
					if i > highestI {
						highestI = i
					}
					break
				}
			} else if set.Contains(e.b) {
				set.Add(e.a)
				added = true
				if i > highestI {
					highestI = i
				}
				break
			}
		}
		if !added {
			panic("should be possible to iterate through every edge and not having everything connected")
		}
		if set.Size() == len(input) {
			break
		}
	}
	lastEdge := edges[highestI]
	return lastEdge.a.x * lastEdge.b.x
}

func p1(input []T) int {
	edges := computeNEdges(input)
	threeBiggest := combineSets(edges)
	sum := 1
	for _, val := range threeBiggest {
		circuitSize := (*val).Size()
		//fmt.Println(circuitSize)
		sum *= circuitSize
	}
	return sum
}

func computeNEdges(input []T) []edge {
	dist := computeEdges(input)
	N := 1000
	if example {
		N = 10
	}
	edges := dist[:N]
	return edges
}

func computeEdges(input []T) []edge {
	dist := make([]edge, 0, len(input)*(len(input)-1)/2)
	for i := 0; i < len(input); i++ {
		c1 := input[i]
		for j := i + 1; j < len(input); j++ {
			c2 := input[j]
			e := edge{c1, c2, Distance(c1, c2)}
			dist = append(dist, e)
		}
	}
	slices.SortFunc(dist, func(e1, e2 edge) int { return e1.length - e2.length })
	return dist
}
func combineSets(edges []edge) []*util.Set[Coordinate] {
	coordSets := make(map[Coordinate]*util.Set[Coordinate])
	allCoords := make(util.Set[Coordinate])
	for _, e := range edges {
		allCoords.Add(e.a)
		allCoords.Add(e.b)
		val1, ok1 := coordSets[e.a]
		val2, ok2 := coordSets[e.b]
		if !ok1 && !ok2 {
			newSet := util.NewSet[Coordinate]()
			newSet.Add(e.a)
			newSet.Add(e.b)
			coordSets[e.a] = &newSet
			coordSets[e.b] = &newSet
			continue
		}
		if ok1 && !ok2 {
			(*val1).Add(e.b)
			coordSets[e.b] = val1
			continue
		}
		if !ok1 && ok2 {
			(*val2).Add(e.a)
			coordSets[e.a] = val2
			continue
		}
		un := (*val1).Union(*val2)
		for val, _ := range un {
			coordSets[val] = &un
		}
		//coordSets[e.a] = &un
		//coordSets[e.b] = &un
	}

	combSet := util.NewSet[*util.Set[Coordinate]]()
	for coord, v := range coordSets {
		if allCoords.Contains(coord) {
			combSet.Add(v)
			allCoords = allCoords.Difference(*v)
		}
	}
	combArr := make([]*util.Set[Coordinate], 0, combSet.Size())
	for circuit, _ := range combSet {
		combArr = append(combArr, circuit)
	}
	slices.SortFunc(combArr,
		func(a, b *util.Set[Coordinate]) int {
			return (*b).Size() - (*a).Size()
		})
	return combArr[:3]
}

func parse(lines []string) []T {
	ret := make([]T, len(lines))
	for i, line := range lines {
		ret[i] = parseLine(line)
	}
	return ret
}
func parseLine(line string) T {
	spl := strings.Split(line, ",")
	nums := sliceAtoi(spl)
	return Coordinate{
		x: nums[0],
		y: nums[1],
		z: nums[2],
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
