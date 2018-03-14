package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"

	"github.com/1lann/sweep"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader  = websocket.Upgrader{}
	mutex     = new(sync.Mutex)
	listeners = make(map[string]chan<- sweep.Scan)
)

func main() {
	f, err := os.Create("./cpu.profile")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			pprof.StopCPUProfile()
			f.Close()
			os.Exit(0)
		}
	}()

	e := echo.New()

	dev, err := sweep.NewDevice("/dev/cu.usbserial-DO0088ZE")
	if err != nil {
		panic(err)
	}

	dev.StopScan()
	dev.Drain()
	// dev.SetMotorSpeed(2)
	// dev.WaitUntilMotorReady()
	// dev.SetSampleRate(sweep.Rate1000)
	fmt.Println("Waiting ready")
	dev.WaitUntilMotorReady()

	fmt.Println("Starting scan")

	_, err = dev.StartScan()
	if err != nil {
		panic(err)
	}

	fmt.Println("Scan started")

	pprof.StartCPUProfile(f)

	// go func() {
	// 	for scan := range scanner {
	// 		mutex.Lock()
	// 		for _, lis := range listeners {
	// 			lis <- scan
	// 		}
	// 		mutex.Unlock()
	// 	}

	// 	mutex.Lock()
	// 	for _, lis := range listeners {
	// 		close(lis)
	// 	}
	// 	mutex.Unlock()
	// }()

	go func() {
		// e.GET("/ws", wsHandler)
		e.File("/", "index.html")
		e.File("/script.js", "script.js")

		e.Start(":9001")
	}()
	select {}
}

func wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	lis := make(chan sweep.Scan, 100)
	id := uuid.New().String()

	mutex.Lock()
	listeners[id] = lis
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(listeners, id)
		mutex.Unlock()
	}()

	for scan := range lis {
		if err := ws.WriteJSON(scan); err != nil {
			return nil
		}
	}

	return nil
}
