package scan

import (
	"encoding/gob"
	"math"

	"github.com/1lann/sweep"
	"github.com/rickqu/Vivid21Software/exhibit%20controller/geo"
)

var maxArc = 30.0
var maxTotalWidth = 50.0
var maxTotalDiff = 20
var maxDiff = 20
var distComp = 0

// Scanner represents a simplified, people scanner.
type Scanner struct {
	minAngle float64
	maxAngle float64
	device   *sweep.Device
	scan     <-chan sweep.Scan
}

// PointMeta represents metaadata regarding a point that is used for the Data
// field of a *geo.Point.
// The pointer representation will be used in the Data field. For instance:
//     meta := p.Data.(*scan.PointMeta)
type PointMeta struct {
	Strength int
}

func init() {
	gob.Register(PointMeta{})
}

// SetupScanner sets up the scanner at the given path.
func SetupScanner(path string, minAngle, maxAngle float64) (*Scanner, error) {
	dev, err := sweep.NewDevice(path)
	if err != nil {
		return nil, err
	}

	dev.Reset()
	dev.WaitUntilMotorReady()
	dev.Close()

	dev, err = sweep.NewDevice(path)
	if err != nil {
		return nil, err
	}

	dev.WaitUntilMotorReady()
	dev.SetMotorSpeed(2)
	dev.WaitUntilMotorReady()
	dev.SetSampleRate(sweep.Rate500)
	dev.WaitUntilMotorReady()

	scan, err := dev.StartScan()
	if err != nil {
		return nil, err
	}

	return &Scanner{minAngle, maxAngle, dev, scan}, nil
}

// ScanPeople scans people into the given map.
func (s *Scanner) ScanPeople(crowd *geo.Map) {
	var lastPoint *sweep.ScanSample

	var aggregate struct {
		startAngle float64
		startDist  int

		count    int
		strength int
		sumX     float64
		sumY     float64
	}

	scan := <-s.scan
	crowd.Lock()
	defer crowd.Unlock()
	crowd.Clear()

	for _, point := range scan {
		if point.Angle > s.maxAngle || point.Angle < s.minAngle {
			continue
		}

		if lastPoint == nil {
			aggregate.count = 1
			aggregate.strength = int(point.SignalStrength)
			aggregate.sumX, aggregate.sumY = point.Cartesian()
			aggregate.startDist = point.Distance
			aggregate.startAngle = point.Rad()
			lastPoint = point
			continue
		}

		if abs(point.Distance-lastPoint.Distance) > maxDiff ||
			math.Abs(lastPoint.Rad()-point.Rad())*float64(point.Distance) > maxArc ||
			abs(point.Distance-aggregate.startDist) > maxTotalDiff ||
			math.Abs(point.Rad()-aggregate.startAngle)*float64(point.Distance) > maxTotalWidth {
			// fmt.Println(aggregate.count)
			if distComp-(lastPoint.Distance/100) < aggregate.count {
				crowd.Add(geo.NewPoint(
					int(aggregate.sumX/float64(aggregate.count)),
					int(aggregate.sumY/float64(aggregate.count)),
					&PointMeta{
						Strength: aggregate.strength / aggregate.count,
					},
				))
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

	if lastPoint == nil {
		return
	}

	// fmt.Println("strength:", aggregate.strength)

	if distComp-(lastPoint.Distance/100) < aggregate.count {
		crowd.Add(geo.NewPoint(
			int(aggregate.sumX/float64(aggregate.count)),
			int(aggregate.sumY/float64(aggregate.count)),
			&PointMeta{
				Strength: aggregate.strength / aggregate.count,
			},
		))
	}
}

func abs(a int) int {
	if a < 0 {
		return -1 * a
	}
	return a
}
