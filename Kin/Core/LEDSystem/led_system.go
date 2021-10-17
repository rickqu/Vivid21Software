package LEDSystem

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/StartCode"
)

type LEDSystem interface {
	Setup(chan<- LightCommand, StartCode.StartCode)
	Start()
	Stop()
}
