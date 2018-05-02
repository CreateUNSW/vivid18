package lighting

import (
	"math"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/pul-s4r/vivid18/akari/geo"
)

// Demo represents a demo effect.
type Demo struct {
	priority int
	start    time.Time
	deadline time.Time
	fern     *Fern
	data     *geo.Map
	loc      *geo.Point
}

// NewDemo returns a new Demo effect.
func NewDemo(fern *Fern, data *geo.Map, loc *geo.Point) *Demo {
	return &Demo{
		priority: 1,
		deadline: time.Now().Add(time.Hour * 8000),
		start:    time.Now(),
		fern:     fern,
		data:     data,
		loc:      loc,
	}
}

// Start returns the start time of the demo effect.
func (d *Demo) Start() time.Time {
	return d.start
}

// Deadline returns the deadline of the demo effect.
func (d *Demo) Deadline() time.Time {
	return d.deadline
}

// Priority returns the priority of the demo effect.
func (d *Demo) Priority() int {
	return d.priority
}

// Run runs.
func (d *Demo) Run(s *System) {
	t := time.Since(d.start).Seconds()

	d.data.Lock()
	defer d.data.Unlock()

	circle := math.Mod(t*100, 360)

	points := d.data.Within(d.loc, 300)
	r := 300.0
	if len(points) > 0 {
		r = math.Sqrt(float64(points[0].SquareDist(d.loc)))
	}

	pos := math.Mod(t*(250.0-(r-50))/250.0*10, 5)

	for _, arm := range d.fern.Arms {
		for i, led := range arm {
			if i != int(pos) {
				led.R = 0
				led.G = 0
				led.B = 0

				continue
			}

			r, g, b := colorful.Hsv(circle, 1, (250.0-(r-50))/250.0).RGB255()
			led.R = r
			led.G = g
			led.B = b
		}
	}
}
