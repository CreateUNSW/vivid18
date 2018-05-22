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
	priority  		int
	active    		bool
	start     		time.Time
	mainColor     	color.Color
	startFern 		*Fern
	speed	  		int // nanoseconds per led
					// (how many nanoseconds it takes for the pulse to move over a single led)
	colorfulColor 	colorful.Color	// Potentially change to hue and chroma for efficiency
	effectRadius 	int
}

// NeuralStepTime represents the amount of time it takes for the neural pulse to move
// one LED.
const NeuralStepTime = 50 * time.Millisecond
// Defines radius of effect in # of LEDs
const NeuralEffectRadius = 15

// NewNeural returns a new Neural effect.
func NewNeural(col color.Color, startFern *Fern, priority int, speed int, radius int) *Neural {
	return &Neural{
		priority:  priority,
		start:     time.Now(),
		active:    true,
		mainColor:     col,
		startFern: startFern,
		speed:	   speed
		colorfulColor: colorful.MakeColor(col),
		effectRadius: radius
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
// gets value from [0..1] for brightness
func (n *Neural) f(x int) float64 {
	if (math.Abs(x) > n.effectRadius) return 0	// if led is outside the radius of effect, it's 0
	return math.Sin(float64(x)*math.Pi/(2*n.effectRadius)+(math.Pi/2))
}

// TODO: Both runFern and runLinear override the LED value - fix by blending?
// - to blend each led requires storing its HCL color, 
// - or repeated conversion from RGB -> colorful.Color -> Hcl

// fernDist is how many leds away this fern is from the starting fern
func (n *Neural) runFern(fernDist int, effectDisplacement int, f *Fern) {
	armLength := len(f.Arms[0])
	// no op if effect doesn't affect this fern
	if (effectDisplacement + n.effectRadius < fernDist 
		|| effectDisplacement - n.effectRadius > fernDist + armLength) return;
	// for each led in an arm
	for i := 0; i < armLength; i++ {
		ledDistance := fernDist + i;
		distFromEffect := ledDistance - effectDisplacement
		colorRGBA = n.getColorFromDisplacement(distFromEffect)		
		for _, arm := range f.Arms {
			col := color.RGBA{
				R:	colorRBGA.R,
				G:	colorRGBA.G,
				B:	colorRGBA.B,
			}
			arm[i] = col	// TODO: blend color?
		}
	}
}

// fernDist is how many leds away this fern is from the starting fern
func (n *Neural) runLinear(linearDist int, effectDisplacement int, linear *Linear, outwards bool) {
	linearLength := len(linear.LEDs)
	// no op if effect doesn't affect this Linear
	if (effectDisplacement + n.effectRadius < linearDist 
		|| effectDisplacement - n.effectRadius > linearDist + linearLength) return;

	// if outwards, iterate from 0'th led outwards to increment distance correctly
	if (outwards) {
		for i := 0; i < linearLength; i++ {
			distFromEffect := linearDist - effectDisplacement
			linear.LEDs[i] = n.getColorFromDisplacement(distFromEffect)
			linearDist++
		}
	} else {
		for i := linearLength - 1; i >= 0; i-- {
			distFromEffect := linearDist - effectDisplacement
			linear.LEDs[i] = n.getColorFromDisplacement(distFromEffect)
			linearDist++			
		}
	}
}

// TODO: blend color with current led ? based on priority?

// Gets the value transformed effect color
// - or  blend between led's current color and effect's color?
func (n *Neural) getColor(value float64) color.Color {
	h, c, l := n.colorfulColor.Hcl()	// Get rid of this function call? 
										// (Store HCL in Neural rather than a colorful.Color)
	return colorful.Hcl(h, value * c, l * value)
}

// Returns a RGBA struct for an led calculated by distance between led and effect
func (n *Neural) getColorFromDisplacement(distFromEffect int) color.RGBA {
	ledVal := n.f(distFromEffect)
	ledColor := n.getColor(ledVal)
	r, g, b, _ = ledColor.RGBA();
	col := color.RGBA{
		R:	r,
		G:	g,
		B:	b,
	}
	return col
}

// From a fern apply the effect to the fern, and the inner linear if not outwards
// or outer linears if outwards
func (n *Neural) recursiveApply(ledDist int, effectDist int, fern *Fern, outwards bool) {
	if (fern == nil) return
	n.runFern(ledDist, effectDist, fern)
	if (outwards) {
		for _,Linear := range fern.OuterLinears {
			n.runLinear(ledDist, effectDist, Linear, outwards)
			n.recursiveApply(ledDist + len(Linear.LEDs), effectDist, Linear.OuterFern, outwards)
		}
	} else {
		n.runLinear(ledDist, effectDist, fern.InnerLinear, outwards)
		n.recursiveApply(ledDist + len(fern.InnerLinear.LEDs), effectDist, Linear.InnerFern, outwards)
	}
}

// Run runs.
func (n *Neural) Run(s *System) {
	duration := (time.Duration) s.CurrTime.Sub(n.start)	// duration since effect started
	effectDisplacement := (int) duration.Nanoseconds / n.speed; // how many leds effect has moved
	
	// Run effect on starting fern
	n.runFern(0, effectDisplacement, n.startFern)

	// Run the effect outwards on outer linears and recursively outwards
	for _, Linear := range n.startFern.OuterLinears {
		n.runLinear(0, effectDist, Linear, true)
		n.recursiveApply(len(Linear.LEDs), effectDist, Linear.OuterFern, true)
	}
	// Run the effect on inner linear and recursively inwards
	// TODO: With current component system - breaks at tree
	n.runLinear(0, effectDist, n.startFern.InnerLinear, false)
	n.recursiveApply(len(fern.InnerLinear.LEDs), effectDist, Linear.InnerFern, false)
}
