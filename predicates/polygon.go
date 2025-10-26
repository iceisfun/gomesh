package predicates

import "github.com/iceisfun/gomesh/types"

// PointInPolygonRayCast tests if a point is inside a polygon using ray casting.
func PointInPolygonRayCast(p types.Point, poly []types.Point, eps float64) bool {
	n := len(poly)
	if n == 0 {
		return false
	}

	// Boundary check first.
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if PointOnSegment(p, poly[i], poly[j], eps) {
			return true
		}
	}

	inside := false
	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		iP := poly[i]
		jP := poly[j]
		diff := (iP.Y > p.Y) != (jP.Y > p.Y)
		if diff {
			t := (p.Y - iP.Y) / (jP.Y - iP.Y)
			x := iP.X + t*(jP.X-iP.X)
			if x > p.X {
				inside = !inside
			}
		}
	}

	return inside
}

// PolygonAABBIntersect tests if a polygon intersects an AABB.
func PolygonAABBIntersect(poly []types.Point, box types.AABB, eps float64) bool {
	n := len(poly)
	if n == 0 {
		return false
	}

	for _, v := range poly {
		if PointInAABB(v, box, eps) {
			return true
		}
	}

	corners := []types.Point{
		{X: box.Min.X, Y: box.Min.Y},
		{X: box.Max.X, Y: box.Min.Y},
		{X: box.Max.X, Y: box.Max.Y},
		{X: box.Min.X, Y: box.Max.Y},
	}

	for _, corner := range corners {
		if PointInPolygonRayCast(corner, poly, eps) {
			return true
		}
	}

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if SegmentAABBIntersect(poly[i], poly[j], box, eps) {
			return true
		}
	}

	return false
}
