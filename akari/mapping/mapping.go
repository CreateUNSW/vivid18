package mapping

import (
	"bytes"
	"image/color"
	"log"
	"net"
	"strconv"

	"github.com/pul-s4r/vivid18/akari/lighting"
)

// Device represents a remote network device with LEDs (i.e. the Arduino).
type Device struct {
	Addr *net.UDPAddr
	LEDs [][50]*color.RGBA

	conn     *net.UDPConn
	response chan<- []byte
}

func parseAddr(addr string) *net.UDPAddr {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		panic("mapping: parseAddr: invalid port number")
	}

	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: portNum,
	}
	if udpAddr.IP == nil {
		panic("mapping: parseAddr: invalid IP address")
	}

	return udpAddr
}

func getAddr(id int) string {
	return "192.168.2." + strconv.Itoa(id) + ":6969"
}

// NewDevice initializes and returns a new device given its address.
func NewDevice(id int, numChains int) *Device {
	conn, err := net.ListenUDP("udp", parseAddr("0.0.0.0:0"))
	if err != nil {
		panic(err)
	}

	remoteAddr := parseAddr(getAddr(id))
	response := make(chan []byte, 1000)

	go func() {
		buffer := make([]byte, 1500)
		for {
			n, readAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("mapping: error reading UDP from %s: %v\n",
					remoteAddr.IP.String(), err)
				return
			}

			if !readAddr.IP.Equal(remoteAddr.IP) {
				log.Printf("mapping: rogue message from %s\n", readAddr.IP.String())
				continue
			}

			result := make([]byte, n)
			copy(result, buffer)
			response <- result
		}
	}()

	d := &Device{
		Addr: remoteAddr,
		LEDs: make([][50]*color.RGBA, numChains),

		conn:     conn,
		response: response,
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

	_, err := d.conn.WriteToUDP(buf.Bytes(), d.Addr)
	return err
}

// AsFern returns a fern mapped to the device's first pin.
func (d *Device) AsFern() *lighting.Fern {
	fern := &lighting.Fern{}

	for i := 0; i < len(fern.Arms); i++ {
		offset := 5 * i
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
