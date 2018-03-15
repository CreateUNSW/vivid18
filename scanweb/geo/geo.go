package geo

// Map represents a map.
type Map struct {
	Points []*Point
}

// NewMap returns a new map.
func NewMap() *Map {
	return &Map{}
}

// Within returns points within a radius of p.
func (m *Map) Within(p *Point, r int) []*Point {
	sqR := r * r
	var results []*Point
	for _, mp := range m.Points {
		dx := (mp.X - p.X) * (mp.X - p.X)
		if dx > sqR {
			continue
		}
		if dx+(mp.Y-p.Y)*(mp.Y-p.Y) > sqR {
			continue
		}
		results = append(results, mp)
	}

	return results
}

// Add adds a point to the map.
func (m *Map) Add(p *Point) {
	m.Points = append(m.Points, p)
}

// Point represents a point.
type Point struct {
	X, Y int
	Data interface{}
}

// Add adds two points together, and inherits the data from b.
func (p *Point) Add(b *Point) *Point {
	return &Point{
		X:    p.X + b.X,
		Y:    p.Y + b.Y,
		Data: b.Data,
	}
}
