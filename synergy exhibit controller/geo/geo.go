package geo

import (
	"encoding/gob"
	"math"
	"sort"
	"sync"
)

// Map represents a map.
type Map struct {
	Points []*Point
	mutex  *sync.Mutex
}

func init() {
	gob.Register(Map{})
}

// NewMap returns a new map.
func NewMap() *Map {
	return &Map{
		mutex: new(sync.Mutex),
	}
}

// ByDistance sorts a slice of points by distance.
type ByDistance struct {
	slice []*Point
	point *Point
}

// Len implements sort.Interface Len
func (b *ByDistance) Len() int {
	return len(b.slice)
}

// Less implements sort.Interface Less
func (b *ByDistance) Less(i, j int) bool {
	return b.point.SquareDist(b.slice[i]) < b.point.SquareDist(b.slice[j])
}

// Swap implements sort.Interface Swap
func (b *ByDistance) Swap(i, j int) {
	b.slice[i], b.slice[j] = b.slice[j], b.slice[i]
}

// Within returns points within a radius of p, sorted in
// ascending order of distance.
func (m *Map) Within(p *Point, r int) []*Point {
	sqR := r * r
	var results []*Point
	for _, mp := range m.Points {
		if mp.SquareDist(p) > sqR {
			continue
		}

		results = append(results, mp)
	}

	sort.Sort(&ByDistance{
		slice: results,
		point: p,
	})

	return results
}

// Add adds a point to the map.
func (m *Map) Add(p *Point) {
	m.Points = append(m.Points, p)
}

// Merge merges the given map into the map with an optional linear translation.
func (m *Map) Merge(b *Map, trans ...*Point) {
	for _, p := range b.Points {
		t := p
		if len(trans) > 0 {
			t = p.Add(trans[0])
		}

		m.Add(&Point{
			X:    t.X,
			Y:    t.Y,
			Data: p.Data,
		})
	}
}

// Lock locks the map for concurrent modification.
func (m *Map) Lock() {
	m.mutex.Lock()
}

// Unlock locks the map for concurrent modification.
func (m *Map) Unlock() {
	m.mutex.Unlock()
}

// Clear clears the map.
func (m *Map) Clear() {
	m.Points = make([]*Point, 0)
}

// Point represents a point.
type Point struct {
	X    int         `json:"x"`
	Y    int         `json:"y"`
	Data interface{} `json:"data"`
}

// NewPoint returns a new point with the given parameters.
func NewPoint(x, y int, data ...interface{}) *Point {
	if len(data) > 0 {
		return &Point{
			X:    x,
			Y:    y,
			Data: data[0],
		}
	}

	return &Point{
		X: x,
		Y: y,
	}
}

// Add adds two points together into a new point, and inherits the data from b.
func (p *Point) Add(b *Point) *Point {
	return &Point{
		X:    p.X + b.X,
		Y:    p.Y + b.Y,
		Data: b.Data,
	}
}

// SquareDist returns the squared distance between two points.
func (p *Point) SquareDist(b *Point) int {
	return (b.X-p.X)*(b.X-p.X) + (b.Y-p.Y)*(b.Y-p.Y)
}

// Angle returns the angle from b relative to p in radians.
func (p *Point) Angle(b *Point) float64 {
	return math.Atan(float64(b.Y-p.Y)/float64(b.X-p.X)) + math.Pi
}
