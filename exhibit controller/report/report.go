// Package report deals with feedback reporting from the boards. They can
// be used to monitor, detect and diagnose issues.
package report

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pul-s4r/vivid18/akari/mapping"
)

// Possible errors.
var (
	ErrBadAck  = errors.New("report: bad ack byte")
	ErrTimeout = errors.New("report: timeout")
)

// DeviceReport represents a report for a device.
type DeviceReport struct {
	LastReports []int
	LastSeen    time.Time
}

// Reporter represents a system reporter.
type Reporter struct {
	conn *net.UDPConn

	seen      map[int]*DeviceReport
	seenMutex *sync.Mutex

	ackListeners  map[int][]chan<- byte
	listenerMutex *sync.Mutex
}

// NewReporter returns a new data reporter from a listener connection.
// It runs a goroutine in the background that listens on conn, so nothing
// else should listen on conn.
func NewReporter(conn *net.UDPConn, logger *logrus.Logger) *Reporter {
	r := &Reporter{
		conn:      conn,
		seen:      make(map[int]*DeviceReport),
		seenMutex: new(sync.Mutex),

		ackListeners:  make(map[int][]chan<- byte),
		listenerMutex: new(sync.Mutex),
	}

	go func() {
		buf := make([]byte, 1600)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				break
			}

			if n < 2 {
				continue
			}

			id := int(addr.IP.To4()[3])
			r.seenMutex.Lock()
			if report, found := r.seen[id]; found {
				report.LastSeen = time.Now()
			} else {
				r.seen[id] = &DeviceReport{
					LastSeen:    time.Now(),
					LastReports: make([]int, 10),
				}
				fmt.Println("found:", id)
			}

			device := r.seen[id]

			switch buf[0] {
			case 'C':
				device.LastReports = append(device.LastReports[1:], int(buf[1]))
			case 'A':
				r.listenerMutex.Lock()
				for _, listen := range r.ackListeners[id] {
					listen <- buf[1]
				}
				r.listenerMutex.Unlock()
			default:
				logger.WithFields(logrus.Fields{
					"ip":     addr.IP.String(),
					"packet": hex.EncodeToString(buf[:n]),
				}).Warn("unexpected packet from device")
			}

			r.seenMutex.Unlock()
		}
	}()

	return r
}

// SeenIDs returns all the seen IDs.
func (r *Reporter) SeenIDs() []int {
	results := make([]int, 0, len(r.seen))

	r.seenMutex.Lock()
	for id := range r.seen {
		results = append(results, id)
	}
	r.seenMutex.Unlock()

	return results
}

// GetReport returns the report of an ID. Returns nil if the ID has never
// been seen.
func (r *Reporter) GetReport(id int) *DeviceReport {
	var result DeviceReport

	r.seenMutex.Lock()
	defer r.seenMutex.Unlock()

	report, found := r.seen[id]
	if !found {
		return nil
	}

	result.LastReports = make([]int, len(report.LastReports))
	copy(result.LastReports, report.LastReports)
	result.LastSeen = report.LastSeen

	return &result
}

// Ping pings the device at the given ID, and returns nil if a correct
// reply is received.
func (r *Reporter) Ping(id int) error {
	listener := make(chan byte, 1)

	r.listenerMutex.Lock()
	r.ackListeners[id] = append(r.ackListeners[id], listener)
	r.listenerMutex.Unlock()

	rbyte := byte(rand.Intn(256))
	r.conn.WriteToUDP([]byte{'S', rbyte}, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 2, byte(id)),
		Port: mapping.DevicePort,
	})

	select {
	case b := <-listener:
		if b != rbyte {
			return ErrBadAck
		}

		return nil
	case <-time.After(time.Second):
		return ErrTimeout
	}
}
