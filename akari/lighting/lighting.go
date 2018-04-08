package lighting

import "time"

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
}

// RunEffects runs all of the effects in the system.
func (s *System) RunEffects() {
	for _, effect := range s.RunningEffects {
		effect.Run(s)
	}
}
