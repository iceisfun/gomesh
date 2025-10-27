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
	// Ensure both triangles are in counter-clockwise order
	// The Sutherland-Hodgman algorithm requires CCW winding for correct "inside" tests
	triA := ensureCCW([]types.Point{a1, a2, a3})
	triB := ensureCCW([]types.Point{b1, b2, b3})

	// Start with triangle A as the subject polygon
	subject := triA

	// Clip against each edge of triangle B
	clipEdges := [][2]types.Point{
		{triB[0], triB[1]},
		{triB[1], triB[2]},
		{triB[2], triB[0]},
	}

	for _, edge := range clipEdges {
		subject = sutherlandHodgmanClip(subject, edge[0], edge[1], eps)
		if len(subject) == 0 {
			return nil
		}
	}

	return subject
}

// ensureCCW ensures a polygon is in counter-clockwise winding order.
// Returns a new slice with vertices in CCW order.
func ensureCCW(poly []types.Point) []types.Point {
	if len(poly) < 3 {
		return poly
	}

	// Calculate signed area using shoelace formula
	// Positive area = CCW, negative area = CW
	signedArea := 0.0
	n := len(poly)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		signedArea += (poly[j].X - poly[i].X) * (poly[j].Y + poly[i].Y)
	}

	// If area is positive (CW by our formula convention), reverse to make CCW
	if signedArea > 0 {
		result := make([]types.Point, n)
		for i := 0; i < n; i++ {
			result[i] = poly[n-1-i]
		}
		return result
	}

	// Already CCW, return copy
	result := make([]types.Point, n)
	copy(result, poly)
	return result
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
