package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/CreateUNSW/vivid18/akari/geo"
	"github.com/CreateUNSW/vivid18/akari/lighting"
	"github.com/CreateUNSW/vivid18/akari/mapping"
	"github.com/CreateUNSW/vivid18/akari/report"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	_ "github.com/CreateUNSW/vivid18/akari/scan"
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
		73, 34, 32, 12, 15, 22, 13, 66, 33, 65, 87, 23,
	}

	devices := make(map[int]*mapping.Device)
	ferns := make(map[int]*lighting.Fern)
	for _, deviceID := range stdDevices {
		devices[deviceID] = mapping.NewStandardDevice(deviceID)
		ferns[deviceID] = devices[deviceID].AsFern(0)
	}

	// TODO: add fern here that has 2 roots

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

	system.AddEffect("breathing", lighting.NewBreathing(1))
	// system.AddEffect("blank", lighting.NewBlank(1))

	reporter := report.NewReporter(mapping.Conn, logger)

	var edgeFerns = []int{66, 33, 23}

	go func() {
		treeLast := time.Now()
		neuralLast := time.Now()
		for range time.Tick(33 * time.Millisecond) {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println(r)
					}
				}()

				if time.Since(treeLast) > 2200*time.Millisecond {
					system.AddEffect(uuid.New().String(), lighting.NewNeural(color.RGBA{
						R: 0,
						G: 0xff,
						B: 0,
					}, ferns[73], 4, lighting.NeuralStepTime, lighting.NeuralEffectRadius, true))
					treeLast = time.Now()
				}

				if time.Since(neuralLast) > 1000*time.Millisecond {
					n := rand.Intn(len(edgeFerns))
					system.AddEffect(uuid.New().String(), lighting.NewNeural(color.RGBA{
						R: 0xff,
						G: 0xff,
						B: 0xff,
					}, ferns[n], 4, lighting.NeuralStepTime, lighting.NeuralEffectRadius, false))
				}

				// 	neuralLast = time.Now()
				// }

				system.Run()
				for _, dev := range devices {
					report := reporter.GetReport(int(dev.Addr.IP.To4()[3]))
					if report != nil && time.Since(report.LastSeen) < 7*time.Second {
						dev.Render()
					}
				}
			}()
			// crowd.Lock()
			// payload := &Payload{
			// 	Ferns:  []*Fern{},
			// 	Sensor: crowd.Points,
			// }
			// crowd.Unlock()
			// for _, lis := range listeners {
			// 	lis <- payload
			// }
		}
	}()

	// effectID := 1
	// scan := bufio.NewScanner(os.Stdin)
	// for {
	// 	fmt.Print("> ")
	// 	if !scan.Scan() {
	// 		fmt.Println("Goodbye!")
	// 		break
	// 	}

	// 	args := strings.Split(scan.Text(), " ")
	// 	if len(args) < 1 {
	// 		printHelp()
	// 		continue
	// 	}

	// 	switch args[0] {
	// 	case "?", "help":
	// 		printHelp()
	// 	case "state":
	// 		fmt.Println(len(system.RunningEffects))
	// 	case "add_neural":
	// 		if len(args) != 3 {
	// 			fmt.Println("expected `add_neural` + 2 arguments")
	// 			break
	// 		}

	// 		fernid, err := strconv.Atoi(args[1])
	// 		if err != nil {
	// 			fmt.Println("Invalid fern id")
	// 			fmt.Println(err)
	// 			break
	// 		}
	// 		startFern, ok := ferns[fernid]
	// 		if !ok {
	// 			//do something here
	// 			fmt.Println("Fern id is invalid")
	// 			break
	// 		}

	// 		dec, err := hex.DecodeString(args[2])
	// 		if err != nil {
	// 			fmt.Println("Invalid hex")
	// 			fmt.Println(err)
	// 			break
	// 		}

	// 		if len(dec) != 3 {
	// 			fmt.Println("Invalid hex: hex must be 3 bytes")
	// 			break
	// 		}

	// 		colorRGBA := color.RGBA{
	// 			R: dec[0],
	// 			G: dec[1],
	// 			B: dec[2],
	// 		}
	// 		neuralEffect := lighting.NewNeural(colorRGBA, startFern, 5,
	// 			lighting.NeuralStepTime, lighting.NeuralEffectRadius, false)
	// 		system.AddEffect(strconv.Itoa(effectID), neuralEffect)
	// 	default:
	// 		fmt.Println("Unknown command, type `help` for help")
	// 	}

	// }

	// TODO: add proper translations
	// receiver, err := netscan.Receive(logger, []*geo.Point{})
	// if err != nil {
	// 	panic(err)
	// }

	// go func() {
	// 	e := echo.New()
	// 	e.GET("/ws", wsHandler)
	// 	e.File("/", "index.html")
	// 	e.File("/script.js", "script.js")
	// 	e.Start(":9000")
	// }()

	// for {
	// 	receiver.ScanPeople(crowd)
	// 	fmt.Println("scan")
	// 	results := receiver.GetAll()

	// 	if results[2] != nil {
	// 		payload := &Payload{
	// 			Ferns:  []*Fern{},
	// 			Sensor: results[2].Points,
	// 		}
	// 		lisMutex.Lock()
	// 		for _, lis := range listeners {
	// 			lis <- payload
	// 		}
	// 		lisMutex.Unlock()
	// 	}

	// 	for i := 2; i <= 5; i++ {
	// 		if results[i] == nil {
	// 			continue
	// 		}

	// 		if len(results[i].Within(&geo.Point{X: 0, Y: 0}, 170)) > 0 {
	// 			fmt.Println("activated!")
	// 			activate[i] = true
	// 		} else {
	// 			activate[i] = false
	// 		}
	// 	}
	// }
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
		"A": &lighting.Linear{
			InnerFern: ferns[73],
			OuterFern: ferns[34],
			LEDs:      devices[73].LEDs[1][0:15],
		},
		"B": &lighting.Linear{
			InnerFern: ferns[34],
			OuterFern: ferns[32],
			LEDs:      devices[73].LEDs[1][15 : 15+21],
		},
		"C": &lighting.Linear{
			InnerFern: ferns[32],
			OuterFern: ferns[33],
			LEDs:      devices[73].LEDs[1][15+21 : 15+21+10],
		},
		"D": &lighting.Linear{
			InnerFern: ferns[12],
			OuterFern: ferns[34],
			LEDs:      devices[12].LEDs[1][0:12],
		},
		"E": &lighting.Linear{
			InnerFern: ferns[34],
			OuterFern: ferns[66],
			LEDs:      devices[12].LEDs[1][12 : 12+21],
		},
		"F": &lighting.Linear{
			InnerFern: ferns[15],
			OuterFern: ferns[22],
			LEDs:      devices[15].LEDs[1][0:12],
		},
		"G": &lighting.Linear{
			InnerFern: ferns[22],
			OuterFern: ferns[13],
			LEDs:      devices[15].LEDs[1][12 : 12+13],
		},
		"H": &lighting.Linear{
			InnerFern: ferns[13],
			OuterFern: ferns[66],
			LEDs:      devices[15].LEDs[1][12+13 : 12+13+11],
		},
		"I": &lighting.Linear{
			InnerFern: ferns[33],
			OuterFern: ferns[65],
			LEDs:      devices[33].LEDs[1][0:7],
		},
		"J": &lighting.Linear{
			InnerFern: ferns[65],
			OuterFern: ferns[87],
			LEDs:      devices[65].LEDs[1][7 : 7+9],
		},
		"K": &lighting.Linear{
			InnerFern: ferns[87],
			OuterFern: ferns[23],
			LEDs:      devices[33].LEDs[1][7+9 : 7+9+19],
		},
	}

	for _, linear := range linears {
		if linear.InnerFern != nil {
			linear.InnerFern.OuterLinears = append(linear.InnerFern.OuterLinears,
				linear)
		}

		if linear.OuterFern != nil {
			linear.OuterFern.InnerLinear = linear
		}
	}

	system.Root = []*lighting.Linear{
		linears["A"],
	}

	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[0]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[1]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[2]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[3]...)

	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[0]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[1]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[2]...)
	// system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[3]...)

	// system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[0]...)
	// system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[1]...)
	// system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[2]...)
	// system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[3]...)
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
