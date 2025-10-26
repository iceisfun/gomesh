package predicates

import (
	"math"

	"github.com/iceisfun/gomesh/types"
)

// Dist2 returns the squared Euclidean distance between two points.
func Dist2(a, b types.Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

// SegmentsIntersect tests if two line segments intersect.
func SegmentsIntersect(a1, a2, b1, b2 types.Point, eps float64) (bool, bool) {
	o1 := Orient(a1, a2, b1, eps)
	o2 := Orient(a1, a2, b2, eps)
	o3 := Orient(b1, b2, a1, eps)
	o4 := Orient(b1, b2, a2, eps)

	// Proper intersection if segments straddle each other.
	if o1*o2 < 0 && o3*o4 < 0 {
		return true, true
	}

	// Check special cases: endpoints and collinear overlaps.
	if o1 == 0 && PointOnSegment(b1, a1, a2, eps) {
		return true, false
	}
	if o2 == 0 && PointOnSegment(b2, a1, a2, eps) {
		return true, false
	}
	if o3 == 0 && PointOnSegment(a1, b1, b2, eps) {
		return true, false
	}
	if o4 == 0 && PointOnSegment(a2, b1, b2, eps) {
		return true, false
	}

	return false, false
}

// SegmentIntersectionPoint computes the intersection point of two segments.
func SegmentIntersectionPoint(a1, a2, b1, b2 types.Point, eps float64) (types.Point, types.IntersectionType) {
	intersects, proper := SegmentsIntersect(a1, a2, b1, b2, eps)
	if !intersects {
		return types.Point{}, types.IntersectNone
	}

	if proper {
		p := lineIntersectionPoint(a1, a2, b1, b2)
		return p, types.IntersectProper
	}

	// Handle collinear overlaps or touching endpoints.
	if isCollinear(a1, a2, b1, eps) && isCollinear(a1, a2, b2, eps) {
		length, point := collinearOverlapPoint(a1, a2, b1, b2, eps)
		if length > eps {
			return point, types.IntersectCollinearOverlap
		}
		return point, types.IntersectTouching
	}

	// Otherwise one of the endpoints lies on the other segment.
	if PointOnSegment(a1, b1, b2, eps) {
		return a1, types.IntersectTouching
	}
	if PointOnSegment(a2, b1, b2, eps) {
		return a2, types.IntersectTouching
	}
	if PointOnSegment(b1, a1, a2, eps) {
		return b1, types.IntersectTouching
	}
	if PointOnSegment(b2, a1, a2, eps) {
		return b2, types.IntersectTouching
	}

	return types.Point{}, types.IntersectNone
}

// PointOnSegment tests if a point lies on a line segment within tolerance.
func PointOnSegment(p, a, b types.Point, eps float64) bool {
	area := math.Abs(Area2(a, b, p))
	segmentLen := math.Sqrt(Dist2(a, b))
	if segmentLen == 0 {
		return Dist2(p, a) <= eps*eps
	}
	if area > (segmentLen * eps) {
		return false
	}

	minX := math.Min(a.X, b.X) - eps
	maxX := math.Max(a.X, b.X) + eps
	minY := math.Min(a.Y, b.Y) - eps
	maxY := math.Max(a.Y, b.Y) + eps

	return p.X >= minX && p.X <= maxX && p.Y >= minY && p.Y <= maxY
}

func isCollinear(a1, a2, p types.Point, eps float64) bool {
	return math.Abs(Area2(a1, a2, p)) <= eps
}

func lineIntersectionPoint(a1, a2, b1, b2 types.Point) types.Point {
	r := types.Point{X: a2.X - a1.X, Y: a2.Y - a1.Y}
	s := types.Point{X: b2.X - b1.X, Y: b2.Y - b1.Y}
	d := cross(r, s)
	if d == 0 {
		// Parallel or collinear; caller handles before.
		return a1
	}
	t := cross(types.Point{X: b1.X - a1.X, Y: b1.Y - a1.Y}, s) / d
	return types.Point{X: a1.X + t*r.X, Y: a1.Y + t*r.Y}
}

func collinearOverlapPoint(a1, a2, b1, b2 types.Point, eps float64) (float64, types.Point) {
	dx := a2.X - a1.X
	dy := a2.Y - a1.Y
	useX := math.Abs(dx) >= math.Abs(dy)

	coord := func(p types.Point) float64 {
		if useX {
			return p.X
		}
		return p.Y
	}

	aMin, aMax := ordered(coord(a1), coord(a2))
	bMin, bMax := ordered(coord(b1), coord(b2))
	overlapMin := math.Max(aMin, bMin)
	overlapMax := math.Min(aMax, bMax)
	length := overlapMax - overlapMin

	mid := overlapMin + length/2

	var point types.Point
	if useX {
		var t float64
		if math.Abs(dx) < eps {
			t = 0
		} else {
			t = (mid - a1.X) / dx
		}
		point = types.Point{X: mid, Y: a1.Y + t*dy}
	} else {
		var t float64
		if math.Abs(dy) < eps {
			t = 0
		} else {
			t = (mid - a1.Y) / dy
		}
		point = types.Point{X: a1.X + t*dx, Y: mid}
	}

	return math.Abs(length), point
}

func ordered(a, b float64) (float64, float64) {
	if a < b {
		return a, b
	}
	return b, a
}

func cross(a, b types.Point) float64 {
	return a.X*b.Y - a.Y*b.X
}
