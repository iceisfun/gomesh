package validation

import (
	"math"
	"strings"
	"testing"

	"github.com/iceisfun/gomesh/predicates"
	"github.com/iceisfun/gomesh/types"
)

// TestPolygonSelfIntersects tests self-intersection detection
func TestPolygonSelfIntersects(t *testing.T) {
	tests := []struct {
		name     string
		polygon  []types.Point
		expected bool
	}{
		{
			name: "valid square - no self-intersection",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			expected: false,
		},
		{
			name: "self-intersecting bowtie",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 0, Y: 10},
				{X: 10, Y: 10},
			},
			expected: true,
		},
		{
			name: "valid triangle",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 5, Y: 10},
			},
			expected: false,
		},
		{
			name: "valid pentagon",
			polygon: []types.Point{
				{X: 5, Y: 0},
				{X: 10, Y: 3},
				{X: 8, Y: 9},
				{X: 2, Y: 9},
				{X: 0, Y: 3},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicates.PolygonSelfIntersects(tt.polygon, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonSelfIntersects() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPolygonContainsPolygon tests polygon containment
func TestPolygonContainsPolygon(t *testing.T) {
	tests := []struct {
		name     string
		outer    []types.Point
		inner    []types.Point
		expected bool
	}{
		{
			name: "outer contains inner - positive case",
			outer: []types.Point{
				{X: 0, Y: 0},
				{X: 20, Y: 0},
				{X: 20, Y: 20},
				{X: 0, Y: 20},
			},
			inner: []types.Point{
				{X: 5, Y: 5},
				{X: 15, Y: 5},
				{X: 15, Y: 15},
				{X: 5, Y: 15},
			},
			expected: true,
		},
		{
			name: "inner not contained - outside",
			outer: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			inner: []types.Point{
				{X: 15, Y: 15},
				{X: 25, Y: 15},
				{X: 25, Y: 25},
				{X: 15, Y: 25},
			},
			expected: false,
		},
		{
			name: "inner partially outside",
			outer: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			inner: []types.Point{
				{X: 5, Y: 5},
				{X: 15, Y: 5},
				{X: 15, Y: 15},
				{X: 5, Y: 15},
			},
			expected: false,
		},
		{
			name: "inner touches outer boundary",
			outer: []types.Point{
				{X: 0, Y: 0},
				{X: 20, Y: 0},
				{X: 20, Y: 20},
				{X: 0, Y: 20},
			},
			inner: []types.Point{
				{X: 0, Y: 5},
				{X: 10, Y: 5},
				{X: 10, Y: 15},
				{X: 0, Y: 15},
			},
			expected: true,
		},
		{
			name: "triangle inside square",
			outer: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			inner: []types.Point{
				{X: 2, Y: 2},
				{X: 8, Y: 2},
				{X: 5, Y: 8},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicates.PolygonContainsPolygon(tt.outer, tt.inner, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonContainsPolygon() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPolygonsIntersect tests polygon intersection detection
func TestPolygonsIntersect(t *testing.T) {
	tests := []struct {
		name     string
		poly1    []types.Point
		poly2    []types.Point
		expected bool
	}{
		{
			name: "overlapping polygons",
			poly1: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			poly2: []types.Point{
				{X: 5, Y: 5},
				{X: 15, Y: 5},
				{X: 15, Y: 15},
				{X: 5, Y: 15},
			},
			expected: true,
		},
		{
			name: "non-overlapping polygons",
			poly1: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			poly2: []types.Point{
				{X: 20, Y: 20},
				{X: 30, Y: 20},
				{X: 30, Y: 30},
				{X: 20, Y: 30},
			},
			expected: false,
		},
		{
			name: "touching at edge",
			poly1: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			poly2: []types.Point{
				{X: 10, Y: 0},
				{X: 20, Y: 0},
				{X: 20, Y: 10},
				{X: 10, Y: 10},
			},
			expected: true,
		},
		{
			name: "one contains the other",
			poly1: []types.Point{
				{X: 0, Y: 0},
				{X: 20, Y: 0},
				{X: 20, Y: 20},
				{X: 0, Y: 20},
			},
			poly2: []types.Point{
				{X: 5, Y: 5},
				{X: 15, Y: 5},
				{X: 15, Y: 15},
				{X: 5, Y: 15},
			},
			expected: true,
		},
		{
			name: "edges cross",
			poly1: []types.Point{
				{X: 0, Y: 5},
				{X: 10, Y: 5},
				{X: 10, Y: 15},
				{X: 0, Y: 15},
			},
			poly2: []types.Point{
				{X: 5, Y: 0},
				{X: 15, Y: 0},
				{X: 15, Y: 10},
				{X: 5, Y: 10},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicates.PolygonsIntersect(tt.poly1, tt.poly2, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonsIntersect() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestValidatePolygon_MinArea tests minimum area validation
func TestValidatePolygon_MinArea(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		minArea   float64
		shouldErr bool
	}{
		{
			name: "area meets minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			minArea:   100,
			shouldErr: false,
		},
		{
			name: "area exceeds minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			minArea:   50,
			shouldErr: false,
		},
		{
			name: "area below minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 2, Y: 0},
				{X: 2, Y: 2},
				{X: 0, Y: 2},
			},
			minArea:   50,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMinArea(tt.minArea))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_MaxArea tests maximum area validation
func TestValidatePolygon_MaxArea(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		maxArea   float64
		shouldErr bool
	}{
		{
			name: "area within maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 5, Y: 0},
				{X: 5, Y: 5},
				{X: 0, Y: 5},
			},
			maxArea:   100,
			shouldErr: false,
		},
		{
			name: "area exceeds maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 20, Y: 0},
				{X: 20, Y: 20},
				{X: 0, Y: 20},
			},
			maxArea:   100,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMaxArea(tt.maxArea))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_MinWidth tests minimum width validation
func TestValidatePolygon_MinWidth(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		minWidth  float64
		shouldErr bool
	}{
		{
			name: "width meets minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 5},
				{X: 0, Y: 5},
			},
			minWidth:  10,
			shouldErr: false,
		},
		{
			name: "width below minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 3, Y: 0},
				{X: 3, Y: 10},
				{X: 0, Y: 10},
			},
			minWidth:  5,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMinWidth(tt.minWidth))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_MaxWidth tests maximum width validation
func TestValidatePolygon_MaxWidth(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		maxWidth  float64
		shouldErr bool
	}{
		{
			name: "width within maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 8, Y: 0},
				{X: 8, Y: 5},
				{X: 0, Y: 5},
			},
			maxWidth:  10,
			shouldErr: false,
		},
		{
			name: "width exceeds maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 15, Y: 0},
				{X: 15, Y: 5},
				{X: 0, Y: 5},
			},
			maxWidth:  10,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMaxWidth(tt.maxWidth))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_MinHeight tests minimum height validation
func TestValidatePolygon_MinHeight(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		minHeight float64
		shouldErr bool
	}{
		{
			name: "height meets minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 5, Y: 0},
				{X: 5, Y: 10},
				{X: 0, Y: 10},
			},
			minHeight: 10,
			shouldErr: false,
		},
		{
			name: "height below minimum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 3},
				{X: 0, Y: 3},
			},
			minHeight: 5,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMinHeight(tt.minHeight))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_MaxHeight tests maximum height validation
func TestValidatePolygon_MaxHeight(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		maxHeight float64
		shouldErr bool
	}{
		{
			name: "height within maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 5, Y: 0},
				{X: 5, Y: 8},
				{X: 0, Y: 8},
			},
			maxHeight: 10,
			shouldErr: false,
		},
		{
			name: "height exceeds maximum",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 5, Y: 0},
				{X: 5, Y: 15},
				{X: 0, Y: 15},
			},
			maxHeight: 10,
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon, WithPolygonMaxHeight(tt.maxHeight))
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygon_AllowSelfIntersection tests self-intersection allowance
func TestValidatePolygon_AllowSelfIntersection(t *testing.T) {
	selfIntersecting := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
	}

	t.Run("self-intersection not allowed (default)", func(t *testing.T) {
		err := ValidatePolygon(selfIntersecting)
		if err == nil {
			t.Error("Expected error for self-intersecting polygon")
		}
	})

	t.Run("self-intersection allowed", func(t *testing.T) {
		err := ValidatePolygon(selfIntersecting, WithAllowSelfIntersection(true))
		if err != nil {
			t.Errorf("Should not error when self-intersection is allowed: %v", err)
		}
	})
}

// TestValidatePolygon_RequireCCW tests counter-clockwise winding requirement
func TestValidatePolygon_RequireCCW(t *testing.T) {
	ccwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	cwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
		{X: 10, Y: 0},
	}

	t.Run("CCW polygon passes CCW requirement", func(t *testing.T) {
		err := ValidatePolygon(ccwSquare, WithRequireCCW(true))
		if err != nil {
			t.Errorf("CCW polygon should pass CCW requirement: %v", err)
		}
	})

	t.Run("CW polygon fails CCW requirement", func(t *testing.T) {
		err := ValidatePolygon(cwSquare, WithRequireCCW(true))
		if err == nil {
			t.Error("CW polygon should fail CCW requirement")
		}
	})
}

// TestValidatePolygon_RequireCW tests clockwise winding requirement
func TestValidatePolygon_RequireCW(t *testing.T) {
	ccwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	cwSquare := []types.Point{
		{X: 0, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
		{X: 10, Y: 0},
	}

	t.Run("CW polygon passes CW requirement", func(t *testing.T) {
		err := ValidatePolygon(cwSquare, WithRequireCW(true))
		if err != nil {
			t.Errorf("CW polygon should pass CW requirement: %v", err)
		}
	})

	t.Run("CCW polygon fails CW requirement", func(t *testing.T) {
		err := ValidatePolygon(ccwSquare, WithRequireCW(true))
		if err == nil {
			t.Error("CCW polygon should fail CW requirement")
		}
	})
}

// TestValidatePolygon_TooFewVertices tests minimum vertex count
func TestValidatePolygon_TooFewVertices(t *testing.T) {
	tests := []struct {
		name      string
		polygon   []types.Point
		shouldErr bool
	}{
		{
			name:      "no vertices",
			polygon:   []types.Point{},
			shouldErr: true,
		},
		{
			name: "one vertex",
			polygon: []types.Point{
				{X: 0, Y: 0},
			},
			shouldErr: true,
		},
		{
			name: "two vertices",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
			},
			shouldErr: true,
		},
		{
			name: "three vertices (valid)",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 5, Y: 10},
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePolygon(tt.polygon)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidatePolygon() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// TestValidatePolygonDetailed tests detailed validation results
func TestValidatePolygonDetailed(t *testing.T) {
	polygon := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	result := ValidatePolygonDetailed(polygon)

	if !result.Valid {
		t.Error("Expected polygon to be valid")
	}

	if result.VertexCount != 4 {
		t.Errorf("VertexCount = %d, want 4", result.VertexCount)
	}

	if math.Abs(result.Area-100) > 1e-9 {
		t.Errorf("Area = %f, want 100", result.Area)
	}

	if math.Abs(result.Width-10) > 1e-9 {
		t.Errorf("Width = %f, want 10", result.Width)
	}

	if math.Abs(result.Height-10) > 1e-9 {
		t.Errorf("Height = %f, want 10", result.Height)
	}

	if !result.IsCCW {
		t.Error("Expected CCW winding")
	}

	if result.SelfIntersects {
		t.Error("Expected no self-intersection")
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
}

// TestPolygonArea tests area calculation
func TestPolygonArea(t *testing.T) {
	tests := []struct {
		name     string
		polygon  []types.Point
		expected float64
	}{
		{
			name: "unit square CCW",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 0, Y: 1},
			},
			expected: 1.0,
		},
		{
			name: "unit square CW",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 1, Y: 0},
			},
			expected: -1.0,
		},
		{
			name: "10x10 square",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			expected: 100.0,
		},
		{
			name: "triangle",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 5, Y: 10},
			},
			expected: 50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area := predicates.PolygonArea(tt.polygon)
			if math.Abs(area-tt.expected) > 1e-9 {
				t.Errorf("PolygonArea() = %f, want %f", area, tt.expected)
			}
		})
	}
}

// TestPolygonBounds tests bounding box calculation
func TestPolygonBounds(t *testing.T) {
	polygon := []types.Point{
		{X: 5, Y: 3},
		{X: 15, Y: 7},
		{X: 12, Y: 18},
		{X: 2, Y: 10},
	}

	bounds := predicates.PolygonBounds(polygon)

	if bounds.Min.X != 2 || bounds.Min.Y != 3 {
		t.Errorf("Min = (%f, %f), want (2, 3)", bounds.Min.X, bounds.Min.Y)
	}

	if bounds.Max.X != 15 || bounds.Max.Y != 18 {
		t.Errorf("Max = (%f, %f), want (15, 18)", bounds.Max.X, bounds.Max.Y)
	}
}

// TestPolygonIsValid tests quick validity check
func TestPolygonIsValid(t *testing.T) {
	tests := []struct {
		name     string
		polygon  []types.Point
		expected bool
	}{
		{
			name: "valid square",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 10, Y: 10},
				{X: 0, Y: 10},
			},
			expected: true,
		},
		{
			name: "self-intersecting",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
				{X: 0, Y: 10},
				{X: 10, Y: 10},
			},
			expected: false,
		},
		{
			name: "too few vertices",
			polygon: []types.Point{
				{X: 0, Y: 0},
				{X: 10, Y: 0},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PolygonIsValid(tt.polygon, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonIsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestPointInPolygon tests point-in-polygon detection
func TestPointInPolygon(t *testing.T) {
	polygon := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	tests := []struct {
		name     string
		point    types.Point
		expected bool
	}{
		{
			name:     "point inside",
			point:    types.Point{X: 5, Y: 5},
			expected: true,
		},
		{
			name:     "point outside",
			point:    types.Point{X: 15, Y: 15},
			expected: false,
		},
		{
			name:     "point on edge",
			point:    types.Point{X: 0, Y: 5},
			expected: true,
		},
		{
			name:     "point at vertex",
			point:    types.Point{X: 0, Y: 0},
			expected: true,
		},
		{
			name:     "point barely inside",
			point:    types.Point{X: 0.1, Y: 0.1},
			expected: true,
		},
		{
			name:     "point barely outside",
			point:    types.Point{X: 10.1, Y: 5},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := predicates.PointInPolygonRayCast(tt.point, polygon, 1e-9)
			if result != tt.expected {
				t.Errorf("PointInPolygonRayCast(%v) = %v, want %v", tt.point, result, tt.expected)
			}
		})
	}
}

// TestPolygonContains is a user-facing wrapper for point containment
func TestPolygonContains(t *testing.T) {
	polygon := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	tests := []struct {
		name     string
		point    types.Point
		expected bool
	}{
		{
			name:     "center point",
			point:    types.Point{X: 5, Y: 5},
			expected: true,
		},
		{
			name:     "corner",
			point:    types.Point{X: 0, Y: 0},
			expected: true,
		},
		{
			name:     "outside",
			point:    types.Point{X: 20, Y: 20},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PolygonContains(polygon, tt.point, 1e-9)
			if result != tt.expected {
				t.Errorf("PolygonContains(%v) = %v, want %v", tt.point, result, tt.expected)
			}
		})
	}
}

// TestMultipleConstraints tests combining multiple validation options
func TestMultipleConstraints(t *testing.T) {
	polygon := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	// Should pass all constraints
	err := ValidatePolygon(polygon,
		WithPolygonMinArea(50),
		WithPolygonMaxArea(200),
		WithPolygonMinWidth(5),
		WithPolygonMaxWidth(15),
		WithPolygonMinHeight(5),
		WithPolygonMaxHeight(15),
		WithRequireCCW(true),
	)

	if err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}

	// Should fail on area
	err = ValidatePolygon(polygon,
		WithPolygonMinArea(200),
	)

	if err == nil {
		t.Error("Expected validation to fail on minimum area")
	}
}

// TestPolygonValidationResultString tests the String method
func TestPolygonValidationResultString(t *testing.T) {
	// Valid polygon
	validPoly := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 10, Y: 10},
		{X: 0, Y: 10},
	}

	result := ValidatePolygonDetailed(validPoly)
	str := result.String()

	// Should contain key information
	if !strings.Contains(str, "vertices=4") {
		t.Errorf("String should contain vertex count, got: %s", str)
	}
	if !strings.Contains(str, "area=100") {
		t.Errorf("String should contain area, got: %s", str)
	}
	if !strings.Contains(str, "winding=CCW") {
		t.Errorf("String should contain winding direction, got: %s", str)
	}

	// Invalid polygon (too small area)
	result = ValidatePolygonDetailed(validPoly, WithPolygonMinArea(200))
	str = result.String()

	if !strings.Contains(str, "area") && !strings.Contains(str, "minimum") {
		t.Errorf("String should contain error message about area, got: %s", str)
	}

	// Self-intersecting polygon
	bowtie := []types.Point{
		{X: 0, Y: 0},
		{X: 10, Y: 0},
		{X: 0, Y: 10},
		{X: 10, Y: 10},
	}

	result = ValidatePolygonDetailed(bowtie)
	str = result.String()

	if !strings.Contains(str, "self-intersects=true") {
		t.Errorf("String should indicate self-intersection, got: %s", str)
	}

	t.Logf("Valid polygon: %s", ValidatePolygonDetailed(validPoly).String())
	t.Logf("Invalid area: %s", ValidatePolygonDetailed(validPoly, WithPolygonMinArea(200)).String())
	t.Logf("Self-intersecting: %s", ValidatePolygonDetailed(bowtie).String())
}
