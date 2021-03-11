package lighting

import (
	"image/color"
	"time"
)

// Specification:
//
// When there is little traffic, system will "breathe" by changing the brightness
// of the whole system in a breathing pattern at about 50 BPM. As traffic increases
// the pattern will increase in speed and colors will become more bright, up to 80 BPM.
// If beyond 70 BPM, the breathing will change from expanding outwards with brighter colors.
//
// The colors will slowly rotate, and the palette depends on the number of people.

// Blank represents a breathing effect.
type Blank struct {
	priority int
	active   bool
	start    time.Time
}

// NewBlank returns a new Blank effect.
func NewBlank(priority int) *Blank {
	return &Blank{
		priority: priority,
		start:    time.Now(),
		active:   true,
	}
}

// Active returns whether or not the effect is still active.
func (b *Blank) Active() bool {
	return true
}

// Start returns the start time of the Blank effect.
func (b *Blank) Start() time.Time {
	return b.start
}

// Priority returns the priority of the Blank effect.
func (b *Blank) Priority() int {
	return b.priority
}

func (b *Blank) recursiveApply(l *Linear, col *color.RGBA) {
	for _, led := range l.LEDs {
		led.R = 0
		led.G = 0
		led.B = 0
	}

	if l.OuterFern != nil {
		for _, arm := range l.OuterFern.Arms {
			for _, led := range arm {
				led.R = 0
				led.G = 0
				led.B = 0
			}
		}

		for _, child := range l.OuterFern.OuterLinears {
			b.recursiveApply(child, col)
		}
	}
}

// Run runs.
func (b *Blank) Run(s *System) {
	for _, led := range s.TreeTop.LEDs {
		led.R = 0
		led.G = 0
		led.B = 0
	}

	for _, led := range s.TreeBase.LEDs {
		led.R = 0
		led.G = 0
		led.B = 0
	}

	for _, arm := range s.Root {
		b.recursiveApply(arm, &color.RGBA{
			R: 0,
			G: 0,
			B: 0,
		})
	}
}
