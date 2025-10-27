package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestOppositeWindingIntersectionArea tests if the geometric intersection
// algorithm correctly detects opposite-winding triangles as overlapping
func TestOppositeWindingIntersectionArea(t *testing.T) {
	// Triangle with vertices in one order
	a1 := types.Point{X: 120, Y: 65}
	b1 := types.Point{X: 122, Y: 64}
	c1 := types.Point{X: 121, Y: 64}

	// Same triangle with opposite winding (vertices reversed)
	a2 := types.Point{X: 121, Y: 64} // c1
	b2 := types.Point{X: 122, Y: 64} // b1
	c2 := types.Point{X: 120, Y: 65} // a1

	eps := 1e-9

	// Calculate intersection area
	area := predicates.TriangleIntersectionArea(a1, b1, c1, a2, b2, c2, eps)

	t.Logf("Triangle 1: (%v, %v, %v)", a1, b1, c1)
	t.Logf("Triangle 2: (%v, %v, %v)", a2, b2, c2)
	t.Logf("Intersection area: %.10f", area)
	t.Logf("Epsilon: %.10f", eps)

	// Calculate expected triangle area
	triangleArea := predicates.PolygonArea([]types.Point{a1, b1, c1})
	t.Logf("Triangle 1 area: %.10f", triangleArea)

	// For identical triangles (even with opposite winding), the intersection area
	// should be close to the full triangle area
	if area < eps {
		t.Errorf("Intersection area (%.10f) is less than epsilon (%.10f)", area, eps)
		t.Errorf("This indicates the algorithm does not detect opposite-winding overlaps")

		// Test the intersection polygon directly
		poly := predicates.TriangleIntersectionPolygon(a1, b1, c1, a2, b2, c2, eps)
		t.Logf("Intersection polygon has %d vertices: %v", len(poly), poly)
	} else {
		t.Logf("✓ Intersection correctly detected (area > epsilon)")
	}
}
