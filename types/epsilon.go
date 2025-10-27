package types

import "math"

// Epsilon stores absolute and relative tolerances for geometric operations.
//
// The combined tolerance for a coordinate with magnitude |v| is computed as:
//
//	tol(v) = Abs + Rel * |v|
//
// For operations involving multiple coordinates or points, the maximum absolute
// coordinate magnitude is used. Negative tolerance values are automatically
// clamped to zero.
type Epsilon struct {
	Abs float64
	Rel float64
}

// NewEpsilon constructs an Epsilon value with the provided parameters.
func NewEpsilon(abs, rel float64) Epsilon {
	return Epsilon{Abs: abs, Rel: rel}.normalized()
}

// DefaultEpsilon returns a conservative default tolerance that roughly matches
// the mesh package defaults.
func DefaultEpsilon() Epsilon {
	return Epsilon{Abs: 1e-9, Rel: 1e-12}
}

// WithAbs returns a copy with the absolute tolerance replaced.
func (e Epsilon) WithAbs(abs float64) Epsilon {
	e.Abs = abs
	return e.normalized()
}

// WithRel returns a copy with the relative tolerance replaced.
func (e Epsilon) WithRel(rel float64) Epsilon {
	e.Rel = rel
	return e.normalized()
}

// Value computes the combined tolerance for the supplied coordinate magnitude.
//
// This is a low-level helper primarily used by TolForPoints/TolForCoords.
func (e Epsilon) Value(mag float64) float64 {
	e = e.normalized()
	return e.Abs + e.Rel*mag
}

// TolForPoints computes the tolerance to use when comparing any of the given
// points. It takes the maximum absolute coordinate across all points and applies
// the combined tolerance.
func (e Epsilon) TolForPoints(points ...Point) float64 {
	if len(points) == 0 {
		return e.Value(0)
	}

	maxMag := 0.0
	for _, p := range points {
		if mag := math.Abs(p.X); mag > maxMag {
			maxMag = mag
		}
		if mag := math.Abs(p.Y); mag > maxMag {
			maxMag = mag
		}
	}
	return e.Value(maxMag)
}

// TolForCoords computes the tolerance for the supplied coordinate magnitudes.
func (e Epsilon) TolForCoords(values ...float64) float64 {
	if len(values) == 0 {
		return e.Value(0)
	}
	maxMag := 0.0
	for _, v := range values {
		if mag := math.Abs(v); mag > maxMag {
			maxMag = mag
		}
	}
	return e.Value(maxMag)
}

// MergeDistance reports the tolerance used for snapping/merging the supplied
// points. This matches the specification abs + rel*max(|x|,|y|).
func (e Epsilon) MergeDistance(a, b Point) float64 {
	return e.TolForPoints(a, b)
}

func (e Epsilon) normalized() Epsilon {
	if e.Abs < 0 {
		e.Abs = -e.Abs
	}
	if e.Rel < 0 {
		e.Rel = -e.Rel
	}
	return e
}
