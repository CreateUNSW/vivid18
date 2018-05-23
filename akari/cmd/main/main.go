package main

import (
	"image/color"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pul-s4r/vivid18/akari/geo"
	"github.com/pul-s4r/vivid18/akari/lighting"
	"github.com/pul-s4r/vivid18/akari/mapping"
	"github.com/pul-s4r/vivid18/akari/netscan"
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

	stdDevices := []int{
		11, 12, 13, 14, 15, 16,
		21, 22, 23, 24, 25, 26, 31, 32, 33, 34, 35, 36,
		41, 42, 43, 44, 45,
		51, 52, 53,
		61, 62, 63, 64, 65, 66,
		71, 72, 73,
		81, 82, 83, 84, 85, 86, 87,
	}

	devices := make(map[int]*mapping.Device)
	ferns := make(map[int]*lighting.Fern)
	for _, deviceID := range stdDevices {
		devices[deviceID] = mapping.NewStandardDevice(deviceID)
		ferns[deviceID] = devices[deviceID].AsFern(0)
	}

	mapSystem(system, devices, ferns)

	// physicalFerns := []*Fern{
	// 	{
	// 		Location: geo.NewPoint(0, -140),
	// 		LEDs:     system.Root[0].Ferns[0].Fern.Arms,
	// 	},
	// 	{
	// 		Location: geo.NewPoint(-170, -170),
	// 		LEDs:     system.Root[0].Ferns[1].Fern.Arms,
	// 	},
	// }

	crowd := geo.NewMap()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	// for fernID, fern := range physicalFerns {
	// 	// for i := 0; i < len(fern.LEDs); i++ {
	// 	// 	for j := 0; j < len(fern.LEDs[i]); j++ {
	// 	// 		fern.LEDs[i][j] = &color.RGBA{}
	// 	// 	}
	// 	// }

	// 	f := system.Root[0].Ferns[fernID].Fern
	// 	ferns = append(ferns, f)

	// 	system.AddEffect(strconv.Itoa(fernID),
	// 		lighting.NewBlob(f, crowd, fern.Location, 310, 120))
	// }

	go func() {
		for range time.Tick(30 * time.Millisecond) {
			system.Run()
			for _, dev := range devices {
				go dev.Render()
			}
			crowd.Lock()
			payload := &Payload{
				Ferns:  []*Fern{},
				Sensor: crowd.Points,
			}
			crowd.Unlock()
			for _, lis := range listeners {
				lis <- payload
			}
		}
	}()

	// TODO: add proper translations
	receiver, err := netscan.Receive(logger, []*geo.Point{})
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

	for {
		receiver.ScanPeople(crowd)
	}
}

func reverseLEDs(leds []*color.RGBA) []*color.RGBA {
	result := make([]*color.RGBA, len(leds))
	for i := range leds {
		result[len(leds)-1-i] = leds[i]
	}
	return result
}

func mapSystem(system *lighting.System, devices map[int]*mapping.Device, ferns map[int]*lighting.Fern) {
	linears := map[string]*lighting.Linear{
		"A1A": &lighting.Linear{
			OuterFern: ferns[11],
		},
		"A1B": &lighting.Linear{
			InnerFern: ferns[11],
			OuterFern: ferns[12],
		},
		"A1C": &lighting.Linear{
			InnerFern: ferns[12],
			OuterFern: ferns[13],
		},

		"A2A": &lighting.Linear{
			OuterFern: ferns[14],
		},
		"A2B": &lighting.Linear{
			InnerFern: ferns[14],
			OuterFern: ferns[15],
		},
		"A2C": &lighting.Linear{
			InnerFern: ferns[15],
			OuterFern: ferns[16],
		},

		"B1A": &lighting.Linear{
			OuterFern: ferns[21],
		},
		"B1B": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[22],
		},
		"B1C": &lighting.Linear{
			InnerFern: ferns[22],
			OuterFern: ferns[23],
		},
		"B1D": &lighting.Linear{
			InnerFern: ferns[23],
			OuterFern: ferns[24],
		},

		"B2A": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[26],
		},
		"B2B": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[25],
		},

		"C1A": &lighting.Linear{
			OuterFern: ferns[31],
		},
		"C1B": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[33],
		},
		"C1C": &lighting.Linear{
			InnerFern: ferns[33],
			OuterFern: ferns[34],
		},
		"C1D": &lighting.Linear{
			InnerFern: ferns[34],
			OuterFern: ferns[35],
		},

		"C2A": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[32],
		},
		"C2B": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[36],
		},

		"D1A": &lighting.Linear{
			OuterFern: ferns[41],
		},
		"D1B": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[42],
		},

		"D2A": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[43],
		},
		"D2B": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[45],
		},

		"E1A": &lighting.Linear{
			OuterFern: ferns[15],
		},
		"E1B": &lighting.Linear{
			OuterFern: ferns[52],
		},
		"E1C": &lighting.Linear{
			InnerFern: ferns[52],
			OuterFern: ferns[53],
		},

		"F1A": &lighting.Linear{
			OuterFern: ferns[61],
		},
		"F1B": &lighting.Linear{
			InnerFern: ferns[61],
			OuterFern: ferns[63],
		},
		"F1C": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[62],
		},

		"F2A": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[64],
		},
		"F2B": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[65],
		},
		"F2C": &lighting.Linear{
			InnerFern: ferns[65],
			OuterFern: ferns[66],
		},

		"G1A": &lighting.Linear{
			OuterFern: ferns[71],
		},
		"G1B": &lighting.Linear{
			InnerFern: ferns[71],
			OuterFern: ferns[72],
		},
		"G1C": &lighting.Linear{
			InnerFern: ferns[72],
			OuterFern: ferns[73],
		},

		"H1A": &lighting.Linear{
			InnerFern: ferns[84],
			OuterFern: ferns[83],
			LEDs:      reverseLEDs(devices[81].LEDs[1][30 : 30+13]),
		},
		"H1B": &lighting.Linear{
			InnerFern: ferns[83],
			OuterFern: ferns[82],
		},
		"H1C": &lighting.Linear{
			InnerFern: ferns[82],
			OuterFern: ferns[81],
		},
		"H2A": &lighting.Linear{
			InnerFern: ferns[82],
			OuterFern: ferns[85],
		},
		"H2B": &lighting.Linear{
			InnerFern: ferns[85],
			OuterFern: ferns[86],
		},
		"H2C": &lighting.Linear{
			InnerFern: ferns[87],
			OuterFern: ferns[88],
		},
	}

	for _, linear := range linears {
		if linear.InnerFern != nil {
			linear.InnerFern.OuterLinear = append(linear.InnerFern.OuterLinear,
				linear)
		}

		if linear.OuterFern != nil {
			linear.InnerFern.InnerLinear = linear
		}
	}
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

	// go func() {
	// 	for {
	// 		if err := ws.ReadJSON(&scan.DebugPoint); err != nil {
	// 			return
	// 		}
	// 	}
	// }()

	for scan := range lis {
		if err := ws.WriteJSON(scan); err != nil {
			return nil
		}
	}

	return nil
}
