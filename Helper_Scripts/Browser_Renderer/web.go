package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/1lann/sweep"
	"github.com/bugra/kmeans"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	upgrader  = websocket.Upgrader{}
	mutex     = new(sync.Mutex)
	listeners = make(map[string]chan<- []scanData)
)

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

func main() {
	e := echo.New()

	dev, err := sweep.NewDevice("/dev/cu.usbserial-DO0088ZE")
	if err != nil {
		panic(err)
	}

	fmt.Println("Stopping scan...")

	dev.StopScan()
	dev.Drain()
	// dev.SetMotorSpeed(2)
	// dev.SetMotorSpeed(2)
	// dev.WaitUntilMotorReady()
	// dev.SetSampleRate(sweep.Rate500)
	fmt.Println("Waiting ready")
	dev.WaitUntilMotorReady()

	fmt.Println("Starting scan")

	scanner, err := dev.StartScan()
	if err != nil {
		panic(err)
	}

	fmt.Println("Scan started")

	go func() {
		for scan := range scanner {
			result := getHumans(scan)
			mutex.Lock()
			for _, lis := range listeners {
				lis <- result
			}
			mutex.Unlock()
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

func getHumans(scan sweep.Scan) []scanData {
	var data [][]float64

	for _, point := range scan {
		rad := (point.Angle * math.Pi) / 180.0

		if math.Cos(rad)*float64(point.Distance) > 200 {
			continue
		}

		data = append(data, []float64{math.Cos(rad) * float64(point.Distance),
			-math.Sin(rad) * float64(point.Distance)})
	}

	result, err := kmeans.Kmeans(data, 10, kmeans.SquaredEuclideanDistance, 20)
	if err != nil {
		log.Println("trash:", err)
		return nil
	}

	var returnResult []scanData
	for index, res := range result {
		returnResult = append(returnResult, scanData{
			X:        int(data[index][0]),
			Y:        int(data[index][1]),
			Color:    colors[res],
			Strength: int(scan[index].SignalStrength),
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
