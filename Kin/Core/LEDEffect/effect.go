package LEDEffect

import (
	"time"

	"github.com/rickqu/Vivid21Software/Kin/Core/LEDSystem"
)

// Effect represents the effect.
type Effect interface {
	Start() time.Time
	Priority() int
	Active() bool
	Run(matrix *LEDSystem.LEDMatrix)
}
