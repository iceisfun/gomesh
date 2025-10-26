package predicates

import (
	"math"

	"github.com/iceisfun/gomesh/types"
)

// PointInAABB tests if a point is inside or on an AABB.
func PointInAABB(p types.Point, box types.AABB, eps float64) bool {
	minX := math.Min(box.Min.X, box.Max.X) - eps
	maxX := math.Max(box.Min.X, box.Max.X) + eps
	minY := math.Min(box.Min.Y, box.Max.Y) - eps
	maxY := math.Max(box.Min.Y, box.Max.Y) + eps

	return p.X >= minX && p.X <= maxX && p.Y >= minY && p.Y <= maxY
}

// SegmentAABBIntersect tests if a line segment intersects an AABB.
func SegmentAABBIntersect(a, b types.Point, box types.AABB, eps float64) bool {
	if PointInAABB(a, box, eps) || PointInAABB(b, box, eps) {
		return true
	}

	minX := math.Min(box.Min.X, box.Max.X) - eps
	maxX := math.Max(box.Min.X, box.Max.X) + eps
	minY := math.Min(box.Min.Y, box.Max.Y) - eps
	maxY := math.Max(box.Min.Y, box.Max.Y) + eps

	segMinX := math.Min(a.X, b.X)
	segMaxX := math.Max(a.X, b.X)
	segMinY := math.Min(a.Y, b.Y)
	segMaxY := math.Max(a.Y, b.Y)

	// Quick rejection if bounding boxes do not overlap.
	if segMaxX < minX || segMinX > maxX || segMaxY < minY || segMinY > maxY {
		return false
	}

	corners := []types.Point{
		{X: minX, Y: minY},
		{X: maxX, Y: minY},
		{X: maxX, Y: maxY},
		{X: minX, Y: maxY},
	}

	edges := [][2]types.Point{
		{corners[0], corners[1]},
		{corners[1], corners[2]},
		{corners[2], corners[3]},
		{corners[3], corners[0]},
	}

	for _, edge := range edges {
		if hit, _ := SegmentsIntersect(a, b, edge[0], edge[1], eps); hit {
			return true
		}
	}

	return false
}

// TriangleAABBIntersect tests if a triangle intersects an AABB.
func TriangleAABBIntersect(a, b, c types.Point, box types.AABB, eps float64) bool {
	if PointInAABB(a, box, eps) || PointInAABB(b, box, eps) || PointInAABB(c, box, eps) {
		return true
	}

	corners := []types.Point{
		{X: box.Min.X, Y: box.Min.Y},
		{X: box.Max.X, Y: box.Min.Y},
		{X: box.Max.X, Y: box.Max.Y},
		{X: box.Min.X, Y: box.Max.Y},
	}

	for _, corner := range corners {
		if PointInTriangle(corner, a, b, c, eps) {
			return true
		}
	}

	triEdges := [][2]types.Point{
		{a, b},
		{b, c},
		{c, a},
	}

	boxEdges := [][2]types.Point{
		{corners[0], corners[1]},
		{corners[1], corners[2]},
		{corners[2], corners[3]},
		{corners[3], corners[0]},
	}

	for _, te := range triEdges {
		for _, be := range boxEdges {
			if hit, _ := SegmentsIntersect(te[0], te[1], be[0], be[1], eps); hit {
				return true
			}
		}
	}

	return false
}
