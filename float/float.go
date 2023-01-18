package float

import "math"

// Epsilon defines a tolerance under which the difference between two float are considered "equal."
const Epsilon = 0.00001

// Equal returns true if the values are exactly equal or are within a difference of Epsilon.
func Equal(a, b float64) bool {
	if math.Abs(a-b) < Epsilon {
		return true
	}
	return false
}
