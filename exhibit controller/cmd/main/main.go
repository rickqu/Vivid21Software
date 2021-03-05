package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/geo"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/lighting"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/mapping"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/netscan"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/report"

	_ "github.com/rickqu/Vivid21Software/exhibit%20controller/scan"
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

var activate [6]bool
var closeFerns = [][]int{
	{}, {}, {62, 64, 66}, {44, 42, 45}, {26, 24, 25},
}

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

	devices[8] = mapping.NewDevice(8, []int{70, 70, 70, 70})
	devices[9] = mapping.NewDevice(9, []int{120, 70, 70, 70})
	devices[10] = mapping.NewDevice(10, []int{70, 70, 70, 70})

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

	system.AddEffect("breathing", lighting.NewBreathing(1))
	// system.AddEffect("blank", lighting.NewBlank(1))

	reporter := report.NewReporter(mapping.Conn, logger)

	// var edgeFerns = []int{66, 62, 64, 53, 42, 35, 32, 24, 26, 73, 25, 16, 13}

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
					}, ferns[84], 4, lighting.NeuralStepTime, lighting.NeuralEffectRadius, true))
					treeLast = time.Now()
				}

				if time.Since(neuralLast) > 1000*time.Millisecond {
					// n := rand.Intn(len(edgeFerns))
					for i := 2; i <= 5; i++ {
						if activate[i] {
							fmt.Println("activated!")
							for _, fern := range closeFerns[i] {
								system.AddEffect(uuid.New().String(), lighting.NewNeural(color.RGBA{
									R: 0xff,
									G: 0xff,
									B: 0xff,
								}, ferns[fern], 4, lighting.NeuralStepTime, lighting.NeuralEffectRadius, false))
							}
						}
					}

					neuralLast = time.Now()
				}

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
		fmt.Println("scan")
		results := receiver.GetAll()

		if results[2] != nil {
			payload := &Payload{
				Ferns:  []*Fern{},
				Sensor: results[2].Points,
			}
			lisMutex.Lock()
			for _, lis := range listeners {
				lis <- payload
			}
			lisMutex.Unlock()
		}

		for i := 2; i <= 5; i++ {
			if results[i] == nil {
				continue
			}

			if len(results[i].Within(&geo.Point{X: 0, Y: 0}, 170)) > 0 {
				fmt.Println("activated!")
				activate[i] = true
			} else {
				activate[i] = false
			}
		}
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
			LEDs:      reverseLEDs(devices[13].LEDs[1][32 : 32+11]),
		},
		"A1B": &lighting.Linear{
			InnerFern: ferns[11],
			OuterFern: ferns[12],
			LEDs:      reverseLEDs(devices[13].LEDs[1][16 : 16+16]),
		},
		"A1C": &lighting.Linear{
			InnerFern: ferns[12],
			OuterFern: ferns[13],
			LEDs:      reverseLEDs(devices[13].LEDs[1][0:16]),
		},

		"A2A": &lighting.Linear{
			OuterFern: ferns[14],
			LEDs:      reverseLEDs(devices[16].LEDs[1][26:36]),
		},
		"A2B": &lighting.Linear{
			InnerFern: ferns[14],
			OuterFern: ferns[15],
			LEDs:      reverseLEDs(devices[16].LEDs[1][13 : 13+13]),
		},
		"A2C": &lighting.Linear{
			InnerFern: ferns[15],
			OuterFern: ferns[16],
			LEDs:      reverseLEDs(devices[16].LEDs[1][0:13]),
		},

		"B1A": &lighting.Linear{
			OuterFern: ferns[21],
			LEDs:      reverseLEDs(devices[24].LEDs[1][27 : 27+9]),
		},
		"B1B": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[22],
			LEDs:      reverseLEDs(devices[24].LEDs[1][20 : 20+7]),
		},
		"B1C": &lighting.Linear{
			InnerFern: ferns[22],
			OuterFern: ferns[23],
			LEDs:      reverseLEDs(devices[24].LEDs[1][13 : 13+7]),
		},
		"B1D": &lighting.Linear{
			InnerFern: ferns[23],
			OuterFern: ferns[24],
			LEDs:      reverseLEDs(devices[24].LEDs[1][0:13]),
		},

		"B2A": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[26],
			LEDs:      reverseLEDs(devices[25].LEDs[1][0:17]),
		},
		"B2B": &lighting.Linear{
			InnerFern: ferns[21],
			OuterFern: ferns[25],
			LEDs:      devices[25].LEDs[1][17 : 17+18],
		},

		"C1A": &lighting.Linear{
			OuterFern: ferns[31],
			LEDs:      reverseLEDs(devices[35].LEDs[1][27 : 27+6]),
		},
		"C1B": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[33],
			LEDs:      reverseLEDs(devices[35].LEDs[1][18 : 18+9]),
		},
		"C1C": &lighting.Linear{
			InnerFern: ferns[33],
			OuterFern: ferns[34],
			LEDs:      reverseLEDs(devices[35].LEDs[1][11 : 11+7]),
		},
		"C1D": &lighting.Linear{
			InnerFern: ferns[34],
			OuterFern: ferns[35],
			LEDs:      reverseLEDs(devices[35].LEDs[1][0:11]),
		},

		"C2A": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[32],
			LEDs:      devices[36].LEDs[1][16 : 16+17],
		},
		"C2B": &lighting.Linear{
			InnerFern: ferns[31],
			OuterFern: ferns[36],
			LEDs:      reverseLEDs(devices[36].LEDs[1][0:16]),
		},

		"D1A": &lighting.Linear{
			OuterFern: ferns[41],
			LEDs:      reverseLEDs(devices[42].LEDs[1][12 : 12+13]),
		},
		"D1B": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[42],
			LEDs:      reverseLEDs(devices[42].LEDs[1][0:12]),
		},

		"D2FUCK": &lighting.Linear{
			InnerFern: ferns[43],
			OuterFern: ferns[44],
			LEDs:      reverseLEDs(devices[44].LEDs[1][0:13]),
		},
		"D2A": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[43],
			LEDs:      reverseLEDs(devices[44].LEDs[1][13 : 13+14]),
		},
		"D2B": &lighting.Linear{
			InnerFern: ferns[41],
			OuterFern: ferns[45],
			LEDs:      devices[44].LEDs[1][13+14 : 13+14+13],
		},

		"E1A": &lighting.Linear{
			OuterFern: ferns[51],
			LEDs:      reverseLEDs(devices[51].LEDs[1][0:13]),
		},
		"E1B": &lighting.Linear{
			OuterFern: ferns[52],
			LEDs:      devices[51].LEDs[1][13 : 13+16],
		},
		"E1C": &lighting.Linear{
			InnerFern: ferns[52],
			OuterFern: ferns[53],
			LEDs:      devices[51].LEDs[1][13+16 : 13+16+10],
		},

		"F1A": &lighting.Linear{
			OuterFern: ferns[61],
			LEDs:      reverseLEDs(devices[62].LEDs[1][12+14 : 12+14+11]),
		},
		"F1B": &lighting.Linear{
			InnerFern: ferns[61],
			OuterFern: ferns[63],
			LEDs:      reverseLEDs(devices[62].LEDs[1][12 : 12+14]),
		},
		"F1C": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[62],
			LEDs:      reverseLEDs(devices[62].LEDs[1][0:12]),
		},

		"F2A": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[64],
			LEDs:      reverseLEDs(devices[64].LEDs[1][0:12]),
		},
		"F2B": &lighting.Linear{
			InnerFern: ferns[63],
			OuterFern: ferns[65],
			LEDs:      devices[64].LEDs[1][12 : 12+10],
		},
		"F2C": &lighting.Linear{
			InnerFern: ferns[65],
			OuterFern: ferns[66],
			LEDs:      devices[64].LEDs[1][12+10 : 12+10+13],
		},

		"G1A": &lighting.Linear{
			OuterFern: ferns[71],
			LEDs:      reverseLEDs(devices[73].LEDs[1][13+13 : 13+13+13]),
		},
		"G1B": &lighting.Linear{
			InnerFern: ferns[71],
			OuterFern: ferns[72],
			LEDs:      reverseLEDs(devices[73].LEDs[1][13 : 13+13]),
		},
		"G1C": &lighting.Linear{
			InnerFern: ferns[72],
			OuterFern: ferns[73],
			LEDs:      reverseLEDs(devices[73].LEDs[1][0:13]),
		},

		"H1A": &lighting.Linear{
			InnerFern: ferns[84],
			OuterFern: ferns[83],
			LEDs:      reverseLEDs(devices[81].LEDs[1][14+16 : 14+16+13]),
		},
		"H1B": &lighting.Linear{
			InnerFern: ferns[83],
			OuterFern: ferns[82],
			LEDs:      reverseLEDs(devices[81].LEDs[1][16 : 16+14]),
		},
		"H1C": &lighting.Linear{
			InnerFern: ferns[82],
			OuterFern: ferns[81],
			LEDs:      reverseLEDs(devices[81].LEDs[1][0:16]),
		},

		"H2A": &lighting.Linear{
			InnerFern: ferns[84],
			OuterFern: ferns[85],
			LEDs:      reverseLEDs(devices[87].LEDs[1][9+11 : 9+11+12]),
		},
		"H2B": &lighting.Linear{
			InnerFern: ferns[85],
			OuterFern: ferns[86],
			LEDs:      reverseLEDs(devices[87].LEDs[1][9 : 9+11]),
		},
		"H2C": &lighting.Linear{
			InnerFern: ferns[86],
			OuterFern: ferns[87],
			LEDs:      reverseLEDs(devices[87].LEDs[1][0:9]),
		},

		"HFAKE": &lighting.Linear{
			OuterFern: ferns[84],
			LEDs:      nil,
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
		linears["C1A"],
		linears["B1A"],
		linears["A2A"],
		linears["A1A"],
		linears["G1A"],
		linears["F1A"],
		linears["E1A"],
		linears["E1B"],
		linears["D1A"],
		linears["HFAKE"],
	}

	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[0]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[1]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[2]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[8].LEDs[3]...)

	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[0]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[1]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[2]...)
	system.TreeBase.LEDs = append(system.TreeBase.LEDs, devices[10].LEDs[3]...)

	system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[0]...)
	system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[1]...)
	system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[2]...)
	system.TreeTop.LEDs = append(system.TreeTop.LEDs, devices[9].LEDs[3]...)
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
