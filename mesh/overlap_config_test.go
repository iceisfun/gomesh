package mesh

import (
	"testing"

	"github.com/iceisfun/gomesh/types"
)

// TestOpposingWindingWithDifferentOptions demonstrates the difference between
// WithOverlapTriangle and WithTriangleOverlapCheck
func TestOpposingWindingWithDifferentOptions(t *testing.T) {
	testCases := []struct {
		name                  string
		allowOverlap          bool // WithOverlapTriangle(allow) - controls errorOnDuplicateTriangle
		checkGeometricOverlap bool // WithTriangleOverlapCheck(enable) - controls validateTriangleOverlapArea
		expectError           bool
		expectedError         string
	}{
		{
			name:                  "Default - allows duplicates",
			allowOverlap:          true, // ALLOW duplicates (default)
			checkGeometricOverlap: false,
			expectError:           false,
		},
		{
			name:                  "Duplicate check enabled - rejects duplicates",
			allowOverlap:          false, // REJECT duplicates
			checkGeometricOverlap: false,
			expectError:           true,
			expectedError:         "duplicate triangle",
		},
		{
			name:                  "Geometric overlap check - rejects overlaps",
			allowOverlap:          true, // ALLOW duplicates
			checkGeometricOverlap: true, // But check geometric overlaps
			expectError:           true,
			expectedError:         "overlap", // Match either form of error
		},
		{
			name:                  "Both enabled - duplicate check wins (runs first)",
			allowOverlap:          false, // REJECT duplicates
			checkGeometricOverlap: true,
			expectError:           true,
			expectedError:         "duplicate triangle", // Caught first
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMesh(
				WithEpsilon(1e-9),
				WithOverlapTriangle(tc.allowOverlap),
				WithTriangleOverlapCheck(tc.checkGeometricOverlap),
			)

			t.Logf("Config: epsilon=%.10f, validateTriangleOverlapArea=%v", m.cfg.epsilon, m.cfg.validateTriangleOverlapArea)

			// Add vertices
			v0, _ := m.AddVertex(types.Point{X: 120.00, Y: 65.00})
			v1, _ := m.AddVertex(types.Point{X: 122.00, Y: 64.00})
			v2, _ := m.AddVertex(types.Point{X: 121.00, Y: 64.00})

			// Add first triangle
			err := m.AddTriangle(v0, v1, v2)
			if err != nil {
				t.Fatalf("Failed to add first triangle: %v", err)
			}

			// Add second triangle with opposite winding (same 3 vertices, reversed order)
			err = m.AddTriangle(v2, v1, v0)

			// Debug: check what the mesh thinks about overlaps
			if tc.checkGeometricOverlap && err == nil {
				overlaps := m.FindOverlappingTriangles()
				t.Logf("DEBUG: Mesh has %d triangles after adding both", m.NumTriangles())
				t.Logf("DEBUG: FindOverlappingTriangles() found %d overlaps", len(overlaps))
				if len(overlaps) > 0 {
					for i, overlap := range overlaps {
						t.Logf("DEBUG: Overlap %d: area=%.4f, type=%s", i+1, overlap.IntersectionArea, overlap.Type)
					}
				}
			}

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', but got none", tc.expectedError)
				} else {
					errStr := err.Error()
					t.Logf("✓ Correctly rejected: %v", err)
					// Verify it's the expected error type
					if tc.expectedError != "" {
						if !contains(errStr, tc.expectedError) {
							t.Errorf("Expected error containing '%s', got: %s", tc.expectedError, errStr)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				} else {
					t.Logf("✓ Correctly allowed duplicate")
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
