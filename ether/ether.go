package main

import (
	"bytes"
	"log"
	"math"
	"net"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var conn *net.UDPConn
var ardIP *net.UDPAddr

func main() {
	var err error
	ardIP, err = net.ResolveUDPAddr("udp", "192.168.2.30:6969")
	if err != nil {
		panic(err)
	}

	myIP, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}

	conn, err = net.ListenUDP("udp", myIP)

	lightLoop()
}

func lightLoop() {
	t := time.NewTicker(33333 * time.Microsecond)
	circle := 0.0

	i := 0
	drift := 0.0

	for range t.C {
		i++

		circle += 5.0
		if circle >= 360 {
			circle = 0
		}

		drift += 0.1
		if drift >= 6 {
			drift = 0
		}

		buf := new(bytes.Buffer)
		for i := 0; i < 40; i++ {
			c := colorful.Hsv(circle, 1, getBrightness(float64(i), drift))
			r, g, b := c.RGB255()
			buf.Write([]byte{r, g, b})
		}

		_, err := conn.WriteToUDP(buf.Bytes(), ardIP)
		if err != nil {
			log.Println("error writing:", err)
		}
	}
}

func getBrightness(position, drift float64) float64 {
	n := (-1/((math.Sin(((position-2.5+drift)*math.Pi)/3)+1)/2)*5 + 9) / 4
	if n < 0 {
		n = 0
	}
	return n
}
