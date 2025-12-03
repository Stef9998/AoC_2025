package util

// SumInts returns the sum of a slice of ints.
func SumInts(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}

// Count returns the number of elements in arr for which pred returns true.
func Count[T any](arr []T, pred func(T) bool) int {
	count := 0
	for _, value := range arr {
		if pred(value) {
			count++
		}
	}
	return count
}
