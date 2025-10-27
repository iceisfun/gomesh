package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestOverlap_DirectCalculation tests TriangleIntersectionArea directly
func TestOverlap_DirectCalculation(t *testing.T) {
	// Triangle 1: 273,249,272
	t1v1 := types.Point{X: 15.000, Y: 132.000} // v273
	t1v2 := types.Point{X: 7.000, Y: 132.000}  // v249
	t1v3 := types.Point{X: 15.000, Y: 149.000} // v272

	// Triangle 2: 273,250,272
	t2v1 := types.Point{X: 15.000, Y: 132.000} // v273 (shared)
	t2v2 := types.Point{X: 5.000, Y: 132.000}  // v250
	t2v3 := types.Point{X: 15.000, Y: 149.000} // v272 (shared)

	t.Logf("Triangle 1: %v, %v, %v", t1v1, t1v2, t1v3)
	t.Logf("Triangle 2: %v, %v, %v", t2v1, t2v2, t2v3)

	// Calculate areas
	area1 := predicates.PolygonArea([]types.Point{t1v1, t1v2, t1v3})
	area2 := predicates.PolygonArea([]types.Point{t2v1, t2v2, t2v3})

	t.Logf("\nTriangle areas:")
	t.Logf("  Triangle 1 area: %.4f", area1)
	t.Logf("  Triangle 2 area: %.4f", area2)

	eps := 1e-9

	// Calculate intersection both ways
	intersectionArea1 := predicates.TriangleIntersectionArea(t1v1, t1v2, t1v3, t2v1, t2v2, t2v3, eps)
	intersectionArea2 := predicates.TriangleIntersectionArea(t2v1, t2v2, t2v3, t1v1, t1v2, t1v3, eps)

	t.Logf("\nIntersection areas:")
	t.Logf("  Tri1 vs Tri2: %.4f", intersectionArea1)
	t.Logf("  Tri2 vs Tri1: %.4f", intersectionArea2)
	t.Logf("  Maximum: %.4f", max(intersectionArea1, intersectionArea2))

	// Get intersection polygons
	poly1 := predicates.TriangleIntersectionPolygon(t1v1, t1v2, t1v3, t2v1, t2v2, t2v3, eps)
	poly2 := predicates.TriangleIntersectionPolygon(t2v1, t2v2, t2v3, t1v1, t1v2, t1v3, eps)

	t.Logf("\nIntersection polygon 1 (tri1 clipped by tri2): %d vertices", len(poly1))
	for i, p := range poly1 {
		t.Logf("  %d: %v", i, p)
	}

	t.Logf("\nIntersection polygon 2 (tri2 clipped by tri1): %d vertices", len(poly2))
	for i, p := range poly2 {
		t.Logf("  %d: %v", i, p)
	}

	expectedArea := 68.0
	maxArea := max(intersectionArea1, intersectionArea2)

	if maxArea < eps {
		t.Errorf("✗ Expected intersection area ~%.1f, but got %.4f", expectedArea, maxArea)
		t.Logf("  These triangles share edge (15,132)→(15,149)")
		t.Logf("  Triangle 2 should completely contain Triangle 1")
	} else {
		t.Logf("✓ Found intersection area: %.4f", maxArea)
	}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
