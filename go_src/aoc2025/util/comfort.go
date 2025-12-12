package util

import (
	"fmt"
	"strconv"
	"time"
)

func Stopwatch[ResType any](function func() ResType) (ResType, time.Duration) {
	start := time.Now()
	result := function()
	duration := time.Since(start)
	return result, duration
}

func Printms(start time.Time) {
	PrintTime(start, "ms")
}

func PrintTime(start time.Time, timeUnit string) {
	duration := time.Since(start)
	var printInt int64 = -1
	var printFloat float64 = -1
	switch timeUnit {
	case "ns":
		printInt = duration.Nanoseconds()
	case "us":
		printInt = duration.Microseconds()
	case "ms":
		printInt = duration.Milliseconds()
	case "s":
		printFloat = duration.Seconds()
	case "min":
		printFloat = duration.Minutes()
	case "h":
		printFloat = duration.Hours()
	}
	if printInt != -1 {
		fmt.Printf("time: %d%s\n", printInt, timeUnit)
		return
	}
	if printFloat != -1 {
		fmt.Printf("time: %f%s\n", printFloat, timeUnit)
		return
	}
	fmt.Println("Wrong time format")
}

func PanicAtoi(in string) (out int) {
	var err error
	out, err = strconv.Atoi(in)
	if err != nil {
		panic(fmt.Errorf("error parsing string: %s\nWith error: %s", in, err))
	}
	return
}
