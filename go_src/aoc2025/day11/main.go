package main

import (
	"fmt"
	"go_src/util"
	"os"
	"strings"
	"time"

	"github.com/dominikbraun/graph"
)

const (
	day = 11
	//example = true
	example = false
)

func readIn[L string | []rune](reader func(day int, filename string) ([]L, error), textfileName string) []L {
	var textfile string
	if textfileName == "" {
		if example {
			textfile = "example"
		} else {
			textfile = "input"
		}
	} else {
		textfile = textfileName
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
	var parsed PT
	var lines []string
	var result int
	if example {
		lines = readIn(util.ReadLines, "example1")
	} else {
		lines = readIn(util.ReadLines, "")
	}
	parsed = parse(lines)
	result = run(p1, parsed, 1)
	//if !example {
	//if result !=  {
	//	fmt.Println("Wrong result of part 1")
	//}
	//}
	if example {
		lines = readIn(util.ReadLines, "example2")
	} else {
		lines = readIn(util.ReadLines, "")
	}
	parsed = parse(lines)
	result = run(p2, parsed, 2)
	if !example {
		if result != 517315308154944 {
			fmt.Println("wrong result")
		}
	}
	_ = result
}

type PT = graph.Graph[string, string]

func p2(input PT) int {
	count := 1

	s, v1, v2, e := p2order(input)
	fmt.Println(v1, "->", v2)

	adj, _ := input.AdjacencyMap()
	part := pathCount(s, v1, adj)
	count *= part
	//fmt.Println(part)
	part = pathCount(v1, v2, adj)
	count *= part
	//fmt.Println(part)
	part = pathCount(v2, e, adj)
	count *= part
	//fmt.Println(part)

	return count
}
func p2order(input PT) (string, string, string, string) {
	var fftToDac bool
	path, err := graph.ShortestPath(input, "fft", "dac")
	fmt.Println(path)
	if err == nil {
		fftToDac = true
	} else {
		fftToDac = false
	}
	var s, v1, v2, e string
	s = "svr"
	e = "out"
	if fftToDac {
		v1 = "fft"
		v2 = "dac"
	} else {
		v1 = "dac"
		v2 = "fft"
	}
	return s, v1, v2, e
}
func pathCount[K string](s, e K, edgeMap map[K]map[K]graph.Edge[K]) int {
	numPath := make(map[K]int)
	numPath[e] = 1 // TODO see if 0

	var rec func(v K) int
	rec = func(v K) int {
		nextVs := edgeMap[v]
		if len(nextVs) == 0 {
			return 0
		}
		sum := 0
		for nextV, _ := range nextVs {
			val, ok := numPath[nextV]
			if ok {
				sum += val
			} else {
				nextVpath := rec(nextV)
				numPath[nextV] = nextVpath
				sum += nextVpath
			}
		}
		return sum
	}
	return rec(s)
}
func p1(input PT) int {
	paths, err := graph.AllPathsBetween(input, "you", "out")
	if err != nil {
		return -1
	}
	return len(paths)
}

func parse(lines []string) PT {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	for _, line := range lines {
		from, tos := parseLine(line)
		_ = g.AddVertex(from)
		//if err != nil {
		//	fmt.Println(err)
		//	panic("error parsing")
		//}
		for _, to := range tos {
			_ = g.AddVertex(to)
			err := g.AddEdge(from, to)
			if err != nil {
				fmt.Println(err)
				panic("error parsing")
			}
		}
	}
	return g
}
func parseLine(line string) (string, []string) {
	frst := strings.Split(line, ":")
	fromNode := frst[0]
	toNodes := strings.Fields(frst[1])
	return fromNode, toNodes
}
