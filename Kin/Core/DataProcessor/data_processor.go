package DataProcessor

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/LEDSystem"
	"github.com/rickqu/Vivid21Software/Kin/Core/SensorArray"
)

type DataProcessor interface {
	StartProcessor(sensorInput chan SensorArray.SensorDatapoint,
		lightingOutput chan LEDSystem.LightCommand)
	StopProcessor()
}
