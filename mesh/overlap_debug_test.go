package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestOverlapDebug directly calls validateTriangleDoesNotOverlap to see what happens
func TestOverlapDebug(t *testing.T) {
	m := NewMesh(
		WithEpsilon(1e-9),
		WithTriangleOverlapCheck(true),
	)

	// Add first triangle vertices
	v0, _ := m.AddVertex(types.Point{X: 151, Y: 155})
	v1, _ := m.AddVertex(types.Point{X: 151, Y: 146})
	v2, _ := m.AddVertex(types.Point{X: 150, Y: 154})

	// Add first triangle
	tri1 := types.NewTriangle(v0, v1, v2)
	m.triangles = append(m.triangles, tri1)

	t.Logf("Added first triangle manually: %v", tri1)
	t.Logf("Mesh now has %d triangles", len(m.triangles))

	// Add second triangle vertices
	v3, _ := m.AddVertex(types.Point{X: 151, Y: 150})
	v4, _ := m.AddVertex(types.Point{X: 151, Y: 154})

	// Create second triangle (shares v2)
	tri2 := types.NewTriangle(v2, v3, v4)

	t.Logf("\nAttempting to validate second triangle: %v", tri2)

	a := m.vertices[tri2.V1()]
	b := m.vertices[tri2.V2()]
	c := m.vertices[tri2.V3()]

	t.Logf("Second triangle coordinates:")
	t.Logf("  v%d: %v", tri2.V1(), a)
	t.Logf("  v%d: %v", tri2.V2(), b)
	t.Logf("  v%d: %v", tri2.V3(), c)

	// Calculate intersection area directly
	a1 := m.vertices[tri1.V1()]
	b1 := m.vertices[tri1.V2()]
	c1 := m.vertices[tri1.V3()]

	intersectionArea1 := predicates.TriangleIntersectionArea(a, b, c, a1, b1, c1, m.cfg.epsilon)
	intersectionArea2 := predicates.TriangleIntersectionArea(a1, b1, c1, a, b, c, m.cfg.epsilon)

	t.Logf("\nDirect intersection area calculation:")
	t.Logf("  Triangle 1: (%v, %v, %v)", a1, b1, c1)
	t.Logf("  Triangle 2: (%v, %v, %v)", a, b, c)
	t.Logf("  Intersection area (tri2 vs tri1): %.10f", intersectionArea1)
	t.Logf("  Intersection area (tri1 vs tri2): %.10f", intersectionArea2)
	t.Logf("  Epsilon: %.10f", m.cfg.epsilon)
	t.Logf("  Area1 > Epsilon? %v", intersectionArea1 > m.cfg.epsilon)
	t.Logf("  Area2 > Epsilon? %v", intersectionArea2 > m.cfg.epsilon)

	intersectionArea := intersectionArea1

	// Now call the actual validation function
	err := m.validateTriangleDoesNotOverlap(tri2, a, b, c)

	t.Logf("\nValidation result: %v", err)

	if err != nil {
		t.Logf("✓ Validation correctly rejected overlap: %v", err)
	} else {
		t.Errorf("✗ Validation did NOT reject overlap even though area=%.4f > epsilon=%.10f",
			intersectionArea, m.cfg.epsilon)
	}
}
