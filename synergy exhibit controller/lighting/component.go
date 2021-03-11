package lighting

import "image/color"

// Linear represents a linear chain of lights.
//
// The start of a Linear chain will ALWAYS be at Inner, that is, address 0
// when Linear is used is ALWAYS towards Inner.
type Linear struct {
	InnerFern *Fern // Linear's inner fern
	OuterFern *Fern // Linear's outer fern(s)

	// Mapping of LEDs on the chain. This is cleared on every Run().
	LEDs []*color.RGBA
}

// Fern represents a fern.
type Fern struct {
	InnerLinear  *Linear   // A Fern's inner linear
	OuterLinears []*Linear // A fern's outer linear(s)
	Arms         [8][5]*color.RGBA
}

// TODO: Consider TreeBase/TreeTop as well as ferns to inherit from
// a common "Node", which all have OuterLinears
// allows for effects to recursively flow through the tree
// by iterating over the tree's outer linears?
// (but will also probably need a boolean to identify "Node" type)

// Also - for priority to work, will probably need a "PrioritySum"
// property for each led, where it's equal to the sum of the priorities
// of the effects applied

// TreeTop represents the lights on the top of the tree.
type TreeTop struct {
	LEDs []*color.RGBA
}

// TreeBase represents the lights at the base of the tree.
type TreeBase struct {
	LEDs []*color.RGBA
}
