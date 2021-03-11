package SensorArray

type SensorArray interface {
	Calibrate()
	Run()
	Stop()
	GetOutputChannel() chan SensorDatapoint
}

type OutputPusher struct {
	outputChannel chan SensorDatapoint
}
