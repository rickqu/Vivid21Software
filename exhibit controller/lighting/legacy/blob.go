package lighting

import (
	"math"
	"time"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/rickqu/Vivid21Software/exhibit%20controller/geo"
)

// Blob represents a blob effect.
type Blob struct {
	priority int
	start    time.Time
	deadline time.Time
	fern     *Fern
	data     *geo.Map
	loc      *geo.Point

	lastBlends [8][5]float64
2
	baseColor   float64
	accentColor float64
}

// NewBlob returns a new Blob effect.
func NewBlob(fern *Fern, data *geo.Map, loc *geo.Point,
	baseColor float64, accentColor float64) *Blob {
	return &Blob{
		priority: 1,
		deadline: time.Now().Add(time.Hour * 8000),
		start:    time.Now(),
		fern:     fern,
		data:     data,
		loc:      loc,

		baseColor:   baseColor,
		accentColor: accentColor,
	}
}

// Start returns the start time of the blob effect.
func (b *Blob) Start() time.Time {
	return b.start
}

// Deadline returns the deadline of the blob effect.
func (b *Blob) Deadline() time.Time {
	return b.deadline
}

// Priority returns the priority of the blob effect.
func (b *Blob) Priority() int {
	return b.priority
}

func blendRGB(a, b colorful.Color, t float64) colorful.Color {
	rp := 1.0
	gp := 1.0
	bp := 1.0
	if b.R-a.R < 0 {
		rp = -1
	}
	if b.G-a.G < 0 {
		gp = -1
	}
	if b.B-a.B < 0 {
		bp = -1
	}

	return colorful.Color{
		R: a.R + math.Sqrt(t*(b.R-a.R)*(b.R-a.R))*rp,
		G: a.G + math.Sqrt(t*(b.G-a.G)*(b.G-a.G))*gp,
		B: a.B + math.Sqrt(t*(b.B-a.B)*(b.B-a.B))*bp,
	}
}

// Run runs.
func (b *Blob) Run(s *System) {
	t := time.Since(b.start)

	b.data.Lock()
	defer b.data.Unlock()

	base := 12
	delta := 7

	for i, arm := range b.fern.Arms {
		angle := float64(i) * (math.Pi / 4.0)

		mx := math.Cos(angle)
		my := math.Sin(angle)

		lumos := (math.Sin((math.Mod(t.Seconds(), 4.0)/2.0)*math.Pi) + 1.3) / 6.0
		// _ = t
		// lumos := 0.5
		baseColor := colorful.Hsl(b.baseColor, 1, lumos)

		for j, led := range arm {
			d := float64(base + delta*j)

			ledPos := b.loc.Add(geo.NewPoint(int(d*mx), int(d*my), nil))
			// fmt.Println(geo.NewPoint(int(d*mx), int(d*my)))

			points := b.data.Within(ledPos, 200)
			blend := 0.0
			if len(points) > 0 {
				dist := math.Sqrt(float64(ledPos.SquareDist(points[0])))
				dist = math.Max(dist-30, 0)
				blend = (170.0 - dist) / 170.0
			}

			if blend > b.lastBlends[i][j] {
				b.lastBlends[i][j] = math.Min(b.lastBlends[i][j]+0.01, blend)
			} else if blend < b.lastBlends[i][j] {
				b.lastBlends[i][j] = math.Max(b.lastBlends[i][j]-0.01, blend)
			}

			// finalBlend := math.Min(b.lastBlends[i][j], 1.0)

			// finalColor := blendRGB(baseColor, colorful.Hsl(b.accentColor, 1.0, 1.0/3.0), b.lastBlends[i][j])
			finalColor := baseColor.BlendHcl(colorful.Hsl(b.accentColor, 1.0, 1.0/4.0), b.lastBlends[i][j]).Clamped()
			// finalColor := colorful.Hsl(b.accentColor, 1.0, 1.0/2.0)
			// _, _ = blend, baseColor

			r, g, b := finalColor.RGB255()
			led.R = r
			led.G = g
			led.B = b
		}
	}
}
