package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/1lann/rpc"
	"github.com/pul-s4r/vivid18/akari/geo"
	"github.com/pul-s4r/vivid18/akari/scan"
)

var client *rpc.Client

func main() {
	// minAngle, err := strconv.Atoi(os.Args[3])
	// if err != nil {
	// 	panic(err)
	// }

	// maxAngle, err := strconv.Atoi(os.Args[4])
	// if err != nil {
	// 	panic(err)
	// }

	scanner, err := scan.SetupScanner(os.Args[2], float64(-1000), float64(1000))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			geoMap := geo.NewMap()
			scanner.ScanPeople(geoMap)
			if client != nil {
				client.Fire("scan-"+os.Args[1], geoMap)
				fmt.Println("completed scan:", len(geoMap.Points))
			}
		}
	}()

	for {
		conn, err := net.DialTimeout("tcp", "192.168.2.1:5555", 2*time.Second)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}

		client, err = rpc.NewClient(conn)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}

		err = client.Receive()
		if err != nil {
			log.Println(":", err)
		}
	}
}
