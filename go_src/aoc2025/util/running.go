package util

import (
	"fmt"
	"os"
	"time"
)

func ReadIn[L string | []rune](reader func(day int, filename string) ([]L, error), day int, example bool) []L {
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

func Run[PT any](code func(PT) int, parsed PT, partNo int) int {
	start := time.Now()
	result := code(parsed)
	duration := time.Since(start)
	fmt.Printf("Part %d: %d  time: %dms\n", partNo, result, duration.Milliseconds())
	return result
}
