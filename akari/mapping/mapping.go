package mapping

import (
	"bytes"
	"image/color"
	"net"

	"github.com/pul-s4r/vivid18/akari/lighting"
)

var localConn *net.UDPConn

type Device struct {
	Addr *net.UDPAddr
	LEDs [3][]*color.RGBA
}

func (d *Device) Render() error {
	buf := new(bytes.Buffer)

	for _, chain := range d.LEDs {
		for _, col := range chain {
			buf.Write([]byte{col.R, col.G, col.B})
		}
	}

	_, err := localConn.WriteToUDP(buf.Bytes(), d.Addr)
	return err
}

func MapSystem(system *lighting.System) {
	d := &Device{
		Addr:
	}

	system.Root = append(system.Root, &Linear{
		ID: 0,
		Outer:
		Inner:
		LEDs:
	})
}

// ID    int
// Outer []LinearOnLinear // Linear node that is going away from the tree.
// Inner *Linear          // Linear node that is going towards the tree.
// Ferns []FernOnLinear

// // Mapping of LEDs on the chain. This is cleared on every Run().
// LEDs []*color.RGBA

// // Determines address mapping.
// startInner bool
