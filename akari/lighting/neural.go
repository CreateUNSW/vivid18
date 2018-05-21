//+build ignore

/*
implementation thoughts/pseudocoding

               v peak of the pulse
0............///\\\...........2*duration
^ edge led  ^ center led
follow a time template of length 2*duration - bell curve of alpha values hard coded as an array

shift starting index according to distance of led from center
i.e. led's at the center start the pulse at t=0 until t=pulseperiod, whereas
led's at the edge ferns start pulse at t=duration-pulseperiod until t=duration

if no explicit distance values are obtainable for each led, create a heuristic to approximate, based on recursive parent location, distance from center of indiviual fern

-main tree first
-runs along chains/cables
-spouts on each fern, from stem outwards to leaves

possible implementation:
    col.R/G/B = uint8(int(float64(r/g/b)*sequence[duration-distance]))
*/

package lighting

import (
	"image/color"
	"time"
)

// Neural represents a neural effect.
type Neural struct {
	priority int
	start    time.Time
	deadline time.Time
	color    color.Color
	fern     *Fern
}

// NeuralStepTime represents the amount of time it takes for the neural pulse to move
// one LED.
const NeuralStepTime = 50 * time.Millisecond

// NewNeural returns a new Neural effect.
func NewNeural(col color.Color, duration time.Duration, fern *Fern, priority int) 
*Neural {

    parentChain := fern.Linear
    var parentLocation int
    for _, f := range parentChain.Ferns {
        if f.Fern == fern {
            parentLocation = f.Location
        }
    }
	return &Neural{
		priority: priority,
		deadline: time.Now().Add(duration),
		start:    time.Now(),
		color:    col,
		fern:     fern,
	}
}

// Start returns the start time of the Neural effect.
func (n *Neural) Start() time.Time {
	return n.start
}

// Deadline returns the deadline of the Neural effect.
func (n *Neural) Deadline() time.Time {
	return n.deadline
}

// Priority returns the priority of the Neural effect.
func (n *Neural) Priority() int {
	return n.priority
}

func (n *Neural) recursiveApply(l *Linear, col color.RGBA) {
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
			n.recursiveApply(linear.Linear, col)
		}
	}
}

// Run runs.
func (n *Neural) Run(s *System) {
	progress := float64(time.Since(n.start)) / float64(n.deadline.Sub(n.start))

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


