package SensorArray

type SensorArray interface {
	AddSensor(string, Sensor)
	CalibrateAll()
	RunAll()
	StopAll()
	StopSensor(string)
	GetOutputFromChannel() SensorDatapoint
}

const output_chan_max_size int = 10

type sensorArray struct {
	dataOutputChannel chan SensorDatapoint
	sensorMap         map[string]Sensor
}

func (s sensorArray) AddSensor(name string, sensor Sensor) {
	s.sensorMap[name] = sensor
}

func (s sensorArray) CalibrateAll() {
	for _, sensor := range s.sensorMap {
		sensor.Calibrate()
	}
}

func (s sensorArray) RunAll() {
	for _, sensor := range s.sensorMap {
		sensor.Run(s.dataOutputChannel)
	}
}

func (s sensorArray) StopAll() {
	for _, sensor := range s.sensorMap {
		sensor.Stop()
	}
}

func (s sensorArray) StopSensor(name string) {
	s.sensorMap[name].Stop()
}

func (s sensorArray) GetOutputFromChannel() SensorDatapoint {
	return <-s.dataOutputChannel
}

func NewSensorArray() SensorArray {
	newArray := sensorArray{}
	newArray.dataOutputChannel = make(chan SensorDatapoint, output_chan_max_size)
	newArray.sensorMap = make(map[string]Sensor)
	return newArray
}
