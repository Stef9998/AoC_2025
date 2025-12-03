package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadLines(day int, filename string) ([]string, error) {
	path := filepath.Join("..", "..", "input", "day"+fmt.Sprint(day), filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err3 := scanner.Err(); err3 != nil {
		return nil, err

	}
	return lines, nil
}

func ReadChars(day int, filename string) ([][]rune, error) {
	lines, err := ReadLines(day, filename)
	if err != nil {
		return nil, err
	}
	return LinesToChars(lines)
}

func LinesToChars(lines []string) ([][]rune, error) {
	mapval := make([][]rune, len(lines))
	for i, line := range lines {
		mapval[i] = make([]rune, len(line))
		for j, char := range line {
			mapval[i][j] = char
		}
	}
	return mapval, nil
}

func SplitAtEmptyLine(lines []string) ([]string, []string, error) {
	for i, line := range lines {
		if line == "" {
			return lines[:i], lines[i+1:], nil
		}
	}
	return nil, nil, fmt.Errorf("no empty line found")
}

func SeparatedNumbers(line string, sep string) ([]int, error) {
	values := strings.Split(line, sep)
	if len(values) == 1 {
		return nil, errors.New("string needs to contain separator")
	}
	numbers := make([]int, len(values))
	for i, value := range values {
		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		numbers[i] = number
	}
	return numbers, nil
}
