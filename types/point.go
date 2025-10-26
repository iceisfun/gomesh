package types

// Point represents a position in 2D Cartesian space.
//
// Coordinates use float64 precision, suitable for most geometric
// applications with appropriate epsilon tolerance for comparisons.
//
// Example:
//
//	p := types.Point{X: 1.5, Y: 2.3}
//	q := types.Point{X: 0.0, Y: 0.0}
type Point struct {
	X float64 // Horizontal coordinate
	Y float64 // Vertical coordinate
}
