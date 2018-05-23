package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"
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

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("add_neural <fern-id>  <hex> <priority> <speed> radius> <- add neural effect")
}

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

	physicalFerns := []*Fern{
		{
			Location: geo.NewPoint(0, -140),
			LEDs:     system.Root[0].Ferns[0].Fern.Arms,
		},
		{
			Location: geo.NewPoint(-170, -170),
			LEDs:     system.Root[0].Ferns[1].Fern.Arms,
		},
	}

	crowd := geo.NewMap()

	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{}

	for fernID, fern := range physicalFerns {
		// for i := 0; i < len(fern.LEDs); i++ {
		// 	for j := 0; j < len(fern.LEDs[i]); j++ {
		// 		fern.LEDs[i][j] = &color.RGBA{}
		// 	}
		// }

		f := system.Root[0].Ferns[fernID].Fern
		ferns = append(ferns, f)

		system.AddEffect(strconv.Itoa(fernID),
			lighting.NewBlob(f, crowd, fern.Location, 310, 120))
	}

	go func() {
		for range time.Tick(30 * time.Millisecond) {
			system.Run()
			for _, dev := range devices {
				go dev.Render()
			}
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
	
	effectID := 1
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scan.Scan() {
			fmt.Println("Goodbye!")
			break
		}

		args := strings.Split(scan.Text(), " ")
		if len(args) < 1 {
			printHelp()
			continue
		}

		switch args[0] {
		case "?", "help":
			printHelp()
		case "add_neural":
			if len(args) != 6 {
				fmt.Println("expected `add_neural` + 5 arguments")
				break
			}

			fernid, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid fern id")
				fmt.Println(err)
				break
			}
			startFern, ok := ferns[fernid]
			if !ok {
				//do something here
				fmt.Println("Fern id is invalid")
				break
			}

			dec, err := hex.DecodeString(args[2])
			if err != nil {
				fmt.Println("Invalid hex")
				fmt.Println(err)
				break
			}

			if len(dec) != 3 {
				fmt.Println("Invalid hex: hex must be 3 bytes")
				break
			}

			priority, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println("Invalid priority")
				fmt.Println(err)
				break
			}
			if (priority < 0 || priority > 5)
			{
				fmt.Println("Invalid priority: priority must be 0-5")
				break
			}

			speed, err := strconv.Atoi(args[4])
			if err != nil {
				fmt.Println("Invalid speed")
				fmt.Println(err)
				break
			}

			neuralRadius, err := strconv.Atoi(args[5])
			if err != nil {
				fmt.Println("Invalid neural radius")
				fmt.Println(err)
				break
			}
			colorRGBA := color.RGBA{
				R: dec[0],
				G: dec[1],
				B: dec[2],
			}
			neuralEffect := lighting.NewNeural(colorRGBA, startFern, priority, speed, neuralRadius)
			system.AddEffect(strconv.Itoa(effectID), neuralEffect)

		case "ls":
			deviceMutex.Lock()
			for ip, seen := range devices {
				fmt.Println(ip, "-", "last seen", humanize.Time(seen))
			}
			deviceMutex.Unlock()
		default:
			fmt.Println("Unknown command, type `help` for help")
		}

	}

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
		// TODO: complete
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
