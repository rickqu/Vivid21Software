package Core

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/DataProcessor"
	"github.com/rickqu/Vivid21Software/Kin/Core/LEDDriver"
	"github.com/rickqu/Vivid21Software/Kin/Core/LEDSystem"
	"github.com/rickqu/Vivid21Software/Kin/Core/SensorArray"
	"github.com/rickqu/Vivid21Software/Kin/Core/StartCode"
)

const chan_max_size int = 10

type Core struct {
	sensorArray SensorArray.SensorArray
	processor   DataProcessor.DataProcessor
	ledSystem   LEDSystem.LEDSystem
	ledDriver   LEDDriver.LEDDriver

	sensorDataChan chan SensorArray.SensorDatapoint
	lightCommands  chan LEDSystem.LightCommand
}

func NewCore(sensorArray SensorArray.SensorArray,
	processor DataProcessor.DataProcessor,
	ledSystem LEDSystem.LEDSystem,
	ledDriver LEDDriver.LEDDriver) *Core {

	core := new(Core)
	core.sensorArray = sensorArray
	core.processor = processor
	core.ledSystem = ledSystem
	core.ledDriver = ledDriver

	core.sensorDataChan = make(chan SensorArray.SensorDatapoint, chan_max_size)
	core.lightCommands = make(chan LEDSystem.LightCommand, chan_max_size)

	return core
}

func (c *Core) Setup(startCode StartCode.StartCode) {
	// do we need a mapper between driver and system?
	c.ledDriver.Initialise()
	c.ledSystem.Setup(c.lightCommands, startCode)
	c.processor.Setup(c.sensorDataChan, c.lightCommands, startCode)
	c.sensorArray.CalibrateAll()
}

func (c *Core) Start(startCode StartCode.StartCode) {
	go c.ledSystem.Start()
	go c.processor.Start()
	c.sensorArray.RunAll(startCode)
}

func (c *Core) Stop() {

}
