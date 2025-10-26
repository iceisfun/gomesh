package validation_test

import (
	"testing"

	"github.com/iceisfun/gomesh/mesh"
	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
	"github.com/iceisfun/gomesh/validation"
)

// TestPolygonLoopSelfIntersects tests self-intersection with PolygonLoop
func TestPolygonLoopSelfIntersects(t *testing.T) {
	m := mesh.NewMesh()

	// Create a self-intersecting bowtie
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 0, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 10, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	if !predicates.PolygonLoopSelfIntersects(m, loop, 1e-9) {
		t.Error("Expected bowtie to self-intersect")
	}

	// Create a valid square
	m2 := mesh.NewMesh()
	s0, _ := m2.AddVertex(types.Point{X: 0, Y: 0})
	s1, _ := m2.AddVertex(types.Point{X: 10, Y: 0})
	s2, _ := m2.AddVertex(types.Point{X: 10, Y: 10})
	s3, _ := m2.AddVertex(types.Point{X: 0, Y: 10})

	square := types.NewPolygonLoop(s0, s1, s2, s3)

	if predicates.PolygonLoopSelfIntersects(m2, square, 1e-9) {
		t.Error("Expected square not to self-intersect")
	}
}

// TestPolygonLoopContains tests point-in-polygon with PolygonLoop
func TestPolygonLoopContains(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	tests := []struct {
		name     string
		point    types.Point
		expected bool
	}{
		{"center", types.Point{X: 5, Y: 5}, true},
		{"outside", types.Point{X: 15, Y: 15}, false},
		{"on edge", types.Point{X: 0, Y: 5}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicates.PolygonLoopContains(m, loop, tt.point, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonLoopContains(%v) = %v, want %v", tt.point, result, tt.expected)
			}
		})
	}
}

// TestPolygonLoopContainsPolygonLoop tests loop containment
func TestPolygonLoopContainsPolygonLoop(t *testing.T) {
	m := mesh.NewMesh()

	// Outer square
	o0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	o1, _ := m.AddVertex(types.Point{X: 20, Y: 0})
	o2, _ := m.AddVertex(types.Point{X: 20, Y: 20})
	o3, _ := m.AddVertex(types.Point{X: 0, Y: 20})
	outer := types.NewPolygonLoop(o0, o1, o2, o3)

	// Inner square
	i0, _ := m.AddVertex(types.Point{X: 5, Y: 5})
	i1, _ := m.AddVertex(types.Point{X: 15, Y: 5})
	i2, _ := m.AddVertex(types.Point{X: 15, Y: 15})
	i3, _ := m.AddVertex(types.Point{X: 5, Y: 15})
	inner := types.NewPolygonLoop(i0, i1, i2, i3)

	if !predicates.PolygonLoopContainsPolygonLoop(m, outer, inner, 1e-9) {
		t.Error("Expected outer to contain inner")
	}

	if predicates.PolygonLoopContainsPolygonLoop(m, inner, outer, 1e-9) {
		t.Error("Expected inner not to contain outer")
	}
}

// TestPolygonLoopsIntersect tests loop intersection
func TestPolygonLoopsIntersect(t *testing.T) {
	m := mesh.NewMesh()

	// First square
	a0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	a1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	a2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	a3, _ := m.AddVertex(types.Point{X: 0, Y: 10})
	loopA := types.NewPolygonLoop(a0, a1, a2, a3)

	// Overlapping square
	b0, _ := m.AddVertex(types.Point{X: 5, Y: 5})
	b1, _ := m.AddVertex(types.Point{X: 15, Y: 5})
	b2, _ := m.AddVertex(types.Point{X: 15, Y: 15})
	b3, _ := m.AddVertex(types.Point{X: 5, Y: 15})
	loopB := types.NewPolygonLoop(b0, b1, b2, b3)

	// Non-overlapping square
	c0, _ := m.AddVertex(types.Point{X: 20, Y: 20})
	c1, _ := m.AddVertex(types.Point{X: 30, Y: 20})
	c2, _ := m.AddVertex(types.Point{X: 30, Y: 30})
	c3, _ := m.AddVertex(types.Point{X: 20, Y: 30})
	loopC := types.NewPolygonLoop(c0, c1, c2, c3)

	if !predicates.PolygonLoopsIntersect(m, loopA, loopB, 1e-9) {
		t.Error("Expected loopA and loopB to intersect")
	}

	if predicates.PolygonLoopsIntersect(m, loopA, loopC, 1e-9) {
		t.Error("Expected loopA and loopC not to intersect")
	}
}

// TestPolygonLoopArea tests area calculation
func TestPolygonLoopArea(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	area := predicates.PolygonLoopArea(m, loop)
	expected := 100.0

	if area != expected {
		t.Errorf("PolygonLoopArea() = %f, want %f", area, expected)
	}
}

// TestPolygonLoopBounds tests bounds calculation
func TestPolygonLoopBounds(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 5, Y: 3})
	v1, _ := m.AddVertex(types.Point{X: 15, Y: 7})
	v2, _ := m.AddVertex(types.Point{X: 12, Y: 18})
	v3, _ := m.AddVertex(types.Point{X: 2, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	bounds := predicates.PolygonLoopBounds(m, loop)

	if bounds.Min.X != 2 || bounds.Min.Y != 3 {
		t.Errorf("Min = (%f, %f), want (2, 3)", bounds.Min.X, bounds.Min.Y)
	}

	if bounds.Max.X != 15 || bounds.Max.Y != 18 {
		t.Errorf("Max = (%f, %f), want (15, 18)", bounds.Max.X, bounds.Max.Y)
	}
}

// TestValidatePolygonLoop tests validation with PolygonLoop
func TestValidatePolygonLoop(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	// Should pass validation
	err := validation.ValidatePolygonLoop(m, loop,
		validation.WithPolygonMinArea(50),
		validation.WithPolygonMinWidth(5),
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Should fail validation (area too small)
	err = validation.ValidatePolygonLoop(m, loop, validation.WithPolygonMinArea(200))
	if err == nil {
		t.Error("Expected validation to fail")
	}
}

// TestValidatePolygonLoopDetailed tests detailed validation
func TestValidatePolygonLoopDetailed(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	loop := types.NewPolygonLoop(v0, v1, v2, v3)

	result := validation.ValidatePolygonLoopDetailed(m, loop)

	if !result.Valid {
		t.Error("Expected loop to be valid")
	}

	if result.VertexCount != 4 {
		t.Errorf("VertexCount = %d, want 4", result.VertexCount)
	}

	if result.Area != 100.0 {
		t.Errorf("Area = %f, want 100", result.Area)
	}
}

// TestPolygonLoopIsValid tests quick validity check
func TestPolygonLoopIsValid(t *testing.T) {
	// Valid square
	m := mesh.NewMesh()
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})
	square := types.NewPolygonLoop(v0, v1, v2, v3)

	if !validation.PolygonLoopIsValid(m, square, 1e-9) {
		t.Error("Expected square to be valid")
	}

	// Self-intersecting bowtie
	m2 := mesh.NewMesh()
	b0, _ := m2.AddVertex(types.Point{X: 0, Y: 0})
	b1, _ := m2.AddVertex(types.Point{X: 10, Y: 0})
	b2, _ := m2.AddVertex(types.Point{X: 0, Y: 10})
	b3, _ := m2.AddVertex(types.Point{X: 10, Y: 10})
	bowtie := types.NewPolygonLoop(b0, b1, b2, b3)

	if validation.PolygonLoopIsValid(m2, bowtie, 1e-9) {
		t.Error("Expected bowtie to be invalid")
	}
}

// TestToPoints tests conversion to point array
func TestToPoints(t *testing.T) {
	m := mesh.NewMesh()

	v0, _ := m.AddVertex(types.Point{X: 1, Y: 2})
	v1, _ := m.AddVertex(types.Point{X: 3, Y: 4})
	v2, _ := m.AddVertex(types.Point{X: 5, Y: 6})

	loop := types.NewPolygonLoop(v0, v1, v2)
	points := loop.ToPoints(m)

	if len(points) != 3 {
		t.Errorf("len(points) = %d, want 3", len(points))
	}

	expected := []types.Point{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 5, Y: 6},
	}

	for i, p := range points {
		if p.X != expected[i].X || p.Y != expected[i].Y {
			t.Errorf("points[%d] = (%f, %f), want (%f, %f)",
				i, p.X, p.Y, expected[i].X, expected[i].Y)
		}
	}
}

// TestPolygonLoopReversedFlipsWinding tests that reversing flips winding direction
func TestPolygonLoopReversedFlipsWinding(t *testing.T) {
	m := mesh.NewMesh()

	// Create CCW square (positive area)
	v0, _ := m.AddVertex(types.Point{X: 0, Y: 0})
	v1, _ := m.AddVertex(types.Point{X: 10, Y: 0})
	v2, _ := m.AddVertex(types.Point{X: 10, Y: 10})
	v3, _ := m.AddVertex(types.Point{X: 0, Y: 10})

	ccwLoop := types.NewPolygonLoop(v0, v1, v2, v3)
	cwLoop := ccwLoop.Reversed()

	// Compute areas
	ccwArea := predicates.PolygonLoopArea(m, ccwLoop)
	cwArea := predicates.PolygonLoopArea(m, cwLoop)

	// CCW should have positive area
	if ccwArea <= 0 {
		t.Errorf("CCW loop area should be positive, got %f", ccwArea)
	}

	// CW should have negative area
	if cwArea >= 0 {
		t.Errorf("CW loop area should be negative, got %f", cwArea)
	}

	// Areas should have same magnitude but opposite sign
	if ccwArea != -cwArea {
		t.Errorf("Areas should be opposite: CCW=%f, CW=%f", ccwArea, cwArea)
	}

	// Double reversal should give original area
	doubleReversed := cwLoop.Reversed()
	doubleReversedArea := predicates.PolygonLoopArea(m, doubleReversed)
	if doubleReversedArea != ccwArea {
		t.Errorf("Double reversal should restore original area: original=%f, double=%f", ccwArea, doubleReversedArea)
	}

	t.Logf("CCW area: %f (positive)", ccwArea)
	t.Logf("CW area: %f (negative)", cwArea)
	t.Logf("Double reversed area: %f", doubleReversedArea)
}
