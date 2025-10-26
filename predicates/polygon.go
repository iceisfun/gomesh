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

// PolygonSelfIntersects checks if a polygon has any self-intersections.
//
// Returns true if any non-adjacent edges intersect.
// Adjacent edges (sharing a vertex) are allowed to touch.
func PolygonSelfIntersects(poly []types.Point, eps float64) bool {
	n := len(poly)
	if n < 3 {
		return false
	}

	// Check each edge against all other non-adjacent edges
	for i := 0; i < n; i++ {
		next := (i + 1) % n
		a1 := poly[i]
		a2 := poly[next]

		for j := i + 2; j < n; j++ {
			// Skip the edge that wraps around and connects to current edge
			if i == 0 && j == n-1 {
				continue
			}

			nextJ := (j + 1) % n
			b1 := poly[j]
			b2 := poly[nextJ]

			// Check for proper intersection (not just touching at shared vertex)
			intersects, proper := SegmentsIntersect(a1, a2, b1, b2, eps)
			if intersects && proper {
				return true
			}
		}
	}

	return false
}

// PolygonContainsPolygon tests if polygon A completely contains polygon B.
//
// Returns true if all vertices of B are inside A and no edges of B
// intersect edges of A.
func PolygonContainsPolygon(a, b []types.Point, eps float64) bool {
	if len(a) < 3 || len(b) < 3 {
		return false
	}

	// All vertices of B must be inside A
	for _, v := range b {
		if !PointInPolygonRayCast(v, a, eps) {
			return false
		}
	}

	// No edges of B should intersect edges of A
	// (if all vertices are inside and no edge crosses, B is fully contained)
	for i := 0; i < len(b); i++ {
		nextB := (i + 1) % len(b)
		b1 := b[i]
		b2 := b[nextB]

		for j := 0; j < len(a); j++ {
			nextA := (j + 1) % len(a)
			a1 := a[j]
			a2 := a[nextA]

			intersects, proper := SegmentsIntersect(b1, b2, a1, a2, eps)
			if intersects && proper {
				return false
			}
		}
	}

	return true
}

// PolygonsIntersect tests if two polygons intersect (overlap or touch).
//
// Returns true if:
//   - Any vertex of one polygon is inside the other
//   - Any edges intersect
//   - One polygon contains the other
func PolygonsIntersect(a, b []types.Point, eps float64) bool {
	if len(a) < 3 || len(b) < 3 {
		return false
	}

	// Check if any vertex of A is inside B
	for _, v := range a {
		if PointInPolygonRayCast(v, b, eps) {
			return true
		}
	}

	// Check if any vertex of B is inside A
	for _, v := range b {
		if PointInPolygonRayCast(v, a, eps) {
			return true
		}
	}

	// Check if any edges intersect
	for i := 0; i < len(a); i++ {
		nextA := (i + 1) % len(a)
		a1 := a[i]
		a2 := a[nextA]

		for j := 0; j < len(b); j++ {
			nextB := (j + 1) % len(b)
			b1 := b[j]
			b2 := b[nextB]

			intersects, _ := SegmentsIntersect(a1, a2, b1, b2, eps)
			if intersects {
				return true
			}
		}
	}

	return false
}

// PolygonArea computes the signed area of a polygon.
//
// Returns:
//   - Positive area for CCW (counter-clockwise) winding
//   - Negative area for CW (clockwise) winding
//   - Zero for degenerate polygons
func PolygonArea(poly []types.Point) float64 {
	n := len(poly)
	if n < 3 {
		return 0
	}

	area := 0.0
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += poly[i].X * poly[j].Y
		area -= poly[j].X * poly[i].Y
	}

	return area / 2.0
}

// PolygonBounds computes the axis-aligned bounding box of a polygon.
func PolygonBounds(poly []types.Point) types.AABB {
	if len(poly) == 0 {
		return types.AABB{}
	}

	bounds := types.AABB{
		Min: poly[0],
		Max: poly[0],
	}

	for _, p := range poly[1:] {
		if p.X < bounds.Min.X {
			bounds.Min.X = p.X
		}
		if p.Y < bounds.Min.Y {
			bounds.Min.Y = p.Y
		}
		if p.X > bounds.Max.X {
			bounds.Max.X = p.X
		}
		if p.Y > bounds.Max.Y {
			bounds.Max.Y = p.Y
		}
	}

	return bounds
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

// PolygonLoop-specific functions that work with types.VertexProvider

// PolygonLoopSelfIntersects checks if a polygon loop has self-intersections.
//
// The vertex provider (e.g., mesh) is used to resolve vertex coordinates.
//
// Example:
//
//	if PolygonLoopSelfIntersects(mesh, loop, 1e-9) {
//	    // Handle self-intersection
//	}
func PolygonLoopSelfIntersects(vp types.VertexProvider, loop types.PolygonLoop, eps float64) bool {
	points := loop.ToPoints(vp)
	return PolygonSelfIntersects(points, eps)
}

// PolygonLoopContains tests if a point is inside a polygon loop.
//
// Example:
//
//	if PolygonLoopContains(mesh, loop, point, 1e-9) {
//	    // Point is inside
//	}
func PolygonLoopContains(vp types.VertexProvider, loop types.PolygonLoop, point types.Point, eps float64) bool {
	points := loop.ToPoints(vp)
	return PointInPolygonRayCast(point, points, eps)
}

// PolygonLoopContainsPolygonLoop tests if polygon loop A contains loop B.
//
// Example:
//
//	if PolygonLoopContainsPolygonLoop(mesh, outer, inner, 1e-9) {
//	    // Outer contains inner
//	}
func PolygonLoopContainsPolygonLoop(vp types.VertexProvider, a, b types.PolygonLoop, eps float64) bool {
	pointsA := a.ToPoints(vp)
	pointsB := b.ToPoints(vp)
	return PolygonContainsPolygon(pointsA, pointsB, eps)
}

// PolygonLoopsIntersect tests if two polygon loops intersect.
//
// Example:
//
//	if PolygonLoopsIntersect(mesh, loop1, loop2, 1e-9) {
//	    // Loops intersect
//	}
func PolygonLoopsIntersect(vp types.VertexProvider, a, b types.PolygonLoop, eps float64) bool {
	pointsA := a.ToPoints(vp)
	pointsB := b.ToPoints(vp)
	return PolygonsIntersect(pointsA, pointsB, eps)
}

// PolygonLoopArea computes the signed area of a polygon loop.
//
// Example:
//
//	area := PolygonLoopArea(mesh, loop)
func PolygonLoopArea(vp types.VertexProvider, loop types.PolygonLoop) float64 {
	points := loop.ToPoints(vp)
	return PolygonArea(points)
}

// PolygonLoopBounds computes the bounding box of a polygon loop.
//
// Example:
//
//	bounds := PolygonLoopBounds(mesh, loop)
func PolygonLoopBounds(vp types.VertexProvider, loop types.PolygonLoop) types.AABB {
	points := loop.ToPoints(vp)
	return PolygonBounds(points)
}
