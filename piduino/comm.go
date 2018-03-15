package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"
)

// Color ...
type Color struct {
	Hue      byte
	Lumosity byte
}

// Chain ...
type Chain struct {
	Pin    byte
	Colors []*Color
}

func main() {
	host.Init()
	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}

	log.Println("set speed:", bus.SetSpeed(1000000))

	var data []byte
	for i := 0; i < 30; i++ {
		data = append(data, 0x12)
	}

	var errors = 0
	go func() {
		for {
			fmt.Println("errors:", errors)
			errors = 0
			time.Sleep(time.Second)
		}
	}()

	for {
		bus.Tx(8, data, nil)
		time.Sleep(30 * time.Millisecond)
		errors++
	}

	// for {
	// 	var i byte
	// 	for i = 0; i <= 255; i++ {
	// 		c := &Chain{
	// 			Pin: 3,
	// 		}

	// 		for j := 0; j < 10; j++ {
	// 			c.Colors = append(c.Colors, &Color{Hue: i, Lumosity: 255})
	// 		}

	// 		err := bus.Tx(8, c.Bytes(), nil)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		time.Sleep(50 * time.Millisecond)
	// 	}
	// }
}

// Bytes returns a byte representation of the chain.
func (c *Chain) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(c.Pin)
	buf.WriteByte(byte(len(c.Colors)))
	for _, col := range c.Colors {
		buf.Write(col.Bytes())
	}
	return buf.Bytes()
}

// Bytes returns a byte representation of the color.
func (c *Color) Bytes() []byte {
	return []byte{c.Hue, c.Lumosity}
}
