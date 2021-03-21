package SensorArray

type Sensor interface {
	Run(<-chan SensorDatapoint)
	Stop()
	Calibrate()
}
