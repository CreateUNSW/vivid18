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
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"math"
	"time"
)

// Neural represents a neural effect.
type Neural struct {
	priority  int
	active    bool
	start     time.Time
	color     color.Color
	startFern *Fern
	speed	  int // nanoseconds per led
					// (how many nanoseconds it takes for the pulse to move over a single led)
}

// NeuralStepTime represents the amount of time it takes for the neural pulse to move
// one LED.
const NeuralStepTime = 50 * time.Millisecond
// Defines radius of effect in # of LEDs
const NeuralEffectRadius = 15

// NewNeural returns a new Neural effect.
func NewNeural(col color.Color, startFern *Fern, priority int, speed int) *Neural {
	return &Neural{
		priority:  priority,
		start:     time.Now(),
		active:    true,
		color:     col,
		startFern: startFern,
		speed:	   speed
	}
}

// Active returns whether or not the effect is still active.
func (n *Neural) Active() bool {
	return time.Since(start) < (5 * time.Second)
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

// for displacement of led from effect centre point,
// gets value from 0-1 for brightness
func (n *Neural) f(x int) float64 {
	if (math.Abs(x) > 15) return 0	// if led is outside the radius of effect, it's 0	
	return math.Sin(float64(x)*math.Pi/(2*NeuralEffectRadius)+(math.Pi/2))
}

// fernDist is how many leds away this fern is from the starting fern
func (n *Neural) runFern(fernDist int, effectDisplacement int, f *Fern) {
	armLength := len(f.Arms[0])
	// for each led in an arm
	for i := 0; i < len(f.Arms[0]); i++ {
		for _, arm := range f.Arms {
			ledDistance := fernDist + i;
			distFromEffect := ledDistance - effectDisplacement
			ledVal := n.f(distFromEffect)
			arm[i] = getColor(ledVal)
		}
	}
}

func (n *Neural) getColor(value float64) color.Color {
	mainColor = (colorful.Color) colorful.MakeColor(n.color)
	return colorful.Hcl(mainColor.Hcl().h, mainColor.Hcl().c, value)
}

// Run runs.
func (n *Neural) Run(s *System) {

	r, g, b, _ := n.color.RGBA()
	duration := (time.Duration) s.CurrTime.Sub(n.start)	// duration since effect started
	effectDisplacement := (int) duration.Nanoseconds / n.speed; // how many leds effect has moved
	
	col := color.RGBA{
		R: uint8(int(float64(r)*progress) >> 8),
		G: uint8(int(float64(g)*progress) >> 8),
		B: uint8(int(float64(b)*progress) >> 8),
	}

	for _, l := range s.Root {
		n.recursiveApply(l, col)
	}
}
