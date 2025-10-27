package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestOverlap_UserReportedCase tests the specific overlap case reported by user
// Triangle 273,250,272 overlaps with 273,249,272 with area 68.0
func TestOverlap_UserReportedCase(t *testing.T) {
	m := NewMesh(
		WithEpsilon(1e-9),
		WithTriangleOverlapCheck(true), // Enable overlap validation
	)

	// Verify config
	t.Logf("Config: validateTriangleOverlapArea=%v, epsilon=%.10f",
		m.cfg.validateTriangleOverlapArea, m.cfg.epsilon)

	// Add vertices
	v273, _ := m.AddVertex(types.Point{X: 15.000, Y: 132.000})
	v249, _ := m.AddVertex(types.Point{X: 7.000, Y: 132.000})
	v272, _ := m.AddVertex(types.Point{X: 15.000, Y: 149.000})
	v250, _ := m.AddVertex(types.Point{X: 5.000, Y: 132.000})

	t.Logf("Vertices:")
	t.Logf("  v%d: %v", v273, m.GetVertex(v273))
	t.Logf("  v%d: %v", v249, m.GetVertex(v249))
	t.Logf("  v%d: %v", v272, m.GetVertex(v272))
	t.Logf("  v%d: %v", v250, m.GetVertex(v250))

	// Add first triangle: 273,249,272
	t.Logf("\nAdding first triangle: %d,%d,%d", v273, v249, v272)
	err := m.AddTriangle(v273, v249, v272)
	if err != nil {
		t.Fatalf("Failed to add first triangle: %v", err)
	}
	t.Logf("✓ First triangle added")

	// Try to add second triangle: 273,250,272
	t.Logf("\nAttempting to add second triangle: %d,%d,%d", v273, v250, v272)
	err = m.AddTriangle(v273, v250, v272)

	if err != nil {
		t.Logf("✓ Second triangle correctly rejected: %v", err)
	} else {
		t.Errorf("✗ Second triangle was NOT rejected (validation bug!)")
		t.Logf("  Expected: overlap rejection with area ~68.0")

		// Verify overlap exists
		overlaps := m.FindOverlappingTriangles()
		t.Logf("  FindOverlappingTriangles() found %d overlaps", len(overlaps))
		if len(overlaps) > 0 {
			for i, o := range overlaps {
				t.Logf("    Overlap %d: area=%.4f", i+1, o.IntersectionArea)
			}
		}
	}
}
