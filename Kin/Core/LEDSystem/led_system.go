package LEDSystem

import (
	"github.com/rickqu/Vivid21Software/Kin/Core/StartCode"
)

type LEDSystem interface {
	Start(chan LightCommand, StartCode.StartCode)
	Stop()
}
