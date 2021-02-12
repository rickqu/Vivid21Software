//+build ignore

package lighting

import (
	"image/color"
	"time"
)

// Pulse represents a fade effect.
type Pulse struct {
	priority int
	start    time.Time
	deadline time.Time
	color    color.Color
	fern *Fern
}

// PulseStepTime represents the amount of time it takes for the pulse to move
// one LED.
const PulseStepTime = 50 * time.Millisecond

// NewPulse returns a new Pulse effect.
func NewPulse(col color.Color, fern *Fern, priority int) *Fade {
	// First, identify number of steps to the tree to calculate the time required.
	parentChain := fern.Linear
	var parentLocation int
	for _, f := range parentChain.Ferns {
		if f.Fern == fern {
			parentLocation = f.Location
		}
	}

	sum := len(parentChain.LEDs) - parentLocation
	currentChain := parentChain.Outer



	return &Fade{
		priority: priority,
		deadline: time.Now().Add(duration),
		start:    time.Now(),
		color:    col,
		fern: fern,
	}
}

// Start returns the start time of the fade effect.
func (p *Pulse) Start() time.Time {
	return p.start
}

// Deadline returns the deadline of the fade effect.
func (p *Pulse) Deadline() time.Time {
	return p.deadline
}

// Priority returns the priority of the fade effect.
func (p *Pulse) Priority() int {
	return p.priority
}

func (p *Pulse) recursiveApply(l *Linear, col color.RGBA) {
	for _, led := range l.LEDs {
		led.R = col.R
		led.G = col.G
		led.B = col.B
	}

	for _, fern := range l.Ferns {
		for _, arm := range fern.Fern.Arms {
			for _, led := range arm {
				led.R = col.R
				led.G = col.G
				led.B = col.B
			}
		}
	}

	if l.Outer != nil {
		f.recursiveApply(p.Outer, col)
	}
}

// Run runs.
func (p *Pulse) Run(s *System) {
	progress := float64(time.Since(f.start)) / float64(f.deadline.Sub(f.start))

	r, g, b, _ := f.color.RGBA()
	col := color.RGBA{
		R: uint8(int(float64(r)*progress) >> 8),
		G: uint8(int(float64(g)*progress) >> 8),
		B: uint8(int(float64(b)*progress) >> 8),
	}

	for _, l := range s.Root {
		f.recursiveApply(l, col)
	}
}
