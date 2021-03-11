package SensorArray

var lastDatapointId int32

type SensorDatapoint struct {
	datapointId int32
	datapoints  []float64
}

func Init() {
	lastDatapointId = 0
}

func NewSensorDatapoint(datapoint []float64) *SensorDatapoint {
	newDatapoint := SensorDatapoint{lastDatapointId, datapoint}
	lastDatapointId += 1
	return &newDatapoint
}

func (d *SensorDatapoint) NumSamples() int {
	return len(d.datapoints)
}

func (d * SensorDatapoint) GetSamples() int [] {
	return d.datapoints
}
