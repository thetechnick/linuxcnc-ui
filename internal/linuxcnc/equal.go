package linuxcnc

// Checks equallity between two status objects.
func (s *Status) Equal(other *Status) bool {
	if s.EstopEnabled != other.EstopEnabled {
		return false
	}
	if s.InPosition != other.InPosition {
		return false
	}
	if s.MotionPaused != other.MotionPaused {
		return false
	}

	if s.Coolant != other.Coolant {
		return false
	}
	if len(s.Joints) != len(other.Joints) {
		return false
	}
	for i := range s.Joints {
		if s.Joints[i] != other.Joints[i] {
			return false
		}
	}

	if len(s.Spindles) != len(other.Spindles) {
		return false
	}
	for i := range s.Spindles {
		if s.Spindles[i] != other.Spindles[i] {
			return false
		}
	}
	return true
}
