package lighting

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

func makeSafe(c color.Color) colorful.Color {
	r, g, b, _ := c.RGBA()
	return colorful.Color{
		R: float64(r) / 65535.0,
		G: float64(g) / 65535.0,
		B: float64(b) / 65535.0,
	}
}

// Neural represents a neural effect.
type Neural struct {
	priority  int
	active    bool
	start     time.Time
	startFern *Fern
	speed     time.Duration // nanoseconds per led
	// (how many nanoseconds it takes for the pulse to move over a single led)
	mainColor    colorful.Color // Potentially change to hue and chroma for efficiency
	effectRadius float64
	startTree    bool
	startRoot    int
}

// NeuralStepTime represents the amount of time it takes for the neural pulse to move
// one LED.
const NeuralStepTime = 100 * time.Millisecond

// NeuralEffectRadius defines radius of effect in # of LEDs
const NeuralEffectRadius = 2

// NewNeural returns a new Neural effect.
func NewNeural(col color.Color, startFern *Fern, priority int, speed time.Duration, radius float64, startTree bool) *Neural {
	return &Neural{
		priority:     priority,
		start:        time.Now(),
		active:       true,
		startFern:    startFern,
		speed:        speed,
		mainColor:    makeSafe(col),
		effectRadius: radius,
		startTree:    startTree,
		startRoot:    -1,
	}
}

// Active returns whether or not the effect is still active.
func (n *Neural) Active() bool {
	return time.Since(n.start) < (7 * time.Second)
}

// Start returns the start time of the Neural effect.
func (n *Neural) Start() time.Time {
	return n.start
}

// Priority returns the priority of the Neural effect.
func (n *Neural) Priority() int {
	return n.priority
}

// for displacement of led from effect centre point,
// gets value from [0..1] for brightness
func (n *Neural) f(x float64) float64 {
	if math.Abs(x) > n.effectRadius {
		return 0 // if led is outside the radius of effect, it's 0
	}
	return math.Sin(x*math.Pi/(2.0*float64(n.effectRadius)) + (math.Pi / 2.0))
}

// TODO: Both runFern and runLinear override the LED value - fix by blending?
// - to blend each led requires storing its HCL color,
// - or repeated conversion from RGB -> colorful.Color -> Hcl

// fernDist is how many leds away this fern is from the starting fern
func (n *Neural) runFern(fernDist float64, effectDisplacement float64, f *Fern) {
	armLength := len(f.Arms[0])
	// no op if effect doesn't affect this fern
	if effectDisplacement+n.effectRadius < fernDist ||
		effectDisplacement-n.effectRadius > fernDist+float64(armLength*4) {
		return
	}

	for i := 0; i < armLength; i++ {
		ledDistance := fernDist + float64(i)
		distFromEffect := ledDistance - effectDisplacement
		blend := n.f(distFromEffect)
		col := makeSafe(*(f.Arms[0][i]))
		r, g, b := col.BlendHsv(n.mainColor, blend).Clamped().RGB255()

		for _, arm := range f.Arms {
			arm[i].R = r
			arm[i].G = g
			arm[i].B = b
		}
	}
}

// fernDist is how many leds away this fern is from the starting fern
func (n *Neural) runLinear(linearDist float64, effectDisplacement float64, linear *Linear, outwards bool) {
	linearLength := len(linear.LEDs)
	// no op if effect doesn't affect this Linear
	if effectDisplacement+n.effectRadius < linearDist ||
		effectDisplacement-n.effectRadius > linearDist+float64(linearLength*2) {
		return
	}

	// if outwards, iterate from 0'th led outwards to increment distance correctly
	if outwards {
		for i := 0; i < linearLength; i++ {
			distFromEffect := linearDist - effectDisplacement

			blend := n.f(distFromEffect)
			col := makeSafe(*(linear.LEDs[i]))
			r, g, b := col.BlendHsv(n.mainColor, blend).Clamped().RGB255()
			linear.LEDs[i].R = r
			linear.LEDs[i].G = g
			linear.LEDs[i].B = b
			linearDist++
		}
	} else {
		for i := linearLength - 1; i >= 0; i-- {
			distFromEffect := linearDist - effectDisplacement
			blend := n.f(distFromEffect)
			col := makeSafe(*(linear.LEDs[i]))
			r, g, b := col.BlendHsv(n.mainColor, blend).Clamped().RGB255()
			linear.LEDs[i].R = r
			linear.LEDs[i].G = g
			linear.LEDs[i].B = b
			linearDist++
		}
	}
}

// TODO: blend color with current led ? based on priority?

// // Gets the value transformed effect color
// // - or  blend between led's current color and effect's color?
// func (n *Neural) getColor(value float64) color.Color {
// 	h, c, l := n.colorfulColor.Hcl() // Get rid of this function call?
// 	// (Store HCL in Neural rather than a colorful.Color)
// 	return colorful.Hcl(h, value*c, l*value)
// }

// // Returns a RGBA struct for an led calculated by distance between led and effect
// func (n *Neural) getColorFromDisplacement(distFromEffect int) color.RGBA {
// 	ledVal := n.f(distFromEffect)
// 	ledColor := n.getColor(ledVal)
// 	r, g, b, _ = ledColor.RGBA()
// 	col := color.RGBA{
// 		R: r,
// 		G: g,
// 		B: b,
// 	}
// 	return col
// }

// From a fern apply the effect to the fern, and the inner linear if not outwards
// or outer linears if outwards
func (n *Neural) recursiveApply(ledDist float64, effectDist float64, fern *Fern, outwards bool) {
	if fern == nil {
		return
	}
	n.runFern(ledDist, effectDist, fern)
	if outwards {
		for _, Linear := range fern.OuterLinears {
			n.runLinear(ledDist, effectDist, Linear, outwards)
			n.recursiveApply(ledDist+float64(len(Linear.LEDs)), effectDist, Linear.OuterFern, outwards)
		}
	} else {
		n.runLinear(ledDist, effectDist, fern.InnerLinear, outwards)
		n.recursiveApply(ledDist+float64(len(fern.InnerLinear.LEDs)), effectDist, fern.InnerLinear.InnerFern, outwards)
	}
}

// Run runs.
func (n *Neural) Run(s *System) {
	duration := time.Since(n.start)
	effectDisplacement := float64(float64(duration) / float64(n.speed))

	if n.startTree {
		if n.startRoot < 0 {
			n.startRoot = rand.Intn(len(s.Root))
		}

		n.runLinear(0, effectDisplacement, s.Root[n.startRoot], true)
		n.recursiveApply(float64(len(s.Root[n.startRoot].LEDs)), effectDisplacement, s.Root[n.startRoot].OuterFern, true)
		return
	}

	// Run effect on starting fern
	n.runFern(0, effectDisplacement, n.startFern)

	// Run the effect outwards on outer linears and recursively outwards
	// for _, Linear := range n.startFern.OuterLinears {
	// 	n.runLinear(0, effectDisplacement, Linear, true)
	// 	n.recursiveApply(float64(len(Linear.LEDs)), effectDisplacement, Linear.OuterFern, true)
	// }
	// Run the effect on inner linear and recursively inwards
	// TODO: With current component system - breaks at tree
	n.runLinear(0, effectDisplacement, n.startFern.InnerLinear, false)
	// n.recursiveApply(len(n.startFern.InnerLinear.LEDs), effectDisplacement, n.startFern.InnerLinear., false)
	n.recursiveApply(float64(len(n.startFern.InnerLinear.LEDs)), effectDisplacement, n.startFern.InnerLinear.InnerFern, false)
}
