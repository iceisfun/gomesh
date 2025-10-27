package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestOverlapCase_CollinearEdges analyzes the specific overlap case
func TestOverlapCase_CollinearEdges(t *testing.T) {
	// Triangle 1 vertices
	v0 := types.Point{X: 151.00, Y: 155.00}
	v1 := types.Point{X: 151.00, Y: 146.00}
	v2 := types.Point{X: 150.00, Y: 154.00}

	// Triangle 2 vertices
	v3 := types.Point{X: 151.00, Y: 150.00}
	v4 := types.Point{X: 151.00, Y: 154.00}
	// v2 shared: (150, 154)

	t.Logf("Triangle 1 vertices:")
	t.Logf("  v0: %v", v0)
	t.Logf("  v1: %v", v1)
	t.Logf("  v2: %v", v2)

	t.Logf("\nTriangle 2 vertices:")
	t.Logf("  v2: %v (shared)", v2)
	t.Logf("  v3: %v", v3)
	t.Logf("  v4: %v", v4)

	// Analyze the geometry
	t.Logf("\n=== Geometric Analysis ===")

	// Check for collinear points
	t.Logf("\nCollinearity checks:")
	t.Logf("  v0, v1, v4 collinear? All have x=151, y: 155, 146, 154")
	t.Logf("  v0: (151, 155)")
	t.Logf("  v1: (151, 146)")
	t.Logf("  v3: (151, 150)")
	t.Logf("  v4: (151, 154)")
	t.Logf("  → v0, v1 form a vertical edge from y=146 to y=155")
	t.Logf("  → v3, v4 form a vertical edge from y=150 to y=154")
	t.Logf("  → v3-v4 edge is CONTAINED within v1-v0 edge (both on x=151)")

	// Calculate triangle areas
	area1 := predicates.PolygonArea([]types.Point{v0, v1, v2})
	area2 := predicates.PolygonArea([]types.Point{v2, v3, v4})

	t.Logf("\nTriangle areas:")
	t.Logf("  Triangle 1 area: %.4f", area1)
	t.Logf("  Triangle 2 area: %.4f", area2)

	// Calculate intersection
	eps := 1e-9
	intersectionArea := predicates.TriangleIntersectionArea(v0, v1, v2, v2, v3, v4, eps)

	t.Logf("\nIntersection:")
	t.Logf("  Intersection area: %.4f", intersectionArea)
	t.Logf("  Epsilon: %.10f", eps)
	t.Logf("  Intersection > epsilon? %v", intersectionArea > eps)

	// Get intersection polygon
	poly := predicates.TriangleIntersectionPolygon(v0, v1, v2, v2, v3, v4, eps)
	t.Logf("\nIntersection polygon:")
	t.Logf("  Vertices: %d", len(poly))
	for i, p := range poly {
		t.Logf("    %d: %v", i, p)
	}

	// Conclusion
	t.Logf("\n=== Conclusion ===")
	if intersectionArea > eps {
		t.Logf("✗ This IS a real overlap")
		t.Logf("  The triangles share vertex v2 and have overlapping collinear edges")
		t.Logf("  Triangle 2's edge v3-v4 (151,150)→(151,154) lies ON Triangle 1's edge v1-v0 (151,146)→(151,155)")
		t.Logf("  This creates a volumetric overlap with area %.4f", intersectionArea)
		t.Logf("  This SHOULD be rejected by WithTriangleOverlapCheck(true)")
	} else {
		t.Logf("✓ No meaningful overlap (intersection area < epsilon)")
	}

	// Now test with mesh validation
	t.Logf("\n=== Mesh Validation Test ===")

	m := NewMesh(
		WithEpsilon(1e-9),
		WithTriangleOverlapCheck(true), // Enable overlap check
	)

	// Verify config
	t.Logf("Mesh config: validateTriangleOverlapArea=%v, epsilon=%.10f",
		m.cfg.validateTriangleOverlapArea, m.cfg.epsilon)

	mv0, _ := m.AddVertex(v0)
	mv1, _ := m.AddVertex(v1)
	mv2, _ := m.AddVertex(v2)
	mv3, _ := m.AddVertex(v3)
	mv4, _ := m.AddVertex(v4)

	// Add first triangle
	err := m.AddTriangle(mv0, mv1, mv2)
	if err != nil {
		t.Fatalf("Failed to add first triangle: %v", err)
	}
	t.Logf("✓ First triangle added")

	// Try to add second triangle
	err = m.AddTriangle(mv2, mv3, mv4)
	if err != nil {
		t.Logf("✓ Second triangle correctly rejected: %v", err)
	} else {
		t.Errorf("✗ Second triangle was NOT rejected (validation bug!)")

		// Check what FindOverlappingTriangles says
		overlaps := m.FindOverlappingTriangles()
		t.Logf("  FindOverlappingTriangles() found %d overlaps", len(overlaps))
		if len(overlaps) > 0 {
			for i, o := range overlaps {
				t.Logf("    Overlap %d: area=%.4f, type=%s", i+1, o.IntersectionArea, o.Type)
			}
		}
	}
}
