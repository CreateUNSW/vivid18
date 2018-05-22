//+build ignore

// TODO: not completed

// Specification:
//
// When there is little traffic, system will "breathe" by changing the brightness
// of the whole system in a breathing pattern at about 50 BPM. As traffic increases
// the pattern will increase in speed and colors will become more bright, up to 80 BPM.
// If beyond 70 BPM, the breathing will change from expanding outwards with brighter colors.
//
// The colors will slowly rotate, and the palette depends on the number of people.

package lighting

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"time"
)

// Breathing represents a neural effect.
type Breathing struct {
	priority  int
	active    bool
	start     time.Time
	color     color.Color
	startFern *Fern
}

// NeuralStepTime represents the amount of time it takes for the neural pulse to move
// one LED.
const NeuralStepTime = 50 * time.Millisecond

// NewNeural returns a new Breathing effect.
func NewNeural(col color.Color, startFern *Fern, priority int) *Breathing {
	return &Breathing{
		priority:  priority,
		start:     time.Now(),
		active:    true,
		color:     col,
		startFern: startFern,
	}
}

// Active returns whether or not the effect is still active.
func (n *Breathing) Active() bool {
	return time.Since(start) < (5 * time.Second)
}

// Start returns the start time of the Breathing effect.
func (n *Breathing) Start() time.Time {
	return n.start
}

// Deadline returns the deadline of the Breathing effect.
func (n *Breathing) Deadline() time.Time {
	return n.deadline
}

// Priority returns the priority of the Breathing effect.
func (n *Breathing) Priority() int {
	return n.priority
}

func (n *Breathing) f(x int) float64 {
	steps := time.Since(n.start) / NeuralStepTime
	steps -= math.Pi

	if math.Abs(float64(x)-steps) > math.Pi {
		return 0
	}

	return math.Sin(float64(x)+(math.Pi/2)) + 1
}

func (n *Breathing) runFern(d int, f *Fern) {
	for i := 0; i < len(f.Arms[0]); i++ {
		for _, arm := range f.Arms {
			color.
			arm[i].
		}
	}
	f.Arms
}

func (n *Breathing) getColor(d int) color.Color {
	310, 120
	colorful.Hcl(310, 1.0, 0.5).BlendHcl(col2 colorful.Color, t float64)
}

// Run runs.
func (n *Breathing) Run(s *System) {
	n.startFern

	for _,

	r, g, b, _ := n.color.RGBA()
	col := color.RGBA{
		R: uint8(int(float64(r)*progress) >> 8),
		G: uint8(int(float64(g)*progress) >> 8),
		B: uint8(int(float64(b)*progress) >> 8),
	}

	for _, l := range s.Root {
		n.recursiveApply(l, col)
	}
}
