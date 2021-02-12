package lighting

import (
	"time"
)

// Effect represents the effect.
type Effect interface {
	Start() time.Time
	Priority() int
	Active() bool
	Run(system *System)
}

// System represents the system.
type System struct {
	RunningEffects map[string]Effect
	Root           []*Linear
	TreeTop        *TreeTop
	TreeBase       *TreeBase
}

// NewSystem returns a new lighting system.
func NewSystem() *System {
	// TODO: actually take in args, or setup Root, TreeTop, TreeBase
	// or something.

	return &System{
		RunningEffects: make(map[string]Effect),
		TreeTop:        &TreeTop{},
		TreeBase:       &TreeBase{},
	}
}

// AddEffect adds an effect in the system.
func (s *System) AddEffect(id string, effect Effect) {
	s.RunningEffects[id] = effect
}

// RemoveEffect removes an effect in the system.
func (s *System) RemoveEffect(id string) {
	delete(s.RunningEffects, id)
}

// Run runs all of the effects in the system.
func (s *System) Run() {
	// TODO: Set all LEDs to black?

	for priority := 1; priority <= 5; priority++ {
		for key, effect := range s.RunningEffects {
			if !effect.Active() {
				delete(s.RunningEffects, key)
				continue
			}

			if effect.Priority() == priority {
				effect.Run(s)
			}
		}
	}
}
