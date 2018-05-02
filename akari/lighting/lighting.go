package lighting

import (
	"time"
)

// Effect represents the effect.
type Effect interface {
	Start() time.Time
	Deadline() time.Time
	Priority() int
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

// RunEffects runs all of the effects in the system.
func (s *System) RunEffects() {
	for key, effect := range s.RunningEffects {
		if effect.Deadline().Before(time.Now()) {
			delete(s.RunningEffects, key)
			continue
		}

		effect.Run(s)
	}
}
