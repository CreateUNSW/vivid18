package lighting

import "image/color"

// FernOnLinear represents a fern that sits on a linear chain of lights.
type FernOnLinear struct {
	Location int // Address of where the fern physically is sitting on linear.
	Fern     *Fern
}

// LinearOnLinear represents an outgoing linear chain that goes into another
// linear chain (i.e. a fork)
type LinearOnLinear struct {
	Location int
	Linear   *Linear
}

// Linear represents a linear chain of lights.
//
// The start of a Linear chain will ALWAYS be at Inner, that is, address 0
// when Linear is used is ALWAYS towards Inner. This is automatically remapped
// if the physical light does not conform to this using startInner.
type Linear struct {
	ID    int
	Outer []LinearOnLinear // Linear node that is going away from the tree.
	Inner *Linear          // Linear node that is going towards the tree.
	Ferns []FernOnLinear

	// Mapping of LEDs on the chain. This is cleared on every Run().
	LEDs []*color.RGBA

	// Determines address mapping.
	startInner bool
}

// Fern represents a fern.
type Fern struct {
	Linear *Linear
	Arms   [8][]*color.RGBA
}

// TreeTop represents the lights on the top of the tree.
type TreeTop struct{}

// TreeBase represents the lights at the base of the tree.
type TreeBase struct{}
