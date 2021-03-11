package lighting

import (
	"image/color"
	"time"
)

// Fade represents a fade effect.
type Fade struct {
	priority int
	start    time.Time
	deadline time.Time
	color    color.Color
}

// NewFade returns a new Fade effect.
func NewFade(col color.Color, duration time.Duration,
	priority int) *Fade {
	return &Fade{
		priority: priority,
		deadline: time.Now().Add(duration),
		start:    time.Now(),
		color:    col,
	}
}

// Start returns the start time of the fade effect.
func (f *Fade) Start() time.Time {
	return f.start
}

// Deadline returns the deadline of the fade effect.
func (f *Fade) Deadline() time.Time {
	return f.deadline
}

// Priority returns the priority of the fade effect.
func (f *Fade) Priority() int {
	return f.priority
}

func (f *Fade) recursiveApply(l *Linear, col color.RGBA) {
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

	if len(l.Outer) > 0 {
		for _, linear := range l.Outer {
			f.recursiveApply(linear.Linear, col)
		}
	}
}

// Run runs.
func (f *Fade) Run(s *System) {
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
