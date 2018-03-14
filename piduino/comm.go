package main

import (
	"log"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
	_ "periph.io/x/periph/host/rpi"
)

func main() {
	host.Init()
	bus, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}

	log.Println("set speed:", bus.SetSpeed(10000))
	log.Println(bus.Tx(8, []byte("Hello, world!\n"), nil))
}
