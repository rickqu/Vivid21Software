package main

import (
	"image/color"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/pul-s4r/vivid18/akari/geo"
	"github.com/pul-s4r/vivid18/akari/lighting"
	"github.com/pul-s4r/vivid18/akari/scan"
)

// {
// 	ferns: [
// 		{
// 			location: {x:, y:},
// 			leds: [
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 				[{r:, g:, b:}, ...],
// 			]
// 		},
// 		...
// 	],
// 	sensor: [
// 		{x:, y:},
// 		...
// 	]
// }

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
	physicalFerns := []*Fern{
		{
			Location: geo.NewPoint(-150, -150),
		},
		{
			Location: geo.NewPoint(0, 0),
		},
		{
			Location: geo.NewPoint(150, 150),
		},
	}

	var ferns []*lighting.Fern

	crowd := geo.NewMap()
	system := lighting.NewSystem()

	for fernID, fern := range physicalFerns {
		for i := 0; i < len(fern.LEDs); i++ {
			for j := 0; j < len(fern.LEDs[i]); j++ {
				fern.LEDs[i][j] = &color.RGBA{}
			}
		}

		f := &lighting.Fern{
			Arms: fern.LEDs,
		}
		ferns = append(ferns, f)

		system.AddEffect(strconv.Itoa(fernID),
			lighting.NewBlob(f, crowd, fern.Location, 270, 180))
	}

	go func() {
		for range time.Tick(30 * time.Millisecond) {
			system.Run()
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

	for {
		scanner.ScanPeople(crowd)
		// fmt.Println("scan completed")
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
