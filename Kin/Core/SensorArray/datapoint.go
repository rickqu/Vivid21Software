package SensorArray

var lastDatapointId int32

type SensorDatapoint struct {
	datapointId int32
	sensorName  string
	datapoints  []float64
}

func Init() {
	lastDatapointId = 0
}

func NewSensorDatapoint(sensorName string, datapoint []float64) *SensorDatapoint {
	newDatapoint := SensorDatapoint{lastDatapointId, sensorName, datapoint}
	lastDatapointId += 1
	return &newDatapoint
}

func (d *SensorDatapoint) NumSamples() int {
	return len(d.datapoints)
}
