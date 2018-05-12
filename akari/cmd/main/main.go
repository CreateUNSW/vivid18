package main

import (
	"fmt"
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pul-s4r/vivid18/akari/geo"
	"github.com/pul-s4r/vivid18/akari/lighting"
	"github.com/pul-s4r/vivid18/akari/mapping"
	"github.com/pul-s4r/vivid18/akari/scan"
)

type Fern struct {
	Location *geo.Point        `json:"location"`
	LEDs     [8][5]*color.RGBA `json:"leds"`
}

type Payload struct {
	Ferns  []*Fern      `json:"ferns"`
	Sensor []*geo.Point `json:"sensor"`
}

var upgrader = websocket.Upgrader{}
var lisMutex = new(sync.Mutex)
var listeners = make(map[string]chan<- *Payload)

func main() {
	system := lighting.NewSystem()

	devices := []*mapping.Device{
		mapping.NewDevice(10, 2),
		mapping.NewDevice(11, 1),
		mapping.NewDevice(12, 1),
	}

	mapSystem(system, devices)

	physicalFerns := []*Fern{
		// {
		// 	Location: geo.NewPoint(-150, -150),
		// 	LEDs:     system.Root[0].Ferns[0].Fern.Arms,
		// },
		{
			Location: geo.NewPoint(0, 0),
			LEDs:     system.Root[0].Ferns[1].Fern.Arms,
		},
		// {
		// 	Location: geo.NewPoint(150, 150),
		// 	LEDs:     system.Root[0].Ferns[2].Fern.Arms,
		// },
	}

	var ferns []*lighting.Fern

	crowd := geo.NewMap()

	for fernID, fern := range physicalFerns {
		// for i := 0; i < len(fern.LEDs); i++ {
		// 	for j := 0; j < len(fern.LEDs[i]); j++ {
		// 		fern.LEDs[i][j] = &color.RGBA{}
		// 	}
		// }

		f := system.Root[0].Ferns[1].Fern
		ferns = append(ferns, f)

		system.AddEffect(strconv.Itoa(fernID),
			lighting.NewDemo(f, crowd, fern.Location))
	}

	go func() {
		for range time.Tick(30 * time.Millisecond) {
			system.Run()
			devices[1].Render()
			crowd.Lock()
			payload := &Payload{
				Ferns:  physicalFerns,
				Sensor: crowd.Points,
			}
			crowd.Unlock()
			for _, lis := range listeners {
				lis <- payload
			}
		}
	}()

	scanner, err := scan.SetupScanner("/dev/cu.usbserial-DO0088ZE")
	if err != nil {
		panic(err)
	}

	go func() {
		e := echo.New()
		e.GET("/ws", wsHandler)
		e.File("/", "index.html")
		e.File("/script.js", "script.js")
		e.Start(":9000")
	}()

	fmt.Println("lol")

	for {
		scanner.ScanPeople(crowd)
		// fmt.Println("scan completed")
	}
}

func mapSystem(system *lighting.System, devices []*mapping.Device) {
	north := &lighting.Linear{
		Outer: nil,
		Inner: nil,
		LEDs:  make([]*color.RGBA, 50),
	}

	system.Root = append(system.Root, north)
	copy(north.LEDs, devices[0].LEDs[1][:])

	north.AddFern(devices[0].AsFern(), 0)
	north.AddFern(devices[1].AsFern(), 15)
	north.AddFern(devices[2].AsFern(), 30)
	// system.Root = append(system.Root, &Linear{})
}

func wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	lis := make(chan *Payload, 100)
	id := uuid.New().String()

	lisMutex.Lock()
	listeners[id] = lis
	lisMutex.Unlock()

	defer func() {
		lisMutex.Lock()
		delete(listeners, id)
		lisMutex.Unlock()
	}()

	go func() {
		for {
			if err := ws.ReadJSON(&scan.DebugPoint); err != nil {
				return
			}
		}
	}()

	for scan := range lis {
		if err := ws.WriteJSON(scan); err != nil {
			return nil
		}
	}

	return nil
}
