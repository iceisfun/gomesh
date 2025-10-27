package geometry

import (
	"math"

	"github.com/iceisfun/gomesh/algorithm/robust"
	"github.com/iceisfun/gomesh/types"
)

const bboxTol = 1e-12

// Area2 computes twice the signed area of triangle (a,b,c).
func Area2(a, b, c types.Point) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

// PointOnSegment reports whether point p lies on the closed segment [a,b].
func PointOnSegment(p, a, b types.Point) bool {
	if robust.Orient2D(a, b, p) != 0 {
		return false
	}

	minX := math.Min(a.X, b.X) - bboxTol
	maxX := math.Max(a.X, b.X) + bboxTol
	minY := math.Min(a.Y, b.Y) - bboxTol
	maxY := math.Max(a.Y, b.Y) + bboxTol

	return p.X >= minX && p.X <= maxX && p.Y >= minY && p.Y <= maxY
}

// DistancePointSegment computes the shortest distance between a point and a segment.
func DistancePointSegment(p, a, b types.Point) float64 {
	ax := b.X - a.X
	ay := b.Y - a.Y
	length2 := ax*ax + ay*ay
	if length2 == 0 {
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	}

	// Project point p onto segment ab.
	t := ((p.X-a.X)*ax + (p.Y-a.Y)*ay) / length2
	switch {
	case t <= 0:
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	case t >= 1:
		return math.Hypot(p.X-b.X, p.Y-b.Y)
	default:
		proj := types.Point{
			X: a.X + t*ax,
			Y: a.Y + t*ay,
		}
		return math.Hypot(p.X-proj.X, p.Y-proj.Y)
	}
}

// Centroid returns the centroid of triangle (a,b,c).
func Centroid(a, b, c types.Point) types.Point {
	return types.Point{
		X: (a.X + b.X + c.X) / 3,
		Y: (a.Y + b.Y + c.Y) / 3,
	}
}

// BBox computes the axis-aligned bounding box of the supplied loop.
func BBox(loop []types.Point) types.AABB {
	if len(loop) == 0 {
		return types.AABB{}
	}

	minX, maxX := loop[0].X, loop[0].X
	minY, maxY := loop[0].Y, loop[0].Y
	for _, p := range loop[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return types.AABB{
		Min: types.Point{X: minX, Y: minY},
		Max: types.Point{X: maxX, Y: maxY},
	}
}
