package main

import (
	"fmt"
	"net"
	"time"
)

// Ports used
const (
	ServerPort = 5050
	DevicePort = 5151
)

const timeout = 5 * time.Second

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, 1),
		Port: ServerPort,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Discovering devices...")

	foundDevices := make(map[string]bool)

	start := time.Now()
	conn.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1000)
	for {
		_, ip, err := conn.ReadFromUDP(buf)
		if time.Since(start) >= timeout {
			break
		}
		if err != nil {
			continue
		}

		if ip.Port != DevicePort {
			continue
		}

		if _, found := foundDevices[ip.IP.String()]; !found {
			foundDevices[ip.IP.String()] = true
			fmt.Println("Found:", ip.IP.String())
		}
	}

	if len(foundDevices) == 0 {
		fmt.Println("No devices could be found!")
	}
}
