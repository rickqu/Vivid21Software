package lighting

import (
	"image/color"
	"math"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

// Specification:
//
// When there is little traffic, system will "breathe" by changing the brightness
// of the whole system in a breathing pattern at about 50 BPM. As traffic increases
// the pattern will increase in speed and colors will become more bright, up to 80 BPM.
// If beyond 70 BPM, the breathing will change from expanding outwards with brighter colors.
//
// The colors will slowly rotate, and the palette depends on the number of people.

// Breathing represents a breathing effect.
type Breathing struct {
	priority int
	active   bool
	start    time.Time
}

// NewBreathing returns a new Breathing effect.
func NewBreathing(priority int) *Breathing {
	return &Breathing{
		priority: priority,
		start:    time.Now(),
		active:   true,
	}
}

// Active returns whether or not the effect is still active.
func (b *Breathing) Active() bool {
	return true
}

// Start returns the start time of the Breathing effect.
func (b *Breathing) Start() time.Time {
	return b.start
}

// Priority returns the priority of the Breathing effect.
func (b *Breathing) Priority() int {
	return b.priority
}

func (b *Breathing) recursiveApply(l *Linear, col *color.RGBA) {
	for _, led := range l.LEDs {
		led.R = col.R
		led.G = col.G
		led.B = col.B
	}

	if l.OuterFern != nil {
		for _, arm := range l.OuterFern.Arms {
			for _, led := range arm {
				led.R = col.R
				led.G = col.G
				led.B = col.B
			}
		}

		for _, child := range l.OuterFern.OuterLinears {
			b.recursiveApply(child, col)
		}
	}
}

// Run runs.
func (b *Breathing) Run(s *System) {
	t := time.Since(b.start)

	h := math.Sin((t.Seconds()*math.Pi)/10.0+1)*60 + 300
	lumos := (math.Sin((math.Mod(t.Seconds(), 4.0)/2.0)*math.Pi) + 1.3) / 10.0
	c := colorful.Hsl(h, 1.0, lumos)
	treeC := colorful.Hsl(h, 1.0, lumos*2)

	tr, tg, tb := treeC.RGB255()
	red, gre, blu := c.RGB255()

	for _, led := range s.TreeTop.LEDs {
		led.R = tr
		led.G = tg
		led.B = tb
	}

	for _, led := range s.TreeBase.LEDs {
		led.R = tr
		led.G = tg
		led.B = tb
	}

	for _, arm := range s.Root {
		b.recursiveApply(arm, &color.RGBA{
			R: red,
			G: gre,
			B: blu,
		})
	}
}
