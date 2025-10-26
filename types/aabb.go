package types

// AABB represents an axis-aligned bounding box in 2D space.
//
// The bounds are inclusive on all sides. An AABB is valid when
// Min.X <= Max.X and Min.Y <= Max.Y. Empty or inverted AABBs
// should be handled explicitly by the caller.
//
// Example:
//
//	box := types.AABB{
//	    Min: types.Point{X: 0.0, Y: 0.0},
//	    Max: types.Point{X: 10.0, Y: 10.0},
//	}
type AABB struct {
	Min Point // Minimum (bottom-left) corner, inclusive
	Max Point // Maximum (top-right) corner, inclusive
}
