package main

type Coordinate struct {
	x int
	y int
	z int
}

func Distance(a, b Coordinate) int {
	return (a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y) + (a.z-b.z)*(a.z-b.z)
}
