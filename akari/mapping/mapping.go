package mapping

import (
	"bytes"
	"image/color"
	"net"

	"github.com/pul-s4r/vivid18/akari/lighting"
)

var conn *net.UDPConn

// Port values
const (
	DevicePort = 5151
	ServerPort = 5050
)

func init() {
	var err error
	conn, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, 1),
		Port: 5050,
	})
	if err != nil {
		panic(err)
	}
}

// Device represents a remote network device with LEDs (i.e. the Arduino).
type Device struct {
	ID   int
	Addr *net.UDPAddr
	LEDs [][50]*color.RGBA
}

// NewDevice initializes and returns a new device given its address.
func NewDevice(id int, numChains int) *Device {
	if id >= 255 || id <= 0 {
		panic("device: NewDevice: id out of range")
	}

	remoteAddr := &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, byte(id)),
		Port: DevicePort,
	}

	d := &Device{
		ID:   id,
		Addr: remoteAddr,
		LEDs: make([][50]*color.RGBA, numChains),
	}

	for i := 0; i < numChains; i++ {
		for j := 0; j < len(d.LEDs[i]); j++ {
			d.LEDs[i][j] = &color.RGBA{}
		}
	}

	return d
}

// Render renders the lighting data to the device.
func (d *Device) Render() error {
	buf := new(bytes.Buffer)

	for _, chain := range d.LEDs {
		for _, col := range chain {
			buf.Write([]byte{col.R, col.G, col.B})
		}
	}

	_, err := conn.WriteToUDP(buf.Bytes(), d.Addr)
	return err
}

// AsFern returns a fern mapped to the device's first pin.
func (d *Device) AsFern(rotation int) *lighting.Fern {
	fern := &lighting.Fern{}

	for i := 0; i < len(fern.Arms); i++ {
		offset := 5 * ((i + rotation) % 8)
		fern.Arms[i] = [5]*color.RGBA{
			d.LEDs[0][0+offset],
			d.LEDs[0][4+offset],
			d.LEDs[0][1+offset],
			d.LEDs[0][3+offset],
			d.LEDs[0][2+offset],
		}
	}

	return fern
}

// ID    int
// Outer []LinearOnLinear // Linear node that is going away from the tree.
// Inner *Linear          // Linear node that is going towards the tree.
// Ferns []FernOnLinear

// // Mapping of LEDs on the chain. This is cleared on every Run().
// LEDs []*color.RGBA

// // Determines address mapping.
// startInner bool
