package predicates

import (
	"math"

	"github.com/iceisfun/gomesh/types"
)

// TriangleIntersectionArea computes the area of intersection between two triangles.
// Returns 0 if they don't intersect.
func TriangleIntersectionArea(a1, a2, a3, b1, b2, b3 types.Point, eps float64) float64 {
	// Get the intersection polygon using Sutherland-Hodgman clipping
	poly := TriangleIntersectionPolygon(a1, a2, a3, b1, b2, b3, eps)
	if len(poly) < 3 {
		return 0
	}
	return math.Abs(PolygonArea(poly))
}

// TriangleIntersectionPolygon computes the polygon formed by the intersection of two triangles.
// Uses Sutherland-Hodgman polygon clipping algorithm.
// Returns empty slice if triangles don't intersect.
func TriangleIntersectionPolygon(a1, a2, a3, b1, b2, b3 types.Point, eps float64) []types.Point {
	// Start with triangle A as the subject polygon
	subject := []types.Point{a1, a2, a3}

	// Clip against each edge of triangle B
	clipEdges := [][2]types.Point{
		{b1, b2},
		{b2, b3},
		{b3, b1},
	}

	for _, edge := range clipEdges {
		subject = sutherlandHodgmanClip(subject, edge[0], edge[1], eps)
		if len(subject) == 0 {
			return nil
		}
	}

	return subject
}

// sutherlandHodgmanClip clips a polygon against a single edge using Sutherland-Hodgman algorithm.
func sutherlandHodgmanClip(poly []types.Point, edgeStart, edgeEnd types.Point, eps float64) []types.Point {
	if len(poly) == 0 {
		return nil
	}

	var output []types.Point

	for i := 0; i < len(poly); i++ {
		current := poly[i]
		previous := poly[(i+len(poly)-1)%len(poly)]

		currentInside := isLeftOfEdge(current, edgeStart, edgeEnd, eps)
		previousInside := isLeftOfEdge(previous, edgeStart, edgeEnd, eps)

		if currentInside {
			if !previousInside {
				// Entering: add intersection point
				intersection := lineLineIntersection(previous, current, edgeStart, edgeEnd)
				output = append(output, intersection)
			}
			// Add current point
			output = append(output, current)
		} else if previousInside {
			// Exiting: add intersection point only
			intersection := lineLineIntersection(previous, current, edgeStart, edgeEnd)
			output = append(output, intersection)
		}
		// else: both outside, add nothing
	}

	return output
}

// isLeftOfEdge tests if a point is on the left side (or on) an edge.
// The edge is directed from start to end.
func isLeftOfEdge(p, edgeStart, edgeEnd types.Point, eps float64) bool {
	// Use cross product to determine side
	cross := (edgeEnd.X-edgeStart.X)*(p.Y-edgeStart.Y) - (edgeEnd.Y-edgeStart.Y)*(p.X-edgeStart.X)
	return cross >= -eps
}

// lineLineIntersection computes the intersection point of two infinite lines.
// Assumes the lines are not parallel (caller should ensure they intersect).
func lineLineIntersection(a1, a2, b1, b2 types.Point) types.Point {
	// Line A: a1 + t * (a2 - a1)
	// Line B: b1 + s * (b2 - b1)

	dx1 := a2.X - a1.X
	dy1 := a2.Y - a1.Y
	dx2 := b2.X - b1.X
	dy2 := b2.Y - b1.Y

	denominator := dx1*dy2 - dy1*dx2

	if math.Abs(denominator) < 1e-10 {
		// Lines are parallel or coincident, return midpoint
		return types.Point{
			X: (a1.X + a2.X) / 2,
			Y: (a1.Y + a2.Y) / 2,
		}
	}

	t := ((b1.X-a1.X)*dy2 - (b1.Y-a1.Y)*dx2) / denominator

	return types.Point{
		X: a1.X + t*dx1,
		Y: a1.Y + t*dy1,
	}
}
