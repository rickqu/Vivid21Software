package DataProcessor

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/LEDSystem"
	"github.com/rickqu/Vivid21Software/Kin/Core/SensorArray"
	"github.com/rickqu/Vivid21Software/Kin/Core/StartCode"
)

type DataProcessor interface {
	Setup(sensorInput chan<- SensorArray.SensorDatapoint,
		lightingOutput <-chan LEDSystem.LightCommand,
		startCode StartCode.StartCode)
	Start()
	Stop()
}
