package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/1lann/sweep"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/tarm/serial"
)

var (
	upgrader  = websocket.Upgrader{}
	mutex     = new(sync.Mutex)
	listeners = make(map[string]chan<- []scanData)
)

var target int
var ard *serial.Port

type scanData struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
	Strength int    `json:"s"`
}

var colors []string

func init() {
	for i := 0; i < 100; i++ {
		col := colorful.FastHappyColor()
		r, g, b := col.RGB255()
		colors = append(colors, fmt.Sprintf("rgba(%d, %d, %d, ", r, g, b))
		fmt.Printf("rgba(%d, %d, %d, \n", r, g, b)
	}
}

func lightLoop() {
	t := time.NewTicker(50 * time.Millisecond)
	var current int
	circle := 0.0

	i := 0

	for range t.C {
		i++

		// if i%10 == 0 {
		fmt.Println(current, target)
		// }

		circle += 5.0
		if circle > 360 {
			circle = 0
		}

		if target < current {
			current -= 5
			if target > current {
				current = target
			}
		} else if target > current {
			current += 5
			if target < current {
				current = target
			}
		}

		c := colorful.Hsv(circle, 1, float64(current)/255.0)

		for i := 0; i < 40; i++ {
			r, g, b := c.RGB255()
			ard.Write([]byte{r, g, b})
		}
		ard.Flush()
	}
}

func main() {
	// f, err := os.Create("./cpu.profile")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func() {
	// 	for range c {
	// 		pprof.StopCPUProfile()
	// 		f.Close()
	// 		os.Exit(0)
	// 	}
	// }()

	e := echo.New()

	dev, err := sweep.NewDevice("/dev/cu.usbserial-DO0088ZE")
	if err != nil {
		panic(err)
	}

	ard, err = serial.OpenPort(&serial.Config{
		Name:        "/dev/cu.usbmodem1411",
		Baud:        115200,
		ReadTimeout: 10 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			io.Copy(os.Stdout, ard)
		}
	}()

	fmt.Println("Stopping scan...")

	dev.StopScan()
	dev.Drain()
	// dev.SetMotorSpeed(2)
	// dev.SetMotorSpeed(2)
	// dev.WaitUntilMotorReady()
	// dev.SetSampleRate(sweep.Rate500)
	// fmt.Println("Waiting ready")
	// dev.WaitUntilMotorReady()

	fmt.Println("Starting scan")

	scanner, err := dev.StartScan()
	if err != nil {
		panic(err)
	}

	fmt.Println("Scan started")

	go lightLoop()

	go func() {
		for scan := range scanner {
			result := getHumans(scan)
			mutex.Lock()
			for _, lis := range listeners {
				lis <- result
			}
			mutex.Unlock()

			min := 100
			for _, res := range result {
				dist := math.Sqrt(float64(res.X*res.X + res.Y*res.Y))
				if dist < float64(min) {
					min = int(dist)
				}
			}

			target = int(float64(80-(min-20)) * 3.2)
		}

		mutex.Lock()
		for _, lis := range listeners {
			close(lis)
		}
		mutex.Unlock()
	}()

	go func() {
		e.GET("/ws", wsHandler)
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

	lis := make(chan []scanData, 100)
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

func abs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}

func getHumans(scan sweep.Scan) []scanData {
	type Aggregate struct {
		count int
		sumX  float64
		sumY  float64
	}

	maxArc := 30.0
	maxTotalWidth := 50.0
	maxTotalDiff := 20
	maxDiff := 20
	distComp := 0

	var lastPoint *sweep.ScanSample

	var aggregate struct {
		startAngle float64
		startDist  int

		count    int
		strength int
		sumX     float64
		sumY     float64
	}

	var returnResult []scanData

	for _, point := range scan {
		// ax, ay := point.Cartesian()
		// returnResult = append(returnResult, scanData{
		// 	X:        int(ax),
		// 	Y:        int(ay),
		// 	Color:    "rgba(255, 0, 0, ",
		// 	Strength: 0,
		// })
		// continue

		if lastPoint == nil {
			aggregate.count = 1
			aggregate.strength = int(point.SignalStrength)
			aggregate.sumX, aggregate.sumY = point.Cartesian()
			aggregate.startDist = point.Distance
			aggregate.startAngle = point.Rad()
			lastPoint = point
			continue
		}

		// fmt.Println("a:", abs(point.Distance-lastPoint.Distance) > maxDiff)
		// fmt.Println("b:", math.Abs(lastPoint.Rad()-point.Rad())*float64(point.Distance) > maxArc)
		// fmt.Println("c:", abs(point.Distance-aggregate.startDist) > maxTotalDiff)
		// fmt.Println("d:", math.Abs(point.Rad()-aggregate.startAngle)*float64(point.Distance))

		if abs(point.Distance-lastPoint.Distance) > maxDiff ||
			math.Abs(lastPoint.Rad()-point.Rad())*float64(point.Distance) > maxArc ||
			abs(point.Distance-aggregate.startDist) > maxTotalDiff ||
			math.Abs(point.Rad()-aggregate.startAngle)*float64(point.Distance) > maxTotalWidth {
			// fmt.Println(aggregate.count)
			if distComp-(lastPoint.Distance/100) < aggregate.count {
				returnResult = append(returnResult, scanData{
					X:        int(aggregate.sumX / float64(aggregate.count)),
					Y:        int(aggregate.sumY / float64(aggregate.count)),
					Color:    "rgba(255, 0, 0, ",
					Strength: aggregate.strength / aggregate.count,
				})
			}

			aggregate.count = 1
			aggregate.strength = int(point.SignalStrength)
			aggregate.sumX, aggregate.sumY = point.Cartesian()
			aggregate.startDist = point.Distance
			aggregate.startAngle = point.Rad()
			lastPoint = point
			continue
		}

		// fmt.Println("diff:", math.Abs(lastPoint.Angle-point.Angle), float64(point.Distance))

		aggregate.count++
		x, y := point.Cartesian()
		aggregate.sumX += x
		aggregate.sumY += y
		aggregate.strength += int(point.SignalStrength)
		lastPoint = point
	}

	// fmt.Println("strength:", aggregate.strength)

	if distComp-(lastPoint.Distance/100) < aggregate.count {
		returnResult = append(returnResult, scanData{
			X:        int(aggregate.sumX / float64(aggregate.count)),
			Y:        int(aggregate.sumY / float64(aggregate.count)),
			Color:    "rgba(255, 0, 0, ",
			Strength: aggregate.strength / aggregate.count,
		})
	}

	return returnResult
}

// func processScan(scan sweep.Scan) sweep.Scan {
// 	sensorP := &geo.Point{
// 		X: 100,
// 		Y: 100,
// 	}
//
// 	m := geo.NewMap()
// 	for _, scanP := range scan {
// 		rad := (scanP.Angle / 180.0) * math.Pi;
//
// 		m.Add(sensorP.Add(&Point{
// 			X: Math.Cos(rad) *
// 		}))
// 	}
// }
