package SensorArray

import "github.com/rickqu/Vivid21Software/Kin/Core/StartCode"

type SensorArray interface {
	AddSensor(string, Sensor)
	CalibrateAll()
	RunAll(startCode StartCode.StartCode)
	StopAll()
	StopSensor(string)
	GetOutputFromChannel() SensorDatapoint
}

type sensorArray struct {
	dataOutputChannel <-chan SensorDatapoint
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

func (s sensorArray) RunAll(startCode StartCode.StartCode) {
	for _, sensor := range s.sensorMap {
		go sensor.Run(s.dataOutputChannel)
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

func NewSensorArray(dataOutputChannel <-chan SensorDatapoint) SensorArray {
	newArray := sensorArray{}
	newArray.dataOutputChannel = dataOutputChannel
	newArray.sensorMap = make(map[string]Sensor)
	return newArray
}
