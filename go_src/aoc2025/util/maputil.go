package util

type Coordinate struct {
	X int
	Y int
}

func (c *Coordinate) Add(other Coordinate) {
	c.X += other.X
	c.Y += other.Y
}

func (c *Coordinate) Sub(other Coordinate) {
	c.X -= other.X
	c.Y -= other.Y
}

func (c *Coordinate) Neighbors() []Coordinate {
	return []Coordinate{
		{X: c.X + 1, Y: c.Y},
		{X: c.X - 1, Y: c.Y},
		{X: c.X, Y: c.Y + 1},
		{X: c.X, Y: c.Y - 1},
	}
}

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func CharToDir(char rune) Direction {
	switch char {
	case '^':
		return Up
	case '>':
		return Right
	case 'v':
		return Down
	case '<':
		return Left
	default:
		panic("not right char")
	}
}

func (d *Direction) ToChar() rune {
	switch *d {
	case Up:
		return '^'
	case Right:
		return '>'
	case Down:
		return 'v'
	case Left:
		return '<'
	default:
		panic("not possible direction")
	}
}

func (d *Direction) Delta() Coordinate {
	switch *d {
	case Up:
		return Coordinate{X: 0, Y: -1}
	case Right:
		return Coordinate{X: 1, Y: 0}
	case Down:
		return Coordinate{X: 0, Y: 1}
	case Left:
		return Coordinate{X: -1, Y: 0}
	default:
		return Coordinate{X: 0, Y: 0}
	}
}

// Rotate returns the direction rotated by 90 degrees.
// If clockwise is true, rotates right; if false, rotates left.
func Rotate(d Direction, clockwise bool) Direction {
	if clockwise {
		return Direction((int(d) + 1) % 4)
	}
	return Direction((int(d) + 3) % 4)
}

// Rotate rotates the direction by 90 degrees in place.
// If clockwise is true, rotates right; if false, rotates left.
func (d *Direction) Rotate(clockwise bool) {
	*d = Rotate(*d, clockwise)
}

func (d *Direction) RotateRight() {
	d.Rotate(true)
}

func (d *Direction) RotateLeft() {
	d.Rotate(false)
}

type Mapping struct {
	Data   [][]rune
	Height int
	Width  int
}

func NewMapping(data [][]rune) Mapping {
	if len(data) == 0 {
		return Mapping{Data: data, Height: 0, Width: 0}
	}
	return Mapping{Data: data, Height: len(data), Width: len(data[0])}
}

func (m *Mapping) OutOfBounds(c Coordinate) bool {
	return c.X < 0 || c.Y < 0 || c.X >= m.Width || c.Y >= m.Height
}

func (m *Mapping) InBounds(c Coordinate) bool {
	return !m.OutOfBounds(c)
}

func (m *Mapping) Neighbors(c Coordinate) []Coordinate {
	neighbors := c.Neighbors()
	inBounds := make([]Coordinate, 0, len(neighbors))
	for _, coord := range neighbors {
		if m.InBounds(coord) {
			inBounds = append(inBounds, coord)
		}
	}
	return inBounds
}

func AddCoords(a, b Coordinate) Coordinate {
	return Coordinate{X: a.X + b.X, Y: a.Y + b.Y}
}

func SubCoords(a, b Coordinate) Coordinate {
	return Coordinate{X: a.X - b.X, Y: a.Y - b.Y}
}
