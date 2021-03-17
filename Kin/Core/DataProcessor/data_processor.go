package DataProcessor

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/LightController"
	"github.com/rickqu/Vivid21Software/Kin/Core/SensorArray"
)

type DataProcessor interface {
	StartProcessor(sensorInput chan SensorArray.SensorDatapoint,
		lightingOutput chan LightController.LightCommand)
	StopProcessor()
}
